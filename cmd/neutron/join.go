package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/b177y/starship/nebutils"
	"github.com/b177y/starship/wormhole"
	"github.com/b177y/starship/xpaseto"
	log "github.com/sirupsen/logrus"
	"github.com/slackhq/nebula/cert"
	"github.com/teris-io/shortid"
	"gopkg.in/yaml.v2"
)

// config to save the address of a quasar server and the node name
type NeutronConfig struct {
	Quasar   string `yaml:"quasar"`
	Nodename string `yaml:"nodename"`
}

// save the neutron config to a yaml file
func saveNeutronConfig(netname, quasar, nodename string) error {
	nconf := NeutronConfig{
		Quasar:   quasar,
		Nodename: nodename,
	}
	path := fmt.Sprintf("/etc/nebula/%s/neutron.yml", netname)
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	err = f.Truncate(0)
	if err != nil {
		return err
	}
	e := yaml.NewEncoder(f)
	e.Encode(nconf)
	return f.Close()
}

// request to join a Starship network
func signReq(netname, hostname, nodename, qAddr string,
	privkey, pubkey []byte) {
	// create request payload
	t := wormhole.RequestJoinSchema{
		Netname:  netname,
		Nodename: nodename,
		Hostname: hostname,
		PubKey:   string(cert.MarshalX25519PublicKey(pubkey)),
	}

	// turn payload into token
	jsonToken, err := wormhole.NewToken(t)
	if err != nil {
		log.Fatal("Error creating paseto." + err.Error())
	}
	// sign token using nebula private key
	signer := xpaseto.NewSigner(privkey, pubkey)
	token, err := signer.SelfSignPaseto(jsonToken)
	if err != nil {
		log.Fatal("Error signing paseto." + err.Error())
	}

	body := bytes.NewBuffer([]byte(token))
	resp, err := http.Post(qAddr+"/api/neutron/join", "text/plain", body)
	if err != nil {
		log.Fatal("Error contacting quasar." + err.Error())
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading response body from quasar." + err.Error())
	}
	if resp.StatusCode != 200 {
		log.Fatal("Error status from quasar: " + string(respBody))
	}
	pubFingerprint := base64.StdEncoding.EncodeToString(pubkey)
	log.Println("Node Fingerprint: ", pubFingerprint)

}

// get the certificate of the CA of a network
func getCaCert(qAddr string, netName string) {
	path := "/etc/nebula/"
	os.MkdirAll(path, 0775)
	path = path + netName
	err := os.Mkdir(path, 0770)
	if err != nil {
		if os.IsExist(err) {
			log.Warning("Network already exists locally. Existing config may be overwritten!")
		} else if os.IsPermission(err) {
			log.Fatal("Permission denied. Are you running as root?")
		} else {
			log.Fatal("Error: " + err.Error())
		}
	}
	ca_url := qAddr + "/api/networks/" + netName + "/cert"
	resp, err := http.Get(ca_url)
	if err != nil {
		log.Fatal("Could not get ca from Quasar server. " + err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == 404 {
			log.Fatalf("Network %s does not exist.\n", netName)
		}
		log.Fatalf("Bad status code from Quasar: %d", resp.StatusCode)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Could not read HTTP body")
	}
	cacert, _, err := cert.UnmarshalNebulaCertificateFromPEM(body)
	if err != nil {
		log.Fatal("Cert is not a valid Nebula Certificate.")
	}
	certfp, err := cacert.Sha256Sum()
	if err != nil {
		log.Fatal("Cert is not a valid Nebula Certificate. Could not calculate fingerprint.")
	}
	fmt.Printf("Certificate Fingerprint: %s\n", certfp)
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Trust fingerprint? (Y/n) ")
	inp, _ := reader.ReadString('\n')
	if inp == "n\n" || inp == "N\n" {
		log.Fatal("Certificate is not trusted. Exitting.")
	}
	log.Println("Saving certificate")
	dst := path + "/ca.crt"
	f, err := os.Create(dst)
	if err != nil {
		log.Fatalf("Could not create file %s", dst)
	}
	r := bytes.NewReader(body)
	_, err = io.Copy(f, r)
}

func initialise(qAddr string,
	nodeName string,
	netName string,
) {
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}
	// validate commandline arguments
	if qAddr == "" {
		log.Fatal("Quasar address must be given with -quasar")
	}
	if netName == "" {
		log.Fatal("Network name must be given with -network")
	}
	if nodeName == "" {
		// generate node name based on hostname if none provided
		hostid, err := shortid.Generate()
		if err != nil {
			log.Fatal("Could not generate id for nodename" + err.Error())
		}
		nodeName = fmt.Sprintf("%s-%s", hostname, hostid)
	}
	getCaCert(qAddr, netName)
	// create nebula keypair
	pubkey, privkey := nebutils.X25519KeyPair()
	err = nebutils.SaveKey("/etc/nebula/"+netName,
		"neutron.key",
		privkey,
	)
	if err != nil {
		log.Fatal("Error saving key" + err.Error())
	}
	// request server to sign pubkey
	signReq(netName, hostname, nodeName, qAddr,
		privkey, pubkey)
	// save config
	saveNeutronConfig(netName, qAddr, nodeName)
}
