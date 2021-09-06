// HTTP Endpoints Relating to Networks
// These should be used by a client using standard auth
// rather than XPASETO auth with a nebula key

package main

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/slackhq/nebula/cert"
)

type NetSchema struct {
	Name          string   `json:"name"`
	Cidr          string   `json:"cidr"`
	Cipher        string   `json:"cipher"`
	Groups        []string `json:"groups"`
	CaFingerprint string   `json:"ca_fingerprint"`
}

type NetOverviewSchema struct {
	Name string `json:"name"`
	Cidr string `json:"cidr"`
}

// /api/networks/all [GET]
func (s *server) handleGetAllNetworks() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Getting all networks")
		networks, err := s.db.allNetworks()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Println("Got all networks", networks)
		if err := json.NewEncoder(w).Encode(networks); err != nil {
			log.Error(err)
		}
	}
}

type NewNetSchema struct {
	Name string `json:"name"`
	Cidr string `json:"cidr"`
}

// /api/networks/new [POST]
func (s *server) handleNewNetwork() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get name and cidr from request body
		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()
		var newnet NewNetSchema
		err := dec.Decode(&newnet)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Println("Creating new network: ", newnet)

		// generate keys
		pubkey, privkey, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			log.Error("Could not generate keys: " + err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// generate and self-sign cert for ca key
		ip, cidr, err := net.ParseCIDR(newnet.Cidr)
		if err != nil {
			log.Error("Invalid cidr definition: " + newnet.Cidr)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		cidr.IP = ip
		subnet := cidr
		nc := cert.NebulaCertificate{
			Details: cert.NebulaCertificateDetails{
				Name:      "quasar" + newnet.Name,
				Ips:       []*net.IPNet{cidr},
				Groups:    []string{},
				Subnets:   []*net.IPNet{subnet},
				NotBefore: time.Now(),
				NotAfter:  time.Now().Add(time.Duration(time.Hour * 2190)),
				PublicKey: pubkey,
				IsCA:      true,
			},
		}

		err = nc.Sign(privkey)
		if err != nil {
			log.Error("Error while signing ca key: %s", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		certbytes, err := nc.MarshalToPEM()

		// write new network to database
		log.Println("ADDING NETWORK TO DB, priv: ", privkey)
		err = s.db.addNetwork(certbytes,
			privkey,
			newnet.Name,
			newnet.Cidr,
			"chachapoly",
		)
		if err != nil {
			log.Error(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Println("Added network to db: ", newnet.Name, newnet.Cidr)
		fmt.Fprintf(w, "SUCCESS")
	}
}

// /api/networks/{NETWORK}/update [POST]
func (s *server) handleUpdateNetwork() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		netname := vars["NETWORK"]
		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()
		var network NetSchema
		err := dec.Decode(&network)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		cidr := network.Cidr
		if cidr != "" {
			_, ncidr, err := net.ParseCIDR(cidr)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			cidr = ncidr.String()
		}
		err = s.db.updateNetwork(netname, cidr, network.Cipher, network.Groups)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Fprintf(w, "SUCCESS")
	}
}

// /api/networks/{NETWORK}/delete [POST]
func (s *server) handleDeleteNetwork() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		net := vars["NETWORK"]
		log.Println("Deleting network", net)
		err := s.db.deleteNetwork(net)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Fprintf(w, "SUCCESS")
	}
}

// /api/networks/{NETWORK}/info [GET]
func (s *server) handleNetworkInfo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		net := vars["NETWORK"]
		log.Println("Getting network info for", net)
		network, err := s.db.networkInfo(net)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		_, certBytes, err := s.db.getNetworkCA(net)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		cert, _, err := cert.UnmarshalNebulaCertificateFromPEM(certBytes)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fingerprint, err := cert.Sha256Sum()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		network.CaFingerprint = fingerprint
		if err := json.NewEncoder(w).Encode(network); err != nil {
			log.Error(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
}

// /api/networks/{NETWORK}/cert [GET]
func (s *server) handleGetNetworkCert() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// unprotected route - anyone with network name can get ca cert
		// returns 404 if no network or ca cert file
		vars := mux.Vars(r)
		net := vars["NETWORK"]
		log.Printf("/api/networks/%s/cert requested.\n", net)
		cacert, err := s.db.getCert(net, "")
		log.Println("GOT CERT")
		if err != nil || string(cacert) == "" {
			switch errmsg := err.Error(); errmsg {
			case "NONETWORK":
				http.Error(w, "Network does not exist.", 404)
			default:
				http.Error(w, "Internal Server Error.", 500)
			}
			return
		}
		nc, _, err := cert.UnmarshalNebulaCertificateFromPEM(cacert)
		if err != nil {
			http.Error(w, "Could not decode CA Certificate", 500)
		}
		pemcert, err := nc.MarshalToPEM()
		if err != nil {
			http.Error(w, "Could not marshal CA certificate to PEM", 500)
		}
		fmt.Fprintf(w, string(pemcert))
	}
}
