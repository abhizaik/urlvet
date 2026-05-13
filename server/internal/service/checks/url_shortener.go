package checks

import "github.com/abhizaik/urlvet/internal/constants"

func IsUrlShortener(domain string) bool {
	_, ok := constants.URLShorteners[domain]
	return ok
}
