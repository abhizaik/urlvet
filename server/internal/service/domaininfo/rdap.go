package domaininfo

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

// ── IANA RDAP bootstrap cache ──────────────────────────────────────────────

var (
	ianaBootstrapMu     sync.RWMutex
	ianaBootstrapCache  map[string]string // lowercase TLD → base RDAP URL
	ianaBootstrapExpiry time.Time
)

// ianaBootstrapJSON is the shape of https://data.iana.org/rdap/dns.json.
type ianaBootstrapJSON struct {
	Services [][][]string `json:"services"`
}

// loadIANABootstrap fetches and parses the IANA RDAP bootstrap file.
// Results are cached in memory for 24 hours; stale data is returned on error.
func loadIANABootstrap() map[string]string {
	ianaBootstrapMu.RLock()
	if ianaBootstrapCache != nil && time.Now().Before(ianaBootstrapExpiry) {
		m := ianaBootstrapCache
		ianaBootstrapMu.RUnlock()
		return m
	}
	ianaBootstrapMu.RUnlock()

	ianaBootstrapMu.Lock()
	defer ianaBootstrapMu.Unlock()

	// Double-check under write lock.
	if ianaBootstrapCache != nil && time.Now().Before(ianaBootstrapExpiry) {
		return ianaBootstrapCache
	}

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get("https://data.iana.org/rdap/dns.json")
	if err != nil {
		return ianaBootstrapCache // return stale data rather than nothing
	}
	defer resp.Body.Close()

	var bootstrap ianaBootstrapJSON
	if err := json.NewDecoder(resp.Body).Decode(&bootstrap); err != nil {
		return ianaBootstrapCache
	}

	result := make(map[string]string, 1500)
	for _, service := range bootstrap.Services {
		if len(service) != 2 || len(service[1]) == 0 {
			continue
		}
		// Use the first server URL, strip trailing slash.
		server := strings.TrimSuffix(service[1][0], "/")
		for _, tld := range service[0] {
			result[strings.ToLower(tld)] = server
		}
	}

	ianaBootstrapCache = result
	ianaBootstrapExpiry = time.Now().Add(24 * time.Hour)
	return result
}

type rdapResponse struct {
	LDHName     string `json:"ldhName"`
	Nameservers []struct {
		LDHName string `json:"ldhName"`
	} `json:"nameservers"`
	Events []struct {
		Action string    `json:"eventAction"`
		Date   time.Time `json:"eventDate"`
	} `json:"events"`
	Entities []struct {
		Roles      []string `json:"roles"`
		VCardArray []any    `json:"vcardArray"`
	} `json:"entities"`
	Status    []string `json:"status"`
	SecureDNS struct {
		DelegationSigned bool `json:"delegationSigned"`
	} `json:"secureDNS"`
}

// fetchRDAP queries RDAP and returns normalized RegistrationData.
// Deprecated: Use fetchRDAPWithContext instead.
func fetchRDAP(domain string) (*RegistrationData, error) {
	return fetchRDAPWithContext(context.Background(), domain)
}

// fetchRDAPWithContext queries RDAP with context support and timeout.
func fetchRDAPWithContext(ctx context.Context, domain string) (*RegistrationData, error) {
	tld := strings.Split(domain, ".")
	if len(tld) < 2 {
		return nil, fmt.Errorf("invalid domain")
	}

	// Get the TLD (last part of domain)
	tldPart := tld[len(tld)-1]

	// Try to find appropriate RDAP server for the TLD
	rdapURL, err := getRDAPServer(tldPart)
	if err != nil {
		return nil, fmt.Errorf("RDAP not supported for TLD .%s: %w", tldPart, err)
	}

	url := fmt.Sprintf("%s/domain/%s", rdapURL, domain)

	// Create HTTP client with timeout to fail fast
	client := &http.Client{
		Timeout: 3 * time.Second, // Fast timeout for non-.com domains
	}

	// Create request with context
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("RDAP request creation failed: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("RDAP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("RDAP query failed: %s", resp.Status)
	}

	var rd rdapResponse
	if err := json.NewDecoder(resp.Body).Decode(&rd); err != nil {
		return nil, fmt.Errorf("RDAP response parsing failed: %w", err)
	}

	// Extract registrar (may be nil in some RDAP servers)
	registrar := ""
	for _, e := range rd.Entities {
		for _, role := range e.Roles {
			if role == "registrar" {
				// Entities often store registrar name in vCardArray
				if len(e.VCardArray) >= 2 {
					if vcardData, ok := e.VCardArray[1].([]any); ok {
						for _, item := range vcardData {
							if entry, ok := item.([]any); ok && len(entry) >= 3 {
								if key, ok := entry[0].(string); ok && key == "fn" { // "fn" = full name
									if value, ok := entry[3].(string); ok {
										registrar = value
										break
									}
								}
							}
						}
					}
				}
			}
		}
		if registrar != "" {
			break
		}
	}

	var created, updated, expiry time.Time
	for _, e := range rd.Events {
		switch e.Action {
		case "registration":
			created = e.Date
		case "last changed":
			updated = e.Date
		case "expiration":
			expiry = e.Date
		}
	}

	var ns []string
	for _, n := range rd.Nameservers {
		ns = append(ns, n.LDHName)
	}

	raw, _ := json.Marshal(rd)

	return &RegistrationData{
		Domain:      rd.LDHName,
		Registrar:   registrar,
		CreatedDate: created,
		UpdatedDate: updated,
		ExpiryDate:  expiry,
		Nameservers: ns,
		Status:      rd.Status,
		DNSSEC:      rd.SecureDNS.DelegationSigned,
		Raw:         string(raw),
		Source:      "RDAP",
	}, nil
}

// wellKnownRDAP is a fast-path cache for the most common TLDs, avoiding a
// network round-trip to the IANA bootstrap service on every request.
var wellKnownRDAP = map[string]string{
	"com": "https://rdap.verisign.com/com/v1",
	"net": "https://rdap.verisign.com/net/v1",
	"org": "https://rdap.pir.org/rdap/org/v1",
	"io":  "https://rdap.nic.io/v1",
	"co":  "https://rdap.nic.co/v1",
	"me":  "https://rdap.nic.me/v1",
	"tv":  "https://rdap.nic.tv/v1",
	"cc":  "https://rdap.nic.cc/v1",
}

// getRDAPServer returns the RDAP base URL for a TLD.
// It checks the well-known map first, then falls back to the IANA bootstrap.
func getRDAPServer(tld string) (string, error) {
	tld = strings.ToLower(tld)

	if server, ok := wellKnownRDAP[tld]; ok {
		return server, nil
	}

	bootstrap := loadIANABootstrap()
	if bootstrap != nil {
		if server, ok := bootstrap[tld]; ok {
			return server, nil
		}
	}

	return "", fmt.Errorf("no RDAP server found for .%s", tld)
}
