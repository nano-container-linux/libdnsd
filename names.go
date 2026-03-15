package libdnsd

import "strings"

// NormalizeName returns a lowercased, trimmed, dot-suffixed domain name.
func NormalizeName(name string) string {
	n := strings.ToLower(strings.TrimSpace(name))
	if n == "" {
		return "."
	}
	if !strings.HasSuffix(n, ".") {
		n += "."
	}
	return n
}

// AcmeChallengeFQDN returns the canonical _acme-challenge.fqdn. name for domain.
// If domain already starts with _acme-challenge., it is returned as-is (normalized).
func AcmeChallengeFQDN(domain string) string {
	norm := NormalizeName(domain)
	norm = strings.TrimPrefix(norm, "*.")
	if strings.HasPrefix(norm, "_acme-challenge.") {
		return norm
	}
	return "_acme-challenge." + norm
}
