package helpers

import (
	"fmt"
	"net"
	"net/url"
)

// linkLocalBlocks are the CIDR ranges rejected by ValidateUpstreamBaseURL. The API Gateway is
// designed to proxy to internal services (RFC1918 ranges and localhost are legitimate upstream
// targets — see FSD), so only link-local/metadata ranges are blocked: these have no legitimate
// use as a registered upstream and are the classic SSRF target for reaching cloud instance
// metadata endpoints (e.g. 169.254.169.254 on AWS/GCP/Azure).
var linkLocalBlocks = []string{
	"169.254.0.0/16", // IPv4 link-local, includes the cloud metadata endpoint
	"fe80::/10",      // IPv6 link-local
}

// ValidateUpstreamBaseURL parses rawURL and rejects it if it resolves to a link-local address.
// DNS is resolved (not just the literal host string) so a hostname can't be used to bypass the
// check via DNS pointing at a blocked IP.
func ValidateUpstreamBaseURL(rawURL string) error {
	parsed, err := url.Parse(rawURL)
	if err != nil || parsed.Scheme == "" || parsed.Hostname() == "" {
		return &FieldError{Field: "base_url", Message: "URL upstream tidak valid"}
	}

	blocks := make([]*net.IPNet, 0, len(linkLocalBlocks))
	for _, cidr := range linkLocalBlocks {
		_, block, _ := net.ParseCIDR(cidr)
		blocks = append(blocks, block)
	}

	host := parsed.Hostname()
	var ips []net.IP
	if ip := net.ParseIP(host); ip != nil {
		ips = []net.IP{ip}
	} else {
		resolved, err := net.LookupIP(host)
		if err != nil {
			return &FieldError{Field: "base_url", Message: fmt.Sprintf("Tidak dapat me-resolve host: %s", host)}
		}
		ips = resolved
	}

	for _, ip := range ips {
		for _, block := range blocks {
			if block.Contains(ip) {
				return &FieldError{Field: "base_url", Message: "URL upstream tidak boleh mengarah ke alamat link-local/metadata"}
			}
		}
	}

	return nil
}
