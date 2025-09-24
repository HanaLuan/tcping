package i18n

// TraditionalChineseLang implements Language interface for Traditional Chinese
type TraditionalChineseLang struct{}

// Program info
func (t *TraditionalChineseLang) ProgramDescription() string {
	return "TCP/HTTP 連線測試工具"
}

func (t *TraditionalChineseLang) Copyright() string {
	return "Copyright (c) 2025. All rights reserved."
}

// Help and usage
func (t *TraditionalChineseLang) UsageDescription() string {
	return "%s 測試到目標主機的TCP連線性或HTTP/HTTPS服務回應。"
}

func (t *TraditionalChineseLang) UsageTCP() string {
	return "tcping [選項] <主機> [連接埠]                  # TCP模式 (預設連接埠: 80)"
}

func (t *TraditionalChineseLang) UsageHTTP() string {
	return "tcping -H [選項] <URI>                       # HTTP模式"
}

func (t *TraditionalChineseLang) OptionsTitle() string {
	return "選項:"
}

func (t *TraditionalChineseLang) TCPExamplesTitle() string {
	return "TCP模式範例:"
}

func (t *TraditionalChineseLang) HTTPExamplesTitle() string {
	return "HTTP模式範例:"
}

// Command line options
func (t *TraditionalChineseLang) OptForceIPv4() string {
	return "強制使用 IPv4"
}

func (t *TraditionalChineseLang) OptForceIPv6() string {
	return "強制使用 IPv6"
}

func (t *TraditionalChineseLang) OptCount() string {
	return "傳送請求的次數 (預設: 4)"
}

func (t *TraditionalChineseLang) OptPort() string {
	return "指定要連線的連接埠 (預設: 80)"
}

func (t *TraditionalChineseLang) OptInterval() string {
	return "請求間隔（毫秒）(預設: 1000毫秒)"
}

func (t *TraditionalChineseLang) OptTimeout() string {
	return "連線逾時（毫秒）(預設: 1000毫秒)"
}

func (t *TraditionalChineseLang) OptColor() string {
	return "啟用彩色輸出"
}

func (t *TraditionalChineseLang) OptVerbose() string {
	return "啟用詳細模式，顯示更多連線資訊"
}

func (t *TraditionalChineseLang) OptHTTP() string {
	return "啟用HTTP模式，測試HTTP/HTTPS服務"
}

func (t *TraditionalChineseLang) OptInsecure() string {
	return "跳過SSL/TLS憑證驗證（僅在HTTP模式下有效）"
}

func (t *TraditionalChineseLang) OptLanguage() string {
	return "設定語言 (en-US, ja-JP, ko-KR, zh-TW, zh-CN)"
}

func (t *TraditionalChineseLang) OptVersion() string {
	return "顯示版本資訊"
}

func (t *TraditionalChineseLang) OptHelp() string {
	return "顯示此說明資訊"
}

// TCP Examples
func (t *TraditionalChineseLang) ExampleBasic() string {
	return "tcping google.com                    # 基本用法 (預設連接埠 80)"
}

func (t *TraditionalChineseLang) ExampleBasicPort() string {
	return "tcping google.com 80                 # 基本用法指定連接埠"
}

func (t *TraditionalChineseLang) ExamplePortFlag() string {
	return "tcping -p 443 google.com             # 使用-p參數指定連接埠"
}

func (t *TraditionalChineseLang) ExampleIPv4() string {
	return "tcping -4 -n 5 8.8.8.8 443           # IPv4, 5次請求"
}

func (t *TraditionalChineseLang) ExampleColorVerbose() string {
	return "tcping -c -v example.com 443         # 彩色輸出和詳細模式"
}

// HTTP Examples
func (t *TraditionalChineseLang) ExampleHTTPS() string {
	return "tcping -H https://www.google.com     # 測試HTTPS服務"
}

func (t *TraditionalChineseLang) ExampleHTTP() string {
	return "tcping -H http://example.com         # 測試HTTP服務"
}

func (t *TraditionalChineseLang) ExampleHTTPCount() string {
	return "tcping -H -n 10 https://github.com   # 傳送10次HTTP請求"
}

func (t *TraditionalChineseLang) ExampleHTTPVerbose() string {
	return "tcping -H -v https://api.github.com  # 詳細模式，顯示回應資訊"
}

func (t *TraditionalChineseLang) ExampleHTTPInsecure() string {
	return "tcping -H -k https://self-signed.badssl.com  # 跳過SSL憑證驗證"
}

// Version info
func (t *TraditionalChineseLang) VersionFormat() string {
	return "%s 版本 %s\n"
}

// Error messages
func (t *TraditionalChineseLang) ErrorPrefix() string {
	return "錯誤: %v\n"
}

func (t *TraditionalChineseLang) ErrorInvalidPort() string {
	return "連接埠號格式無效"
}

func (t *TraditionalChineseLang) ErrorPortRange() string {
	return "連接埠號必須在 1 到 65535 之間"
}

func (t *TraditionalChineseLang) ErrorIPv6Decimal() string {
	return "IPv6 地址不支援十進位格式"
}

func (t *TraditionalChineseLang) ErrorIPv6Hex() string {
	return "IPv6 地址不支援十六進位格式"
}

func (t *TraditionalChineseLang) ErrorResolve() string {
	return "解析 %s 失敗: %v"
}

func (t *TraditionalChineseLang) ErrorNoIP() string {
	return "未找到 %s 的 IP 地址"
}

func (t *TraditionalChineseLang) ErrorNoIPv4() string {
	return "未找到 %s 的 IPv4 地址"
}

func (t *TraditionalChineseLang) ErrorNoIPv6() string {
	return "未找到 %s 的 IPv6 地址"
}

func (t *TraditionalChineseLang) ErrorBothIPv4IPv6() string {
	return "無法同時使用 -4 和 -6 標誌"
}

func (t *TraditionalChineseLang) ErrorNegativeInterval() string {
	return "間隔時間不能為負值"
}

func (t *TraditionalChineseLang) ErrorNegativeTimeout() string {
	return "逾時時間不能為負值"
}

func (t *TraditionalChineseLang) ErrorHostRequired() string {
	return "需要提供主機參數\n\n用法: tcping [選項] <主機> [連接埠]\n嘗試 'tcping -h' 取得更多資訊"
}

func (t *TraditionalChineseLang) ErrorPortMustBeInRange() string {
	return "連接埠號必須在 1 到 65535 之間"
}

func (t *TraditionalChineseLang) ErrorHTTPModeURIRequired() string {
	return "HTTP模式需要提供URI參數\n\n用法: tcping -H [選項] <URI>\n嘗試 'tcping -h' 取得更多資訊"
}

func (t *TraditionalChineseLang) ErrorInvalidURI() string {
	return "無效的URI格式: %v"
}

func (t *TraditionalChineseLang) ErrorURIMustStartWithHTTP() string {
	return "URI必須以http://或https://開頭"
}

// Runtime messages
func (t *TraditionalChineseLang) MsgTCPPingStart() string {
	return "正在對 %s (%s - %s) 連接埠 %s 執行 TCP Ping\n"
}

func (t *TraditionalChineseLang) MsgHTTPPingStart() string {
	return "正在對 %s 執行 HTTP Ping (User-Agent: tcping/%s.%s)\n"
}

func (t *TraditionalChineseLang) MsgInterrupted() string {
	return "\n操作被中斷。\n"
}

func (t *TraditionalChineseLang) MsgOperationCanceled() string {
	return "\n操作被中斷, 連線嘗試已中止\n"
}

func (t *TraditionalChineseLang) MsgHTTPOperationCanceled() string {
	return "\n操作被中斷, HTTP請求已中止\n"
}

func (t *TraditionalChineseLang) MsgConnectionTimeout() string {
	return "連線逾時"
}

// Connection messages
func (t *TraditionalChineseLang) MsgTCPConnectionFailed() string {
	return "TCP連線失敗 %s:%s: seq=%d 錯誤=%s\n"
}

func (t *TraditionalChineseLang) MsgTCPConnectionSuccess() string {
	return "從 %s:%s 收到回應: seq=%d time=%.2fms\n"
}

func (t *TraditionalChineseLang) MsgHTTPRequestFailed() string {
	return "HTTP請求建立失敗 %s: seq=%d 錯誤=%v\n"
}

func (t *TraditionalChineseLang) MsgHTTPRequestFailedExec() string {
	return "HTTP請求失敗 %s: seq=%d 錯誤=%v\n"
}

func (t *TraditionalChineseLang) MsgHTTPResponseFailed() string {
	return "HTTP回應讀取失敗 %s: seq=%d 錯誤=%v\n"
}

func (t *TraditionalChineseLang) MsgHTTPResponse() string {
	return "HTTP %d %s: seq=%d time=%.2fms size=%d bytes bandwidth=%.2f Mbps\n"
}

// Verbose messages
func (t *TraditionalChineseLang) MsgVerboseDetails() string {
	return "  詳細資訊: 連線嘗試耗時 %.2fms, 目標 %s:%s\n"
}

func (t *TraditionalChineseLang) MsgVerboseConnection() string {
	return "  詳細資訊: 本地地址=%s, 遠端地址=%s:%s\n"
}

func (t *TraditionalChineseLang) MsgVerboseHTTP() string {
	return "  詳細資訊: 狀態=%s, Content-Type=%s, Server=%s\n"
}

func (t *TraditionalChineseLang) MsgVerboseHTTPDetails() string {
	return "  詳細資訊:\n"
}

func (t *TraditionalChineseLang) MsgVerboseHTTPStatus() string {
	return "    狀態: %s\n"
}

func (t *TraditionalChineseLang) MsgVerboseHTTPHeaders() string {
	return "    回應標頭:\n"
}

// Statistics
func (t *TraditionalChineseLang) MsgTCPStatisticsTitle() string {
	return "\n\n--- 目標主機 TCP ping 統計 ---\n"
}

func (t *TraditionalChineseLang) MsgHTTPStatisticsTitle() string {
	return "\n\n--- HTTP ping 統計 ---\n"
}

func (t *TraditionalChineseLang) MsgStatisticsSummary() string {
	return "已傳送 = %d, 已接收 = %d, 遺失 = %d (%.1f%% 遺失)\n"
}

func (t *TraditionalChineseLang) MsgStatisticsRTT() string {
	return "往返時間(RTT): 最小 = %.2fms, 最大 = %.2fms, 平均 = %.2fms\n"
}

func (t *TraditionalChineseLang) MsgStatisticsTotalData() string {
	return "總傳輸資料: %d bytes (%.2f MB)\n"
}

func (t *TraditionalChineseLang) MsgStatisticsBandwidth() string {
	return "估算頻寬: 最小 = %.2f Mbps, 最大 = %.2f Mbps, 平均 = %.2f Mbps\n"
}

// IP type strings
func (t *TraditionalChineseLang) IPv4String() string {
	return "IPv4"
}

func (t *TraditionalChineseLang) IPv6String() string {
	return "IPv6"
}