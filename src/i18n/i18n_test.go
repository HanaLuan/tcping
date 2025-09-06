package i18n

import (
	"os"
	"testing"
)

func TestGetLanguageByCode(t *testing.T) {
	tests := []struct {
		code     string
		expected string
	}{
		{"en", "EnglishLang"},
		{"en-US", "EnglishLang"},
		{"ja", "JapaneseLang"},
		{"ja-JP", "JapaneseLang"},
		{"ko", "KoreanLang"},
		{"ko-KR", "KoreanLang"},
		{"zh-TW", "TraditionalChineseLang"},
		{"zh-tw", "TraditionalChineseLang"},
		{"zh-CN", "SimplifiedChineseLang"},
		{"zh-cn", "SimplifiedChineseLang"},
		{"zh", "SimplifiedChineseLang"},
		{"invalid", "EnglishLang"}, // Should fallback to English
	}

	for _, test := range tests {
		lang := GetLanguageByCode(test.code)
		langType := ""
		
		switch lang.(type) {
		case *EnglishLang:
			langType = "EnglishLang"
		case *JapaneseLang:
			langType = "JapaneseLang"
		case *KoreanLang:
			langType = "KoreanLang"
		case *TraditionalChineseLang:
			langType = "TraditionalChineseLang"
		case *SimplifiedChineseLang:
			langType = "SimplifiedChineseLang"
		}
		
		if langType != test.expected {
			t.Errorf("GetLanguageByCode(%q) returned %s, expected %s", test.code, langType, test.expected)
		}
	}
}

func TestDetectLanguage(t *testing.T) {
	// Save original environment
	originalLang := os.Getenv("LANG")
	originalTcpingLang := os.Getenv("TCPING_LANG")
	originalLcAll := os.Getenv("LC_ALL")
	
	// Clean environment
	os.Unsetenv("LANG")
	os.Unsetenv("TCPING_LANG")
	os.Unsetenv("LC_ALL")
	os.Unsetenv("LC_MESSAGES")
	
	defer func() {
		// Restore environment
		if originalLang != "" {
			os.Setenv("LANG", originalLang)
		}
		if originalTcpingLang != "" {
			os.Setenv("TCPING_LANG", originalTcpingLang)
		}
		if originalLcAll != "" {
			os.Setenv("LC_ALL", originalLcAll)
		}
	}()
	
	// Test default (should be English)
	lang := DetectLanguage()
	if _, ok := lang.(*EnglishLang); !ok {
		t.Error("DetectLanguage() should return EnglishLang by default")
	}
	
	// Test TCPING_LANG priority
	os.Setenv("TCPING_LANG", "ja-JP")
	os.Setenv("LANG", "ko-KR")
	lang = DetectLanguage()
	if _, ok := lang.(*JapaneseLang); !ok {
		t.Error("DetectLanguage() should prioritize TCPING_LANG over LANG")
	}
	
	// Test LC_ALL priority over LANG
	os.Unsetenv("TCPING_LANG")
	os.Setenv("LC_ALL", "ko-KR")
	os.Setenv("LANG", "zh-CN")
	lang = DetectLanguage()
	if _, ok := lang.(*KoreanLang); !ok {
		t.Error("DetectLanguage() should prioritize LC_ALL over LANG")
	}
}

func TestLanguageStrings(t *testing.T) {
	languages := []Language{
		&EnglishLang{},
		&JapaneseLang{},
		&KoreanLang{},
		&TraditionalChineseLang{},
		&SimplifiedChineseLang{},
	}
	
	for _, lang := range languages {
		// Test that all strings are non-empty
		if lang.ProgramDescription() == "" {
			t.Errorf("ProgramDescription() should not be empty for %T", lang)
		}
		if lang.ErrorPrefix() == "" {
			t.Errorf("ErrorPrefix() should not be empty for %T", lang)
		}
		if lang.MsgTCPPingStart() == "" {
			t.Errorf("MsgTCPPingStart() should not be empty for %T", lang)
		}
		if lang.IPv4String() == "" {
			t.Errorf("IPv4String() should not be empty for %T", lang)
		}
		if lang.IPv6String() == "" {
			t.Errorf("IPv6String() should not be empty for %T", lang)
		}
	}
}

func TestInitialize(t *testing.T) {
	// Test explicit language setting
	Initialize("ja-JP")
	if _, ok := GetLanguage().(*JapaneseLang); !ok {
		t.Error("Initialize('ja-JP') should set Japanese language")
	}
	
	// Test empty string (should use detection)
	os.Setenv("LANG", "ko-KR")
	Initialize("")
	if _, ok := GetLanguage().(*KoreanLang); !ok {
		t.Error("Initialize('') should use language detection")
	}
	
	// Reset to default
	Initialize("en-US")
}