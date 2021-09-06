package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/b177y/starship/nebutils"
	"github.com/b177y/starship/wormhole"
	"github.com/b177y/starship/xpaseto"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/slackhq/nebula/cert"
)

type StaticHost struct {
	NebulaAddress string
	Endpoint      []string
}

type NodeConfigSchema struct {
	Address      string       `json:"address"` // do i need this?
	Lighthouses  []string     `json:"lighthouses"`
	AmLighthouse bool         `json:"am_lighthouse"`
	StaticHosts  []StaticHost `json:"static_hosts"`
	ListenPort   int          `json:"listen_port"`

	FirewallInbound  []FirewallRule `json:"firewall_inbound"`
	FirewallOutbound []FirewallRule `json:"firewall_outbound"`
	Cipher           string         `json:"cipher"`
	Cert             string         `json:"cert"`
}

func (s *server) CheckNodeIdentity(netname string,
	nodename string, token string) (err error) {
	log.Printf("Checking identity for node %s (network %s).\n", nodename,
		netname)
	// get node pubkey
	pubkey, err := s.db.getNodePubkey(netname, nodename)
	if err != nil {
		return err
	}
	log.Println("Got node pubkey: ", pubkey)
	signer := xpaseto.NewSigner([]byte{0}, pubkey)
	jsonToken, err := signer.ParsePaseto(token)
	if err != nil {
		return err
	}
	nodeIdentity := *new(wormhole.NodeIdentitySchema)
	err = wormhole.SchemaFromJSONToken(jsonToken, &nodeIdentity)
	if err != nil {
		return err
	}
	if nodeIdentity.Netname != netname || nodeIdentity.Nodename != nodename {
		log.Error("Node identity does not match url params")
		return fmt.Errorf("Node Identity does not match url params.")
	}
	return nil
}

func (s *server) SignPayload(netname string,
	payload interface{}) (signed string, err error) {
	jsonToken, err := wormhole.NewToken(payload)
	if err != nil {
		return "", err
	}
	edprivkey, _, err := s.db.getNetworkCA(netname)
	if err != nil {
		return "", err
	}
	privkey := nebutils.PrivateKeyToCurve25519(edprivkey[:32])
	signer := xpaseto.NewSigner(privkey, []byte{})
	token, err := signer.SignPaseto(jsonToken)
	if err != nil {
		return "", err
	}
	return token, nil
}

// /api/neutron/config?net=NETWORK&node=NODE [GET]
func (s *server) handleGetConfig() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		urlParams := r.URL.Query()
		netname := urlParams["net"][0]
		nodename := urlParams["node"][0]
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Bad Request", 400)
			return
		}
		token := string(b)
		err = s.CheckNodeIdentity(netname, nodename, token)
		if err != nil {
			http.Error(w, err.Error(), 401)
			log.Error(err)
			return
		}
		status, err := s.db.getNodeStatus(netname, nodename)
		if err != nil {
			http.Error(w, "Could not get node status: "+err.Error(), 500)
			log.Error("Could not get node status: ", err)
			return
		}
		if status != "active" {
			http.Error(w, "425 - Node is not active. You may need to have it approved.", 425)
			log.Error("Node is not active, not returning cert.")
			return
		}
		log.Printf("Getting cert for node %s in network %s.\n", nodename, netname)
		err = s.signNodeCert(netname, nodename)
		if err != nil {
			http.Error(w, "Could not sign certificate", 503)
			return
		}
		nodeCert, err := s.db.getCert(netname, nodename)
		if err != nil || nodeCert == nil {
			http.Error(w, "Could not get certificate", 503)
			return
		}
		nc, err := cert.UnmarshalNebulaCertificate(nodeCert)
		if err != nil {
			http.Error(w, "Could not decode CA Certificate", 500)
			return
		}
		pemcert, err := nc.MarshalToPEM()
		if err != nil {
			http.Error(w, "Could not marshal CA certificate to PEM", 500)
			return
		}
		node, err := s.db.getNodeConfig(netname, nodename)
		node.Cert = string(pemcert)
		err = s.db.updateLatestFetch(netname, nodename, time.Now().Format(time.RFC3339))
		signedResponse, err := s.SignPayload(netname, node)
		if err != nil {
			log.Error(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Fprintf(w, signedResponse)
	}
}

// /api/neutron/join [POST]
func (s *server) handleJoinNetwork() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Join network requested")
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Bad Request", 400)
			return
		}
		token := string(b)
		signer := xpaseto.NewSigner([]byte{0}, []byte{0})
		jsonToken, err := signer.ParseSelfSigned(token)
		if err != nil {
			http.Error(w, "Invalid Signature: "+err.Error(), 503)
			return
		}
		joinReq := *new(wormhole.RequestJoinSchema)
		err = wormhole.SchemaFromJSONToken(jsonToken, &joinReq)
		if err != nil {
			http.Error(w, "Internal Server Error", 500)
			log.Error(err)
			return
		}
		pubkey, _, err := cert.UnmarshalX25519PublicKey([]byte(joinReq.PubKey))
		if err != nil {
			http.Error(w, "Internal Server Error", 500)
			log.Error(err)
			return
		}
		address, err := s.db.newAddress(joinReq.Netname)
		if err != nil {
			http.Error(w, "Internal Server Error", 500)
			log.Error(err)
			return
		}
		err = s.db.addJoinRequest(joinReq.Netname,
			joinReq.Nodename,
			joinReq.Hostname,
			address,
			pubkey,
		)
		if err != nil {
			http.Error(w, "Internal Server Error: "+err.Error(), 500)
			log.Error(err)
			return
		}
		fmt.Fprintf(w, "SUCCESS")
	}
}

// /api/neutron/leave?net=NETWORK&node=NODE [POST]
func (s *server) handleLeaveNetwork() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		nodename := vars["NODENAME"]
		log.Printf("/api/neutron/%s/leave requested.\n", nodename)
	}
}
