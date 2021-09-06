// curve25519.go - Mostly copied from https://github.com/signal-golang/textsecure
// signal-golang/textsecure license: GPL-3.0 - https://github.com/signal-golang/textsecure/blob/master/COPYING
// This file contains the Sign and Verify functions for XEdDSA Signatures
// XEdDSA is for using Montgomery keys (traditionally used for X25519 Diffie-Hellman functions) for creating and verifying EdDSA compatible signatures.
// See more here: https://signal.org/docs/specifications/xeddsa/#xeddsa

package xpaseto

import (
	"crypto/sha512"

	"github.com/pkg/errors"
	"github.com/signal-golang/ed25519"
	"github.com/signal-golang/ed25519/edwards25519"
)

// Sign signs a message with an X25519 key and returns a signature.
//
// An error will be returned if an invalid private key is given.
func Sign(privateKey []byte, message []byte, random [64]byte) (signature []byte,
	err error) {
	sig := new([64]byte)
	var privkey [32]byte
	if n := copy(privkey[:], privateKey); n != 32 {
		return signature, errors.Errorf("Invalid Private Key. Cannot Sign Payload. ")
	}

	// Calculate Ed25519 public key from Curve25519 private key
	var A edwards25519.ExtendedGroupElement
	var publicKey [32]byte
	edwards25519.GeScalarMultBase(&A, &privkey)
	A.ToBytes(&publicKey)

	// Calculate r
	diversifier := [32]byte{
		0xFE, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}

	var r [64]byte
	h := sha512.New()
	h.Write(diversifier[:])
	h.Write(privkey[:])
	h.Write(message)
	h.Write(random[:])
	h.Sum(r[:0])

	// Calculate R
	var rReduced [32]byte
	edwards25519.ScReduce(&rReduced, &r)
	var R edwards25519.ExtendedGroupElement
	edwards25519.GeScalarMultBase(&R, &rReduced)

	var encodedR [32]byte
	R.ToBytes(&encodedR)

	// Calculate S = r + SHA2-512(R || A_ed || msg) * a  (mod L)
	var hramDigest [64]byte
	h.Reset()
	h.Write(encodedR[:])
	h.Write(publicKey[:])
	h.Write(message)
	h.Sum(hramDigest[:0])
	var hramDigestReduced [32]byte
	edwards25519.ScReduce(&hramDigestReduced, &hramDigest)

	var s [32]byte
	edwards25519.ScMulAdd(&s, &hramDigestReduced, &privkey, &rReduced)

	copy(sig[:], encodedR[:])
	copy(sig[32:], s[:])
	sig[63] |= publicKey[31] & 0x80

	signature = sig[:]
	return signature, nil
}

// Verify checks whether the message has a valid signature.
//
// Returns true if the signature is valid, otherwise returns false.
func Verify(publicKey []byte, message []byte, signature []byte) bool {

	var sig [64]byte
	if n := copy(sig[:], signature); n != 64 {
		return false
	}
	var pubkey [32]byte
	if n := copy(pubkey[:], publicKey); n != 32 {
		return false
	}
	pubkey[31] &= 0x7F

	var edY, one, montX, montXMinusOne, montXPlusOne edwards25519.FieldElement
	edwards25519.FeFromBytes(&montX, &pubkey)
	edwards25519.FeOne(&one)
	edwards25519.FeSub(&montXMinusOne, &montX, &one)
	edwards25519.FeAdd(&montXPlusOne, &montX, &one)
	edwards25519.FeInvert(&montXPlusOne, &montXPlusOne)
	edwards25519.FeMul(&edY, &montXMinusOne, &montXPlusOne)

	var A_ed [32]byte
	edwards25519.FeToBytes(&A_ed, &edY)

	A_ed[31] |= sig[63] & 0x80
	sig[63] &= 0x7F

	return ed25519.Verify(&A_ed, message, &sig)
}
