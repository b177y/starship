package nebutils

import (
	"crypto/rand"
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"github.com/slackhq/nebula/cert"
	"golang.org/x/crypto/curve25519"
)

// Create a curve25519 keypair
func X25519KeyPair() (kpub []byte, kpriv []byte) {
	var pubkey, privkey [32]byte
	if _, err := io.ReadFull(rand.Reader, privkey[:]); err != nil {
		log.Fatal(err)
	}
	privkey[0] &= 248
	privkey[31] &= 63
	privkey[31] |= 64
	curve25519.ScalarBaseMult(&pubkey, &privkey)
	return pubkey[:], privkey[:]
}

func SaveKey(directory, name string,
	keyBytes []byte) (err error) {
	key_fn := filepath.Join(directory, name)
	log.WithFields(log.Fields{
		"directory": directory,
		"name":      name,
		"priv_fn":   key_fn,
	}).Info("Saving Key")
	err = ioutil.WriteFile(key_fn, cert.MarshalX25519PrivateKey(keyBytes), 0660)
	if err != nil {
		return fmt.Errorf("error while writing out-key: %s", err)
	}
	return nil
}
