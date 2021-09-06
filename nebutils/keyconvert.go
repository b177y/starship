package nebutils

import (
	"crypto/ed25519"
	"crypto/sha512"
	"math/big"

	"golang.org/x/crypto/curve25519"
)

// Convert an Edwards private key to a curve25519 private key
func PrivateKeyToCurve25519(privateKey []byte) (curvePrivate []byte) {
	h := sha512.New()
	h.Write(privateKey)
	digest := h.Sum(nil)

	// key clamping
	digest[0] &= 248
	digest[31] &= 127
	digest[31] |= 64

	return digest[:32]
}

var curve25519P, _ = new(big.Int).SetString("57896044618658097711785492504343953926634992332820282019728792003956564819949", 10)

// from Filo Sottile's 'age' : https://github.com/FiloSottile/age/blob/bbab440e198a4d67ba78591176c7853e62d29e04/internal/age/ssh.go#L174
// See https://blog.filippo.io/using-ed25519-keys-for-encryption.
func Ed25519PublicKeyToCurve25519(pk ed25519.PublicKey) []byte {
	bigEndianY := make([]byte, ed25519.PublicKeySize)
	for i, b := range pk {
		bigEndianY[ed25519.PublicKeySize-i-1] = b
	}
	bigEndianY[0] &= 0b0111_1111

	y := new(big.Int).SetBytes(bigEndianY)
	denom := big.NewInt(1)
	denom.ModInverse(denom.Sub(denom, y), curve25519P)
	u := y.Mul(y.Add(y, big.NewInt(1)), denom)
	u.Mod(u, curve25519P)

	out := make([]byte, curve25519.PointSize)
	uBytes := u.Bytes()
	for i, b := range uBytes {
		out[len(uBytes)-i-1] = b
	}

	return out
}
