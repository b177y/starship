package xpaseto_test

import (
	"crypto/rand"
	"io"
	"testing"

	"github.com/b177y/starship/nebutils"
	"github.com/b177y/starship/xpaseto"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/curve25519"
)

func randBytes(data []byte) {
	if _, err := io.ReadFull(rand.Reader, data); err != nil {
		panic(err)
	}
}

func TestSign(t *testing.T) {
	msg := make([]byte, 200)

	var priv, pub [32]byte
	var random [64]byte

	// Test for random values of the keys, nonce and message
	for i := 0; i < 100; i++ {
		randBytes(priv[:])
		priv[0] &= 248
		priv[31] &= 63
		priv[31] |= 64
		curve25519.ScalarBaseMult(&pub, &priv)
		pub := pub[:]
		priv := priv[:]
		randBytes(random[:])
		randBytes(msg)
		sig, err := xpaseto.Sign(priv, msg, random)
		assert.True(t, err == nil, "Sign must work")
		v := xpaseto.Verify(pub, msg, sig)
		assert.True(t, v, "Verify must work")
	}
}

func TestSignNebkey(t *testing.T) {
	msg := make([]byte, 200)

	var random [64]byte

	// Test for random values of the keys, nonce and message
	for i := 0; i < 100; i++ {
		pub, priv := nebutils.X25519KeyPair()
		randBytes(random[:])
		randBytes(msg)
		sig, err := xpaseto.Sign(priv, msg, random)
		assert.True(t, err == nil, "Sign must work")
		v := xpaseto.Verify(pub, msg, sig)
		assert.True(t, v, "Verify must work")
	}
}
