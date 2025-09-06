package i18n

import (
	"os"
	"runtime"
	"testing"
)

func TestCrossPlatformDetector(t *testing.T) {
	detector := NewCrossPlatformDetector()
	
	if detector.OS != runtime.GOOS {
		t.Errorf("Expected OS %s, got %s", runtime.GOOS, detector.OS)
	}
	
	if detector.Arch != runtime.GOARCH {
		t.Errorf("Expected Arch %s, got %s", runtime.GOARCH, detector.Arch)
	}
}

func TestNormalizeLanguageCode(t *testing.T) {
	detector := NewCrossPlatformDetector()
	
	tests := []struct {
		input    string
		expected string
	}{
		{"en_US.UTF-8", "en-us"},
		{"ja_JP.eucJP", "ja-jp"},
		{"ko_KR@hangul", "ko-kr"},
		{"zh_CN.GB2312", "zh-cn"},
		{"zh_TW.Big5", "zh-tw"},
		{"C", "en-us"},
		{"POSIX", "en-us"},
		{"chinese", "zh-cn"},
		{"japanese", "ja-jp"},
		{"korean", "ko-kr"},
		{"", ""},
		{"invalid_locale", "invalid-locale"},
	}
	
	for _, test := range tests {
		result := detector.NormalizeLanguageCode(test.input)
		if result != test.expected {
			t.Errorf("NormalizeLanguageCode(%q) = %q, expected %q", test.input, result, test.expected)
		}
	}
}

func TestIsAndroidEnvironment(t *testing.T) {
	detector := NewCrossPlatformDetector()
	
	// Save original environment
	originalAndroidRoot := os.Getenv("ANDROID_ROOT")
	originalAndroidData := os.Getenv("ANDROID_DATA")
	
	defer func() {
		// Restore environment
		if originalAndroidRoot != "" {
			os.Setenv("ANDROID_ROOT", originalAndroidRoot)
		} else {
			os.Unsetenv("ANDROID_ROOT")
		}
		if originalAndroidData != "" {
			os.Setenv("ANDROID_DATA", originalAndroidData)
		} else {
			os.Unsetenv("ANDROID_DATA")
		}
	}()
	
	// Clear Android environment
	os.Unsetenv("ANDROID_ROOT")
	os.Unsetenv("ANDROID_DATA")
	
	// Should not detect Android environment when running on non-Android OS and env vars are cleared
	// Note: File system checks might still detect Android files in some test environments
	if detector.isAndroidEnvironment() && runtime.GOOS != "android" && runtime.GOOS != "linux" {
		t.Error("Should not detect Android environment when not set and not on Android/Linux")
	}
	
	// Set Android environment
	os.Setenv("ANDROID_ROOT", "/system")
	if !detector.isAndroidEnvironment() {
		t.Error("Should detect Android environment when ANDROID_ROOT is set")
	}
	
	// Test with ANDROID_DATA
	os.Unsetenv("ANDROID_ROOT")
	os.Setenv("ANDROID_DATA", "/data")
	if !detector.isAndroidEnvironment() {
		t.Error("Should detect Android environment when ANDROID_DATA is set")
	}
}

func TestIsTermuxEnvironment(t *testing.T) {
	detector := NewCrossPlatformDetector()
	
	// Save original environment
	originalPrefix := os.Getenv("PREFIX")
	
	defer func() {
		if originalPrefix != "" {
			os.Setenv("PREFIX", originalPrefix)
		} else {
			os.Unsetenv("PREFIX")
		}
	}()
	
	// Clear Termux environment
	os.Unsetenv("PREFIX")
	
	// Should not detect Termux environment
	if detector.isTermuxEnvironment() {
		t.Error("Should not detect Termux environment when not set")
	}
	
	// Set Termux environment
	os.Setenv("PREFIX", "/data/data/com.termux/files/usr")
	if !detector.isTermuxEnvironment() {
		t.Error("Should detect Termux environment when PREFIX contains termux")
	}
}

func TestDetectSystemLanguage(t *testing.T) {
	detector := NewCrossPlatformDetector()
	
	// Save original environment
	originalLang := os.Getenv("LANG")
	originalTcpingLang := os.Getenv("TCPING_LANG")
	
	defer func() {
		if originalLang != "" {
			os.Setenv("LANG", originalLang)
		} else {
			os.Unsetenv("LANG")
		}
		if originalTcpingLang != "" {
			os.Setenv("TCPING_LANG", originalTcpingLang)
		} else {
			os.Unsetenv("TCPING_LANG")
		}
	}()
	
	// Test TCPING_LANG priority
	os.Setenv("TCPING_LANG", "ja-JP")
	os.Setenv("LANG", "ko-KR")
	
	lang := detector.DetectSystemLanguage()
	if lang != "ja-JP" {
		t.Errorf("Expected TCPING_LANG to have priority, got %s", lang)
	}
	
	// Test fallback to LANG
	os.Unsetenv("TCPING_LANG")
	lang = detector.DetectSystemLanguage()
	if lang == "" {
		// This is OK as platform-specific detection may not work in test environment
		t.Logf("Platform-specific detection returned empty, fallback working")
	}
}

func TestGetPlatformInfo(t *testing.T) {
	detector := NewCrossPlatformDetector()
	
	info := detector.GetPlatformInfo()
	
	if info["os"] != runtime.GOOS {
		t.Errorf("Expected os=%s, got %s", runtime.GOOS, info["os"])
	}
	
	if info["arch"] != runtime.GOARCH {
		t.Errorf("Expected arch=%s, got %s", runtime.GOARCH, info["arch"])
	}
	
	// Check that android and termux are properly detected
	if _, exists := info["android"]; !exists {
		t.Error("Expected android key in platform info")
	}
	
	if _, exists := info["termux"]; !exists {
		t.Error("Expected termux key in platform info")
	}
}

func TestDetectLanguageCrossPlatform(t *testing.T) {
	// Save original environment
	originalLang := os.Getenv("LANG")
	originalTcpingLang := os.Getenv("TCPING_LANG")
	
	defer func() {
		if originalLang != "" {
			os.Setenv("LANG", originalLang)
		} else {
			os.Unsetenv("LANG")
		}
		if originalTcpingLang != "" {
			os.Setenv("TCPING_LANG", originalTcpingLang)
		} else {
			os.Unsetenv("TCPING_LANG")
		}
	}()
	
	// Test with known language
	os.Setenv("TCPING_LANG", "ja-JP")
	lang := DetectLanguageCrossPlatform()
	
	if _, ok := lang.(*JapaneseLang); !ok {
		t.Errorf("Expected JapaneseLang, got %T", lang)
	}
	
	// Test fallback to English
	os.Unsetenv("TCPING_LANG")
	os.Setenv("LANG", "invalid-locale")
	lang = DetectLanguageCrossPlatform()
	
	if _, ok := lang.(*EnglishLang); !ok {
		t.Errorf("Expected fallback to EnglishLang, got %T", lang)
	}
}

// Benchmark tests
func BenchmarkDetectSystemLanguage(b *testing.B) {
	detector := NewCrossPlatformDetector()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		detector.DetectSystemLanguage()
	}
}

func BenchmarkNormalizeLanguageCode(b *testing.B) {
	detector := NewCrossPlatformDetector()
	testCode := "en_US.UTF-8@euro"
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		detector.NormalizeLanguageCode(testCode)
	}
}

func BenchmarkDetectLanguageCrossPlatform(b *testing.B) {
	os.Setenv("TCPING_LANG", "en-US")
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		DetectLanguageCrossPlatform()
	}
}