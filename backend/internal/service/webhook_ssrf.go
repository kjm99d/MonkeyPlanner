package service

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"os"
	"strings"
)

// Errors returned by validateWebhookURL. Exposed so HTTP handlers can map
// them to 400 responses instead of 500s.
var (
	ErrWebhookInvalidScheme   = errors.New("webhook URL must use http or https")
	ErrWebhookMetadataBlocked = errors.New("webhook URL targets a cloud metadata endpoint")
	ErrWebhookPrivateBlocked  = errors.New("webhook URL resolves to a private/loopback address; set MP_WEBHOOK_ALLOW_PRIVATE=1 to override")
	ErrWebhookUnresolved      = errors.New("webhook URL could not be resolved to any IP address")
)

// Cloud-provider instance metadata endpoints. These are the high-value
// SSRF targets — an attacker who persuades MonkeyPlanner to POST to them
// can steal IAM credentials. Blocked unconditionally, no env opt-out.
var metadataHosts = map[string]struct{}{
	"169.254.169.254":            {}, // AWS, Azure, OCI, DigitalOcean
	"metadata.google.internal":   {}, // GCP (by DNS)
	"metadata":                   {}, // GCP short form
	"100.100.100.200":            {}, // Alibaba Cloud
	"fd00:ec2::254":              {}, // AWS IPv6
}

// validateWebhookURL parses rawURL and rejects destinations that should
// never receive a webhook from this server:
//   - anything other than http:// or https://
//   - cloud-provider metadata endpoints (never legitimate)
//   - loopback / private / link-local addresses (unless the operator
//     set MP_WEBHOOK_ALLOW_PRIVATE=1, e.g. for self-tests against a
//     local receiver on the same host)
//
// Every resolved IP must pass; a hostname that resolves to both a public
// and a private IP is rejected.
func validateWebhookURL(rawURL string) error {
	u, err := url.Parse(strings.TrimSpace(rawURL))
	if err != nil {
		return fmt.Errorf("parse webhook URL: %w", err)
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return ErrWebhookInvalidScheme
	}

	host := u.Hostname()
	if host == "" {
		return ErrWebhookInvalidScheme
	}

	lowerHost := strings.ToLower(host)
	if _, meta := metadataHosts[lowerHost]; meta {
		return ErrWebhookMetadataBlocked
	}

	allowPrivate := os.Getenv("MP_WEBHOOK_ALLOW_PRIVATE") == "1"

	// Resolve to IPs. If the hostname is already a literal IP, LookupIP
	// still returns it, so the same checks apply uniformly.
	ips, err := net.LookupIP(host)
	if err != nil || len(ips) == 0 {
		return ErrWebhookUnresolved
	}

	for _, ip := range ips {
		if _, meta := metadataHosts[ip.String()]; meta {
			return ErrWebhookMetadataBlocked
		}
		if allowPrivate {
			continue
		}
		if isBlockedIP(ip) {
			return ErrWebhookPrivateBlocked
		}
	}
	return nil
}

// isBlockedIP returns true if ip is in a range that should never receive
// a webhook from an untrusted-URL context: loopback, link-local, private
// (RFC1918 / ULA), unspecified, or multicast.
func isBlockedIP(ip net.IP) bool {
	return ip.IsLoopback() ||
		ip.IsLinkLocalUnicast() ||
		ip.IsLinkLocalMulticast() ||
		ip.IsMulticast() ||
		ip.IsUnspecified() ||
		ip.IsPrivate()
}
