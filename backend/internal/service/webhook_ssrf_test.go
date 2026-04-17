package service

import (
	"errors"
	"testing"
)

func TestValidateWebhookURL(t *testing.T) {
	cases := []struct {
		name    string
		url     string
		wantErr error
	}{
		// Scheme rejects.
		{"ftp rejected", "ftp://example.com/hook", ErrWebhookInvalidScheme},
		{"file rejected", "file:///etc/passwd", ErrWebhookInvalidScheme},
		{"no scheme rejected", "example.com/hook", ErrWebhookInvalidScheme},
		{"empty rejected", "", ErrWebhookInvalidScheme},

		// Metadata endpoints — always blocked.
		{"aws metadata by ip", "http://169.254.169.254/latest/meta-data/", ErrWebhookMetadataBlocked},
		{"gcp metadata by host", "http://metadata.google.internal/", ErrWebhookMetadataBlocked},
		{"alibaba metadata", "http://100.100.100.200/latest/meta-data/", ErrWebhookMetadataBlocked},

		// Private / loopback — blocked unless MP_WEBHOOK_ALLOW_PRIVATE=1.
		{"loopback v4 blocked", "http://127.0.0.1:8080/hook", ErrWebhookPrivateBlocked},
		{"loopback v6 blocked", "http://[::1]:8080/hook", ErrWebhookPrivateBlocked},
		{"rfc1918 10.x blocked", "http://10.0.0.5/hook", ErrWebhookPrivateBlocked},
		{"rfc1918 192.168.x blocked", "http://192.168.1.1/hook", ErrWebhookPrivateBlocked},
		{"link-local v4 blocked", "http://169.254.1.1/hook", ErrWebhookPrivateBlocked},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateWebhookURL(tc.url)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("url=%s: got %v, want %v", tc.url, err, tc.wantErr)
			}
		})
	}
}

func TestValidateWebhookURL_AllowPrivateOptOut(t *testing.T) {
	t.Setenv("MP_WEBHOOK_ALLOW_PRIVATE", "1")

	// Loopback should now pass (local self-test scenario).
	if err := validateWebhookURL("http://127.0.0.1:8080/hook"); err != nil {
		t.Fatalf("loopback should be allowed with opt-out: %v", err)
	}
	// Metadata endpoints must still be blocked even with opt-out.
	if err := validateWebhookURL("http://169.254.169.254/"); !errors.Is(err, ErrWebhookMetadataBlocked) {
		t.Fatalf("metadata must stay blocked: %v", err)
	}
}
