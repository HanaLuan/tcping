package main

/*
This file demonstrates how to integrate the i18n system into the existing tcping main.go file.

INTEGRATION INSTRUCTIONS:
1. Add this import to main.go:
   import "./i18n"

2. Add language command-line option to setupFlags():
   lang := flag.String("l", "", "Set language (en-US, ja-JP, ko-KR, zh-TW, zh-CN)")
   flag.StringVar(lang, "lang", "", "Set language (en-US, ja-JP, ko-KR, zh-TW, zh-CN)")

3. Add language initialization to main() after setupFlags():
   i18n.Initialize(*lang)

4. Replace hardcoded strings with i18n calls:

EXAMPLES OF REPLACEMENTS:

Original:
   fmt.Fprintf(os.Stderr, "错误: %v\n", err)
Replacement:
   fmt.Fprintf(os.Stderr, i18n.T().ErrorPrefix(), err)

Original:
   fmt.Printf(`%s %s - TCP/HTTP 连接测试工具`, programName, version)
Replacement:
   fmt.Printf(`%s %s - %s`, programName, version, i18n.T().ProgramDescription())

Original:
   fmt.Printf("正在对 %s (%s - %s) 端口 %s 执行 TCP Ping\n", originalHost, ipType, ipAddress, port)
Replacement:
   fmt.Printf(i18n.T().MsgTCPPingStart(), originalHost, ipType, ipAddress, port)

Original:
   fmt.Printf("从 %s:%s 收到响应: seq=%d time=%.2fms\n", ip, port, seq, elapsed)
Replacement:
   fmt.Printf(i18n.T().MsgTCPConnectionSuccess(), ip, port, seq, elapsed)

ENVIRONMENT VARIABLE SUPPORT:

Users can set language via environment variables:
- TCPING_LANG=en-US ./tcping google.com
- LANG=ja-JP ./tcping google.com
- LC_ALL=ko-KR ./tcping google.com

Or via command line:
- ./tcping -l en-US google.com
- ./tcping --lang ja-JP google.com

LANGUAGE DETECTION PRIORITY:
1. Command line flag (-l/--lang)
2. TCPING_LANG environment variable
3. LC_ALL environment variable  
4. LC_MESSAGES environment variable
5. LANG environment variable
6. Default to en-US

SUPPORTED LANGUAGES:
- en-US: English (United States) - Default
- ja-JP: Japanese (Japan)
- ko-KR: Korean (South Korea)
- zh-TW: Traditional Chinese (Taiwan)
- zh-CN: Simplified Chinese (China) - Original language

EXAMPLE USAGE AFTER INTEGRATION:

# English (default)
./tcping google.com

# Japanese
./tcping -l ja-JP google.com
LANG=ja-JP ./tcping google.com

# Korean
./tcping --lang ko-KR google.com
TCPING_LANG=ko-KR ./tcping google.com

# Traditional Chinese
./tcping -l zh-TW google.com

# Simplified Chinese (original)
./tcping -l zh-CN google.com

The system will automatically handle:
- Language detection from environment
- Fallback to English if language not found
- Proper formatting with printf-style placeholders
- All user-facing strings in the selected language
*/

import (
	"fmt"
	"./i18n"
)

// Example of how to integrate i18n into existing functions

func examplePrintHelp() {
	lang := i18n.T()
	fmt.Printf(`%s %s - %s

%s:
    %s

%s: 
    %s
    %s

%s
    -4, --ipv4              %s
    -6, --ipv6              %s
    -n, --count <次数>      %s
    -p, --port <端口>       %s
    -t, --interval <毫秒>   %s
    -w, --timeout <毫秒>    %s
    -c, --color             %s
    -v, --verbose           %s
    -H, --http              %s
    -l, --lang <code>       Set language (en-US, ja-JP, ko-KR, zh-TW, zh-CN)
    -V, --version           %s
    -h, --help              %s

%s
    %s
    %s
    %s
    %s
    %s

%s
    %s
    %s
    %s
    %s

`, "TCPing", "v1.8.0", lang.ProgramDescription(),
		lang.UsageDescription(), "TCPing",
		lang.UsageTCP(),
		lang.UsageHTTP(),
		lang.OptionsTitle(),
		lang.OptForceIPv4(),
		lang.OptForceIPv6(),
		lang.OptCount(),
		lang.OptPort(),
		lang.OptInterval(),
		lang.OptTimeout(),
		lang.OptColor(),
		lang.OptVerbose(),
		lang.OptHTTP(),
		lang.OptVersion(),
		lang.OptHelp(),
		lang.TCPExamplesTitle(),
		lang.ExampleBasic(),
		lang.ExampleBasicPort(),
		lang.ExamplePortFlag(),
		lang.ExampleIPv4(),
		lang.ExampleColorVerbose(),
		lang.HTTPExamplesTitle(),
		lang.ExampleHTTPS(),
		lang.ExampleHTTP(),
		lang.ExampleHTTPCount(),
		lang.ExampleHTTPVerbose())
}

func examplePrintVersion() {
	lang := i18n.T()
	fmt.Printf(lang.VersionFormat(), "TCPing", "v1.8.0")
	fmt.Printf("Git commit: %s\n", "11ae0ba")
	fmt.Printf("Build time: %s\n", "2025-01-01 00:00:00 UTC")
	fmt.Println(lang.Copyright())
}

func exampleHandleError(err error, exitCode int) {
	if err != nil {
		lang := i18n.T()
		fmt.Fprintf(os.Stderr, lang.ErrorPrefix(), err)
		os.Exit(exitCode)
	}
}