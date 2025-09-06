package i18n

import (
	"os"
	"runtime"
	"strings"
)

// CrossPlatformLanguageDetector provides platform-specific language detection
type CrossPlatformLanguageDetector struct {
	OS   string
	Arch string
}

// NewCrossPlatformDetector creates a new cross-platform language detector
func NewCrossPlatformDetector() *CrossPlatformLanguageDetector {
	return &CrossPlatformLanguageDetector{
		OS:   runtime.GOOS,
		Arch: runtime.GOARCH,
	}
}

// DetectSystemLanguage detects language using platform-specific methods
func (c *CrossPlatformLanguageDetector) DetectSystemLanguage() string {
	// Priority order:
	// 1. TCPING_LANG (application-specific)
	// 2. Command line argument (handled elsewhere)
	// 3. Platform-specific detection
	// 4. Standard environment variables
	// 5. Default to English

	// Check application-specific environment variable first
	if lang := os.Getenv("TCPING_LANG"); lang != "" {
		return lang
	}

	// Platform-specific detection
	var platformLang string
	switch c.OS {
	case "windows":
		platformLang = c.detectWindowsLanguage()
	case "android":
		platformLang = c.detectAndroidLanguage()
	case "linux":
		// Check if we're in Android environment
		if c.isAndroidEnvironment() {
			platformLang = c.detectAndroidLanguage()
		} else {
			platformLang = c.detectLinuxLanguage()
		}
	case "freebsd", "openbsd", "netbsd", "dragonfly":
		platformLang = c.detectBSDLanguage()
	case "darwin":
		platformLang = c.detectDarwinLanguage()
	default:
		platformLang = c.detectUnixLanguage()
	}

	if platformLang != "" {
		return platformLang
	}

	// Fallback to standard environment variables
	return c.detectFromEnvironment()
}

// detectWindowsLanguage is implemented in platform_windows.go

// detectAndroidLanguage is implemented in platform_android.go

// Platform-specific detection methods are implemented in platform_*.go files

// detectFromEnvironment is implemented in platform-specific files

// isAndroidEnvironment checks if we're running in an Android environment
func (c *CrossPlatformLanguageDetector) isAndroidEnvironment() bool {
	// Check for Android-specific indicators
	if _, err := os.Stat("/system/build.prop"); err == nil {
		return true
	}
	if _, err := os.Stat("/android_root"); err == nil {
		return true
	}
	if os.Getenv("ANDROID_ROOT") != "" {
		return true
	}
	if os.Getenv("ANDROID_DATA") != "" {
		return true
	}
	return false
}

// isTermuxEnvironment checks if we're running in Termux
func (c *CrossPlatformLanguageDetector) isTermuxEnvironment() bool {
	return os.Getenv("PREFIX") != "" && strings.Contains(os.Getenv("PREFIX"), "termux")
}

// Windows-specific methods implemented in platform_windows.go

// Platform-specific methods are implemented in platform_*.go files

// Locale parsing methods are implemented in platform-specific files

// NormalizeLanguageCode normalizes language code across platforms
func (c *CrossPlatformLanguageDetector) NormalizeLanguageCode(code string) string {
	if code == "" {
		return ""
	}
	
	// Remove charset and modifiers (en_US.UTF-8@euro -> en_US)
	if idx := strings.Index(code, "."); idx != -1 {
		code = code[:idx]
	}
	if idx := strings.Index(code, "@"); idx != -1 {
		code = code[:idx]
	}
	
	// Convert underscores to hyphens and lowercase
	code = strings.ToLower(strings.Replace(code, "_", "-", -1))
	
	// Handle common variations and mappings
	switch code {
	case "c", "posix":
		return "en-us"
	case "chinese", "china":
		return "zh-cn"
	case "chinese-traditional", "taiwan":
		return "zh-tw"
	case "japanese", "japan":
		return "ja-jp"
	case "korean", "korea":
		return "ko-kr"
	case "english":
		return "en-us"
	}
	
	return code
}

// GetPlatformInfo returns platform-specific information for debugging
func (c *CrossPlatformLanguageDetector) GetPlatformInfo() map[string]string {
	info := make(map[string]string)
	
	info["os"] = c.OS
	info["arch"] = c.Arch
	info["android"] = "false"
	info["termux"] = "false"
	
	if c.isAndroidEnvironment() {
		info["android"] = "true"
	}
	
	if c.isTermuxEnvironment() {
		info["termux"] = "true"
	}
	
	// Add environment variables for debugging
	envVars := []string{"LANG", "LC_ALL", "LC_MESSAGES", "TCPING_LANG", 
		"ANDROID_ROOT", "ANDROID_DATA", "PREFIX", "TERMUX_LANG"}
	
	for _, env := range envVars {
		if val := os.Getenv(env); val != "" {
			info[env] = val
		}
	}
	
	return info
}

// Windows API methods are implemented in platform_windows.go

// DetectLanguageCrossPlatform provides the main cross-platform detection function
func DetectLanguageCrossPlatform() Language {
	detector := NewCrossPlatformDetector()
	
	// Get language code using platform-specific detection
	langCode := detector.DetectSystemLanguage()
	
	// Normalize the language code
	normalizedCode := detector.NormalizeLanguageCode(langCode)
	
	// Get language instance
	if lang := GetLanguageByCode(normalizedCode); lang != nil {
		return lang
	}
	
	// Default to English
	return &EnglishLang{}
}