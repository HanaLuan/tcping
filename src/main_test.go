package main

import (
	"testing"
	"time"
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

// 测试HTTP统计更新
func TestStatisticsUpdateHTTP(t *testing.T) {
	var s Statistics
	
	// 测试第一次更新
	s.updateHTTP(100, 1024, true)
	sent, responded, totalBytes, minTime, maxTime, avgTime, minBW, maxBW, avgBW := s.getHTTPStats()
	
	if sent != 1 {
		t.Errorf("sent = %d, want 1", sent)
	}
	if responded != 1 {
		t.Errorf("responded = %d, want 1", responded)
	}
	if totalBytes != 1024 {
		t.Errorf("totalBytes = %d, want 1024", totalBytes)
	}
	if minTime != 100 {
		t.Errorf("minTime = %v, want 100", minTime)
	}
	if maxTime != 100 {
		t.Errorf("maxTime = %v, want 100", maxTime)
	}
	if avgTime != 100 {
		t.Errorf("avgTime = %v, want 100", avgTime)
	}
	
	// 计算期望的带宽: (1024 * 8) / (100 * 1000) = 0.08192 Mbps
	expectedBW := float64(1024*8) / (100 * 1000)
	if minBW != expectedBW {
		t.Errorf("minBW = %v, want %v", minBW, expectedBW)
	}
	if maxBW != expectedBW {
		t.Errorf("maxBW = %v, want %v", maxBW, expectedBW)
	}
	if avgBW != expectedBW {
		t.Errorf("avgBW = %v, want %v", avgBW, expectedBW)
	}
	
	// 测试失败的请求
	s.updateHTTP(50, 0, false)
	sent, responded, _, _, _, _, _, _, _ = s.getHTTPStats()
	if sent != 2 {
		t.Errorf("sent = %d, want 2", sent)
	}
	if responded != 1 {
		t.Errorf("responded = %d, want 1 (failed request should not increment)", responded)
	}
	
	// 测试多次成功请求
	s.updateHTTP(200, 2048, true)
	sent, responded, totalBytes, minTime, maxTime, avgTime, _, _, _ = s.getHTTPStats()
	if sent != 3 {
		t.Errorf("sent = %d, want 3", sent)
	}
	if responded != 2 {
		t.Errorf("responded = %d, want 2", responded)
	}
	if totalBytes != 3072 {
		t.Errorf("totalBytes = %d, want 3072", totalBytes)
	}
	if minTime != 100 {
		t.Errorf("minTime = %v, want 100", minTime)
	}
	if maxTime != 200 {
		t.Errorf("maxTime = %v, want 200", maxTime)
	}
	if avgTime != 150 {
		t.Errorf("avgTime = %v, want 150", avgTime)
	}
}

// 测试并发更新
func TestStatisticsConcurrentUpdateHTTP(t *testing.T) {
	var s Statistics
	done := make(chan bool)
	
	// 启动多个goroutine并发更新
	for i := 0; i < 10; i++ {
		go func(id int) {
			time.Sleep(time.Millisecond * time.Duration(id))
			s.updateHTTP(float64(100+id*10), int64(1024*(id+1)), true)
			done <- true
		}(i)
	}
	
	// 等待所有goroutine完成
	for i := 0; i < 10; i++ {
		<-done
	}
	
	sent, responded, totalBytes, _, _, _, _, _, _ := s.getHTTPStats()
	if sent != 10 {
		t.Errorf("sent = %d, want 10", sent)
	}
	if responded != 10 {
		t.Errorf("responded = %d, want 10", responded)
	}
	
	// 总字节数应该是 1024 + 2048 + 3072 + ... + 10240 = 1024 * (1+2+...+10) = 1024 * 55
	expectedBytes := int64(1024 * 55)
	if totalBytes != expectedBytes {
		t.Errorf("totalBytes = %d, want %d", totalBytes, expectedBytes)
	}
}
