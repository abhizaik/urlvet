package checks

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"
)

func SupportsHSTS(rawURL string) (bool, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return false, fmt.Errorf("invalid URL: %v", err)
	}

	if parsedURL.Scheme != "https" {
		parsedURL.Scheme = "https"
	}

	req, err := http.NewRequest("HEAD", parsedURL.String(), nil)
	if err != nil {
		return false, err
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			DialContext:     ssrfSafeDialContext(&net.Dialer{Timeout: 5 * time.Second}),
		},
		Timeout: 5 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	_, ok := resp.Header["Strict-Transport-Security"]
	return ok, nil
}
