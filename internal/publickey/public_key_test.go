package publickey

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	t.Parallel()

	// Generate RSA key.
	rsaKey, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.NoError(t, err)

	rsaPublicKeyBytes, err := x509.MarshalPKIXPublicKey(&rsaKey.PublicKey)
	assert.NoError(t, err)

	rsaPublicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: rsaPublicKeyBytes,
	})

	// Generate ECDSA key.
	ecdsaKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	assert.NoError(t, err)

	ecdsaPublicKeyBytes, err := x509.MarshalPKIXPublicKey(&ecdsaKey.PublicKey)
	assert.NoError(t, err)

	ecdsaPublicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: ecdsaPublicKeyBytes,
	})

	tests := []struct {
		name          string
		input         []byte
		expectedError error
		setup         func([]byte) []byte
	}{
		{
			name:          "Valid public key",
			input:         rsaPublicKeyPEM,
			expectedError: nil,
		},
		{
			name:          "Invalid PEM data",
			input:         []byte("invalid data"),
			expectedError: ErrCannotDecode,
		},
		{
			name:  "Wrong key type",
			input: rsaPublicKeyPEM,
			setup: func(data []byte) []byte {
				block, _ := pem.Decode(data)
				return pem.EncodeToMemory(&pem.Block{
					Type:  "PRIVATE KEY",
					Bytes: block.Bytes,
				})
			},
			expectedError: ErrNotPublicKey,
		},
		{
			name:  "Invalid public key bytes",
			input: rsaPublicKeyPEM,
			setup: func(data []byte) []byte {
				return pem.EncodeToMemory(&pem.Block{
					Type:  "PUBLIC KEY",
					Bytes: []byte("invalid key bytes"),
				})
			},
			expectedError: nil,
		},
		{
			name:          "Unexpected key type",
			input:         ecdsaPublicKeyPEM,
			expectedError: ErrUnexpectedType,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			input := test.input

			if test.setup != nil {
				input = test.setup(input)
			}

			key, err := Parse(input)
			if test.expectedError != nil {
				assert.ErrorIs(t, err, test.expectedError)
			}

			if test.name == "Valid public key" {
				assert.NoError(t, err)
				assert.Equal(t, rsaKey.PublicKey, key)
				return
			}

			assert.Error(t, err)
		})
	}
}

func TestVerifyEventSignature(t *testing.T) {
	t.Parallel()

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.NoError(t, err)

	message := []byte("test")
	hash := sha256.Sum256(message)

	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hash[:])
	assert.NoError(t, err)

	validSignatureBase64 := base64.StdEncoding.EncodeToString(signature)

	tests := []struct {
		name        string
		body        []byte
		signature   string
		expectError bool
	}{
		{
			name:        "Valid signature",
			body:        message,
			signature:   validSignatureBase64,
			expectError: false,
		},
		{
			name:        "Invalid base64 signature",
			body:        message,
			signature:   "invalid-base64",
			expectError: true,
		},
		{
			name:        "Invalid signature",
			body:        message,
			signature:   base64.StdEncoding.EncodeToString([]byte("invalid-signature")),
			expectError: true,
		},
		{
			name:        "Different message",
			body:        []byte("different message"),
			signature:   validSignatureBase64,
			expectError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := VerifyEventSignature(&privateKey.PublicKey, test.body, []byte(test.signature))
			if test.expectError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
		})
	}
}
