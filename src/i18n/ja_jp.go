package i18n

// JapaneseLang implements Language interface for Japanese
type JapaneseLang struct{}

// Program info
func (j *JapaneseLang) ProgramDescription() string {
	return "TCP/HTTP接続テストツール"
}

func (j *JapaneseLang) Copyright() string {
	return "Copyright (c) 2025. All rights reserved."
}

// Help and usage
func (j *JapaneseLang) UsageDescription() string {
	return "%sはターゲットホストへのTCP接続性やHTTP/HTTPSサービスレスポンスをテストします。"
}

func (j *JapaneseLang) UsageTCP() string {
	return "tcping [オプション] <ホスト> [ポート]           # TCPモード (デフォルトポート: 80)"
}

func (j *JapaneseLang) UsageHTTP() string {
	return "tcping -H [オプション] <URI>                   # HTTPモード"
}

func (j *JapaneseLang) OptionsTitle() string {
	return "オプション:"
}

func (j *JapaneseLang) TCPExamplesTitle() string {
	return "TCPモードの例:"
}

func (j *JapaneseLang) HTTPExamplesTitle() string {
	return "HTTPモードの例:"
}

// Command line options
func (j *JapaneseLang) OptForceIPv4() string {
	return "IPv4を強制使用"
}

func (j *JapaneseLang) OptForceIPv6() string {
	return "IPv6を強制使用"
}

func (j *JapaneseLang) OptCount() string {
	return "送信するリクエスト数 (デフォルト: 4)"
}

func (j *JapaneseLang) OptPort() string {
	return "接続するポートを指定 (デフォルト: 80)"
}

func (j *JapaneseLang) OptInterval() string {
	return "リクエスト間隔（ミリ秒）(デフォルト: 1000ms)"
}

func (j *JapaneseLang) OptTimeout() string {
	return "接続タイムアウト（ミリ秒）(デフォルト: 1000ms)"
}

func (j *JapaneseLang) OptColor() string {
	return "カラー出力を有効化"
}

func (j *JapaneseLang) OptVerbose() string {
	return "詳細モードを有効化、より多くの接続情報を表示"
}

func (j *JapaneseLang) OptHTTP() string {
	return "HTTPモードを有効化、HTTP/HTTPSサービスをテスト"
}

func (j *JapaneseLang) OptInsecure() string {
	return "SSL/TLS証明書の検証をスキップ（HTTPモードのみ有効）"
}

func (j *JapaneseLang) OptLanguage() string {
	return "言語を設定 (en-US, ja-JP, ko-KR, zh-TW, zh-CN)"
}

func (j *JapaneseLang) OptVersion() string {
	return "バージョン情報を表示"
}

func (j *JapaneseLang) OptHelp() string {
	return "このヘルプ情報を表示"
}

// TCP Examples
func (j *JapaneseLang) ExampleBasic() string {
	return "tcping google.com                    # 基本的な使用法 (デフォルトポート 80)"
}

func (j *JapaneseLang) ExampleBasicPort() string {
	return "tcping google.com 80                 # ポートを指定した基本的な使用法"
}

func (j *JapaneseLang) ExamplePortFlag() string {
	return "tcping -p 443 google.com             # -pフラグでポートを指定"
}

func (j *JapaneseLang) ExampleIPv4() string {
	return "tcping -4 -n 5 8.8.8.8 443           # IPv4, 5回リクエスト"
}

func (j *JapaneseLang) ExampleColorVerbose() string {
	return "tcping -c -v example.com 443         # カラー出力と詳細モード"
}

// HTTP Examples
func (j *JapaneseLang) ExampleHTTPS() string {
	return "tcping -H https://www.google.com     # HTTPSサービスをテスト"
}

func (j *JapaneseLang) ExampleHTTP() string {
	return "tcping -H http://example.com         # HTTPサービスをテスト"
}

func (j *JapaneseLang) ExampleHTTPCount() string {
	return "tcping -H -n 10 https://github.com   # 10回HTTPリクエストを送信"
}

func (j *JapaneseLang) ExampleHTTPVerbose() string {
	return "tcping -H -v https://api.github.com  # 詳細モード、レスポンス詳細を表示"
}

func (j *JapaneseLang) ExampleHTTPInsecure() string {
	return "tcping -H -k https://self-signed.badssl.com  # SSL証明書検証をスキップ"
}

// Version info
func (j *JapaneseLang) VersionFormat() string {
	return "%s バージョン %s\n"
}

// Error messages
func (j *JapaneseLang) ErrorPrefix() string {
	return "エラー: %v\n"
}

func (j *JapaneseLang) ErrorInvalidPort() string {
	return "ポート番号の形式が無効です"
}

func (j *JapaneseLang) ErrorPortRange() string {
	return "ポート番号は1から65535の間である必要があります"
}

func (j *JapaneseLang) ErrorIPv6Decimal() string {
	return "IPv6アドレスは10進数形式をサポートしていません"
}

func (j *JapaneseLang) ErrorIPv6Hex() string {
	return "IPv6アドレスは16進数形式をサポートしていません"
}

func (j *JapaneseLang) ErrorResolve() string {
	return "%sの解決に失敗しました: %v"
}

func (j *JapaneseLang) ErrorNoIP() string {
	return "%sのIPアドレスが見つかりませんでした"
}

func (j *JapaneseLang) ErrorNoIPv4() string {
	return "%sのIPv4アドレスが見つかりませんでした"
}

func (j *JapaneseLang) ErrorNoIPv6() string {
	return "%sのIPv6アドレスが見つかりませんでした"
}

func (j *JapaneseLang) ErrorBothIPv4IPv6() string {
	return "-4と-6フラグの両方を同時に使用することはできません"
}

func (j *JapaneseLang) ErrorNegativeInterval() string {
	return "間隔時間は負の値にできません"
}

func (j *JapaneseLang) ErrorNegativeTimeout() string {
	return "タイムアウト時間は負の値にできません"
}

func (j *JapaneseLang) ErrorHostRequired() string {
	return "ホストパラメータが必要です\n\n使用法: tcping [オプション] <ホスト> [ポート]\n詳細は 'tcping -h' を実行してください"
}

func (j *JapaneseLang) ErrorPortMustBeInRange() string {
	return "ポート番号は1から65535の間である必要があります"
}

func (j *JapaneseLang) ErrorHTTPModeURIRequired() string {
	return "HTTPモードにはURIパラメータが必要です\n\n使用法: tcping -H [オプション] <URI>\n詳細は 'tcping -h' を実行してください"
}

func (j *JapaneseLang) ErrorInvalidURI() string {
	return "無効なURI形式: %v"
}

func (j *JapaneseLang) ErrorURIMustStartWithHTTP() string {
	return "URIはhttp://またはhttps://で始まる必要があります"
}

// Runtime messages
func (j *JapaneseLang) MsgTCPPingStart() string {
	return "%s (%s - %s) ポート%sにTCP Ping実行中\n"
}

func (j *JapaneseLang) MsgHTTPPingStart() string {
	return "%sにHTTP Ping実行中 (User-Agent: tcping/%s.%s)\n"
}

func (j *JapaneseLang) MsgInterrupted() string {
	return "\n操作が中断されました。\n"
}

func (j *JapaneseLang) MsgOperationCanceled() string {
	return "\n操作が中断されました、接続試行を中止します\n"
}

func (j *JapaneseLang) MsgHTTPOperationCanceled() string {
	return "\n操作が中断されました、HTTPリクエストを中止します\n"
}

func (j *JapaneseLang) MsgConnectionTimeout() string {
	return "接続タイムアウト"
}

// Connection messages
func (j *JapaneseLang) MsgTCPConnectionFailed() string {
	return "TCP接続失敗 %s:%s: seq=%d エラー=%s\n"
}

func (j *JapaneseLang) MsgTCPConnectionSuccess() string {
	return "%s:%sからレスポンス: seq=%d time=%.2fms\n"
}

func (j *JapaneseLang) MsgHTTPRequestFailed() string {
	return "HTTPリクエスト作成失敗 %s: seq=%d エラー=%v\n"
}

func (j *JapaneseLang) MsgHTTPRequestFailedExec() string {
	return "HTTPリクエスト失敗 %s: seq=%d エラー=%v\n"
}

func (j *JapaneseLang) MsgHTTPResponseFailed() string {
	return "HTTPレスポンス読み取り失敗 %s: seq=%d エラー=%v\n"
}

func (j *JapaneseLang) MsgHTTPResponse() string {
	return "HTTP %d %s: seq=%d time=%.2fms size=%d bytes bandwidth=%.2f Mbps\n"
}

// Verbose messages
func (j *JapaneseLang) MsgVerboseDetails() string {
	return "  詳細: 接続試行時間 %.2fms、ターゲット %s:%s\n"
}

func (j *JapaneseLang) MsgVerboseConnection() string {
	return "  詳細: ローカルアドレス=%s、リモートアドレス=%s:%s\n"
}

func (j *JapaneseLang) MsgVerboseHTTP() string {
	return "  詳細: ステータス=%s、Content-Type=%s、Server=%s\n"
}

func (j *JapaneseLang) MsgVerboseHTTPDetails() string {
	return "  詳細:\n"
}

func (j *JapaneseLang) MsgVerboseHTTPStatus() string {
	return "    ステータス: %s\n"
}

func (j *JapaneseLang) MsgVerboseHTTPHeaders() string {
	return "    レスポンスヘッダー:\n"
}

// Statistics
func (j *JapaneseLang) MsgTCPStatisticsTitle() string {
	return "\n\n--- TCP ping統計 ---\n"
}

func (j *JapaneseLang) MsgHTTPStatisticsTitle() string {
	return "\n\n--- HTTP ping統計 ---\n"
}

func (j *JapaneseLang) MsgStatisticsSummary() string {
	return "送信 = %d、受信 = %d、損失 = %d (%.1f%% 損失)\n"
}

func (j *JapaneseLang) MsgStatisticsRTT() string {
	return "往復時間(RTT): 最小 = %.2fms、最大 = %.2fms、平均 = %.2fms\n"
}

func (j *JapaneseLang) MsgStatisticsTotalData() string {
	return "総転送データ: %d bytes (%.2f MB)\n"
}

func (j *JapaneseLang) MsgStatisticsBandwidth() string {
	return "推定帯域幅: 最小 = %.2f Mbps、最大 = %.2f Mbps、平均 = %.2f Mbps\n"
}

// IP type strings
func (j *JapaneseLang) IPv4String() string {
	return "IPv4"
}

func (j *JapaneseLang) IPv6String() string {
	return "IPv6"
}