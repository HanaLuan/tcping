package i18n

// KoreanLang implements Language interface for Korean
type KoreanLang struct{}

// Program info
func (k *KoreanLang) ProgramDescription() string {
	return "TCP/HTTP 연결 테스트 도구"
}

func (k *KoreanLang) Copyright() string {
	return "Copyright (c) 2025. All rights reserved."
}

// Help and usage
func (k *KoreanLang) UsageDescription() string {
	return "%s는 대상 호스트에 대한 TCP 연결성 또는 HTTP/HTTPS 서비스 응답을 테스트합니다."
}

func (k *KoreanLang) UsageTCP() string {
	return "tcping [옵션] <호스트> [포트]                  # TCP 모드 (기본 포트: 80)"
}

func (k *KoreanLang) UsageHTTP() string {
	return "tcping -H [옵션] <URI>                        # HTTP 모드"
}

func (k *KoreanLang) OptionsTitle() string {
	return "옵션:"
}

func (k *KoreanLang) TCPExamplesTitle() string {
	return "TCP 모드 예시:"
}

func (k *KoreanLang) HTTPExamplesTitle() string {
	return "HTTP 모드 예시:"
}

// Command line options
func (k *KoreanLang) OptForceIPv4() string {
	return "IPv4 강제 사용"
}

func (k *KoreanLang) OptForceIPv6() string {
	return "IPv6 강제 사용"
}

func (k *KoreanLang) OptCount() string {
	return "전송할 요청 수 (기본값: 4)"
}

func (k *KoreanLang) OptPort() string {
	return "연결할 포트 지정 (기본값: 80)"
}

func (k *KoreanLang) OptInterval() string {
	return "요청 간격(밀리초) (기본값: 1000ms)"
}

func (k *KoreanLang) OptTimeout() string {
	return "연결 타임아웃(밀리초) (기본값: 1000ms)"
}

func (k *KoreanLang) OptColor() string {
	return "컬러 출력 활성화"
}

func (k *KoreanLang) OptVerbose() string {
	return "상세 모드 활성화, 더 많은 연결 정보 표시"
}

func (k *KoreanLang) OptHTTP() string {
	return "HTTP 모드 활성화, HTTP/HTTPS 서비스 테스트"
}

func (k *KoreanLang) OptInsecure() string {
	return "SSL/TLS 인증서 검증 건너뛰기 (HTTP 모드만 유효)"
}

func (k *KoreanLang) OptLanguage() string {
	return "언어 설정 (en-US, ja-JP, ko-KR, zh-TW, zh-CN)"
}

func (k *KoreanLang) OptVersion() string {
	return "버전 정보 표시"
}

func (k *KoreanLang) OptHelp() string {
	return "이 도움말 정보 표시"
}

// TCP Examples
func (k *KoreanLang) ExampleBasic() string {
	return "tcping google.com                    # 기본 사용법 (기본 포트 80)"
}

func (k *KoreanLang) ExampleBasicPort() string {
	return "tcping google.com 80                 # 포트를 지정한 기본 사용법"
}

func (k *KoreanLang) ExamplePortFlag() string {
	return "tcping -p 443 google.com             # -p 플래그로 포트 지정"
}

func (k *KoreanLang) ExampleIPv4() string {
	return "tcping -4 -n 5 8.8.8.8 443           # IPv4, 5회 요청"
}

func (k *KoreanLang) ExampleColorVerbose() string {
	return "tcping -c -v example.com 443         # 컬러 출력과 상세 모드"
}

// HTTP Examples
func (k *KoreanLang) ExampleHTTPS() string {
	return "tcping -H https://www.google.com     # HTTPS 서비스 테스트"
}

func (k *KoreanLang) ExampleHTTP() string {
	return "tcping -H http://example.com         # HTTP 서비스 테스트"
}

func (k *KoreanLang) ExampleHTTPCount() string {
	return "tcping -H -n 10 https://github.com   # 10회 HTTP 요청 전송"
}

func (k *KoreanLang) ExampleHTTPVerbose() string {
	return "tcping -H -v https://api.github.com  # 상세 모드, 응답 세부사항 표시"
}

func (k *KoreanLang) ExampleHTTPInsecure() string {
	return "tcping -H -k https://self-signed.badssl.com  # SSL 인증서 검증 건너뛰기"
}

// Version info
func (k *KoreanLang) VersionFormat() string {
	return "%s 버전 %s\n"
}

// Error messages
func (k *KoreanLang) ErrorPrefix() string {
	return "오류: %v\n"
}

func (k *KoreanLang) ErrorInvalidPort() string {
	return "잘못된 포트 번호 형식"
}

func (k *KoreanLang) ErrorPortRange() string {
	return "포트 번호는 1부터 65535 사이여야 합니다"
}

func (k *KoreanLang) ErrorIPv6Decimal() string {
	return "IPv6 주소는 10진수 형식을 지원하지 않습니다"
}

func (k *KoreanLang) ErrorIPv6Hex() string {
	return "IPv6 주소는 16진수 형식을 지원하지 않습니다"
}

func (k *KoreanLang) ErrorResolve() string {
	return "%s 해석 실패: %v"
}

func (k *KoreanLang) ErrorNoIP() string {
	return "%s에 대한 IP 주소를 찾을 수 없습니다"
}

func (k *KoreanLang) ErrorNoIPv4() string {
	return "%s에 대한 IPv4 주소를 찾을 수 없습니다"
}

func (k *KoreanLang) ErrorNoIPv6() string {
	return "%s에 대한 IPv6 주소를 찾을 수 없습니다"
}

func (k *KoreanLang) ErrorBothIPv4IPv6() string {
	return "-4와 -6 플래그를 동시에 사용할 수 없습니다"
}

func (k *KoreanLang) ErrorNegativeInterval() string {
	return "간격 시간은 음수일 수 없습니다"
}

func (k *KoreanLang) ErrorNegativeTimeout() string {
	return "타임아웃은 음수일 수 없습니다"
}

func (k *KoreanLang) ErrorHostRequired() string {
	return "호스트 매개변수가 필요합니다\n\n사용법: tcping [옵션] <호스트> [포트]\n자세한 정보는 'tcping -h'를 입력하세요"
}

func (k *KoreanLang) ErrorPortMustBeInRange() string {
	return "포트 번호는 1부터 65535 사이여야 합니다"
}

func (k *KoreanLang) ErrorHTTPModeURIRequired() string {
	return "HTTP 모드에는 URI 매개변수가 필요합니다\n\n사용법: tcping -H [옵션] <URI>\n자세한 정보는 'tcping -h'를 입력하세요"
}

func (k *KoreanLang) ErrorInvalidURI() string {
	return "잘못된 URI 형식: %v"
}

func (k *KoreanLang) ErrorURIMustStartWithHTTP() string {
	return "URI는 http:// 또는 https://로 시작해야 합니다"
}

// Runtime messages
func (k *KoreanLang) MsgTCPPingStart() string {
	return "%s (%s - %s) 포트 %s에 TCP Ping 실행 중\n"
}

func (k *KoreanLang) MsgHTTPPingStart() string {
	return "%s에 HTTP Ping 실행 중 (User-Agent: tcping/%s.%s)\n"
}

func (k *KoreanLang) MsgInterrupted() string {
	return "\n작업이 중단되었습니다.\n"
}

func (k *KoreanLang) MsgOperationCanceled() string {
	return "\n작업이 중단되었습니다, 연결 시도를 중단합니다\n"
}

func (k *KoreanLang) MsgHTTPOperationCanceled() string {
	return "\n작업이 중단되었습니다, HTTP 요청을 중단합니다\n"
}

func (k *KoreanLang) MsgConnectionTimeout() string {
	return "연결 타임아웃"
}

// Connection messages
func (k *KoreanLang) MsgTCPConnectionFailed() string {
	return "TCP 연결 실패 %s:%s: seq=%d 오류=%s\n"
}

func (k *KoreanLang) MsgTCPConnectionSuccess() string {
	return "%s:%s에서 응답: seq=%d time=%.2fms\n"
}

func (k *KoreanLang) MsgHTTPRequestFailed() string {
	return "HTTP 요청 생성 실패 %s: seq=%d 오류=%v\n"
}

func (k *KoreanLang) MsgHTTPRequestFailedExec() string {
	return "HTTP 요청 실패 %s: seq=%d 오류=%v\n"
}

func (k *KoreanLang) MsgHTTPResponseFailed() string {
	return "HTTP 응답 읽기 실패 %s: seq=%d 오류=%v\n"
}

func (k *KoreanLang) MsgHTTPResponse() string {
	return "HTTP %d %s: seq=%d time=%.2fms size=%d bytes bandwidth=%.2f Mbps\n"
}

// Verbose messages
func (k *KoreanLang) MsgVerboseDetails() string {
	return "  세부사항: 연결 시도 소요시간 %.2fms, 대상 %s:%s\n"
}

func (k *KoreanLang) MsgVerboseConnection() string {
	return "  세부사항: 로컬 주소=%s, 원격 주소=%s:%s\n"
}

func (k *KoreanLang) MsgVerboseHTTP() string {
	return "  세부사항: 상태=%s, Content-Type=%s, Server=%s\n"
}

func (k *KoreanLang) MsgVerboseHTTPDetails() string {
	return "  세부사항:\n"
}

func (k *KoreanLang) MsgVerboseHTTPStatus() string {
	return "    상태: %s\n"
}

func (k *KoreanLang) MsgVerboseHTTPHeaders() string {
	return "    응답 헤더:\n"
}

// Statistics
func (k *KoreanLang) MsgTCPStatisticsTitle() string {
	return "\n\n--- TCP ping 통계 ---\n"
}

func (k *KoreanLang) MsgHTTPStatisticsTitle() string {
	return "\n\n--- HTTP ping 통계 ---\n"
}

func (k *KoreanLang) MsgStatisticsSummary() string {
	return "전송 = %d, 수신 = %d, 손실 = %d (%.1f%% 손실)\n"
}

func (k *KoreanLang) MsgStatisticsRTT() string {
	return "왕복시간(RTT): 최소 = %.2fms, 최대 = %.2fms, 평균 = %.2fms\n"
}

func (k *KoreanLang) MsgStatisticsTotalData() string {
	return "총 전송 데이터: %d bytes (%.2f MB)\n"
}

func (k *KoreanLang) MsgStatisticsBandwidth() string {
	return "추정 대역폭: 최소 = %.2f Mbps, 최대 = %.2f Mbps, 평균 = %.2f Mbps\n"
}

// IP type strings
func (k *KoreanLang) IPv4String() string {
	return "IPv4"
}

func (k *KoreanLang) IPv6String() string {
	return "IPv6"
}