package i18n

// SimplifiedChineseLang implements Language interface for Simplified Chinese (original language)
type SimplifiedChineseLang struct{}

// Program info
func (s *SimplifiedChineseLang) ProgramDescription() string {
	return "TCP/HTTP 连接测试工具"
}

func (s *SimplifiedChineseLang) Copyright() string {
	return "Copyright (c) 2025. All rights reserved."
}

// Help and usage
func (s *SimplifiedChineseLang) UsageDescription() string {
	return "%s 测试到目标主机的TCP连接性或HTTP/HTTPS服务响应。"
}

func (s *SimplifiedChineseLang) UsageTCP() string {
	return "tcping [选项] <主机> [端口]                  # TCP模式 (默认端口: 80)"
}

func (s *SimplifiedChineseLang) UsageHTTP() string {
	return "tcping -H [选项] <URI>                       # HTTP模式"
}

func (s *SimplifiedChineseLang) OptionsTitle() string {
	return "选项:"
}

func (s *SimplifiedChineseLang) TCPExamplesTitle() string {
	return "TCP模式示例:"
}

func (s *SimplifiedChineseLang) HTTPExamplesTitle() string {
	return "HTTP模式示例:"
}

// Command line options
func (s *SimplifiedChineseLang) OptForceIPv4() string {
	return "强制使用 IPv4"
}

func (s *SimplifiedChineseLang) OptForceIPv6() string {
	return "强制使用 IPv6"
}

func (s *SimplifiedChineseLang) OptCount() string {
	return "发送请求的次数 (默认: 4)"
}

func (s *SimplifiedChineseLang) OptPort() string {
	return "指定要连接的端口 (默认: 80)"
}

func (s *SimplifiedChineseLang) OptInterval() string {
	return "请求间隔（毫秒）(默认: 1000毫秒)"
}

func (s *SimplifiedChineseLang) OptTimeout() string {
	return "连接超时（毫秒）(默认: 1000毫秒)"
}

func (s *SimplifiedChineseLang) OptColor() string {
	return "启用彩色输出"
}

func (s *SimplifiedChineseLang) OptVerbose() string {
	return "启用详细模式，显示更多连接信息"
}

func (s *SimplifiedChineseLang) OptHTTP() string {
	return "启用HTTP模式，测试HTTP/HTTPS服务"
}

func (s *SimplifiedChineseLang) OptInsecure() string {
	return "跳过SSL/TLS证书验证（仅在HTTP模式下有效）"
}

func (s *SimplifiedChineseLang) OptVersion() string {
	return "显示版本信息"
}

func (s *SimplifiedChineseLang) OptHelp() string {
	return "显示此帮助信息"
}

// TCP Examples
func (s *SimplifiedChineseLang) ExampleBasic() string {
	return "tcping google.com                    # 基本用法 (默认端口 80)"
}

func (s *SimplifiedChineseLang) ExampleBasicPort() string {
	return "tcping google.com 80                 # 基本用法指定端口"
}

func (s *SimplifiedChineseLang) ExamplePortFlag() string {
	return "tcping -p 443 google.com             # 使用-p参数指定端口"
}

func (s *SimplifiedChineseLang) ExampleIPv4() string {
	return "tcping -4 -n 5 8.8.8.8 443           # IPv4, 5次请求"
}

func (s *SimplifiedChineseLang) ExampleColorVerbose() string {
	return "tcping -c -v example.com 443         # 彩色输出和详细模式"
}

// HTTP Examples
func (s *SimplifiedChineseLang) ExampleHTTPS() string {
	return "tcping -H https://www.google.com     # 测试HTTPS服务"
}

func (s *SimplifiedChineseLang) ExampleHTTP() string {
	return "tcping -H http://example.com         # 测试HTTP服务"
}

func (s *SimplifiedChineseLang) ExampleHTTPCount() string {
	return "tcping -H -n 10 https://github.com   # 发送10次HTTP请求"
}

func (s *SimplifiedChineseLang) ExampleHTTPVerbose() string {
	return "tcping -H -v https://api.github.com  # 详细模式，显示响应信息"
}

func (s *SimplifiedChineseLang) ExampleHTTPInsecure() string {
	return "tcping -H -k https://self-signed.badssl.com  # 跳过SSL证书验证"
}

// Version info
func (s *SimplifiedChineseLang) VersionFormat() string {
	return "%s 版本 %s\n"
}

// Error messages
func (s *SimplifiedChineseLang) ErrorPrefix() string {
	return "错误: %v\n"
}

func (s *SimplifiedChineseLang) ErrorInvalidPort() string {
	return "端口号格式无效"
}

func (s *SimplifiedChineseLang) ErrorPortRange() string {
	return "端口号必须在 1 到 65535 之间"
}

func (s *SimplifiedChineseLang) ErrorIPv6Decimal() string {
	return "IPv6 地址不支持十进制格式"
}

func (s *SimplifiedChineseLang) ErrorIPv6Hex() string {
	return "IPv6 地址不支持十六进制格式"
}

func (s *SimplifiedChineseLang) ErrorResolve() string {
	return "解析 %s 失败: %v"
}

func (s *SimplifiedChineseLang) ErrorNoIP() string {
	return "未找到 %s 的 IP 地址"
}

func (s *SimplifiedChineseLang) ErrorNoIPv4() string {
	return "未找到 %s 的 IPv4 地址"
}

func (s *SimplifiedChineseLang) ErrorNoIPv6() string {
	return "未找到 %s 的 IPv6 地址"
}

func (s *SimplifiedChineseLang) ErrorBothIPv4IPv6() string {
	return "无法同时使用 -4 和 -6 标志"
}

func (s *SimplifiedChineseLang) ErrorNegativeInterval() string {
	return "间隔时间不能为负值"
}

func (s *SimplifiedChineseLang) ErrorNegativeTimeout() string {
	return "超时时间不能为负值"
}

func (s *SimplifiedChineseLang) ErrorHostRequired() string {
	return "需要提供主机参数\n\n用法: tcping [选项] <主机> [端口]\n尝试 'tcping -h' 获取更多信息"
}

func (s *SimplifiedChineseLang) ErrorPortMustBeInRange() string {
	return "端口号必须在 1 到 65535 之间"
}

func (s *SimplifiedChineseLang) ErrorHTTPModeURIRequired() string {
	return "HTTP模式需要提供URI参数\n\n用法: tcping -H [选项] <URI>\n尝试 'tcping -h' 获取更多信息"
}

func (s *SimplifiedChineseLang) ErrorInvalidURI() string {
	return "无效的URI格式: %v"
}

func (s *SimplifiedChineseLang) ErrorURIMustStartWithHTTP() string {
	return "URI必须以http://或https://开头"
}

// Runtime messages
func (s *SimplifiedChineseLang) MsgTCPPingStart() string {
	return "正在对 %s (%s - %s) 端口 %s 执行 TCP Ping\n"
}

func (s *SimplifiedChineseLang) MsgHTTPPingStart() string {
	return "正在对 %s 执行 HTTP Ping (User-Agent: tcping/%s.%s)\n"
}

func (s *SimplifiedChineseLang) MsgInterrupted() string {
	return "\n操作被中断。\n"
}

func (s *SimplifiedChineseLang) MsgOperationCanceled() string {
	return "\n操作被中断, 连接尝试已中止\n"
}

func (s *SimplifiedChineseLang) MsgHTTPOperationCanceled() string {
	return "\n操作被中断, HTTP请求已中止\n"
}

func (s *SimplifiedChineseLang) MsgConnectionTimeout() string {
	return "连接超时"
}

// Connection messages
func (s *SimplifiedChineseLang) MsgTCPConnectionFailed() string {
	return "TCP连接失败 %s:%s: seq=%d 错误=%s\n"
}

func (s *SimplifiedChineseLang) MsgTCPConnectionSuccess() string {
	return "从 %s:%s 收到响应: seq=%d time=%.2fms\n"
}

func (s *SimplifiedChineseLang) MsgHTTPRequestFailed() string {
	return "HTTP请求创建失败 %s: seq=%d 错误=%v\n"
}

func (s *SimplifiedChineseLang) MsgHTTPRequestFailedExec() string {
	return "HTTP请求失败 %s: seq=%d 错误=%v\n"
}

func (s *SimplifiedChineseLang) MsgHTTPResponseFailed() string {
	return "HTTP响应读取失败 %s: seq=%d 错误=%v\n"
}

func (s *SimplifiedChineseLang) MsgHTTPResponse() string {
	return "HTTP %d %s: seq=%d time=%.2fms size=%d bytes bandwidth=%.2f Mbps\n"
}

// Verbose messages
func (s *SimplifiedChineseLang) MsgVerboseDetails() string {
	return "  详细信息: 连接尝试耗时 %.2fms, 目标 %s:%s\n"
}

func (s *SimplifiedChineseLang) MsgVerboseConnection() string {
	return "  详细信息: 本地地址=%s, 远程地址=%s:%s\n"
}

func (s *SimplifiedChineseLang) MsgVerboseHTTP() string {
	return "  详细信息: 状态=%s, Content-Type=%s, Server=%s\n"
}

// Statistics
func (s *SimplifiedChineseLang) MsgTCPStatisticsTitle() string {
	return "\n\n--- 目标主机 TCP ping 统计 ---\n"
}

func (s *SimplifiedChineseLang) MsgHTTPStatisticsTitle() string {
	return "\n\n--- HTTP ping 统计 ---\n"
}

func (s *SimplifiedChineseLang) MsgStatisticsSummary() string {
	return "已发送 = %d, 已接收 = %d, 丢失 = %d (%.1f%% 丢失)\n"
}

func (s *SimplifiedChineseLang) MsgStatisticsRTT() string {
	return "往返时间(RTT): 最小 = %.2fms, 最大 = %.2fms, 平均 = %.2fms\n"
}

func (s *SimplifiedChineseLang) MsgStatisticsTotalData() string {
	return "总传输数据: %d bytes (%.2f MB)\n"
}

func (s *SimplifiedChineseLang) MsgStatisticsBandwidth() string {
	return "估算带宽: 最小 = %.2f Mbps, 最大 = %.2f Mbps, 平均 = %.2f Mbps\n"
}

// IP type strings
func (s *SimplifiedChineseLang) IPv4String() string {
	return "IPv4"
}

func (s *SimplifiedChineseLang) IPv6String() string {
	return "IPv6"
}