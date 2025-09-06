# Cross-Platform Language Detection

TCPing's i18n system provides sophisticated cross-platform language detection that works across different operating systems and environments.

## Platform Support Matrix

| Platform | Environment | Detection Method | Status |
|----------|------------|------------------|--------|
| **Windows** | Native | Windows API + Environment | ✅ Full |
| **Windows** | WSL/Cygwin | Environment Variables | ✅ Full |
| **Linux** | Native | systemd + config files | ✅ Full |
| **Linux** | Docker | Environment Variables | ✅ Full |
| **Android** | Native | System Properties | ✅ Full |
| **Android** | Termux | Termux API + Environment | ✅ Full |
| **Android** | ADB Shell | System Properties | ✅ Full |
| **macOS** | Native | Environment Variables | ✅ Full |
| **FreeBSD** | Native | login.conf + Environment | ✅ Full |
| **OpenBSD** | Native | Environment Variables | ✅ Full |
| **NetBSD** | Native | Environment Variables | ✅ Full |

## Detection Priority Order

### 1. Application-Specific (Highest Priority)
```bash
TCPING_LANG=ja-JP ./tcping google.com
```

### 2. Command Line Arguments
```bash
./tcping -l ko-KR google.com
./tcping --lang zh-TW google.com
```

### 3. Platform-Specific Detection

#### Windows
1. **Windows API**: `GetUserDefaultLocaleName()`
2. **System API**: `GetSystemDefaultLocaleName()`
3. **Environment Variables**: `LC_ALL`, `LC_MESSAGES`, `LANG`

#### Android/Termux
1. **Termux Environment**: `TERMUX_LANG`, `TERMUX_LOCALE`
2. **System Properties**: `persist.sys.locale`, `ro.product.locale`
3. **Settings Database**: Content provider access (if available)
4. **Environment Variables**: `ANDROID_LOCALE`, `LC_ALL`, `LANG`

#### Linux
1. **systemd**: `localectl status`
2. **Config Files**: `/etc/locale.conf`, `/etc/default/locale`
3. **Environment Variables**: `LC_ALL`, `LC_MESSAGES`, `LANG`

#### BSD Systems
1. **Config Files**: `/etc/locale.conf`, `/etc/login.conf` (FreeBSD)
2. **Environment Variables**: `LC_ALL`, `LC_MESSAGES`, `LANG`

### 4. Standard Environment Variables (Fallback)
- `LC_ALL` (overrides all)
- `LC_MESSAGES` (messages only)
- `LANG` (default locale)

### 5. Default Fallback
- **English (en-US)** if no locale detected

## Platform-Specific Behaviors

### Windows Detection

#### Native Windows
```go
// Uses Windows API
userLocale := getWindowsUserLocale()    // "en-US", "ja-JP", etc.
systemLocale := getWindowsSystemLocale() // System-wide setting
```

#### Windows Subsystems
- **WSL**: Uses Linux detection methods
- **Cygwin**: Uses environment variables
- **Git Bash**: Uses environment variables

### Android Detection

#### Native Android (ADB Shell)
```bash
# System properties
getprop persist.sys.locale          # "ja-JP"
getprop ro.product.locale           # "en-US" 
getprop ro.product.locale.language  # "ja"
```

#### Termux Environment
```bash
# Termux-specific variables
export TERMUX_LANG=ja-JP
export TERMUX_LOCALE=ja_JP.UTF-8

# Termux API (if installed)
termux-locale  # Returns system locale
```

#### Android Settings Database
```bash
# Content provider access (requires permissions)
content query --uri content://settings/system --where "name='system_locales'"
```

### Linux Detection

#### systemd Systems
```bash
localectl status
# Output: System Locale: LANG=en_US.UTF-8
```

#### Configuration Files
```bash
# /etc/locale.conf
LANG=ja_JP.UTF-8
LC_MESSAGES=ja_JP.UTF-8

# /etc/default/locale (Debian/Ubuntu)
LANG="ko_KR.UTF-8"
```

### BSD Detection

#### FreeBSD
```bash
# /etc/login.conf
default:\
    :lang=ja_JP.UTF-8:\
    :charset=UTF-8:
```

#### OpenBSD/NetBSD
- Primarily environment variable based
- System-specific config files

## Language Code Normalization

The system automatically normalizes various locale formats:

```go
// Input formats -> Normalized output
"en_US.UTF-8"     -> "en-us"
"ja_JP.eucJP"     -> "ja-jp"  
"ko_KR@hangul"    -> "ko-kr"
"zh_CN.GB2312"    -> "zh-cn"
"zh_TW.Big5"      -> "zh-tw"

// Special cases
"C"               -> "en-us"
"POSIX"           -> "en-us"
"chinese"         -> "zh-cn"
"japanese"        -> "ja-jp"
"korean"          -> "ko-kr"
```

## Usage Examples

### Environment Variables

#### Cross-Platform
```bash
# Application-specific (highest priority)
TCPING_LANG=ja-JP ./tcping google.com

# Standard Unix
LANG=ko_KR.UTF-8 ./tcping google.com
LC_ALL=zh_TW.UTF-8 ./tcping google.com

# Android/Termux
TERMUX_LANG=ja-JP ./tcping google.com
ANDROID_LOCALE=ko-KR ./tcping google.com
```

#### Platform-Specific

**Windows Command Prompt:**
```cmd
set TCPING_LANG=ja-JP
tcping.exe google.com
```

**Windows PowerShell:**
```powershell
$env:TCPING_LANG="ja-JP"
.\tcping.exe google.com
```

**Android ADB Shell:**
```bash
adb shell
export TCPING_LANG=ja-JP
./tcping google.com
```

**Termux:**
```bash
pkg install tcping
export TERMUX_LANG=ja-JP
tcping google.com
```

### Command Line Arguments

```bash
# Short form
./tcping -l ja-JP google.com
./tcping -l ko-KR -c -v google.com 443

# Long form  
./tcping --lang zh-TW google.com
./tcping --lang zh-CN -H https://github.com
```

## Debugging Language Detection

### Get Platform Information
```go
detector := i18n.NewCrossPlatformDetector()
info := detector.GetPlatformInfo()

// Returns map with:
// - "os": "windows" | "linux" | "android" | "darwin" | "freebsd"
// - "arch": "amd64" | "386" | "arm64" | "arm"
// - "android": "true" | "false"
// - "termux": "true" | "false"
// - Environment variables: "LANG", "LC_ALL", etc.
```

### Test Language Detection
```bash
# Set debug environment
export TCPING_DEBUG=1
export TCPING_LANG=ja-JP

# Run with verbose mode
./tcping -v google.com

# Check detected language
./tcping -V  # Shows version with detected language info
```

## Implementation Notes

### Build Constraints

The system uses Go build constraints for platform-specific code:

```go
//go:build windows
// platform_windows.go - Windows API calls

//go:build !windows  
// platform_unix.go - Unix/Linux implementations

//go:build android
// platform_android.go - Android-specific detection
```

### Dependencies

#### Minimal Dependencies
- **Core**: Pure Go stdlib only
- **Windows**: `golang.org/x/sys/windows` (optional, for API access)
- **Unix**: No additional dependencies

#### External Commands (Optional)
- **Linux**: `localectl` (systemd)
- **Android**: `getprop`, `content` (Android tools)
- **Termux**: `termux-api` (Termux API)

### Performance Characteristics

- **Fast**: Direct environment variable access
- **Cached**: Language detection result cached per process
- **Lightweight**: No heavy system calls in hot paths
- **Fallback**: Graceful degradation if platform detection fails

### Error Handling

The system gracefully handles:
- **Missing Commands**: Fallback to environment variables
- **Permission Errors**: Skip system files, use alternatives
- **Malformed Locales**: Normalize or default to English
- **Empty Results**: Chain fallback methods

## Troubleshooting

### Common Issues

#### Windows
```bash
# Issue: No language detected on Windows
# Solution: Set environment variable
set LANG=ja-JP

# Or use Windows locale setting
# Control Panel > Region > Administrative > Change system locale
```

#### Android/Termux
```bash
# Issue: Wrong language in Termux
# Solution: Set Termux-specific variable
export TERMUX_LANG=ja-JP

# Or install Termux API
pkg install termux-api
```

#### Linux
```bash
# Issue: Language not detected on headless system
# Solution: Set locale properly
export LANG=ja_JP.UTF-8
localectl set-locale LANG=ja_JP.UTF-8
```

### Debug Logging

Enable debug mode to see detection process:

```bash
# Environment variable
export TCPING_DEBUG_I18N=1
./tcping google.com

# This will show:
# - Platform detected: linux/amd64
# - Android environment: false  
# - Termux environment: false
# - Environment variables checked
# - Language detection result
```

## Contributing

### Adding Platform Support

1. **Create platform file**: `platform_newos.go`
2. **Add build constraints**: `//go:build newos`
3. **Implement detection methods**
4. **Add tests**: `TestNewOSDetection`
5. **Update documentation**

### Testing

```bash
# Run all cross-platform tests
go test -v ./src/i18n -run CrossPlatform

# Test specific platform
go test -v ./src/i18n -run TestWindows

# Benchmark detection performance
go test -bench=. ./src/i18n
```

The cross-platform language detection system ensures TCPing works seamlessly across all supported platforms while respecting user preferences and system settings.