package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/boltdb/bolt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type database interface {
	connect(filepath string) error
	addNetwork(
		cacert []byte,
		capriv []byte,
		name string,
		cidr string,
		cipher string,
	) error
	getCert(network, host string) ([]byte, error)
	updateLatestFetch(netname, nodename, timestamp string) error
	addJoinRequest(netname string,
		nodename string,
		hostname string,
		address string,
		pubkey []byte,
	) error
	allNetworks() (
		networks []NetOverviewSchema,
		err error,
	)
	deleteNetwork(netname string) error
	networkInfo(netname string) (
		network NetSchema,
		err error,
	)
	updateNetwork(netname, cidr, cipher string,
		groups []string) error
	allNodes(netname string) (
		nodes []NodeOverviewSchema,
		err error,
	)
	updateNodeStatus(netname string,
		nodename string,
		status string) error
	updateNodeInfo(netname string,
		nodename string,
		node NodeSchema) error
	getNodeStatus(netname string,
		nodename string) (
		status string,
		err error,
	)
	getNodeInfo(netname string, nodename string) (nodeinfo NodeSchema,
		err error)
	getNodeConfig(netname, nodename string) (config NodeConfigSchema,
		err error)
	getNodePubkey(netname string,
		nodename string) (pubkey []byte,
		err error)
	getNetworkCA(netname string) (privkey []byte,
		cert []byte,
		err error)
	saveNodeCert(netname string,
		nodename string,
		cert []byte) error
	newAddress(netname string) (address string,
		err error)
}

// boltdb interface
type boltdbi struct {
	db *bolt.DB
}

func (b *boltdbi) addNetwork(cacert []byte,
	capriv []byte,
	name string,
	cidr string,
	cipher string,
) (err error) {
	log.Info("Bolt adding network: ", name)
	err = b.db.Update(func(tx *bolt.Tx) error {
		nb, err := tx.CreateBucket([]byte(name))
		if err != nil {
			return err
		}

		err = nb.Put([]byte("NET_NAME"), []byte(name))
		if err != nil {
			return err
		}

		err = nb.Put([]byte("CA_PRIV_KEY"), capriv)
		if err != nil {
			return err
		}

		err = nb.Put([]byte("CA_CERT"), cacert)
		if err != nil {
			return err
		}

		err = nb.Put([]byte("CIDR"), []byte(cidr))
		if err != nil {
			return err
		}

		groups, err := json.Marshal([]string{})
		if err != nil {
			return err
		}
		err = nb.Put([]byte("GROUPS"), groups)
		if err != nil {
			return err
		}

		err = nb.Put([]byte("CIPHER"), []byte(cipher))
		if err != nil {
			return err
		}

		return nil

	})
	return err
}

func (b *boltdbi) connect(filepath string) (err error) {
	log.Info("Bolt connecting ", filepath)
	b.db, err = bolt.Open(filepath,
		0600,
		&bolt.Options{Timeout: 3 * time.Second},
	)
	return nil
}

func (b *boltdbi) getCert(network, host string) (cert []byte,
	err error) {
	err = b.db.View(func(tx *bolt.Tx) error {
		bkt := tx.Bucket([]byte(network))
		if bkt == nil {
			return fmt.Errorf("NONETWORK")
		}
		if host != "" {
			nodeBkt := bkt.Bucket([]byte(host))
			if nodeBkt == nil {
				return fmt.Errorf("NOHOST")
			}
			cert = nodeBkt.Get([]byte("cert"))
			return nil
		}
		cert = bkt.Get([]byte("CA_CERT"))
		return nil
	})
	return cert, err
}

func (b *boltdbi) addJoinRequest(netname string,
	nodename string,
	hostname string,
	address string,
	pubkey []byte,
) (err error) {
	err = b.db.Update(func(tx *bolt.Tx) error {
		netBkt := tx.Bucket([]byte(netname))
		if netBkt == nil {
			return errors.Errorf("Network does not exist.")
		}
		nodeBkt, err := netBkt.CreateBucket([]byte(nodename))
		if err != nil {
			return errors.Errorf("Node exists in network.")
		}
		err = nodeBkt.Put([]byte("hostname"), []byte(hostname))
		if err != nil {
			return err
		}
		err = nodeBkt.Put([]byte("address"), []byte(address))
		if err != nil {
			return err
		}
		err = nodeBkt.Put([]byte("latest_fetch"), []byte("NEVER"))
		if err != nil {
			return err
		}
		err = nodeBkt.Put([]byte("pubkey"), pubkey)
		if err != nil {
			return err
		}
		err = nodeBkt.Put([]byte("status"), []byte("pending"))
		if err != nil {
			return err
		}
		err = nodeBkt.Put([]byte("listen_port"), []byte("0"))
		if err != nil {
			return err
		}
		groups, err := json.Marshal([]string{})
		if err != nil {
			return err
		}
		err = nodeBkt.Put([]byte("groups"), groups)
		if err != nil {
			return err
		}
		err = nodeBkt.Put([]byte("is_lighthouse"), []byte("false"))
		if err != nil {
			return err
		}
		err = nodeBkt.Put([]byte("static_address"), []byte(""))
		if err != nil {
			return err
		}
		inbound, outbound := defaultRules()
		log.Println("Converting to bytes", inbound, outbound)
		inboundBytes, err := json.Marshal(inbound)
		if err != nil {
			return err
		}
		outboundBytes, err := json.Marshal(outbound)
		if err != nil {
			return err
		}
		err = json.Unmarshal(outboundBytes, &outbound)
		if err != nil {
			return err
		}
		err = nodeBkt.Put([]byte("firewall_outbound"), outboundBytes)
		if err != nil {
			return err
		}
		err = nodeBkt.Put([]byte("firewall_inbound"), inboundBytes)
		return err
	})
	return err
}

func (b *boltdbi) allNetworks() (networks []NetOverviewSchema, err error) {
	networks = make([]NetOverviewSchema, 0)
	err = b.db.View(func(tx *bolt.Tx) error {
		tx.ForEach(func(name []byte, b *bolt.Bucket) error {
			netcidr := b.Get([]byte("CIDR"))
			n := NetOverviewSchema{
				Name: string(name),
				Cidr: string(netcidr),
			}
			networks = append(networks, n)
			return nil
		})
		return nil
	})
	return networks, err
}

func (b *boltdbi) deleteNetwork(netname string) (err error) {
	err = b.db.Update(func(tx *bolt.Tx) error {
		return tx.DeleteBucket([]byte(netname))
	})
	return err
}

func (b *boltdbi) networkInfo(netname string) (network NetSchema,
	err error) {
	err = b.db.View(func(tx *bolt.Tx) error {
		bkt := tx.Bucket([]byte(netname))
		if bkt == nil {
			return fmt.Errorf("NONETWORK")
		}
		netcidr := bkt.Get([]byte("CIDR"))
		cipher := bkt.Get([]byte("CIPHER"))
		groupsBytes := bkt.Get([]byte("GROUPS"))
		var groups []string
		err := json.Unmarshal(groupsBytes, &groups)
		if err != nil {
			return err
		}
		network = NetSchema{
			Name:   netname,
			Cidr:   string(netcidr),
			Cipher: string(cipher),
			Groups: groups,
		}
		return nil
	})
	return network, err
}

func (b *boltdbi) allNodes(netname string) (nodes []NodeOverviewSchema,
	err error) {
	nodes = make([]NodeOverviewSchema, 0)
	err = b.db.View(func(tx *bolt.Tx) error {
		bkt := tx.Bucket([]byte(netname))
		if bkt == nil {
			return fmt.Errorf("NONETWORK")
		}
		err := bkt.ForEach(func(key, val []byte) error {
			if val == nil {
				// keyval is bucket so is node
				nb := bkt.Bucket(key)
				hostName := nb.Get([]byte("hostname"))
				status := nb.Get([]byte("status"))
				address := nb.Get([]byte("address"))
				pubkeyBytes := nb.Get([]byte("pubkey"))
				pubkey := base64.StdEncoding.EncodeToString([]byte(pubkeyBytes))
				latest_fetch := nb.Get([]byte("latest_fetch"))
				node := NodeOverviewSchema{
					Nodename:    string(key),
					Hostname:    string(hostName),
					Status:      string(status),
					Address:     string(address),
					LatestFetch: string(latest_fetch),
					PubKey:      pubkey,
				}
				nodes = append(nodes, node)
			}
			return err
		})
		return err
	})
	return nodes, err

}

func (b *boltdbi) updateNodeStatus(netname string,
	nodename string,
	status string,
) (err error) {
	err = b.db.Update(func(tx *bolt.Tx) error {
		netBkt := tx.Bucket([]byte(netname))
		if netBkt == nil {
			return errors.Errorf("Network does not exist.")
		}
		nodeBkt := netBkt.Bucket([]byte(nodename))
		if nodeBkt == nil {
			return errors.Errorf("Node does not exist in network.")
		}
		err = nodeBkt.Put([]byte("status"), []byte(status))
		return err
	})
	return err

}

func (b *boltdbi) updateNodeInfo(netname string,
	nodename string,
	node NodeSchema,
) (err error) {
	err = b.db.Update(func(tx *bolt.Tx) error {
		netBkt := tx.Bucket([]byte(netname))
		if netBkt == nil {
			return errors.Errorf("Network does not exist.")
		}
		nodeBkt := netBkt.Bucket([]byte(nodename))
		if nodeBkt == nil {
			return errors.Errorf("Node does not exist in network.")
		}
		inboundBytes, err := json.Marshal(node.FirewallInbound)
		if err != nil {
			return err
		}
		outboundBytes, err := json.Marshal(node.FirewallOutbound)
		if err != nil {
			return err
		}
		err = nodeBkt.Put([]byte("firewall_outbound"), outboundBytes)
		if err != nil {
			return err
		}
		err = nodeBkt.Put([]byte("firewall_inbound"), inboundBytes)
		if err != nil {
			return err
		}
		if node.StaticAddress != "" {
			err = nodeBkt.Put([]byte("static_address"), []byte(node.StaticAddress))
			if err != nil {
				return err
			}
		}
		err = nodeBkt.Put([]byte("is_lighthouse"), []byte(strconv.FormatBool(node.Lighthouse)))
		if err != nil {
			return err
		}
		err = nodeBkt.Put([]byte("listen_port"), []byte(fmt.Sprint(node.ListenPort)))
		if err != nil {
			return err
		}

		log.Println("Adding groups to db:", node.Groups)
		groupsBytes, err := json.Marshal(node.Groups)
		if err != nil {
			return err
		}
		err = nodeBkt.Put([]byte("groups"), groupsBytes)
		return err
	})
	return err

}

func (b *boltdbi) getNodeStatus(netname string,
	nodename string) (status string,
	err error) {
	err = b.db.View(func(tx *bolt.Tx) error {
		netBkt := tx.Bucket([]byte(netname))
		if netBkt == nil {
			return errors.Errorf("Network does not exist.")
		}
		nodeBkt := netBkt.Bucket([]byte(nodename))
		if nodeBkt == nil {
			return errors.Errorf("Node does not exist in network.")
		}
		statusBytes := nodeBkt.Get([]byte("status"))
		status = string(statusBytes)
		return err
	})
	return status, err
}

func (b *boltdbi) getNodeInfo(netname string,
	nodename string) (node NodeSchema,
	err error) {
	err = b.db.View(func(tx *bolt.Tx) error {
		netBkt := tx.Bucket([]byte(netname))
		if netBkt == nil {
			return errors.Errorf("Network does not exist.")
		}
		nodeBkt := netBkt.Bucket([]byte(nodename))
		if nodeBkt == nil {
			return errors.Errorf("Node does not exist in network.")
		}
		statusBytes := nodeBkt.Get([]byte("status"))
		hostnameBytes := nodeBkt.Get([]byte("hostname"))
		addressBytes := nodeBkt.Get([]byte("address"))
		listenPortBytes := nodeBkt.Get([]byte("listen_port"))
		listenPort, err := strconv.Atoi(string(listenPortBytes))
		if err != nil {
			return err
		}
		staticAddressBytes := nodeBkt.Get([]byte("static_address"))
		groupsBytes := nodeBkt.Get([]byte("groups"))
		var groups []string
		err = json.Unmarshal(groupsBytes, &groups)
		if err != nil {
			return err
		}
		inboundBytes := nodeBkt.Get([]byte("firewall_inbound"))
		outboundBytes := nodeBkt.Get([]byte("firewall_outbound"))
		var inbound []FirewallRule
		var outbound []FirewallRule
		err = json.Unmarshal(inboundBytes, &inbound)
		if err != nil {
			return err
		}
		err = json.Unmarshal(outboundBytes, &outbound)
		if err != nil {
			return err
		}
		lighthouseBytes := nodeBkt.Get([]byte("is_lighthouse"))
		is_lighthouse, err := strconv.ParseBool(string(lighthouseBytes))
		if err != nil {
			return err
		}
		node = NodeSchema{
			Nodename:         nodename,
			Hostname:         string(hostnameBytes),
			Status:           string(statusBytes),
			Address:          string(addressBytes),
			StaticAddress:    string(staticAddressBytes),
			ListenPort:       listenPort,
			Lighthouse:       is_lighthouse,
			Groups:           groups,
			FirewallInbound:  inbound,
			FirewallOutbound: outbound,
		}
		return err
	})
	return node, err
}

func getLighthouses(bkt *bolt.Bucket) (lighthouses []string,
	err error) {
	lighthouses = []string{}
	err = bkt.ForEach(func(key, val []byte) error {
		if val == nil {
			nb := bkt.Bucket(key)
			log.Println("Checking if node is lighthouse", string(key))
			lhBytes := nb.Get([]byte("is_lighthouse"))
			isLighthouse, err := strconv.ParseBool(string(lhBytes))
			log.Println("lighthouse: ", isLighthouse, string(lhBytes))
			if err != nil {
				return err
			}
			if isLighthouse {
				addressBytes := nb.Get([]byte("address"))
				lighthouses = append(lighthouses, string(addressBytes))
			}
		}
		return err
	})
	return lighthouses, err
}

func getStaticHosts(nodename string, bkt *bolt.Bucket) (hosts []StaticHost,
	err error) {
	hosts = []StaticHost{}
	err = bkt.ForEach(func(key, val []byte) error {
		if val == nil && string(key) != nodename {
			nb := bkt.Bucket(key)
			addressBytes := nb.Get([]byte("address"))
			staticAddressBytes := nb.Get([]byte("static_address"))
			staticAddress := string(staticAddressBytes)
			portBytes := nb.Get([]byte("listen_port"))
			port := string(portBytes)
			endpoint := fmt.Sprintf("%s:%s", staticAddress, port)
			if staticAddress != "" {
				host := StaticHost{
					NebulaAddress: string(addressBytes),
					Endpoint:      []string{endpoint},
				}
				hosts = append(hosts, host)
			}
		}
		return err
	})
	return hosts, err
}

func (b *boltdbi) getNodeConfig(netname string,
	nodename string) (config NodeConfigSchema,
	err error) {
	err = b.db.View(func(tx *bolt.Tx) error {
		netBkt := tx.Bucket([]byte(netname))
		if netBkt == nil {
			return errors.Errorf("Network does not exist.")
		}
		cipher := netBkt.Get([]byte("CIPHER"))
		nodeBkt := netBkt.Bucket([]byte(nodename))
		if nodeBkt == nil {
			return errors.Errorf("Node does not exist in network.")
		}
		addressBytes := nodeBkt.Get([]byte("address"))
		listenPortBytes := nodeBkt.Get([]byte("listen_port"))
		listenPort, err := strconv.Atoi(string(listenPortBytes))
		groupsBytes := nodeBkt.Get([]byte("groups"))
		var groups []string
		err = json.Unmarshal(groupsBytes, &groups)
		lighthouseBytes := nodeBkt.Get([]byte("is_lighthouse"))
		is_lighthouse, err := strconv.ParseBool(string(lighthouseBytes))
		var lighthouses []string
		if !is_lighthouse {
			lighthouses, err = getLighthouses(netBkt)
			if err != nil {
				return err
			}
		} else {
			lighthouses = []string{}
		}
		if err != nil {
			return err
		}
		staticHosts, err := getStaticHosts(nodename, netBkt)
		inboundBytes := nodeBkt.Get([]byte("firewall_inbound"))
		outboundBytes := nodeBkt.Get([]byte("firewall_outbound"))
		var inbound []FirewallRule
		var outbound []FirewallRule
		err = json.Unmarshal(inboundBytes, &inbound)
		if err != nil {
			return err
		}
		err = json.Unmarshal(outboundBytes, &outbound)
		if err != nil {
			return err
		}
		config = NodeConfigSchema{
			Address:          string(addressBytes),
			AmLighthouse:     is_lighthouse,
			Cipher:           string(cipher),
			Lighthouses:      lighthouses,
			StaticHosts:      staticHosts,
			ListenPort:       listenPort,
			FirewallInbound:  inbound,
			FirewallOutbound: outbound,
		}
		return err
	})
	return config, err
}

func (b *boltdbi) getNodePubkey(netname string,
	nodename string) (pubkey []byte,
	err error) {
	err = b.db.View(func(tx *bolt.Tx) error {
		netBkt := tx.Bucket([]byte(netname))
		if netBkt == nil {
			return errors.Errorf("Network does not exist.")
		}
		nodeBkt := netBkt.Bucket([]byte(nodename))
		if nodeBkt == nil {
			return errors.Errorf("Node does not exist in network.")
		}
		pubkey = nodeBkt.Get([]byte("pubkey"))
		return err
	})
	return pubkey, err
}

func (b *boltdbi) getNetworkCA(netname string) (privkey []byte,
	cert []byte,
	err error) {
	err = b.db.View(func(tx *bolt.Tx) error {
		bkt := tx.Bucket([]byte(netname))
		if bkt == nil {
			return errors.Errorf("Network does not exist.")
		}
		privkey = bkt.Get([]byte("CA_PRIV_KEY"))
		cert = bkt.Get([]byte("CA_CERT"))
		return err
	})
	return privkey, cert, err
}

func (b *boltdbi) saveNodeCert(netname string,
	nodename string,
	cert []byte,
) (err error) {
	err = b.db.Update(func(tx *bolt.Tx) error {
		netBkt := tx.Bucket([]byte(netname))
		if netBkt == nil {
			return errors.Errorf("Network does not exist.")
		}
		nodeBkt := netBkt.Bucket([]byte(nodename))
		if nodeBkt == nil {
			return errors.Errorf("Node does not exist in network.")
		}
		err = nodeBkt.Put([]byte("cert"), cert)
		return err
	})
	return err
}

func (b *boltdbi) updateLatestFetch(netname string,
	nodename string,
	timestamp string,
) (err error) {
	err = b.db.Update(func(tx *bolt.Tx) error {
		netBkt := tx.Bucket([]byte(netname))
		if netBkt == nil {
			return errors.Errorf("Network does not exist.")
		}
		nodeBkt := netBkt.Bucket([]byte(nodename))
		if nodeBkt == nil {
			return errors.Errorf("Node does not exist in network.")
		}
		err = nodeBkt.Put([]byte("latest_fetch"), []byte(timestamp))
		return err
	})
	return err
}

func (b *boltdbi) updateNetwork(netname, cidr,
	cipher string, groups []string) (err error) {
	err = b.db.Update(func(tx *bolt.Tx) error {
		bkt := tx.Bucket([]byte(netname))
		if bkt == nil {
			return errors.Errorf("Network does not exist.")
		}
		groupsBytes, err := json.Marshal(groups)
		if err != nil {
			return err
		}
		err = bkt.Put([]byte("GROUPS"), groupsBytes)
		if err != nil {
			return err
		}
		if cidr != "" {
			err = bkt.Put([]byte("CIDR"), []byte(cidr))
		}
		if cipher != "" {
			err = bkt.Put([]byte("CIPHER"), []byte(cipher))
		}
		return err
	})
	return err
}

func (b *boltdbi) newAddress(netname string) (address string,
	err error) {
	var used []string
	var cidr string
	err = b.db.View(func(tx *bolt.Tx) error {
		bkt := tx.Bucket([]byte(netname))
		if bkt == nil {
			return errors.Errorf("Network does not exist.")
		}
		cidr = string(bkt.Get([]byte("CIDR")))
		err := bkt.ForEach(func(key, val []byte) error {
			if val == nil {
				// keyval is bucket so is node
				nb := bkt.Bucket(key)
				addr := nb.Get([]byte("address"))
				used = append(used, string(addr))
			}
			return err
		})
		return err
	})
	return newAddress(cidr, used)
}
