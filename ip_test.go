package main

import (
	"net"
	"testing"
)

func TestIPv4(t *testing.T) {
	// https://en.wikipedia.org/wiki/Google_Public_DNS#Service
	conn, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		t.Skipf("IPv4 is not available: %s", err.Error())
	}
	if err := conn.Close(); err != nil {
		t.Errorf("error closing connection: %s", err.Error())
	}

	ip, err := IPv4()
	if err != nil {
		t.Errorf("IPv4() failed with error: %s", err.Error())
	}
	if net.ParseIP(ip) == nil {
		t.Errorf("IPv4() returned invalid IP: %s", ip)
	}
}

func TestIPv6(t *testing.T) {
	// https://en.wikipedia.org/wiki/Google_Public_DNS#Service
	conn, err := net.Dial("udp", "[2001:4860:4860::8888]:53")
	if err != nil {
		t.Skipf("IPv6 is not available: %s", err.Error())
	}
	if err := conn.Close(); err != nil {
		t.Errorf("error closing connection: %s", err.Error())
	}

	ip, err := IPv6()
	if err != nil {
		t.Errorf("IPv6() failed with error: %s", err.Error())
	}
	if net.ParseIP(ip) == nil {
		t.Errorf("IPv6() returned invalid IP: %s", ip)
	}
}
