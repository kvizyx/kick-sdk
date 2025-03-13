package kicksdk

import (
	"fmt"
	"github.com/glichtv/kick-sdk/internal/publickey"
	"net/http"
)

// ExtractWebhookEventHeader extracts all Kick's event-specific headers from the provided request
// and returns it as a WebhookEventHeader.
//
// Reference: https://docs.kick.com/events/webhook-security#headers
func ExtractWebhookEventHeader(request *http.Request) WebhookEventHeader {
	return WebhookEventHeader{
		MessageID:        request.Header.Get("Kick-Event-Message-Id"),
		SubscriptionID:   request.Header.Get("Kick-Event-Subscription-Id"),
		Signature:        request.Header.Get("Kick-Event-Signature"),
		MessageTimestamp: request.Header.Get("Kick-Event-Message-Timestamp"),
		EventType:        request.Header.Get("Kick-Event-Type"),
		EventVersion:     request.Header.Get("Kick-Event-Version"),
	}
}

// VerifyWebhookEvent verifies webhook event signature to ensure that event with provided header and body
// was actually sent from the Kick server.
//
// Reference: https://docs.kick.com/events/webhook-security#webhook-sender-validation
func VerifyWebhookEvent(header WebhookEventHeader, publicKey string, body []byte) error {
	signature := []byte(fmt.Sprintf("%s.%s.%s", header.MessageID, header.MessageTimestamp, body))

	rsaPublicKey, err := publickey.Parse([]byte(publicKey))
	if err != nil {
		return fmt.Errorf("parse public key: %w", err)
	}

	if err = publickey.VerifyEventSignature(&rsaPublicKey, signature, []byte(header.Signature)); err != nil {
		return fmt.Errorf("verify event signature: %w", err)
	}

	return nil
}
