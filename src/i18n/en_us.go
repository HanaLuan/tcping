package i18n

// EnglishLang implements Language interface for English (US)
type EnglishLang struct{}

// Program info
func (e *EnglishLang) ProgramDescription() string {
	return "TCP/HTTP Connection Test Tool"
}

func (e *EnglishLang) Copyright() string {
	return "Copyright (c) 2025. All rights reserved."
}

// Help and usage
func (e *EnglishLang) UsageDescription() string {
	return "%s tests TCP connectivity or HTTP/HTTPS service response to target hosts."
}

func (e *EnglishLang) UsageTCP() string {
	return "tcping [options] <host> [port]                  # TCP mode (default port: 80)"
}

func (e *EnglishLang) UsageHTTP() string {
	return "tcping -H [options] <URI>                       # HTTP mode"
}

func (e *EnglishLang) OptionsTitle() string {
	return "Options:"
}

func (e *EnglishLang) TCPExamplesTitle() string {
	return "TCP Mode Examples:"
}

func (e *EnglishLang) HTTPExamplesTitle() string {
	return "HTTP Mode Examples:"
}

// Command line options
func (e *EnglishLang) OptForceIPv4() string {
	return "Force IPv4"
}

func (e *EnglishLang) OptForceIPv6() string {
	return "Force IPv6"
}

func (e *EnglishLang) OptCount() string {
	return "Number of requests to send (default: 4)"
}

func (e *EnglishLang) OptPort() string {
	return "Specify the port to connect to (default: 80)"
}

func (e *EnglishLang) OptInterval() string {
	return "Request interval in milliseconds (default: 1000ms)"
}

func (e *EnglishLang) OptTimeout() string {
	return "Connection timeout in milliseconds (default: 1000ms)"
}

func (e *EnglishLang) OptColor() string {
	return "Enable colored output"
}

func (e *EnglishLang) OptVerbose() string {
	return "Enable verbose mode, show more connection details"
}

func (e *EnglishLang) OptHTTP() string {
	return "Enable HTTP mode to test HTTP/HTTPS services"
}

func (e *EnglishLang) OptInsecure() string {
	return "Skip SSL/TLS certificate verification (HTTP mode only)"
}

func (e *EnglishLang) OptVersion() string {
	return "Show version information"
}

func (e *EnglishLang) OptHelp() string {
	return "Show this help information"
}

// TCP Examples
func (e *EnglishLang) ExampleBasic() string {
	return "tcping google.com                    # Basic usage (default port 80)"
}

func (e *EnglishLang) ExampleBasicPort() string {
	return "tcping google.com 80                 # Basic usage with port specified"
}

func (e *EnglishLang) ExamplePortFlag() string {
	return "tcping -p 443 google.com             # Use -p flag to specify port"
}

func (e *EnglishLang) ExampleIPv4() string {
	return "tcping -4 -n 5 8.8.8.8 443           # IPv4, 5 requests"
}

func (e *EnglishLang) ExampleColorVerbose() string {
	return "tcping -c -v example.com 443         # Colored output and verbose mode"
}

// HTTP Examples
func (e *EnglishLang) ExampleHTTPS() string {
	return "tcping -H https://www.google.com     # Test HTTPS service"
}

func (e *EnglishLang) ExampleHTTP() string {
	return "tcping -H http://example.com         # Test HTTP service"
}

func (e *EnglishLang) ExampleHTTPCount() string {
	return "tcping -H -n 10 https://github.com   # Send 10 HTTP requests"
}

func (e *EnglishLang) ExampleHTTPVerbose() string {
	return "tcping -H -v https://api.github.com  # Verbose mode, show response details"
}

func (e *EnglishLang) ExampleHTTPInsecure() string {
	return "tcping -H -k https://self-signed.badssl.com  # Skip SSL certificate verification"
}

// Version info
func (e *EnglishLang) VersionFormat() string {
	return "%s version %s\n"
}

// Error messages
func (e *EnglishLang) ErrorPrefix() string {
	return "Error: %v\n"
}

func (e *EnglishLang) ErrorInvalidPort() string {
	return "Invalid port number format"
}

func (e *EnglishLang) ErrorPortRange() string {
	return "Port number must be between 1 and 65535"
}

func (e *EnglishLang) ErrorIPv6Decimal() string {
	return "IPv6 addresses do not support decimal format"
}

func (e *EnglishLang) ErrorIPv6Hex() string {
	return "IPv6 addresses do not support hexadecimal format"
}

func (e *EnglishLang) ErrorResolve() string {
	return "Failed to resolve %s: %v"
}

func (e *EnglishLang) ErrorNoIP() string {
	return "No IP address found for %s"
}

func (e *EnglishLang) ErrorNoIPv4() string {
	return "No IPv4 address found for %s"
}

func (e *EnglishLang) ErrorNoIPv6() string {
	return "No IPv6 address found for %s"
}

func (e *EnglishLang) ErrorBothIPv4IPv6() string {
	return "Cannot use both -4 and -6 flags"
}

func (e *EnglishLang) ErrorNegativeInterval() string {
	return "Interval time cannot be negative"
}

func (e *EnglishLang) ErrorNegativeTimeout() string {
	return "Timeout cannot be negative"
}

func (e *EnglishLang) ErrorHostRequired() string {
	return "Host parameter is required\n\nUsage: tcping [options] <host> [port]\nTry 'tcping -h' for more information"
}

func (e *EnglishLang) ErrorPortMustBeInRange() string {
	return "Port number must be between 1 and 65535"
}

func (e *EnglishLang) ErrorHTTPModeURIRequired() string {
	return "HTTP mode requires URI parameter\n\nUsage: tcping -H [options] <URI>\nTry 'tcping -h' for more information"
}

func (e *EnglishLang) ErrorInvalidURI() string {
	return "Invalid URI format: %v"
}

func (e *EnglishLang) ErrorURIMustStartWithHTTP() string {
	return "URI must start with http:// or https://"
}

// Runtime messages
func (e *EnglishLang) MsgTCPPingStart() string {
	return "TCPing to %s (%s - %s) port %s\n"
}

func (e *EnglishLang) MsgHTTPPingStart() string {
	return "HTTP Ping to %s (User-Agent: tcping/%s.%s)\n"
}

func (e *EnglishLang) MsgInterrupted() string {
	return "\nOperation interrupted.\n"
}

func (e *EnglishLang) MsgOperationCanceled() string {
	return "\nOperation interrupted, connection attempt aborted\n"
}

func (e *EnglishLang) MsgHTTPOperationCanceled() string {
	return "\nOperation interrupted, HTTP request aborted\n"
}

func (e *EnglishLang) MsgConnectionTimeout() string {
	return "Connection timeout"
}

// Connection messages
func (e *EnglishLang) MsgTCPConnectionFailed() string {
	return "TCP connection failed %s:%s: seq=%d error=%s\n"
}

func (e *EnglishLang) MsgTCPConnectionSuccess() string {
	return "Response from %s:%s: seq=%d time=%.2fms\n"
}

func (e *EnglishLang) MsgHTTPRequestFailed() string {
	return "HTTP request creation failed %s: seq=%d error=%v\n"
}

func (e *EnglishLang) MsgHTTPRequestFailedExec() string {
	return "HTTP request failed %s: seq=%d error=%v\n"
}

func (e *EnglishLang) MsgHTTPResponseFailed() string {
	return "HTTP response read failed %s: seq=%d error=%v\n"
}

func (e *EnglishLang) MsgHTTPResponse() string {
	return "HTTP %d %s: seq=%d time=%.2fms size=%d bytes bandwidth=%.2f Mbps\n"
}

// Verbose messages
func (e *EnglishLang) MsgVerboseDetails() string {
	return "  Details: Connection attempt took %.2fms, target %s:%s\n"
}

func (e *EnglishLang) MsgVerboseConnection() string {
	return "  Details: Local address=%s, Remote address=%s:%s\n"
}

func (e *EnglishLang) MsgVerboseHTTP() string {
	return "  Details: Status=%s, Content-Type=%s, Server=%s\n"
}

// Statistics
func (e *EnglishLang) MsgTCPStatisticsTitle() string {
	return "\n\n--- TCP ping statistics ---\n"
}

func (e *EnglishLang) MsgHTTPStatisticsTitle() string {
	return "\n\n--- HTTP ping statistics ---\n"
}

func (e *EnglishLang) MsgStatisticsSummary() string {
	return "Sent = %d, Received = %d, Lost = %d (%.1f%% loss)\n"
}

func (e *EnglishLang) MsgStatisticsRTT() string {
	return "Round-trip times: Min = %.2fms, Max = %.2fms, Avg = %.2fms\n"
}

func (e *EnglishLang) MsgStatisticsTotalData() string {
	return "Total data transferred: %d bytes (%.2f MB)\n"
}

func (e *EnglishLang) MsgStatisticsBandwidth() string {
	return "Estimated bandwidth: Min = %.2f Mbps, Max = %.2f Mbps, Avg = %.2f Mbps\n"
}

// IP type strings
func (e *EnglishLang) IPv4String() string {
	return "IPv4"
}

func (e *EnglishLang) IPv6String() string {
	return "IPv6"
}