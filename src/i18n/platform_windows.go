//go:build windows

package i18n

import (
	"os"
	"syscall"
	"unsafe"
)

var (
	kernel32                      = syscall.NewLazyDLL("kernel32.dll")
	procGetUserDefaultLocaleName  = kernel32.NewProc("GetUserDefaultLocaleName")
	procGetSystemDefaultLocaleName = kernel32.NewProc("GetSystemDefaultLocaleName")
)

// getWindowsUserLocale gets the Windows user locale
func getWindowsUserLocale() string {
	buf := make([]uint16, 85) // LOCALE_NAME_MAX_LENGTH
	ret, _, _ := procGetUserDefaultLocaleName.Call(uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))
	
	if ret != 0 {
		return syscall.UTF16ToString(buf)
	}
	
	return ""
}

// getWindowsSystemLocale gets the Windows system locale
func getWindowsSystemLocale() string {
	buf := make([]uint16, 85) // LOCALE_NAME_MAX_LENGTH
	ret, _, _ := procGetSystemDefaultLocaleName.Call(uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))
	
	if ret != 0 {
		return syscall.UTF16ToString(buf)
	}
	
	return ""
}

// Enhanced Windows detection
func (c *CrossPlatformLanguageDetector) detectWindowsLanguage() string {
	// 1. Check environment variables first (for compatibility)
	for _, env := range []string{"LC_ALL", "LC_MESSAGES", "LANG"} {
		if lang := os.Getenv(env); lang != "" {
			return lang
		}
	}
	
	// 2. Try Windows user locale
	if userLocale := getWindowsUserLocale(); userLocale != "" {
		return userLocale
	}
	
	// 3. Try Windows system locale as fallback
	if systemLocale := getWindowsSystemLocale(); systemLocale != "" {
		return systemLocale
	}
	
	// 4. Check Windows-specific environment variables
	if lang := os.Getenv("USERPROFILE"); lang != "" {
		// Could implement registry reading here
	}
	
	return ""
}

// detectFromEnvironment detects language from environment variables (Windows implementation)
func (c *CrossPlatformLanguageDetector) detectFromEnvironment() string {
	// Standard environment variable priority (Windows)
	for _, env := range []string{"LC_ALL", "LC_MESSAGES", "LANG"} {
		if lang := os.Getenv(env); lang != "" {
			return lang
		}
	}
	return ""
}

// detectAndroidLanguage stub for Windows systems
func (c *CrossPlatformLanguageDetector) detectAndroidLanguage() string {
	// Not applicable on Windows systems
	return ""
}

// detectLinuxLanguage stub for Windows systems
func (c *CrossPlatformLanguageDetector) detectLinuxLanguage() string {
	// Not applicable on Windows systems - use Windows detection
	return c.detectWindowsLanguage()
}

// detectBSDLanguage stub for Windows systems
func (c *CrossPlatformLanguageDetector) detectBSDLanguage() string {
	// Not applicable on Windows systems - use Windows detection
	return c.detectWindowsLanguage()
}

// detectDarwinLanguage stub for Windows systems
func (c *CrossPlatformLanguageDetector) detectDarwinLanguage() string {
	// Not applicable on Windows systems - use Windows detection
	return c.detectWindowsLanguage()
}

// detectUnixLanguage stub for Windows systems
func (c *CrossPlatformLanguageDetector) detectUnixLanguage() string {
	// Not applicable on Windows systems - use Windows detection
	return c.detectWindowsLanguage()
}