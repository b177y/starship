package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/b177y/starship/nebutils"
	"github.com/b177y/starship/wormhole"
	"github.com/b177y/starship/xpaseto"
	log "github.com/sirupsen/logrus"
	"github.com/slackhq/nebula/cert"
	"gopkg.in/yaml.v2"
)

// create signed XPASETO token to authenticate to Quasar server
func getIdentityToken(netname, nodename string,
	privkey []byte) (token string,
	err error) {
	t := wormhole.NodeIdentitySchema{
		Netname:  netname,
		Nodename: nodename,
	}
	jsonToken, err := wormhole.NewToken(t)
	if err != nil {
		log.Fatal("Error creating paseto." + err.Error())
	}
	signer := xpaseto.NewSigner(privkey, []byte{})
	token, err = signer.SelfSignPaseto(jsonToken)
	if err != nil {
		log.Fatal("Error signing paseto." + err.Error())
	}
	return token, nil
}

// open private key to use for signing requests
func getKey(netname string) (key []byte, err error) {
	privpem, err := ioutil.ReadFile(fmt.Sprintf("/etc/nebula/%s/neutron.key",
		netname))
	if err != nil {
		log.Error("Error reading private key from /etc/nebula")
		return nil, err
	}
	key, _, err = cert.UnmarshalX25519PrivateKey(privpem)
	if err != nil {
		log.Error("Error decoding key from /etc/nebula")
		return nil, err
	}
	return key, nil
}

// Save certificate from fetched config to neutron network directory
func saveCert(netname, cert string) error {
	path := fmt.Sprintf("/etc/nebula/%s/neutron.crt", netname)
	err := ioutil.WriteFile(path, []byte(cert), 0660)
	return err
}

// get public key from CA certificate and validate cert
func getPubkey(netname string) (pubkey []byte, err error) {
	cacert, err := ioutil.ReadFile(fmt.Sprintf("/etc/nebula/%s/ca.crt",
		netname))
	if err != nil {
		return pubkey, err
	}
	nc, _, err := cert.UnmarshalNebulaCertificateFromPEM(cacert)
	caPool := cert.NewCAPool()
	// validate certificate
	_, err = caPool.AddCACertificate(cacert)
	if err != nil {
		log.Fatalf("Cert for network %s is not valid: %s.\n", netname, err.Error())
	}
	// get public key and convert from edwards format to curve25519
	edpub := nc.Details.PublicKey
	pubkey = nebutils.Ed25519PublicKeyToCurve25519(edpub)
	return pubkey, nil

}

func getConfig(netname, nodename, qAddr string) {
	endPoint := fmt.Sprintf("%s/api/neutron/config?net=%s&node=%s",
		qAddr, netname, nodename)
	privkey, err := getKey(netname)
	token, err := getIdentityToken(netname, nodename, privkey)
	body := bytes.NewBuffer([]byte(token))
	req, err := http.NewRequest("GET", endPoint, body)
	if err != nil {
		log.Fatal("Error creating request." + err.Error())
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error contacting Quasar." + err.Error())
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode != 200 {
		if resp.StatusCode == 425 {
			log.Error("Node is not enabled - please enable it from the frontend!")
			return
		}
		log.Error("Error from Quasar: " + string(b))
		return
	}
	pubkey, err := getPubkey(netname)
	if err != nil {
		log.Fatal("Could not get CA pubkey: " + err.Error())
	}
	signer := xpaseto.NewSigner([]byte{}, pubkey)
	jsonToken, err := signer.ParsePaseto(string(b))
	if err != nil {
		log.Fatal("Could not decode response token: " + err.Error())
	}
	config := *new(NodeConfigSchema)
	err = wormhole.SchemaFromJSONToken(jsonToken, &config)
	if err != nil {
		log.Error("Can't decode node config: ", err)
		return
	}
	config.Netname = netname
	err = saveCert(netname, config.Cert)
	if err != nil {
		log.Error(err)
		return
	}
	err = genConfig(config)
	if err != nil {
		log.Error(err)
		return
	}
	log.Println("Successfully updated config.")
}

func loadNeutronConfig(netname string) (config NeutronConfig, err error) {
	path := fmt.Sprintf("/etc/nebula/%s/neutron.yml", netname)
	config = NeutronConfig{}
	file, err := os.Open(path)
	if err != nil {
		return config, err
	}
	defer file.Close()

	d := yaml.NewDecoder(file)

	if err := d.Decode(&config); err != nil {
		return config, err
	}

	return config, err
}
func update(netname string) {
	config, err := loadNeutronConfig(netname)
	if err != nil {
		log.Fatal(err)
	}
	getConfig(netname, config.Nodename, config.Quasar)
}
