//go:build !windows

package i18n

import (
	"bufio"
	"os"
	"os/exec"
	"strings"
)

// getWindowsUserLocale is not available on Unix systems
func getWindowsUserLocale() string {
	return ""
}

// getWindowsSystemLocale is not available on Unix systems  
func getWindowsSystemLocale() string {
	return ""
}

// detectBSDLanguage detects language on BSD systems
func (c *CrossPlatformLanguageDetector) detectBSDLanguage() string {
	// Check BSD-specific locations
	if lang := c.readFileContent("/etc/locale.conf"); lang != "" {
		return c.parseLocaleConf(lang)
	}
	
	// FreeBSD specific
	if c.OS == "freebsd" {
		if lang := c.readFileContent("/etc/login.conf"); lang != "" {
			return c.parseLoginConf(lang)
		}
	}
	
	// Fallback to standard Unix detection
	return c.detectUnixLanguage()
}

// detectLinuxLanguage detects language on Linux systems
func (c *CrossPlatformLanguageDetector) detectLinuxLanguage() string {
	// Try systemd locale
	if lang := c.getSystemdLocale(); lang != "" {
		return lang
	}
	
	// Try locale files
	if lang := c.readFileContent("/etc/locale.conf"); lang != "" {
		return c.parseLocaleConf(lang)
	}
	
	if lang := c.readFileContent("/etc/default/locale"); lang != "" {
		return c.parseDefaultLocale(lang)
	}
	
	// Fallback to standard Unix detection
	return c.detectUnixLanguage()
}

// detectDarwinLanguage detects language on macOS systems  
func (c *CrossPlatformLanguageDetector) detectDarwinLanguage() string {
	// macOS-specific detection could be added here
	// For now, use Unix detection
	return c.detectUnixLanguage()
}

// detectUnixLanguage detects language using standard Unix methods
func (c *CrossPlatformLanguageDetector) detectUnixLanguage() string {
	return c.detectFromEnvironment()
}

// detectFromEnvironment detects language from environment variables (Unix implementation)
func (c *CrossPlatformLanguageDetector) detectFromEnvironment() string {
	// Standard environment variable priority
	for _, env := range []string{"LC_ALL", "LC_MESSAGES", "LANG"} {
		if lang := os.Getenv(env); lang != "" {
			return lang
		}
	}
	return ""
}

// parseLoginConf parses FreeBSD login.conf format
func (c *CrossPlatformLanguageDetector) parseLoginConf(content string) string {
	// Parse login.conf for locale settings
	// This is BSD-specific - simplified implementation
	return ""
}

// parseDefaultLocale parses /etc/default/locale format  
func (c *CrossPlatformLanguageDetector) parseDefaultLocale(content string) string {
	return c.parseLocaleConf(content) // Same format as locale.conf
}

// parseLocaleConf parses locale.conf file format
func (c *CrossPlatformLanguageDetector) parseLocaleConf(content string) string {
	// Parse LANG=en_US.UTF-8 format
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "LANG=") {
			return strings.TrimSpace(strings.TrimPrefix(line, "LANG="))
		}
	}
	return ""
}

// detectAndroidLanguage stub for non-Android Unix systems
func (c *CrossPlatformLanguageDetector) detectAndroidLanguage() string {
	// If we detect Android environment on non-Android build, fallback to Unix
	return c.detectUnixLanguage()
}

// Enhanced Unix/Linux detection
func (c *CrossPlatformLanguageDetector) detectWindowsLanguage() string {
	// Not applicable on Unix systems
	return ""
}

// readFileContent safely reads file content
func (c *CrossPlatformLanguageDetector) readFileContent(filename string) string {
	file, err := os.Open(filename)
	if err != nil {
		return ""
	}
	defer file.Close()
	
	var content strings.Builder
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		content.WriteString(scanner.Text())
		content.WriteString("\n")
	}
	
	return content.String()
}

// getAndroidSystemProperty executes getprop command
func (c *CrossPlatformLanguageDetector) getAndroidSystemProperty(prop string) string {
	// Check if getprop exists
	if _, err := os.Stat("/system/bin/getprop"); err != nil {
		return ""
	}
	
	cmd := exec.Command("getprop", prop)
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	
	result := strings.TrimSpace(string(output))
	if result == "" || result == "undefined" {
		return ""
	}
	
	return result
}

// getSystemdLocale executes localectl command
func (c *CrossPlatformLanguageDetector) getSystemdLocale() string {
	cmd := exec.Command("localectl", "status")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "LANG=") {
			return strings.TrimPrefix(line, "LANG=")
		}
		if strings.Contains(line, "System Locale:") && strings.Contains(line, "LANG=") {
			parts := strings.Split(line, "LANG=")
			if len(parts) > 1 {
				return strings.Fields(parts[1])[0]
			}
		}
	}
	
	return ""
}

// getTermuxLanguage gets language from Termux environment
func (c *CrossPlatformLanguageDetector) getTermuxLanguage() string {
	// Check Termux-specific environment variables
	for _, env := range []string{"TERMUX_LOCALE", "TERMUX_LANG"} {
		if lang := os.Getenv(env); lang != "" {
			return lang
		}
	}
	
	// Try termux-api if available
	if _, err := exec.LookPath("termux-api"); err == nil {
		// Could use termux-api to get system locale
		// For now, return empty
	}
	
	return ""
}