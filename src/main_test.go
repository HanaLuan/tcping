package main

import (
	"testing"
)

func TestValidatePort(t *testing.T) {
	tests := []struct {
		port    string
		wantErr bool
	}{
		{"80", false},
		{"0", true},
		{"65536", true},
		{"22", false},
		{"abc", true},
	}
	for _, tt := range tests {
		err := validatePort(tt.port)
		if (err != nil) != tt.wantErr {
			t.Errorf("validatePort(%q) error = %v, wantErr %v", tt.port, err, tt.wantErr)
		}
	}
}

func TestParseNumericIPv4(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"3232235777", "192.168.1.1"},
		{"0xc0a80101", "192.168.1.1"},
		{"0x08080808", "8.8.8.8"},
		{"134744072", "8.8.8.8"},
		{"notanip", ""},
	}
	for _, tt := range tests {
		ip := parseNumericIPv4(tt.input)
		if tt.want == "" && ip != nil {
			t.Errorf("parseNumericIPv4(%q) = %v, want nil", tt.input, ip)
		}
		if tt.want != "" && (ip == nil || ip.String() != tt.want) {
			t.Errorf("parseNumericIPv4(%q) = %v, want %v", tt.input, ip, tt.want)
		}
	}
}

func TestIsIPv4(t *testing.T) {
	if !isIPv4("8.8.8.8") {
		t.Error("isIPv4 failed for 8.8.8.8")
	}
	if !isIPv4("134744072") {
		t.Error("isIPv4 failed for decimal IPv4")
	}
	if isIPv4("::1") {
		t.Error("isIPv4 failed for IPv6")
	}
}

func TestIsIPv6(t *testing.T) {
	if !isIPv6("::1") {
		t.Error("isIPv6 failed for ::1")
	}
	if isIPv6("8.8.8.8") {
		t.Error("isIPv6 failed for 8.8.8.8")
	}
}

func TestColorize(t *testing.T) {
	txt := "hello"
	colored := colorize(txt, "31", true)
	if colored == txt {
		t.Error("colorize should add color codes when useColor is true")
	}
	if colorize(txt, "31", false) != txt {
		t.Error("colorize should not add color codes when useColor is false")
	}
}

func TestStatisticsUpdateAndGet(t *testing.T) {
	var s Statistics
	s.update(10, true)
	s.update(20, true)
	s.update(30, false)
	s.update(5, true)
	sent, responded, min, max, avg := s.getStats()
	if sent != 4 {
		t.Errorf("sent = %d, want 4", sent)
	}
	if responded != 3 {
		t.Errorf("responded = %d, want 3", responded)
	}
	if min != 5 {
		t.Errorf("min = %v, want 5", min)
	}
	if max != 20 {
		t.Errorf("max = %v, want 20", max)
	}
	if avg < 11.6 || avg > 11.7 {
		t.Errorf("avg = %v, want about 11.666", avg)
	}
}

func TestResolveAddress_IPv4(t *testing.T) {
	ip, err := resolveAddress("8.8.8.8", true, false)
	if err != nil {
		t.Errorf("resolveAddress IPv4 error: %v", err)
	}
	if ip != "8.8.8.8" {
		t.Errorf("resolveAddress IPv4 got %v", ip)
	}
}

func TestResolveAddress_IPv6(t *testing.T) {
	// ::1 is always available as loopback
	ip, err := resolveAddress("::1", false, true)
	if err != nil {
		t.Errorf("resolveAddress IPv6 error: %v", err)
	}
	if ip != "[::1]" {
		t.Errorf("resolveAddress IPv6 got %v", ip)
	}
}

func TestResolveAddress_Invalid(t *testing.T) {
	_, err := resolveAddress("notarealhost.invalid", true, false)
	if err == nil {
		t.Error("resolveAddress should fail for invalid host")
	}
}

func TestValidateOptions(t *testing.T) {
	opts := &Options{UseIPv4: true, UseIPv6: true}
	_, _, err := validateOptions(opts, []string{"8.8.8.8"})
	if err == nil {
		t.Error("validateOptions should fail when both IPv4 and IPv6 are set")
	}
	opts = &Options{UseIPv4: true, UseIPv6: false, Interval: -1}
	_, _, err = validateOptions(opts, []string{"8.8.8.8"})
	if err == nil {
		t.Error("validateOptions should fail for negative interval")
	}
	opts = &Options{UseIPv4: true, UseIPv6: false, Interval: 1000, Timeout: -1}
	_, _, err = validateOptions(opts, []string{"8.8.8.8"})
	if err == nil {
		t.Error("validateOptions should fail for negative timeout")
	}
	opts = &Options{UseIPv4: true, UseIPv6: false, Interval: 1000, Timeout: 1000}
	_, _, err = validateOptions(opts, []string{})
	if err == nil {
		t.Error("validateOptions should fail for missing host")
	}
	opts = &Options{UseIPv4: true, UseIPv6: false, Interval: 1000, Timeout: 1000, Port: 70000}
	_, _, err = validateOptions(opts, []string{"8.8.8.8"})
	if err == nil {
		t.Error("validateOptions should fail for invalid port")
	}
	opts = &Options{UseIPv4: true, UseIPv6: false, Interval: 1000, Timeout: 1000, Port: 80}
	_, _, err = validateOptions(opts, []string{"8.8.8.8"})
	if err != nil {
		t.Errorf("validateOptions failed: %v", err)
	}
}
