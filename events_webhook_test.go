package kicksdk

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVerifyWebhookEvent(t *testing.T) {
	t.Parallel()

	t.Run("Successful verification", func(t *testing.T) {
		var (
			body   = []byte("test-body")
			header = WebhookEventHeader{
				MessageID:        "test-id",
				MessageTimestamp: "2023-01-01T00:00:00Z",
			}
		)

		publicKey, signature, err := generateMockKeyAndSignature(header.MessageID, header.MessageTimestamp, body)
		assert.NoError(t, err)

		header.Signature = signature

		err = VerifyWebhookEvent(header, publicKey, body)
		assert.NoError(t, err)
	})

	t.Run("Public key parsing error", func(t *testing.T) {
		var (
			body   = []byte("test-body")
			header = WebhookEventHeader{
				MessageID:        "test-id",
				MessageTimestamp: "2023-01-01T00:00:00Z",
				Signature:        "valid-signature",
			}
		)

		err := VerifyWebhookEvent(header, "invalid-public-key", body)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "parse public key")
	})

	t.Run("Signature verification error", func(t *testing.T) {
		var (
			body   = []byte("test-body")
			header = WebhookEventHeader{
				MessageID:        "test-id",
				MessageTimestamp: "2023-01-01T00:00:00Z",
				Signature:        "invalid-base64-signature",
			}
		)

		publicKey, _, err := generateMockKeyAndSignature(header.MessageID, header.MessageTimestamp, body)
		assert.NoError(t, err)

		err = VerifyWebhookEvent(header, publicKey, body)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "verify event signature")
	})
}

func generateMockKeyAndSignature(messageID, timestamp string, body []byte) (string, string, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", "", fmt.Errorf("generate key: %w", err)
	}

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return "", "", fmt.Errorf("marshal public key: %w", err)
	}

	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	})

	var (
		signaturePayload = []byte(fmt.Sprintf("%s.%s.%s", messageID, timestamp, body))
		hash             = sha256.Sum256(signaturePayload)
	)

	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hash[:])
	if err != nil {
		return "", "", fmt.Errorf("sign data: %w", err)
	}

	signatureBase64 := base64.StdEncoding.EncodeToString(signature)

	return string(publicKeyPEM), signatureBase64, nil
}
