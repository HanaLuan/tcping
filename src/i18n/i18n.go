package i18n

import (
	"strings"
)

// Language interface defines all translatable strings
type Language interface {
	// Program info
	ProgramDescription() string
	Copyright() string
	
	// Help and usage
	UsageDescription() string
	UsageTCP() string
	UsageHTTP() string
	OptionsTitle() string
	TCPExamplesTitle() string
	HTTPExamplesTitle() string
	
	// Command line options
	OptForceIPv4() string
	OptForceIPv6() string
	OptCount() string
	OptPort() string
	OptInterval() string
	OptTimeout() string
	OptColor() string
	OptVerbose() string
	OptHTTP() string
	OptInsecure() string
	OptLanguage() string
	OptVersion() string
	OptHelp() string
	
	// TCP Examples
	ExampleBasic() string
	ExampleBasicPort() string
	ExamplePortFlag() string
	ExampleIPv4() string
	ExampleColorVerbose() string
	
	// HTTP Examples
	ExampleHTTPS() string
	ExampleHTTP() string
	ExampleHTTPCount() string
	ExampleHTTPVerbose() string
	ExampleHTTPInsecure() string
	
	// Version info
	VersionFormat() string
	
	// Error messages
	ErrorPrefix() string
	ErrorInvalidPort() string
	ErrorPortRange() string
	ErrorIPv6Decimal() string
	ErrorIPv6Hex() string
	ErrorResolve() string
	ErrorNoIP() string
	ErrorNoIPv4() string
	ErrorNoIPv6() string
	ErrorBothIPv4IPv6() string
	ErrorNegativeInterval() string
	ErrorNegativeTimeout() string
	ErrorHostRequired() string
	ErrorPortMustBeInRange() string
	ErrorHTTPModeURIRequired() string
	ErrorInvalidURI() string
	ErrorURIMustStartWithHTTP() string
	
	// Runtime messages
	MsgTCPPingStart() string          // "正在对 %s (%s - %s) 端口 %s 执行 TCP Ping\n"
	MsgHTTPPingStart() string         // "正在对 %s 执行 HTTP Ping (User-Agent: tcping/%s.%s)\n"
	MsgInterrupted() string           // "\n操作被中断。\n"
	MsgOperationCanceled() string     // "\n操作被中断, 连接尝试已中止\n"
	MsgHTTPOperationCanceled() string // "\n操作被中断, HTTP请求已中止\n"
	MsgConnectionTimeout() string     // "连接超时"
	
	// Connection messages
	MsgTCPConnectionFailed() string   // "TCP连接失败 %s:%s: seq=%d 错误=%s\n"
	MsgTCPConnectionSuccess() string  // "从 %s:%s 收到响应: seq=%d time=%.2fms\n"
	MsgHTTPRequestFailed() string     // "HTTP请求创建失败 %s: seq=%d 错误=%v\n"
	MsgHTTPRequestFailedExec() string // "HTTP请求失败 %s: seq=%d 错误=%v\n"
	MsgHTTPResponseFailed() string    // "HTTP响应读取失败 %s: seq=%d 错误=%v\n"
	MsgHTTPResponse() string          // "HTTP %d %s: seq=%d time=%.2fms size=%d bytes bandwidth=%.2f Mbps\n"
	
	// Verbose messages
	MsgVerboseDetails() string        // "  详细信息: 连接尝试耗时 %.2fms, 目标 %s:%s\n"
	MsgVerboseConnection() string     // "  详细信息: 本地地址=%s, 远程地址=%s:%s\n"
	MsgVerboseHTTP() string          // "  详细信息: 状态=%s, Content-Type=%s, Server=%s\n"
	MsgVerboseHTTPDetails() string   // "  Details:\n"
	MsgVerboseHTTPStatus() string    // "    Status: %s\n"
	MsgVerboseHTTPHeaders() string   // "    Response Headers:\n"
	
	// Statistics
	MsgTCPStatisticsTitle() string    // "\n\n--- 目标主机 TCP ping 统计 ---\n"
	MsgHTTPStatisticsTitle() string   // "\n\n--- HTTP ping 统计 ---\n"
	MsgStatisticsSummary() string     // "已发送 = %d, 已接收 = %d, 丢失 = %d (%.1f%% 丢失)\n"
	MsgStatisticsRTT() string         // "往返时间(RTT): 最小 = %.2fms, 最大 = %.2fms, 平均 = %.2fms\n"
	MsgStatisticsTotalData() string   // "总传输数据: %d bytes (%.2f MB)\n"
	MsgStatisticsBandwidth() string   // "估算带宽: 最小 = %.2f Mbps, 最大 = %.2f Mbps, 平均 = %.2f Mbps\n"
	
	// IP type strings
	IPv4String() string               // "IPv4"
	IPv6String() string               // "IPv6"
}

// Global language instance
var currentLang Language

// Supported languages
const (
	LangEnUS = "en-US"
	LangJaJP = "ja-JP"
	LangKoKR = "ko-KR"
	LangZhTW = "zh-TW" 
	LangZhCN = "zh-CN"
)

// GetLanguage returns the current language instance
func GetLanguage() Language {
	if currentLang == nil {
		// Default to English if not set
		currentLang = &EnglishLang{}
	}
	return currentLang
}

// SetLanguage sets the current language
func SetLanguage(lang Language) {
	currentLang = lang
}

// DetectLanguage detects system language and returns appropriate Language instance
func DetectLanguage() Language {
	// Use cross-platform detection
	return DetectLanguageCrossPlatform()
}

// GetLanguageByCode returns Language instance by language code
func GetLanguageByCode(code string) Language {
	// Normalize the code - handle various formats
	code = strings.ToLower(strings.Replace(code, "_", "-", -1))
	
	// Extract primary language if full locale (e.g., "en-us.utf-8" -> "en-us")
	if dotIndex := strings.Index(code, "."); dotIndex != -1 {
		code = code[:dotIndex]
	}
	
	switch code {
	case "en", "en-us":
		return &EnglishLang{}
	case "ja", "ja-jp":
		return &JapaneseLang{}
	case "ko", "ko-kr":
		return &KoreanLang{}
	case "zh-tw", "zh-hant":
		return &TraditionalChineseLang{}
	case "zh", "zh-cn", "zh-hans":
		return &SimplifiedChineseLang{}
	default:
		return &EnglishLang{} // Default fallback
	}
}

// Initialize sets up the language system
func Initialize(langCode string) {
	var lang Language
	if langCode != "" {
		lang = GetLanguageByCode(langCode)
	} else {
		lang = DetectLanguage()
	}
	SetLanguage(lang)
}

// Helper function to get current language strings
func T() Language {
	return GetLanguage()
}