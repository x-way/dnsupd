package main

import (
	"fmt"
	"net/http"
	"testing"
)

func TestValidHostname(t *testing.T) {
	tt := []struct {
		hostname string
		zone     string
		hosts    []HostConfig
		expected bool
	}{
		{"name.zone", "zone", []HostConfig{{"name", "", ""}}, true},
		{"name.otherzone", "zone", []HostConfig{{"name", "", ""}}, false},
		{"namewithnozone", "zone", []HostConfig{{"name", "", ""}}, false},
		{"sub.name.zone", "zone", []HostConfig{{"name", "", ""}}, false},
	}

	for _, tc := range tt {
		config.Hosts = tc.hosts
		config.Zone = tc.zone
		got := validHostname(tc.hostname)
		if got != tc.expected {
			t.Errorf("validHostname failed for '%s', got %v, expected %v", tc.hostname, got, tc.expected)
		}
	}
}

func TestValidAuth(t *testing.T) {
	tt := []struct {
		hostname string
		user     string
		password string
		zone     string
		hosts    []HostConfig
		expected bool
	}{
		{"name.zone", "user", "pass", "zone", []HostConfig{{"name", "user", "pass"}}, true},
		{"name.zone", "otheruser", "pass", "zone", []HostConfig{{"name", "user", "pass"}}, false},
		{"name.zone", "otheruser", "pass", "zone", []HostConfig{{"name", "user", "pass"}, {"name2", "otheruser", "pass"}}, false},
		{"name2.zone", "otheruser", "pass", "zone", []HostConfig{{"name", "user", "pass"}, {"name2", "otheruser", "pass"}}, true},
		{"name.zone", "user", "otherpass", "zone", []HostConfig{{"name", "user", "pass"}}, false},
		{"name.otherzone", "user", "pass", "zone", []HostConfig{{"name", "user", "pass"}}, false},
		{"namewithnozone", "user", "pass", "zone", []HostConfig{{"name", "user", "pass"}}, false},
		{"sub.name.zone", "user", "pass", "zone", []HostConfig{{"name", "user", "pass"}}, false},
	}

	for _, tc := range tt {
		config.Hosts = tc.hosts
		config.Zone = tc.zone
		got := validAuth(tc.hostname, tc.user, tc.password)
		if got != tc.expected {
			t.Errorf("validAuth failed for '%s', '%s', '%s' got %v, expected %v", tc.hostname, tc.user, tc.password, got, tc.expected)
		}
	}
}

func TestGetParameters(t *testing.T) {
	tt := []struct {
		ipheader         string
		myip             string
		myipheaderheader string
		myipheader       string
		hostname         string
		expectedHostname string
		expectedMyipstr  string
		expectedRrtype   string
		expectedOk       bool
	}{
		{"ip-in-header-only", "", "ip-in-header-only", "1.2.3.4", "myhostname", "myhostname", "1.2.3.4", "A", true},
		{"ip-not-in-header", "4.5.6.7", "ip-not-in-header", "1.2.3.4", "myhostname", "myhostname", "4.5.6.7", "A", true},
		{"ip-not-in-header", "4.5.6.7", "some-header", "", "myhostname", "myhostname", "4.5.6.7", "A", true},
		{"ipv6-not-in-header", "2001:db8::1234", "some-header", "", "myhostname", "myhostname", "2001:db8::1234", "AAAA", true},
		{"ipv6-in-header-only", "", "ipv6-in-header-only", "2001:db8::5", "myhostname", "myhostname", "2001:db8::5", "AAAA", true},
		{"non-ip-value", "xyz", "some-header", "", "myhostname", "myhostname", "xyz", "", false},
		{"non-ip-header-value", "", "non-ip-header-value", "xyz", "myhostname", "myhostname", "xyz", "", false},
		{"", "2001:db8::1234", "some-header", "1.2.3.4", "myhostname", "myhostname", "2001:db8::1234", "AAAA", true},
		{"", "4.5.6.7", "some-header", "1.2.3.4", "myhostname", "myhostname", "4.5.6.7", "A", true},
		{"", "", "some-header", "4.5.6.7", "myhostname", "myhostname", "", "", false},
	}

	for _, tc := range tt {
		config.IPHeader = tc.ipheader

		url := fmt.Sprintf("https://foo.bar/?myip=%s&hostname=%s", tc.myip, tc.hostname)
		req, err := http.NewRequest(http.MethodGet, url, nil)
		req.Header.Set(tc.myipheaderheader, tc.myipheader)
		if err != nil {
			t.Errorf("getParameters failed to build http.NewRequest for '%s', '%s', '%s' got error: %v", tc.ipheader, tc.myip, tc.hostname, err)
		}

		gotHostname, gotMyipstr, gotRrtype, gotOk := getParameters(req)
		if gotHostname != tc.expectedHostname {
			t.Errorf("getParameters hostname failed for '%s', '%s', '%s' got %v, expected %v", tc.ipheader, tc.myip, tc.hostname, gotHostname, tc.expectedHostname)
		}
		if gotMyipstr != tc.expectedMyipstr {
			t.Errorf("getParameters myipstr failed for '%s', '%s', '%s' got %v, expected %v", tc.ipheader, tc.myip, tc.hostname, gotMyipstr, tc.expectedMyipstr)
		}
		if gotRrtype != tc.expectedRrtype {
			t.Errorf("getParameters rrtype failed for '%s', '%s', '%s' got %v, expected %v", tc.ipheader, tc.myip, tc.hostname, gotRrtype, tc.expectedRrtype)
		}
		if gotOk != tc.expectedOk {
			t.Errorf("getParameters ok failed for '%s', '%s', '%s' got %v, expected %v", tc.ipheader, tc.myip, tc.hostname, gotOk, tc.expectedOk)
		}
	}
}
