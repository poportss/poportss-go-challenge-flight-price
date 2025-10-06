package util

import (
	"crypto/tls"
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

func EnvOr(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}
