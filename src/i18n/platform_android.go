//go:build android

package i18n

import (
	"os"
	"os/exec"
	"strings"
)

// detectAndroidLanguage detects language in Android environments (Termux, adb shell)
func (c *CrossPlatformLanguageDetector) detectAndroidLanguage() string {
	// 1. Check Termux-specific environment first
	if c.isTermuxEnvironment() {
		if lang := c.getTermuxLanguage(); lang != "" {
			return lang
		}
	}
	
	// 2. Try Android system properties
	androidProps := []string{
		"persist.sys.locale",
		"ro.product.locale",
		"ro.product.locale.language",
		"persist.sys.language",
		"ro.config.locale",
	}
	
	for _, prop := range androidProps {
		if lang := c.getAndroidSystemProperty(prop); lang != "" {
			return lang
		}
	}
	
	// 3. Check Android-specific environment variables
	for _, env := range []string{"ANDROID_LOCALE", "ANDROID_LANG"} {
		if lang := os.Getenv(env); lang != "" {
			return lang
		}
	}
	
	// 4. Try to read Android settings database (if accessible)
	if lang := c.getAndroidSettingsLocale(); lang != "" {
		return lang
	}
	
	// 5. Fallback to standard Unix detection
	return c.detectUnixLanguage()
}

// getAndroidSettingsLocale attempts to read Android settings
func (c *CrossPlatformLanguageDetector) getAndroidSettingsLocale() string {
	// Try to access Android settings database
	// This may require root or specific permissions
	
	// Try content provider interface if available
	if _, err := exec.LookPath("content"); err == nil {
		cmd := exec.Command("content", "query", "--uri", "content://settings/system", "--where", "name='system_locales'")
		output, err := cmd.Output()
		if err == nil && len(output) > 0 {
			// Parse content provider output
			result := strings.TrimSpace(string(output))
			if result != "" {
				return c.parseAndroidContentOutput(result)
			}
		}
	}
	
	return ""
}

// parseAndroidContentOutput parses Android content provider output
func (c *CrossPlatformLanguageDetector) parseAndroidContentOutput(output string) string {
	// Parse Android content provider output format
	// This is a simplified parser - real implementation would be more robust
	
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.Contains(line, "value=") {
			parts := strings.Split(line, "value=")
			if len(parts) > 1 {
				value := strings.TrimSpace(parts[1])
				value = strings.Trim(value, "\"'")
				if value != "" {
					return value
				}
			}
		}
	}
	
	return ""
}

// isTermuxEnvironment enhanced check for Termux
func (c *CrossPlatformLanguageDetector) isTermuxEnvironment() bool {
	// Check multiple Termux indicators
	indicators := []string{
		"PREFIX",
		"TERMUX_VERSION",
		"TERMUX_APP_PID",
	}
	
	for _, indicator := range indicators {
		if val := os.Getenv(indicator); val != "" {
			if strings.Contains(val, "termux") || strings.Contains(val, "/data/data/com.termux") {
				return true
			}
		}
	}
	
	// Check if we're running in Termux directory structure
	if prefix := os.Getenv("PREFIX"); prefix != "" {
		if strings.Contains(prefix, "/data/data/com.termux") {
			return true
		}
	}
	
	// Check for Termux-specific files
	termuxFiles := []string{
		"/data/data/com.termux/files/usr/bin/termux-info",
		"/data/data/com.termux/files/usr/etc/termux.conf",
	}
	
	for _, file := range termuxFiles {
		if _, err := os.Stat(file); err == nil {
			return true
		}
	}
	
	return false
}

// detectUnixLanguage fallback for Android
func (c *CrossPlatformLanguageDetector) detectUnixLanguage() string {
	return c.detectFromEnvironment()
}

// detectFromEnvironment detects language from environment variables
func (c *CrossPlatformLanguageDetector) detectFromEnvironment() string {
	// Standard environment variable priority
	for _, env := range []string{"LC_ALL", "LC_MESSAGES", "LANG"} {
		if lang := os.Getenv(env); lang != "" {
			return lang
		}
	}
	return ""
}