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
	sync.RWMutex         // 使用读写锁以允许并发读取
	sentCount      int64 // 使用int64确保原子操作安全
	respondedCount int64
	minTime        float64
	maxTime        float64
	avgTime        float64
	totalTime      float64 // 添加总时间以简化平均计算
	totalBytes     int64  // HTTP模式下的总字节数
	minBandwidth   float64 // 最小带宽 (Mbps)
	maxBandwidth   float64 // 最大带宽 (Mbps)
	avgBandwidth   float64 // 平均带宽 (Mbps)
}

func (s *Statistics) update(elapsed float64, success bool) {
	// 原子操作增加发送计数，无需加锁
	atomic.AddInt64(&s.sentCount, 1)

	if !success {
		return
	}

	// 对于成功响应的更新需要加锁
	s.Lock()
	defer s.Unlock()

	newCount := atomic.AddInt64(&s.respondedCount, 1)
	s.totalTime += elapsed
	s.avgTime = s.totalTime / float64(newCount)

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

// 新增：更新HTTP统计信息
func (s *Statistics) updateHTTP(elapsed float64, bytes int64, success bool) {
	// 原子操作增加发送计数
	atomic.AddInt64(&s.sentCount, 1)

	if !success {
		return
	}

	// 对于成功响应的更新需要加锁
	s.Lock()
	defer s.Unlock()

	newCount := atomic.AddInt64(&s.respondedCount, 1)
	s.totalTime += elapsed
	s.avgTime = s.totalTime / float64(newCount)
	atomic.AddInt64(&s.totalBytes, bytes)

	// 计算带宽 (Mbps) = (bytes * 8) / (elapsed_ms * 1000)
	bandwidth := float64(bytes*8) / (elapsed * 1000)

	// 首次响应特殊处理
	if newCount == 1 {
		s.minTime = elapsed
		s.maxTime = elapsed
		s.minBandwidth = bandwidth
		s.maxBandwidth = bandwidth
		s.avgBandwidth = bandwidth
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
	// 计算平均带宽
	totalBandwidth := s.avgBandwidth * float64(newCount-1)
	s.avgBandwidth = (totalBandwidth + bandwidth) / float64(newCount)
}

// 添加新的方法，获取统计信息，使用读锁减少阻塞
func (s *Statistics) getStats() (sent, responded int64, min, max, avg float64) {
	s.RLock()
	defer s.RUnlock()

	return s.sentCount, s.respondedCount, s.minTime, s.maxTime, s.avgTime
}

// 获取HTTP统计信息
func (s *Statistics) getHTTPStats() (sent, responded, totalBytes int64, minTime, maxTime, avgTime, minBW, maxBW, avgBW float64) {
	s.RLock()
	defer s.RUnlock()

	return s.sentCount, s.respondedCount, s.totalBytes, s.minTime, s.maxTime, s.avgTime, s.minBandwidth, s.maxBandwidth, s.avgBandwidth
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
}

func handleError(err error, exitCode int) {
	if err != nil {
		_, err := fmt.Fprintf(os.Stderr, "错误: %v\n", err)
		if err != nil {
			return
		}
		os.Exit(exitCode)
	}
}

func printHelp() {
	fmt.Printf(`%s %s - TCP/HTTP 连接测试工具

描述:
    %s 测试到目标主机的TCP连接性或HTTP/HTTPS服务响应。

用法: 
    tcping [选项] <主机> [端口]                  # TCP模式 (默认端口: 80)
    tcping -H [选项] <URI>                       # HTTP模式

选项:
    -4, --ipv4              强制使用 IPv4
    -6, --ipv6              强制使用 IPv6
    -n, --count <次数>      发送请求的次数 (默认: 4)
    -p, --port <端口>       指定要连接的端口 (默认: 80)
    -t, --interval <毫秒>   请求间隔 (默认: 1000毫秒)
    -w, --timeout <毫秒>    连接超时 (默认: 1000毫秒)
    -c, --color             启用彩色输出
    -v, --verbose           启用详细模式，显示更多连接信息
    -H, --http              启用HTTP模式，测试HTTP/HTTPS服务
    -k, --insecure          跳过SSL/TLS证书验证（仅在HTTP模式下有效）
    -V, --version           显示版本信息
    -h, --help              显示此帮助信息

TCP模式示例:
    tcping google.com                    # 基本用法 (默认端口 80)
    tcping google.com 80                 # 基本用法指定端口
    tcping -p 443 google.com             # 使用-p参数指定端口
    tcping -4 -n 5 8.8.8.8 443           # IPv4, 5次请求
    tcping -c -v example.com 443         # 彩色输出和详细模式

HTTP模式示例:
    tcping -H https://www.google.com     # 测试HTTPS服务
    tcping -H http://example.com         # 测试HTTP服务  
    tcping -H -n 10 https://github.com   # 发送10次HTTP请求
    tcping -H -v https://api.github.com  # 详细模式，显示响应信息
    tcping -H -k https://self-signed.badssl.com  # 跳过SSL证书验证

`, programName, version, programName)
}

func printVersion() {
	fmt.Printf("%s 版本 %s\n", programName, version)
	if gitHash != "unknown" {
		fmt.Printf("Git commit: %s\n", gitHash)
	}
	if buildTime != "unknown" {
		fmt.Printf("构建时间: %s\n", buildTime)
	}
	fmt.Println(copyright)
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

// 修改函数签名，添加context参数
func pingOnce(ctx context.Context, address, port string, timeout int, stats *Statistics, seq int, ip string,
	opts *Options) {
	// 创建可取消的连接上下文，继承父上下文
	dialCtx, dialCancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Millisecond)

	// 确保在函数返回时取消上下文，防止资源泄漏
	defer dialCancel()

	// 创建完成通道，用于确保所有操作完成
	done := make(chan struct{})
	var conn net.Conn
	var err error
	var elapsed float64

	// 启动协程执行连接操作
	go func() {
		start := time.Now()
		var d net.Dialer
		conn, err = d.DialContext(dialCtx, "tcp", address+":"+port)
		elapsed = float64(time.Since(start).Microseconds()) / 1000.0
		close(done)
	}()

	// 等待连接完成或上下文取消
	select {
	case <-dialCtx.Done():
		// 如果是主上下文取消，显示中断信息
		if errors.Is(ctx.Err(), context.Canceled) {
			msg := "\n操作被中断, 连接尝试已中止\n"
			fmt.Print(infoText(msg, opts.ColorOutput))
			return
		}
		// 超时错误处理
		<-done // 等待连接协程完成
		err = fmt.Errorf("连接超时")
	case <-done:
		// 连接完成，继续处理
	}

	success := err == nil
	stats.update(elapsed, success)

	if !success {
		// 修改错误消息格式，去除重复的IP:端口信息
		errMsg := fmt.Sprintf("%v", err)

		// 清理错误信息
		targetAddr := address + ":" + port
		if strings.Contains(errMsg, targetAddr) {
			errParts := strings.Split(errMsg, targetAddr)
			if len(errParts) > 1 && strings.HasPrefix(errMsg, "dial ") {
				prefix := strings.Split(errMsg, targetAddr)[0]
				suffix := strings.Join(strings.Split(errMsg, targetAddr)[1:], "")
				errMsg = prefix + suffix
			}
		}

		msg := fmt.Sprintf("TCP连接失败 %s:%s: seq=%d 错误=%s\n", ip, port, seq, errMsg)
		fmt.Print(errorText(msg, opts.ColorOutput))

		if opts.VerboseMode {
			fmt.Printf("  详细信息: 连接尝试耗时 %.2fms, 目标 %s:%s\n", elapsed, address, port)
		}
		return
	}

	// 确保连接被关闭
	if conn != nil {
		defer conn.Close()
	}
	msg := fmt.Sprintf("从 %s:%s 收到响应: seq=%d time=%.2fms\n", ip, port, seq, elapsed)
	fmt.Print(successText(msg, opts.ColorOutput))

	if opts.VerboseMode && conn != nil {
		localAddr := conn.LocalAddr().String()
		fmt.Printf("  详细信息: 本地地址=%s, 远程地址=%s:%s\n", localAddr, ip, port)
	}
}

// HTTP ping功能
func httpPingOnce(ctx context.Context, uri string, timeout int, stats *Statistics, seq int, opts *Options) {
	// 创建自定义Transport支持超时和TLS
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: opts.InsecureSSL,
		},
		DisableKeepAlives: true,
	}

	// 创建HTTP客户端
	client := &http.Client{
		Transport: transport,
		Timeout:   time.Duration(timeout) * time.Millisecond,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse // 不自动跟随重定向
		},
	}

	// 设置User-Agent
	userAgent := fmt.Sprintf("tcping/%s.%s", version, gitHash)

	// 创建请求
	req, err := http.NewRequestWithContext(ctx, "GET", uri, nil)
	if err != nil {
		stats.updateHTTP(0, 0, false)
		msg := fmt.Sprintf("HTTP请求创建失败 %s: seq=%d 错误=%v\n", uri, seq, err)
		fmt.Print(errorText(msg, opts.ColorOutput))
		return
	}
	req.Header.Set("User-Agent", userAgent)

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
			msg := "\n操作被中断, HTTP请求已中止\n"
			fmt.Print(infoText(msg, opts.ColorOutput))
			return
		}
		stats.updateHTTP(elapsed, 0, false)
		msg := fmt.Sprintf("HTTP请求失败 %s: seq=%d 错误=%v\n", uri, seq, err)
		fmt.Print(errorText(msg, opts.ColorOutput))
		return
	}
	defer resp.Body.Close()

	// 读取响应体以计算传输字节数
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		stats.updateHTTP(elapsed, 0, false)
		msg := fmt.Sprintf("HTTP响应读取失败 %s: seq=%d 错误=%v\n", uri, seq, err)
		fmt.Print(errorText(msg, opts.ColorOutput))
		return
	}

	// 计算总字节数（包括响应头的估算）
	totalBytes := int64(len(body))
	headerSize := 0
	for key, values := range resp.Header {
		headerSize += len(key)
		for _, value := range values {
			headerSize += len(value)
		}
	}
	totalBytes += int64(headerSize)

	// 更新统计
	stats.updateHTTP(elapsed, totalBytes, true)

	// 计算带宽 (Mbps)
	bandwidth := float64(totalBytes*8) / (elapsed * 1000)

	// 输出结果
	msg := fmt.Sprintf("HTTP %d %s: seq=%d time=%.2fms size=%d bytes bandwidth=%.2f Mbps\n",
		resp.StatusCode, uri, seq, elapsed, totalBytes, bandwidth)
	
	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		fmt.Print(successText(msg, opts.ColorOutput))
	} else {
		fmt.Print(errorText(msg, opts.ColorOutput))
	}

	if opts.VerboseMode {
		fmt.Printf("  详细信息: 状态=%s, Content-Type=%s, Server=%s\n",
			resp.Status,
			resp.Header.Get("Content-Type"),
			resp.Header.Get("Server"))
	}
}

func printTCPingStatistics(stats *Statistics) {
	sent, responded, statMin, statMax, avg := stats.getStats()

	fmt.Printf("\n\n--- 目标主机 TCP ping 统计 ---\n")

	if sent > 0 {
		lossRate := float64(sent-responded) / float64(sent) * 100
		fmt.Printf("已发送 = %d, 已接收 = %d, 丢失 = %d (%.1f%% 丢失)\n",
			sent, responded, sent-responded, lossRate)

		if responded > 0 {
			fmt.Printf("往返时间(RTT): 最小 = %.2fms, 最大 = %.2fms, 平均 = %.2fms\n",
				statMin, statMax, avg)
		}
	}
}

// HTTP统计打印
func printHTTPStatistics(stats *Statistics) {
	sent, responded, totalBytes, minTime, maxTime, avgTime, minBW, maxBW, avgBW := stats.getHTTPStats()

	fmt.Printf("\n\n--- HTTP ping 统计 ---\n")

	if sent > 0 {
		lossRate := float64(sent-responded) / float64(sent) * 100
		fmt.Printf("已发送 = %d, 已接收 = %d, 丢失 = %d (%.1f%% 丢失)\n",
			sent, responded, sent-responded, lossRate)

		if responded > 0 {
			fmt.Printf("往返时间(RTT): 最小 = %.2fms, 最大 = %.2fms, 平均 = %.2fms\n",
				minTime, maxTime, avgTime)
			fmt.Printf("总传输数据: %d bytes (%.2f MB)\n", totalBytes, float64(totalBytes)/1024/1024)
			fmt.Printf("估算带宽: 最小 = %.2f Mbps, 最大 = %.2f Mbps, 平均 = %.2f Mbps\n",
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

		fmt.Printf("正在对 %s 执行 HTTP Ping (User-Agent: tcping/%s.%s)\n", uri, version, gitHash)

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
			fmt.Printf("\n操作被中断。\n")
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
		fmt.Printf("\n操作被中断。\n")
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
