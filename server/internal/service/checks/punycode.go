package checks

import (
	"strings"
)

func ContainsPunycode(rawURL string) (bool, error) {
	host, err := GetHost(rawURL)
	if err != nil {
		return false, err
	}

	for _, label := range strings.Split(host, ".") {
		if strings.HasPrefix(label, "xn--") {
			return true, nil
		}
		// Raw Unicode IDN (not yet punycode-encoded) — any non-ASCII rune flags it
		for _, r := range label {
			if r > 127 {
				return true, nil
			}
		}
	}
	return false, nil
}
