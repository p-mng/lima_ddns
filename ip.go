package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

func IP(mode string) (string, error) {
	if mode != "4" && mode != "6" {
		return "", fmt.Errorf("invalid IP mode: IPv%s", mode)
	}

	var res *http.Response
	var err error

	if mode == "4" {
		res, err = http.Get("https://ipv4.icanhazip.com")
	} else {
		res, err = http.Get("https://ipv6.icanhazip.com")
	}
	if err != nil {
		return "", err
	}
	defer func(r io.ReadCloser) {
		if err := r.Close(); err != nil {
			log.Printf("Warning: Error closing icanhazip API response: %s", err.Error())
		}
	}(res.Body)

	if res.StatusCode != 200 {
		return "", fmt.Errorf("API request failed with status code %d", res.StatusCode)
	}
	if res.Body == nil {
		return "", fmt.Errorf("API request returned no body")
	}

	raw, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	ip := strings.TrimSpace(string(raw))
	if ip == "" {
		return "", fmt.Errorf("API request returned no IP")
	}
	if (mode == "4" && strings.Contains(ip, ":")) || (mode == "6" && strings.Contains(ip, ".")) {
		return "", fmt.Errorf("API request returned invalid IPv%s: %s", mode, ip)
	}

	return ip, nil
}

func IPv4() (string, error) { return IP("4") }

func IPv6() (string, error) { return IP("6") }
