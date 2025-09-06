package main

import (
	"fmt"
	"os"
	"tcping/src/i18n"
)

func main() {
	// Test different language codes
	languages := []string{"en-US", "ja-JP", "ko-KR", "zh-TW", "zh-CN"}
	
	for _, langCode := range languages {
		fmt.Printf("\n=== Testing Language: %s ===\n", langCode)
		
		// Initialize with specific language
		i18n.Initialize(langCode)
		
		// Test various strings
		lang := i18n.T()
		
		fmt.Printf("Program Description: %s\n", lang.ProgramDescription())
		fmt.Printf("Version Format: %s\n", fmt.Sprintf(lang.VersionFormat(), "TCPing", "v1.8.0"))
		fmt.Printf("HTTP Mode Option: %s\n", lang.OptHTTP())
		fmt.Printf("SSL Insecure Option: %s\n", lang.OptInsecure())
		fmt.Printf("TCP Connection Success: %s\n", fmt.Sprintf(lang.MsgTCPConnectionSuccess(), "8.8.8.8", "80", 1, 15.23))
		fmt.Printf("Error Prefix: %s\n", fmt.Sprintf(lang.ErrorPrefix(), "Test error message"))
		fmt.Printf("Statistics Title: %s\n", lang.MsgTCPStatisticsTitle())
	}
	
	// Test environment variable detection
	fmt.Printf("\n=== Testing Environment Detection ===\n")
	
	// Test with TCPING_LANG
	os.Setenv("TCPING_LANG", "ja-JP")
	i18n.Initialize("")
	fmt.Printf("TCPING_LANG=ja-JP detected: %s\n", i18n.T().ProgramDescription())
	
	// Test with LANG fallback
	os.Unsetenv("TCPING_LANG")
	os.Setenv("LANG", "ko-KR")
	i18n.Initialize("")
	fmt.Printf("LANG=ko-KR detected: %s\n", i18n.T().ProgramDescription())
	
	// Test cross-platform detector
	fmt.Printf("\n=== Testing Cross-Platform Detection ===\n")
	detector := i18n.NewCrossPlatformDetector()
	info := detector.GetPlatformInfo()
	
	fmt.Printf("Platform Info:\n")
	for key, value := range info {
		fmt.Printf("  %s: %s\n", key, value)
	}
	
	// Test language code normalization
	fmt.Printf("\n=== Testing Language Code Normalization ===\n")
	testCodes := []string{
		"en_US.UTF-8",
		"ja_JP.eucJP@japanese", 
		"ko_KR.UTF-8",
		"zh_CN.GB2312",
		"C",
		"POSIX",
		"chinese",
	}
	
	for _, code := range testCodes {
		normalized := detector.NormalizeLanguageCode(code)
		fmt.Printf("%s -> %s\n", code, normalized)
	}
}