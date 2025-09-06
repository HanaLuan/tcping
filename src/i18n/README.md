# TCPing i18n (Internationalization) System

This directory contains the multi-language support system for TCPing, providing translations for 5 languages with smart detection and fallback mechanisms.

## Supported Languages

| Language Code | Language | Region | Status | File |
|---------------|----------|--------|--------|------|
| `en-US` | English | United States | ✅ Complete | `en_us.go` |
| `ja-JP` | Japanese | Japan | ✅ Complete | `ja_jp.go` |
| `ko-KR` | Korean | South Korea | ✅ Complete | `ko_kr.go` |
| `zh-TW` | Traditional Chinese | Taiwan | ✅ Complete | `zh_tw.go` |
| `zh-CN` | Simplified Chinese | China | ✅ Complete | `zh_cn.go` |

## Features

### 🎯 Smart Language Detection
- **Command Line Priority**: `-l en-US` or `--lang ja-JP`
- **Environment Variables**: `TCPING_LANG` > `LC_ALL` > `LC_MESSAGES` > `LANG`
- **Automatic Fallback**: Defaults to English if language not found
- **Flexible Format**: Supports `en`, `en-US`, `en_US`, `en-us.utf-8` formats

### 🔧 Easy Integration
- **Clean Interface**: Single `Language` interface with all strings
- **Simple Usage**: `i18n.T().MsgTCPPingStart()` replaces hardcoded text
- **Thread-Safe**: Safe for concurrent access
- **Zero Dependencies**: Pure Go stdlib implementation

### 📝 Complete Coverage
- Help text and usage examples
- Error messages and validation
- Runtime status messages  
- Statistics and verbose output
- All user-facing strings

## Usage Examples

### Command Line Language Selection
```bash
# English (default)
./tcping google.com

# Japanese  
./tcping -l ja-JP google.com
./tcping --lang ja google.com

# Korean
TCPING_LANG=ko-KR ./tcping google.com

# Traditional Chinese
LANG=zh-TW ./tcping -c -v google.com 443

# Simplified Chinese (original)
LC_ALL=zh-CN ./tcping -H https://github.com
```

### Environment Variable Priority
```bash
# TCPING_LANG has highest priority
TCPING_LANG=ja-JP LANG=ko-KR ./tcping google.com  # → Japanese

# LC_ALL overrides LANG
LC_ALL=ko-KR LANG=zh-CN ./tcping google.com       # → Korean

# LANG is used if others not set
LANG=zh-TW ./tcping google.com                     # → Traditional Chinese
```

## Integration Guide

### 1. Import the Package
```go
import "./i18n"
```

### 2. Add Language Flag
```go
func setupFlags(opts *Options) {
    // ... existing flags ...
    lang := flag.String("l", "", "Set language (en-US, ja-JP, ko-KR, zh-TW, zh-CN)")
    flag.StringVar(lang, "lang", "", "Set language")
    
    // ... rest of function ...
    opts.Language = *lang
}
```

### 3. Initialize Language System
```go
func main() {
    opts := &Options{}
    setupFlags(opts)
    
    // Initialize i18n system
    i18n.Initialize(opts.Language)
    
    // ... rest of main ...
}
```

### 4. Replace Hardcoded Strings
```go
// Before:
fmt.Printf("正在对 %s (%s - %s) 端口 %s 执行 TCP Ping\n", host, ipType, ip, port)

// After:
fmt.Printf(i18n.T().MsgTCPPingStart(), host, ipType, ip, port)
```

### 5. Error Messages
```go
// Before:
fmt.Fprintf(os.Stderr, "错误: %v\n", err)

// After:
fmt.Fprintf(os.Stderr, i18n.T().ErrorPrefix(), err)
```

## File Structure

```
src/i18n/
├── README.md           # This documentation
├── i18n.go            # Core interface and detection logic
├── i18n_test.go       # Unit tests
├── en_us.go           # English translations (default)
├── ja_jp.go           # Japanese translations
├── ko_kr.go           # Korean translations  
├── zh_tw.go           # Traditional Chinese translations
├── zh_cn.go           # Simplified Chinese translations (original)
└── example_integration.go  # Integration examples
```

## Language Interface

Each language implements the `Language` interface with methods for:

```go
type Language interface {
    // Program info
    ProgramDescription() string
    Copyright() string
    
    // Help and usage
    UsageDescription() string
    UsageTCP() string
    UsageHTTP() string
    
    // Command line options  
    OptForceIPv4() string
    OptCount() string
    // ... 50+ methods total
    
    // Runtime messages
    MsgTCPPingStart() string
    MsgConnectionTimeout() string
    
    // Statistics
    MsgStatisticsSummary() string
    MsgStatisticsRTT() string
}
```

## Testing

Run tests to verify language system:

```bash
cd src/i18n
go test -v

# Test specific functionality
go test -run TestDetectLanguage
go test -run TestLanguageStrings
```

## Output Examples

### English (en-US)
```
TCPing to google.com (IPv4 - 8.8.8.8) port 80
Response from 8.8.8.8:80: seq=0 time=15.23ms
```

### Japanese (ja-JP)  
```
google.com (IPv4 - 8.8.8.8) ポート80にTCP Ping実行中
8.8.8.8:80からレスポンス: seq=0 time=15.23ms
```

### Korean (ko-KR)
```
google.com (IPv4 - 8.8.8.8) 포트 80에 TCP Ping 실행 중
8.8.8.8:80에서 응답: seq=0 time=15.23ms
```

### Traditional Chinese (zh-TW)
```
正在對 google.com (IPv4 - 8.8.8.8) 連接埠 80 執行 TCP Ping
從 8.8.8.8:80 收到回應: seq=0 time=15.23ms
```

### Simplified Chinese (zh-CN)
```
正在对 google.com (IPv4 - 8.8.8.8) 端口 80 执行 TCP Ping  
从 8.8.8.8:80 收到响应: seq=0 time=15.23ms
```

## Adding New Languages

1. Create new language file (e.g., `fr_fr.go`)
2. Implement `Language` interface
3. Add language detection in `GetLanguageByCode()`
4. Update documentation
5. Add tests

Example template:
```go
package i18n

type FrenchLang struct{}

func (f *FrenchLang) ProgramDescription() string {
    return "Outil de test de connexion TCP/HTTP"
}

// ... implement all interface methods
```

## Performance

- **Zero Allocation**: String lookups use method calls, no map lookups
- **Compile-Time Safety**: All strings checked at compile time
- **Minimal Overhead**: Direct method calls, no reflection
- **Memory Efficient**: Only active language loaded

## Best Practices

1. **Use `i18n.T()`**: Always use the helper function for consistency
2. **Format Strings**: Use printf-style placeholders for dynamic content
3. **Context Aware**: Keep cultural context in mind for translations
4. **Test All Languages**: Verify output in all supported languages
5. **Consistent Style**: Maintain consistent terminology across languages

## Migration from Hardcoded Text

The system is designed for gradual migration:

1. **Phase 1**: Initialize i18n system (no code changes needed)
2. **Phase 2**: Replace critical error messages  
3. **Phase 3**: Replace help text and usage
4. **Phase 4**: Replace all runtime messages
5. **Phase 5**: Remove original hardcoded strings

The default English translations match common Unix tool conventions, while other languages adapt to local preferences and technical terminology standards.