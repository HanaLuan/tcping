# SSL/TLS Certificate Verification Bypass Feature

## Overview

TCPing now supports bypassing SSL/TLS certificate verification for HTTP/HTTPS requests, which is useful for testing self-signed certificates, internal services, or development environments.

## New Command Line Option

### Short and Long Forms
```bash
-k, --insecure          跳过SSL/TLS证书验证（仅在HTTP模式下有效）
```

**Note:** This option only works in HTTP mode (`-H` flag) and has no effect in TCP mode.

## Usage Examples

### Basic SSL Bypass
```bash
# Test self-signed certificate site
tcping -H -k https://self-signed.badssl.com

# Test with custom timeout and verbose output
tcping -H -k -v -w 5000 https://internal-server.company.com
```

### Development and Testing
```bash
# Test local development server with self-signed cert
tcping -H -k -n 5 https://localhost:8443

# Test internal API with invalid certificate
tcping -H -k -v https://internal-api.local/health
```

### Production Monitoring (with warning)
```bash
# Monitor service while ignoring cert issues
tcping -H -k -n 10 -t 5000 https://legacy-system.company.com/status
```

## Security Warning

⚠️ **SECURITY WARNING**: The `-k/--insecure` option disables SSL/TLS certificate verification, making the connection vulnerable to man-in-the-middle attacks. Use only for:

- **Development/Testing**: Testing self-signed certificates
- **Internal Networks**: Trusted internal services  
- **Debugging**: Troubleshooting certificate issues
- **Legacy Systems**: Old systems with expired certificates

**DO NOT USE** in production for external services or when security is critical.

## Verbose Mode Integration

When using both verbose mode (`-v`) and insecure mode (`-k`), TCPing displays a warning:

```bash
$ tcping -H -v -k https://self-signed.badssl.com
正在对 https://self-signed.badssl.com 执行 HTTP Ping (User-Agent: tcping/v1.8.0.unknown)
  警告: SSL/TLS证书验证已禁用
HTTP 200 https://self-signed.badssl.com: seq=0 time=150.23ms size=1024 bytes bandwidth=0.05 Mbps
  详细信息: 状态=200 OK, Content-Type=text/html, Server=nginx
```

## Multi-Language Support

The new option is fully integrated with TCPing's i18n system:

### English (en-US)
```
-k, --insecure          Skip SSL/TLS certificate verification (HTTP mode only)
tcping -H -k https://self-signed.badssl.com  # Skip SSL certificate verification
```

### Japanese (ja-JP)  
```
-k, --insecure          SSL/TLS証明書の検証をスキップ（HTTPモードのみ有効）
tcping -H -k https://self-signed.badssl.com  # SSL証明書検証をスキップ
```

### Korean (ko-KR)
```
-k, --insecure          SSL/TLS 인증서 검증 건너뛰기 (HTTP 모드만 유효)
tcping -H -k https://self-signed.badssl.com  # SSL 인증서 검증 건너뛰기
```

### Traditional Chinese (zh-TW)
```
-k, --insecure          跳過SSL/TLS證書驗證（僅在HTTP模式下有效）
tcping -H -k https://self-signed.badssl.com  # 跳過SSL證書驗證
```

### Simplified Chinese (zh-CN)
```
-k, --insecure          跳过SSL/TLS证书验证（仅在HTTP模式下有效）
tcping -H -k https://self-signed.badssl.com  # 跳过SSL证书验证
```

## Implementation Details

### Code Changes

1. **Options Structure**: Added `InsecureSSL bool` field to `Options` struct
2. **HTTP Transport**: Modified `TLSClientConfig` to use `opts.InsecureSSL`
3. **Command Line**: Added `-k/--insecure` flag parsing
4. **Help Text**: Updated usage examples and descriptions
5. **Verbose Mode**: Added warning message when SSL verification is disabled
6. **i18n Integration**: Added translations for all supported languages

### Technical Implementation
```go
// In httpPingOnce function
transport := &http.Transport{
    TLSClientConfig: &tls.Config{
        InsecureSkipVerify: opts.InsecureSSL,  // Previously hardcoded to false
    },
    DisableKeepAlives: true,
}
```

## Common Use Cases

### 1. Development Testing
```bash
# Test local development server
tcping -H -k https://dev.localhost:3000/api/health

# Test staging environment with self-signed cert  
tcping -H -k -v https://staging.internal/status
```

### 2. Certificate Troubleshooting
```bash
# Compare response with and without certificate validation
tcping -H https://expired.badssl.com          # Will fail
tcping -H -k https://expired.badssl.com       # Will succeed

# Test certificate chain issues
tcping -H -k -v https://incomplete-chain.badssl.com
```

### 3. Internal Service Monitoring
```bash
# Monitor internal service with self-signed certificate
tcping -H -k -n 60 -t 60000 https://internal.company.local/health

# Check API response time ignoring certificate issues  
tcping -H -k -v https://internal-api.local/version
```

### 4. Legacy System Integration
```bash
# Connect to old system with outdated SSL configuration
tcping -H -k https://legacy-system.company.com:8443/status

# Test old SSL/TLS protocols
tcping -H -k -w 10000 https://old.ssl-server.com
```

## Comparison with Other Tools

### curl Equivalent
```bash
# TCPing
tcping -H -k https://self-signed.badssl.com

# curl equivalent  
curl -k -w "@curl-format.txt" https://self-signed.badssl.com
```

### wget Equivalent
```bash
# TCPing
tcping -H -k https://self-signed.badssl.com

# wget equivalent
wget --no-check-certificate https://self-signed.badssl.com
```

### openssl s_client
```bash
# TCPing (tests full HTTP response)
tcping -H -k -v https://self-signed.badssl.com

# openssl (tests SSL handshake only)
openssl s_client -connect self-signed.badssl.com:443 -verify_return_error
```

## Migration Notes

### From Previous Versions
- **Behavior Change**: Previously, SSL verification was always enabled
- **Backward Compatibility**: Default behavior unchanged (SSL verification enabled)
- **New Flag Required**: Must explicitly use `-k` to disable verification

### Integration with Existing Scripts
```bash
# Old script (always verified SSL)
tcping -H https://internal.company.com

# New script (can bypass SSL verification)
tcping -H -k https://internal.company.com  # Add -k flag
```

## Best Practices

### 1. Use Sparingly
Only use `-k` when absolutely necessary. Always prefer proper SSL certificates.

### 2. Document Usage
When using in scripts or documentation, clearly indicate why SSL verification is disabled:
```bash
# Disable SSL verification for development environment only
tcping -H -k https://dev.internal.company.com
```

### 3. Combine with Verbose Mode
Use `-v` flag to see SSL warnings and response details:
```bash
tcping -H -k -v https://self-signed.badssl.com
```

### 4. Monitor Certificate Expiration
Use the insecure mode temporarily while fixing certificate issues:
```bash
# Temporary bypass while renewing certificates
tcping -H -k https://service.company.com
```

### 5. Environment-Specific Configuration
Different configurations for different environments:
```bash
# Production (strict SSL verification)
tcping -H https://api.production.com

# Development (allow self-signed certificates)  
tcping -H -k https://api.development.com
```

The SSL bypass feature provides flexibility for testing and development while maintaining security by default and providing clear warnings when certificate verification is disabled.