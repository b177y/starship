package xpaseto_test

import (
	"crypto/rand"
	"io"
	"testing"

	"github.com/b177y/starship/nebutils"
	"github.com/b177y/starship/xpaseto"
	"github.com/stretchr/testify/assert"
)

func TestVerifyAndSign(t *testing.T) {
	xv2 := xpaseto.NewXV2()
	pubkey, privkey := nebutils.X25519KeyPair()

	testMsg := "Hello World!"

	// Test Sign
	token, err := xv2.Sign(privkey, testMsg, "test_footer")
	assert.NoError(t, err, "XV2.Sign should not return an error")
	assert.IsType(t, *new(string), token, "token should be a string")

	// Test Verify
	var payload string
	var footer string
	err = xv2.Verify(token, pubkey, &payload, &footer)
	assert.NoError(t, err, "XV2.Verify should succeed.")
	assert.Equal(t, testMsg, payload)

}

func TestIncorrectPubkey(t *testing.T) {
	xv2 := xpaseto.NewXV2()
	_, privkey := nebutils.X25519KeyPair()
	differentPubkey, _ := nebutils.X25519KeyPair()

	testMsg := "Hello World!"

	// Sign
	token, err := xv2.Sign(privkey, testMsg, "test_footer")

	// Test Verify
	var payload string
	var footer string
	err = xv2.Verify(token, differentPubkey, &payload, &footer)
	assert.Error(t, err, "XV2.Verify should fail with incorrect public key.")

}

func TestInvalidPrivKey(t *testing.T) {
	xv2 := xpaseto.NewXV2()
	var privkey [33]byte
	io.ReadFull(rand.Reader, privkey[:])
	var privkey2 [31]byte
	io.ReadFull(rand.Reader, privkey[:])
	testMsg := "Hello World!"

	_, err := xv2.Sign(privkey[:], testMsg, "test_footer")
	assert.Error(t, err, "XV2.Sign should fail with invalid private key")
	_, err = xv2.Sign(privkey2[:], testMsg, "test_footer")
	assert.Error(t, err, "XV2.Sign should fail with invalid private key")
}
