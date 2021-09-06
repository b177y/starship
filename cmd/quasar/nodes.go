// This contains endpoints for managing nodes through the API
// using standard auth (not XPASETO auth)

package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/slackhq/nebula/cert"
)

type NodeOverviewSchema struct {
	Nodename    string `json:"name"`
	Hostname    string `json:"hostname"`
	LatestFetch string `json:"latest_fetch"`
	Status      string `json:"status"`
	Address     string `json:"address"`
	PubKey      string `json:"pubkey"`
}

type NodeSchema struct {
	Nodename         string         `json:"name"`
	Hostname         string         `json:"hostname"`
	Status           string         `json:"status"`
	Address          string         `json:"address"`
	StaticAddress    string         `json:"static_address"`
	ListenPort       int            `json:"listen_port"`
	Lighthouse       bool           `json:"is_lighthouse"`
	Groups           []string       `json:"groups"`
	FirewallOutbound []FirewallRule `json:"firewall_outbound"`
	FirewallInbound  []FirewallRule `json:"firewall_inbound"`
}

// /api/networks/{NETWORK}/nodes/all [GET]
func (s *server) handleGetAllNodes() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		net := vars["NETWORK"]
		log.Println("Getting nodes in network", net)
		nodes, err := s.db.allNodes(net)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := json.NewEncoder(w).Encode(nodes); err != nil {
			log.Error(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
}

func (s *server) signNodeCert(netname string, nodename string) error {
	log.Printf("Signing cert for node %s in network %s.", nodename, netname)
	pubkey, err := s.db.getNodePubkey(netname, nodename)
	if err != nil {
		return err
	}
	node, err := s.db.getNodeInfo(netname, nodename)
	if err != nil {
		return err
	}
	network, err := s.db.networkInfo(netname)
	if err != nil {
		return err
	}
	quasarPrivKey, certBytes, err := s.db.getNetworkCA(netname)
	log.Println("Trying to unmarshal cert")
	quasarCert, _, err := cert.UnmarshalNebulaCertificateFromPEM(certBytes)
	if err != nil {
		return err
	}
	log.Println("Trying to get issuer")
	issuer, err := quasarCert.Sha256Sum()
	if err != nil {
		return err
	}
	log.Println("cidr stuffs.")
	ip, cidr, err := net.ParseCIDR(network.Cidr)
	ip = net.ParseIP(node.Address)
	if err != nil {
		return err
	}
	cidr.IP = ip
	subnet := cidr
	exp := time.Until(quasarCert.Details.NotAfter) - time.Second*1
	nc := cert.NebulaCertificate{
		Details: cert.NebulaCertificateDetails{
			Name:      nodename,
			Ips:       []*net.IPNet{cidr},
			Groups:    node.Groups,
			Subnets:   []*net.IPNet{subnet},
			NotBefore: time.Now(),
			NotAfter:  time.Now().Add(exp),
			PublicKey: pubkey,
			IsCA:      false,
			Issuer:    issuer,
		},
	}
	err = nc.Sign(quasarPrivKey)
	if err != nil {
		return err
	}
	signedCertBytes, err := nc.Marshal()
	if err != nil {
		return err
	}
	err = s.db.saveNodeCert(netname, nodename, signedCertBytes)
	if err != nil {
		return err
	}
	pemBytes, _ := nc.MarshalToPEM()
	log.Println("Saved cert", string(pemBytes))
	return nil
}

// /api/networks/{NETWORK}/nodes/{NODENAME}/approve [POST]
func (s *server) handleApproveNode() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		net := vars["NETWORK"]
		node := vars["NODENAME"]
		log.Printf("Approving node %s in network %s.\n", node, net)
		// sign pubkey and create cert
		err := s.signNodeCert(net, node)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// update status to active
		err = s.db.updateNodeStatus(net, node, "active")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Fprintf(w, "SUCCESS")
	}
}

// /api/networks/{NETWORK}/nodes/{NODENAME}/update [POST]
func (s *server) handleUpdateNode() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		netname := vars["NETWORK"]
		nodename := vars["NODENAME"]
		log.Printf("Updating node %s in network %s.\n", nodename, netname)
		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()
		var node NodeSchema
		err := dec.Decode(&node)
		if err != nil {
			log.Error("Could not decode json", err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Println("UPDATING NODE WITH", node)
		err = s.db.updateNodeInfo(netname, nodename, node)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Fprintf(w, "SUCCESS")
	}
}

// /api/networks/{NETWORK}/nodes/{NODENAME}/info [GET]
func (s *server) handleNodeInfo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		netname := vars["NETWORK"]
		nodename := vars["NODENAME"]
		node, err := s.db.getNodeInfo(netname, nodename)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := json.NewEncoder(w).Encode(node); err != nil {
			log.Error(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
}

// /api/networks/{NETWORK}/nodes/{NODENAME}/disable [POST]
func (s *server) handleDisableNode() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		net := vars["NETWORK"]
		node := vars["NODENAME"]
		log.Printf("Disabling node %s in network %s.\n", node, net)
		err := s.db.updateNodeStatus(net, node, "disabled")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Fprintf(w, "SUCCESS")
	}
}
