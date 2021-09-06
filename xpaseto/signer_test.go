package xpaseto_test

import (
	"testing"

	"github.com/b177y/starship/nebutils"
	"github.com/b177y/starship/xpaseto"
	"github.com/o1egl/paseto"
	"github.com/stretchr/testify/assert"
)

func TestNewSigner(t *testing.T) {
	pub, priv := nebutils.X25519KeyPair()
	signer := xpaseto.NewSigner(priv, pub)
	assert.IsType(t, *new(xpaseto.Signer), signer, "signer should be of type xpaseto.Signer")
}

func TestSignAndVerifyToken(t *testing.T) {
	jsonToken := paseto.JSONToken{
		Audience: "example_audience",
		Issuer:   "example_issuer",
	}
	jsonToken.Set("example", "value")
	jsonToken.Set("anotherexample", "anothervalue")

	pub, priv := nebutils.X25519KeyPair()
	signer := xpaseto.NewSigner(priv, pub)

	// Test that we can sign
	token, err := signer.SignPaseto(jsonToken)
	assert.NoError(t, err)
	assert.IsType(t, *new(string), token, "token should be a string")

	// Test that we can parse the token
	jt, err := signer.ParsePaseto(token)
	assert.NoError(t, err)
	assert.IsType(t, *new(paseto.JSONToken), jt, "parsed object should be paseto.JSONToken")

	assert.Equal(t, jsonToken, jt, "Parsed json token should be the same as the orignal json")

}
