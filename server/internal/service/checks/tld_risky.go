package checks

import "github.com/abhizaik/urlvet/internal/constants"

func IsRiskyTld(domain string) (bool, bool, string) {
	tld, icann := GetTld(domain)
	_, ok := constants.RiskyTLDs[tld]

	return ok, icann, tld
}
