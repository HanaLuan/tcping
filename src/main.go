package main

import (
	"context"
	"crypto/tls"
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
	
	"tcping/src/i18n"
)

// 这些变量将在编译时通过 -ldflags 注入
var (
	version     = "dev"      // 默认开发版本
	gitHash     = "unknown"  // 默认未知hash
	buildTime   = "unknown"  // 构建时间
)

const (
	copyright   = "Copyright (c) 2025. All rights reserved."
	programName = "TCPing"
)

type Statistics struct {
	// 使用原子操作的计数器，减少锁竞争
	sentCount      int64 // 原子操作
	respondedCount int64 // 原子操作
	totalBytes     int64 // 原子操作
	
	// 只有需要更复杂操作的字段使用锁
	sync.RWMutex
	minTime        float64
	maxTime        float64
	totalTime      float64 // 用于计算平均值
	minBandwidth   float64
	maxBandwidth   float64
	totalBandwidth float64 // 用于计算平均带宽
}

func (s *Statistics) update(elapsed float64, success bool) {
	// 原子操作增加发送计数，无需加锁
	atomic.AddInt64(&s.sentCount, 1)

	if !success {
		return
	}

	// 原子操作增加成功计数
	newCount := atomic.AddInt64(&s.respondedCount, 1)

	// 只在更新复杂统计时加锁
	s.Lock()
	defer s.Unlock()

	s.totalTime += elapsed

	// 首次响应特殊处理
	if newCount == 1 {
		s.minTime = elapsed
		s.maxTime = elapsed
		return
	}

	// 更新最小和最大时间
	if elapsed < s.minTime {
		s.minTime = elapsed
	}
	if elapsed > s.maxTime {
		s.maxTime = elapsed
	}
}

// 优化的HTTP统计更新函数
func (s *Statistics) updateHTTP(elapsed float64, bytes int64, success bool) {
	// 原子操作增加发送计数
	atomic.AddInt64(&s.sentCount, 1)

	if !success {
		return
	}

	// 原子操作增加成功计数和字节数
	newCount := atomic.AddInt64(&s.respondedCount, 1)
	atomic.AddInt64(&s.totalBytes, bytes)

	// 计算带宽 (Mbps)
	bandwidth := float64(bytes*8) / (elapsed * 1000)

	// 只在更新复杂统计时加锁
	s.Lock()
	defer s.Unlock()

	s.totalTime += elapsed
	s.totalBandwidth += bandwidth

	// 首次响应特殊处理
	if newCount == 1 {
		s.minTime = elapsed
		s.maxTime = elapsed
		s.minBandwidth = bandwidth
		s.maxBandwidth = bandwidth
		return
	}

	// 更新时间统计
	if elapsed < s.minTime {
		s.minTime = elapsed
	}
	if elapsed > s.maxTime {
		s.maxTime = elapsed
	}

	// 更新带宽统计
	if bandwidth < s.minBandwidth {
		s.minBandwidth = bandwidth
	}
	if bandwidth > s.maxBandwidth {
		s.maxBandwidth = bandwidth
	}
}

// 优化的统计获取函数，减少锁使用
func (s *Statistics) getStats() (sent, responded int64, min, max, avg float64) {
	// 原子操作读取计数器，无需加锁
	sent = atomic.LoadInt64(&s.sentCount)
	responded = atomic.LoadInt64(&s.respondedCount)
	
	// 只在读取复杂统计时加读锁
	s.RLock()
	min, max = s.minTime, s.maxTime
	if responded > 0 {
		avg = s.totalTime / float64(responded)
	}
	s.RUnlock()
	
	return
}

// 优化的HTTP统计获取函数
func (s *Statistics) getHTTPStats() (sent, responded, totalBytes int64, minTime, maxTime, avgTime, minBW, maxBW, avgBW float64) {
	// 原子操作读取计数器，无需加锁
	sent = atomic.LoadInt64(&s.sentCount)
	responded = atomic.LoadInt64(&s.respondedCount)
	totalBytes = atomic.LoadInt64(&s.totalBytes)
	
	// 只在读取复杂统计时加读锁
	s.RLock()
	minTime, maxTime = s.minTime, s.maxTime
	minBW, maxBW = s.minBandwidth, s.maxBandwidth
	if responded > 0 {
		avgTime = s.totalTime / float64(responded)
		avgBW = s.totalBandwidth / float64(responded)
	}
	s.RUnlock()
	
	return
}

type Options struct {
	UseIPv4     bool
	UseIPv6     bool
	Count       int
	Interval    int // 请求间隔（毫秒）
	Timeout     int
	ColorOutput bool
	VerboseMode bool
	ShowVersion bool
	ShowHelp    bool
	Port        int
	HTTPMode    bool // HTTP模式
	InsecureSSL bool // 跳过SSL/TLS证书验证
	Language    string // 语言设置
}

func handleError(err error, exitCode int) {
	if err != nil {
		_, err := fmt.Fprintf(os.Stderr, i18n.T().ErrorPrefix(), err)
		if err != nil {
			return
		}
		os.Exit(exitCode)
	}
}

func printHelp() {
	lang := i18n.T()
	fmt.Printf(`%s %s - %s

%s

%s
%s

%s:
    -4, --ipv4              %s
    -6, --ipv6              %s
    -n, --count <count>     %s
    -p, --port <port>       %s
    -t, --interval <ms>     %s
    -w, --timeout <ms>      %s
    -c, --color             %s
    -v, --verbose           %s
    -H, --http              %s
    -k, --insecure          %s
    -l, --language <code>   %s
    -V, --version           %s
    -h, --help              %s

%s:
    %s
    %s
    %s
    %s
    %s

%s:
    %s
    %s
    %s
    %s
    %s

`, programName, version, lang.ProgramDescription(),
		fmt.Sprintf(lang.UsageDescription(), programName),
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
		lang.OptInsecure(),
		lang.OptLanguage(),
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
		lang.ExampleHTTPVerbose(),
		lang.ExampleHTTPInsecure())
}

func printVersion() {
	lang := i18n.T()
	fmt.Printf(lang.VersionFormat(), programName, version)
	if gitHash != "unknown" {
		fmt.Printf("Git commit: %s\n", gitHash)
	}
	if buildTime != "unknown" {
		fmt.Printf("Build time: %s\n", buildTime)
	}
	fmt.Println(lang.Copyright())
}

func validatePort(port string) error {
	portNum, err := strconv.Atoi(port)
	if err != nil {
		return fmt.Errorf("端口号格式无效")
	}
	if portNum <= 0 || portNum > 65535 {
		return fmt.Errorf("端口号必须在 1 到 65535 之间")
	}
	return nil
}

func parseNumericIPv4(address string) net.IP {
	// Try decimal first
	if decIP, err := strconv.ParseUint(address, 10, 32); err == nil {
		return net.IPv4(
			byte(decIP>>24),
			byte(decIP>>16),
			byte(decIP>>8),
			byte(decIP),
		).To4()
	}

	// Try hexadecimal (with or without 0x prefix)
	addr := strings.ToLower(address)
	addr = strings.TrimPrefix(addr, "0x")
	if hexIP, err := strconv.ParseUint(addr, 16, 32); err == nil {
		return net.IPv4(
			byte(hexIP>>24),
			byte(hexIP>>16),
			byte(hexIP>>8),
			byte(hexIP),
		).To4()
	}

	return nil
}

func resolveAddress(address string, useIPv4, useIPv6 bool) (string, error) {
	// 检查IPv6数字格式
	if useIPv6 {
		if _, err := strconv.ParseUint(address, 10, 32); err == nil {
			return "", errors.New("IPv6 地址不支持十进制格式")
		}
		lowerAddr := strings.ToLower(address)
		if strings.HasPrefix(lowerAddr, "0x") {
			if _, err := strconv.ParseUint(strings.TrimPrefix(lowerAddr, "0x"), 16, 32); err == nil {
				return "", errors.New("IPv6 地址不支持十六进制格式")
			}
		}
	}

	// 尝试解析数字格式IPv4地址
	if useIPv4 || !useIPv6 {
		if ip := parseNumericIPv4(address); ip != nil {
			return ip.String(), nil
		}
	}

	// 尝试标准IP解析
	if ip := net.ParseIP(address); ip != nil {
		isV4 := ip.To4() != nil
		if useIPv4 && !isV4 {
			return "", fmt.Errorf("地址 %s 不是 IPv4 地址", address)
		}
		if useIPv6 && isV4 {
			return "", fmt.Errorf("地址 %s 不是 IPv6 地址", address)
		}
		if !isV4 {
			return "[" + ip.String() + "]", nil
		}
		return ip.String(), nil
	}

	// 最后尝试DNS解析
	ipList, err := net.LookupIP(address)
	if err != nil {
		return "", fmt.Errorf("解析 %s 失败: %v", address, err)
	}

	if len(ipList) == 0 {
		return "", fmt.Errorf("未找到 %s 的 IP 地址", address)
	}

	if useIPv4 {
		for _, ip := range ipList {
			if ip.To4() != nil {
				return ip.String(), nil
			}
		}
		return "", fmt.Errorf("未找到 %s 的 IPv4 地址", address)
	}

	if useIPv6 {
		for _, ip := range ipList {
			if ip.To4() == nil {
				return "[" + ip.String() + "]", nil
			}
		}
		return "", fmt.Errorf("未找到 %s 的 IPv6 地址", address)
	}

	ip := ipList[0]
	if ip.To4() == nil {
		return "[" + ip.String() + "]", nil
	}
	return ip.String(), nil
}

func isIPv4(address string) bool {
	if parseNumericIPv4(address) != nil {
		return true
	}
	return net.ParseIP(address) != nil && strings.Count(address, ":") == 0
}

func isIPv6(address string) bool {
	return strings.Count(address, ":") >= 2
}

// 优化的TCP连接函数，减少goroutine开销和内存分配
func pingOnce(ctx context.Context, address, port string, timeout int, stats *Statistics, seq int, ip string,
	opts *Options) {
	// 创建可取消的连接上下文，继承父上下文
	dialCtx, dialCancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Millisecond)
	defer dialCancel()

	// 直接在当前goroutine中执行连接，避免不必要的goroutine创建
	start := time.Now()
	
	// 创建带超时的dialer，复用连接配置
	dialer := &net.Dialer{
		Timeout: time.Duration(timeout) * time.Millisecond,
	}
	
	conn, err := dialer.DialContext(dialCtx, "tcp", address+":"+port)
	elapsed := float64(time.Since(start).Microseconds()) / 1000.0

	// 检查上下文取消
	if errors.Is(ctx.Err(), context.Canceled) {
		fmt.Print(infoText(i18n.T().MsgOperationCanceled(), opts.ColorOutput))
		return
	}

	success := err == nil
	stats.update(elapsed, success)

	if !success {
		// 优化错误消息处理，减少字符串操作
		errMsg := err.Error()
		
		// 使用strings.Builder减少内存分配
		var msgBuilder strings.Builder
		msgBuilder.Grow(64) // 预分配合理大小
		msgBuilder.WriteString("TCP连接失败 ")
		msgBuilder.WriteString(ip)
		msgBuilder.WriteByte(':')
		msgBuilder.WriteString(port)
		msgBuilder.WriteString(": seq=")
		msgBuilder.WriteString(strconv.Itoa(seq))
		msgBuilder.WriteString(" 错误=")
		msgBuilder.WriteString(errMsg)
		msgBuilder.WriteByte('\n')
		
		fmt.Print(errorText(msgBuilder.String(), opts.ColorOutput))

		if opts.VerboseMode {
			fmt.Printf(i18n.T().MsgVerboseDetails(), elapsed, address, port)
		}
		return
	}

	// 确保连接被关闭
	if conn != nil {
		defer conn.Close()
	}
	
	// 使用strings.Builder优化成功消息构建
	var msgBuilder strings.Builder
	msgBuilder.Grow(48) // 预分配合理大小
	msgBuilder.WriteString("从 ")
	msgBuilder.WriteString(ip)
	msgBuilder.WriteByte(':')
	msgBuilder.WriteString(port)
	msgBuilder.WriteString(" 收到响应: seq=")
	msgBuilder.WriteString(strconv.Itoa(seq))
	msgBuilder.WriteString(" time=")
	msgBuilder.WriteString(fmt.Sprintf("%.2f", elapsed))
	msgBuilder.WriteString("ms\n")
	
	fmt.Print(successText(msgBuilder.String(), opts.ColorOutput))

	if opts.VerboseMode && conn != nil {
		localAddr := conn.LocalAddr().String()
		fmt.Printf(i18n.T().MsgVerboseConnection(), localAddr, ip, port)
	}
}

// 全局HTTP客户端池，复用连接以提高性能
var httpClientPool = sync.Pool{
	New: func() interface{} {
		transport := &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 10,
			IdleConnTimeout:     30 * time.Second,
			TLSHandshakeTimeout: 10 * time.Second,
			DisableCompression:  false,
		}
		return &http.Client{
			Transport: transport,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}
	},
}

// HTTP ping功能 - 优化版本，使用连接池
func httpPingOnce(ctx context.Context, uri string, timeout int, stats *Statistics, seq int, opts *Options) {
	// 从池中获取HTTP客户端
	client := httpClientPool.Get().(*http.Client)
	defer httpClientPool.Put(client)
	
	// 动态设置超时和TLS配置
	transport := client.Transport.(*http.Transport)
	transport.TLSClientConfig = &tls.Config{
		InsecureSkipVerify: opts.InsecureSSL,
	}
	client.Timeout = time.Duration(timeout) * time.Millisecond

	// 创建请求
	req, err := http.NewRequestWithContext(ctx, "GET", uri, nil)
	if err != nil {
		stats.updateHTTP(0, 0, false)
		// 使用strings.Builder优化错误消息
		var msgBuilder strings.Builder
		msgBuilder.WriteString("HTTP请求创建失败 ")
		msgBuilder.WriteString(uri)
		msgBuilder.WriteString(": seq=")
		msgBuilder.WriteString(strconv.Itoa(seq))
		msgBuilder.WriteString(" 错误=")
		msgBuilder.WriteString(err.Error())
		msgBuilder.WriteByte('\n')
		fmt.Print(errorText(msgBuilder.String(), opts.ColorOutput))
		return
	}
	
	// 设置优化的User-Agent，避免重复字符串拼接
	req.Header.Set("User-Agent", "tcping/"+version+"."+gitHash)

	// 显示SSL验证警告（仅在详细模式下）
	if opts.VerboseMode && opts.InsecureSSL {
		fmt.Printf("  警告: SSL/TLS证书验证已禁用\n")
	}

	// 执行请求并计时
	start := time.Now()
	resp, err := client.Do(req)
	elapsed := float64(time.Since(start).Microseconds()) / 1000.0

	if err != nil {
		// 检查是否是上下文取消
		if errors.Is(ctx.Err(), context.Canceled) {
			fmt.Print(infoText(i18n.T().MsgHTTPOperationCanceled(), opts.ColorOutput))
			return
		}
		stats.updateHTTP(elapsed, 0, false)
		// 使用i18n格式化错误消息
		msg := fmt.Sprintf(i18n.T().MsgHTTPRequestFailedExec(), uri, seq, err)
		fmt.Print(errorText(msg, opts.ColorOutput))
		return
	}
	defer resp.Body.Close()

	// 使用缓冲读取优化内存使用
	var totalBytes int64
	buf := make([]byte, 4096) // 4KB缓冲区
	for {
		n, err := resp.Body.Read(buf)
		totalBytes += int64(n)
		if err == io.EOF {
			break
		}
		if err != nil {
			stats.updateHTTP(elapsed, 0, false)
			// 使用strings.Builder优化错误消息
			var msgBuilder strings.Builder
			msgBuilder.WriteString("HTTP响应读取失败 ")
			msgBuilder.WriteString(uri)
			msgBuilder.WriteString(": seq=")
			msgBuilder.WriteString(strconv.Itoa(seq))
			msgBuilder.WriteString(" 错误=")
			msgBuilder.WriteString(err.Error())
			msgBuilder.WriteByte('\n')
			fmt.Print(errorText(msgBuilder.String(), opts.ColorOutput))
			return
		}
	}

	// 计算响应头大小（优化版本）
	headerSize := int64(0)
	for key, values := range resp.Header {
		headerSize += int64(len(key) + 2) // key + ": "
		for _, value := range values {
			headerSize += int64(len(value) + 2) // value + "\r\n"
		}
	}
	totalBytes += headerSize

	// 更新统计
	stats.updateHTTP(elapsed, totalBytes, true)

	// 计算带宽 (Mbps) - 避免重复计算
	bandwidth := float64(totalBytes*8) / (elapsed * 1000)

	// 使用strings.Builder优化输出消息构建
	var msgBuilder strings.Builder
	msgBuilder.WriteString("HTTP ")
	msgBuilder.WriteString(strconv.Itoa(resp.StatusCode))
	msgBuilder.WriteByte(' ')
	msgBuilder.WriteString(uri)
	msgBuilder.WriteString(": seq=")
	msgBuilder.WriteString(strconv.Itoa(seq))
	msgBuilder.WriteString(" time=")
	msgBuilder.WriteString(fmt.Sprintf("%.2f", elapsed))
	msgBuilder.WriteString("ms size=")
	msgBuilder.WriteString(strconv.FormatInt(totalBytes, 10))
	msgBuilder.WriteString(" bytes bandwidth=")
	msgBuilder.WriteString(fmt.Sprintf("%.2f", bandwidth))
	msgBuilder.WriteString(" Mbps\n")
	
	msg := msgBuilder.String()
	
	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		fmt.Print(successText(msg, opts.ColorOutput))
	} else {
		fmt.Print(errorText(msg, opts.ColorOutput))
	}

	if opts.VerboseMode {
		// Display response details with proper formatting
		fmt.Print(i18n.T().MsgVerboseHTTPDetails())
		fmt.Printf(i18n.T().MsgVerboseHTTPStatus(), resp.Status)
		
		// Display key headers
		if contentType := resp.Header.Get("Content-Type"); contentType != "" {
			fmt.Printf("    Content-Type: %s\n", contentType)
		}
		if server := resp.Header.Get("Server"); server != "" {
			fmt.Printf("    Server: %s\n", server)
		}
		if contentLength := resp.Header.Get("Content-Length"); contentLength != "" {
			fmt.Printf("    Content-Length: %s\n", contentLength)
		}
		if lastModified := resp.Header.Get("Last-Modified"); lastModified != "" {
			fmt.Printf("    Last-Modified: %s\n", lastModified)
		}
		if cacheControl := resp.Header.Get("Cache-Control"); cacheControl != "" {
			fmt.Printf("    Cache-Control: %s\n", cacheControl)
		}
		
		// Display all response headers in organized format
		fmt.Print(i18n.T().MsgVerboseHTTPHeaders())
		headerNames := make([]string, 0, len(resp.Header))
		for name := range resp.Header {
			headerNames = append(headerNames, name)
		}
		// Sort headers for consistent display
		for i := 0; i < len(headerNames); i++ {
			for j := i + 1; j < len(headerNames); j++ {
				if headerNames[i] > headerNames[j] {
					headerNames[i], headerNames[j] = headerNames[j], headerNames[i]
				}
			}
		}
		
		// Display each header with proper line breaks for long values
		for _, name := range headerNames {
			values := resp.Header[name]
			for _, value := range values {
				// Break long header values into multiple lines if needed
				if len(value) > 60 {
					fmt.Printf("    %s:\n      %s\n", name, value)
				} else {
					fmt.Printf("    %s: %s\n", name, value)
				}
			}
		}
	}
}

func printTCPingStatistics(stats *Statistics) {
	sent, responded, statMin, statMax, avg := stats.getStats()

	fmt.Print(i18n.T().MsgTCPStatisticsTitle())

	if sent > 0 {
		lossRate := float64(sent-responded) / float64(sent) * 100
		fmt.Printf(i18n.T().MsgStatisticsSummary(),
			sent, responded, sent-responded, lossRate)

		if responded > 0 {
			fmt.Printf(i18n.T().MsgStatisticsRTT(),
				statMin, statMax, avg)
		}
	}
}

// HTTP统计打印
func printHTTPStatistics(stats *Statistics) {
	sent, responded, totalBytes, minTime, maxTime, avgTime, minBW, maxBW, avgBW := stats.getHTTPStats()

	fmt.Print(i18n.T().MsgHTTPStatisticsTitle())

	if sent > 0 {
		lossRate := float64(sent-responded) / float64(sent) * 100
		fmt.Printf(i18n.T().MsgStatisticsSummary(),
			sent, responded, sent-responded, lossRate)

		if responded > 0 {
			fmt.Printf(i18n.T().MsgStatisticsRTT(),
				minTime, maxTime, avgTime)
			fmt.Printf(i18n.T().MsgStatisticsTotalData(), totalBytes, float64(totalBytes)/1024/1024)
			fmt.Printf(i18n.T().MsgStatisticsBandwidth(),
				minBW, maxBW, avgBW)
		}
	}
}

func colorize(text string, colorCode string, useColor bool) string {
	if !useColor {
		return text
	}
	return "\033[" + colorCode + "m" + text + "\033[0m"
}

func successText(text string, useColor bool) string {
	return colorize(text, "32", useColor) // 绿色
}

func errorText(text string, useColor bool) string {
	return colorize(text, "31", useColor) // 红色
}

func infoText(text string, useColor bool) string {
	return colorize(text, "36", useColor) // 青色
}

// 处理短选项和长选项映射的函数
func setupFlags(opts *Options) {
	// 定义命令行标志，同时设置短选项和长选项
	ipv4 := flag.Bool("4", false, "使用 IPv4 地址")
	ipv6 := flag.Bool("6", false, "使用 IPv6 地址")
	count := flag.Int("n", -1, "发送请求次数 (默认: 4)") // 默认-1，后续判断
	interval := flag.Int("t", 1000, "请求间隔（毫秒）")
	timeout := flag.Int("w", 1000, "连接超时（毫秒）")
	port := flag.Int("p", 0, "指定要连接的端口 (默认: 80)")
	color := flag.Bool("c", false, "启用彩色输出")
	verbose := flag.Bool("v", false, "启用详细模式")
	httpMode := flag.Bool("H", false, "启用HTTP模式")
	insecure := flag.Bool("k", false, "跳过SSL/TLS证书验证")
	language := flag.String("l", "", "设置语言 (en-US, ja-JP, ko-KR, zh-TW, zh-CN)")
	version := flag.Bool("V", false, "显示版本信息")
	help := flag.Bool("h", false, "显示帮助信息")

	// 设置长选项别名
	flag.BoolVar(ipv4, "ipv4", false, "使用 IPv4 地址")
	flag.BoolVar(ipv6, "ipv6", false, "使用 IPv6 地址")
	flag.IntVar(count, "count", -1, "发送请求次数 (默认: 4)")
	flag.IntVar(interval, "interval", 1000, "请求间隔（毫秒）")
	flag.IntVar(timeout, "timeout", 1000, "连接超时（毫秒）")
	flag.IntVar(port, "port", 0, "指定要连接的端口 (默认: 80)")
	flag.BoolVar(color, "color", false, "启用彩色输出")
	flag.BoolVar(verbose, "verbose", false, "启用详细模式")
	flag.BoolVar(httpMode, "http", false, "启用HTTP模式")
	flag.BoolVar(insecure, "insecure", false, "跳过SSL/TLS证书验证")
	flag.StringVar(language, "language", "", "设置语言 (en-US, ja-JP, ko-KR, zh-TW, zh-CN)")
	flag.BoolVar(version, "version", false, "显示版本信息")
	flag.BoolVar(help, "help", false, "显示帮助信息")

	// 解析命令行参数
	flag.Parse()

	// 设置选项结构
	opts.UseIPv4 = *ipv4
	opts.UseIPv6 = *ipv6
	// 关键变更：如果未指定 -n/--count，则默认4次
	if *count == -1 {
		opts.Count = 4
	} else {
		opts.Count = *count
	}
	opts.Interval = *interval
	opts.Timeout = *timeout
	opts.Port = *port
	opts.ColorOutput = *color
	opts.VerboseMode = *verbose
	opts.HTTPMode = *httpMode
	opts.InsecureSSL = *insecure
	opts.Language = *language
	opts.ShowVersion = *version
	opts.ShowHelp = *help
}

// 新增集中的参数验证函数
func validateOptions(opts *Options, args []string) (string, string, error) {
	// 验证基本选项
	if opts.UseIPv4 && opts.UseIPv6 {
		return "", "", errors.New("无法同时使用 -4 和 -6 标志")
	}

	if opts.Interval < 0 {
		return "", "", errors.New("间隔时间不能为负值")
	}

	if opts.Timeout < 0 {
		return "", "", errors.New("超时时间不能为负值")
	}

	// 验证主机参数
	if len(args) < 1 {
		return "", "", errors.New("需要提供主机参数\n\n用法: tcping [选项] <主机> [端口]\n尝试 'tcping -h' 获取更多信息")
	}

	host := args[0]
	port := "80" // 默认端口为 80

	// 优先级：命令行直接指定的端口 > -p参数指定的端口 > 默认端口80
	if len(args) > 1 {
		port = args[1]
	} else if opts.Port > 0 {
		// 如果通过-p参数指定了端口且命令行没有直接指定端口，则使用-p参数的值
		if opts.Port > 65535 {
			return "", "", errors.New("端口号必须在 1 到 65535 之间")
		}
		port = strconv.Itoa(opts.Port)
	}

	// 验证端口
	if err := validatePort(port); err != nil {
		return "", "", err
	}

	return host, port, nil
}

func main() {
	// 创建选项结构
	opts := &Options{}

	// 设置和解析命令行参数
	setupFlags(opts)
	
	// 初始化国际化系统
	i18n.Initialize(opts.Language)

	// 处理帮助和版本信息选项，这些选项优先级最高
	if opts.ShowHelp {
		printHelp()
		os.Exit(0)
	}

	if opts.ShowVersion {
		printVersion()
		os.Exit(0)
	}

	stats := &Statistics{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 创建信号捕获通道
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	// 使用 WaitGroup 来确保后台 goroutine 正确退出
	var wg sync.WaitGroup
	wg.Add(1)

	// 创建错误通道
	errChan := make(chan error, 1)

	// HTTP模式处理
	if opts.HTTPMode {
		// HTTP模式下验证URI参数
		if len(flag.Args()) < 1 {
			handleError(errors.New("HTTP模式需要提供URI参数\n\n用法: tcping -H [选项] <URI>\n尝试 'tcping -h' 获取更多信息"), 1)
		}

		uri := flag.Args()[0]
		
		// 验证URI格式
		parsedURL, err := url.Parse(uri)
		if err != nil {
			handleError(fmt.Errorf("无效的URI格式: %v", err), 1)
		}
		if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
			handleError(errors.New("URI必须以http://或https://开头"), 1)
		}

		fmt.Printf(i18n.T().MsgHTTPPingStart(), uri, version, gitHash)

		// 启动HTTP ping协程
		go func() {
			defer wg.Done()
			defer signal.Stop(interrupt)

		httpPingLoop:
			for i := 0; opts.Count == 0 || i < opts.Count; i++ {
				select {
				case <-ctx.Done():
					return
				default:
				}

				// 执行HTTP ping
				httpPingOnce(ctx, uri, opts.Timeout, stats, i, opts)

				if opts.Count != 0 && i == opts.Count-1 {
					break httpPingLoop
				}

				select {
				case <-ctx.Done():
					return
				case <-time.After(time.Duration(opts.Interval) * time.Millisecond):
				}
			}
			select {
			case errChan <- nil:
			default:
			}
		}()

		// 等待中断信号或完成
		select {
		case <-interrupt:
			fmt.Print(i18n.T().MsgInterrupted())
			cancel()
		case err := <-errChan:
			if err != nil {
				handleError(err, 1)
			}
		}

		wg.Wait()
		printHTTPStatistics(stats)
		return
	}

	// TCP模式处理（原有逻辑）
	// 集中验证所有参数
	host, port, err := validateOptions(opts, flag.Args())
	if err != nil {
		handleError(err, 1)
	}

	// 确定使用IPv4还是IPv6
	useIPv4 := opts.UseIPv4 || (!opts.UseIPv6 && isIPv4(host))
	useIPv6 := opts.UseIPv6 || isIPv6(host)

	// 保存原始主机名用于显示
	originalHost := host

	// 解析IP地址
	address, err := resolveAddress(host, useIPv4, useIPv6)
	if err != nil {
		handleError(err, 1)
	}

	// 提取IP地址用于显示
	ipType := "IPv4"
	ipAddress := address
	if strings.HasPrefix(address, "[") && strings.HasSuffix(address, "]") {
		ipType = "IPv6"
		ipAddress = address[1 : len(address)-1]
	}

	fmt.Printf("正在对 %s (%s - %s) 端口 %s 执行 TCP Ping\n", originalHost, ipType, ipAddress, port)

	// 启动ping协程
	go func() {
		defer wg.Done()
		defer signal.Stop(interrupt) // 停止信号捕获

	pingLoop:
		for i := 0; opts.Count == 0 || i < opts.Count; i++ {
			// 检查上下文是否已取消
			select {
			case <-ctx.Done():
				return
			default:
				// 继续执行
			}

			// 执行ping
			pingOnce(ctx, address, port, opts.Timeout, stats, i, ipAddress, opts)

			// 检查是否完成所有请求
			if opts.Count != 0 && i == opts.Count-1 {
				break pingLoop
			}

			// 等待下一次ping的间隔
			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Duration(opts.Interval) * time.Millisecond):
				// 继续下一次ping
			}
		}
		// 所有ping完成，发送nil到错误通道表示正常完成
		select {
		case errChan <- nil:
		default:
		}
	}()

	// 等待中断信号或完成
	select {
	case <-interrupt:
		fmt.Print(i18n.T().MsgInterrupted())
		cancel() // 取消上下文
	case err := <-errChan:
		if err != nil {
			handleError(err, 1)
		}
		// 正常完成
	}

	// 等待ping协程完成
	wg.Wait()
	printTCPingStatistics(stats)
}
