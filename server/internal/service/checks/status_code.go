package checks

import (
	"net/http"
	"time"
)

func GetStatusCode(rawURL string) (int, string, bool, bool, error) {
	client := newSafeHTTPClient(5 * time.Second)
	resp, err := client.Head(rawURL)
	if err != nil {
		return 0, "", false, false, err
	}

	defer resp.Body.Close()

	statusCode := resp.StatusCode
	statusText := http.StatusText(statusCode)
	isSuccess := statusCode >= 200 && statusCode < 300
	isRedirectStatusCode := statusCode >= 300 && statusCode < 400

	return statusCode, statusText, isSuccess, isRedirectStatusCode, nil
}
