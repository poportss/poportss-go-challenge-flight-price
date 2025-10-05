package util

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"time"
)

func NewHTTPClient(timeout time.Duration) *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
	}
	return &http.Client{Timeout: timeout, Transport: tr}
}

func Ensure2xx(resp *http.Response) error {
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("upstream status %d", resp.StatusCode)
	}
	return nil
}

func EnvOr(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}
