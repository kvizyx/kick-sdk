package publickey

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
)

var (
	ErrCannotDecode   = errors.New("public key can't be decoded")
	ErrNotPublicKey   = errors.New("payload is not a public key")
	ErrUnexpectedType = errors.New("unexpected public key type")
)

// Parse ...
func Parse(payload []byte) (rsa.PublicKey, error) {
	block, _ := pem.Decode(payload)
	if block == nil {
		return rsa.PublicKey{}, ErrCannotDecode
	}

	if block.Type != "PUBLIC KEY" {
		return rsa.PublicKey{}, ErrNotPublicKey
	}

	parsed, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return rsa.PublicKey{}, fmt.Errorf("parse PKI public key: %w", err)
	}

	publicKey, ok := parsed.(*rsa.PublicKey)
	if !ok {
		return rsa.PublicKey{}, ErrUnexpectedType
	}

	return *publicKey, nil
}

// VerifyEventSignature ...
func VerifyEventSignature(publicKey *rsa.PublicKey, body []byte, signature []byte) error {
	decoded := make([]byte, base64.StdEncoding.DecodedLen(len(signature)))

	at, err := base64.StdEncoding.Decode(decoded, signature)
	if err != nil {
		return fmt.Errorf("decode signature: %w", err)
	}

	signature = decoded[:at]
	hashedBody := sha256.Sum256(body)

	return rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hashedBody[:], signature)
}
