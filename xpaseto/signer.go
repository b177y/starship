package xpaseto

import (
	"time"

	paseto "github.com/o1egl/paseto/v2"
	"github.com/slackhq/nebula/cert"
)

// Signer is a struct which is used to Sign and Verify XPASETO tokens
// Privkey is a montgomery x25519 private key
// Pubkey is the public counterpart to the Privkey
type Signer struct {
	paseto  *XV2
	Privkey []byte
	Pubkey  []byte
}

// NewSigner returns a new Signer, initialised with the given
// private and public x25519 keys
func NewSigner(privkey []byte,
	pubkey []byte) Signer {
	return Signer{
		paseto:  NewXV2(),
		Privkey: privkey,
		Pubkey:  pubkey,
	}
}

// SignPaseto returns a token from the given paseto.JSONToken
func (s *Signer) SignPaseto(jsonToken paseto.JSONToken) (token string, err error) {
	token, err = s.paseto.Sign(s.Privkey, jsonToken, "STARSHIP")
	if err != nil {
		return "", err
	}
	return token, nil
}

// ParsePaseto returns a paseto.JSONToken from a raw token string
// If the token signature is invalid, an error will be returned
func (s *Signer) ParsePaseto(token string) (jsonToken paseto.JSONToken, err error) {
	var payload paseto.JSONToken
	var footer string
	err = s.paseto.Verify(token, s.Pubkey, &payload, &footer)
	if err != nil {
		return payload, err
	}
	err = payload.Validate(paseto.ValidAt(time.Now()))
	if err != nil {
		return payload, err
	}
	return payload, nil
}

// SelfSignPaseto sets the "pubkey" additional claim to the value of the public part
// of the keypair used to sign the XPASETO token
func (s *Signer) SelfSignPaseto(jsonToken paseto.JSONToken) (token string, err error) {
	pubpem := string(cert.MarshalX25519PublicKey(s.Pubkey))
	jsonToken.Set("pubkey", pubpem)
	return s.SignPaseto(jsonToken)
}

// ParseSelfSigned parses and validates an XPASETO token against the 'pubkey'
// in the additional claims.
func (s *Signer) ParseSelfSigned(token string) (jsonToken paseto.JSONToken,
	err error) {
	data, _, err := splitToken([]byte(token), headerXV2Public)
	payloadBytes := data[:len(data)-64]
	err = fillValue(payloadBytes, &jsonToken)
	if err != nil {
		return jsonToken, err
	}
	var pubkey string
	err = jsonToken.Get("pubkey", &pubkey)

	if err != nil {
		return jsonToken, err
	}
	s.Pubkey, _, err = cert.UnmarshalX25519PublicKey([]byte(pubkey))
	if err != nil {
		return jsonToken, err
	}
	return s.ParsePaseto(token)
}
