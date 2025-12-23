package health

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"time"
)

// TCPResult represents TCP port check result
type TCPResult struct {
	Host     string
	Port     int
	Open     bool
	Duration time.Duration
	Error    error
}

// CheckTCP checks if a TCP port is open
func CheckTCP(host string, port int, timeout time.Duration) TCPResult {
	result := TCPResult{Host: host, Port: port}
	start := time.Now()

	address := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.DialTimeout("tcp", address, timeout)
	result.Duration = time.Since(start)

	if err != nil {
		result.Error = err
		result.Open = false
		return result
	}
	defer conn.Close()

	result.Open = true
	return result
}

// DNSResult represents DNS resolution result
type DNSResult struct {
	Host     string
	IPs      []string
	Duration time.Duration
	Error    error
	Resolved bool
}

// CheckDNS resolves a hostname and returns IPs
func CheckDNS(host string, timeout time.Duration) DNSResult {
	result := DNSResult{Host: host}
	start := time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	resolver := net.Resolver{}
	ips, err := resolver.LookupHost(ctx, host)
	result.Duration = time.Since(start)

	if err != nil {
		result.Error = err
		result.Resolved = false
		return result
	}

	result.IPs = ips
	result.Resolved = len(ips) > 0
	return result
}

// SSLResult represents SSL certificate check result
type SSLResult struct {
	Host      string
	Port      int
	Valid     bool
	Issuer    string
	Subject   string
	NotBefore time.Time
	NotAfter  time.Time
	DaysLeft  int
	Duration  time.Duration
	Error     error
}

// CheckSSL checks SSL certificate validity and expiry
func CheckSSL(host string, port int, timeout time.Duration) SSLResult {
	result := SSLResult{Host: host, Port: port}
	start := time.Now()

	if port == 0 {
		port = 443
	}

	address := fmt.Sprintf("%s:%d", host, port)

	dialer := &net.Dialer{Timeout: timeout}
	conn, err := tls.DialWithDialer(dialer, "tcp", address, &tls.Config{
		InsecureSkipVerify: false,
	})
	result.Duration = time.Since(start)

	if err != nil {
		result.Error = err
		result.Valid = false
		return result
	}
	defer conn.Close()

	certs := conn.ConnectionState().PeerCertificates
	if len(certs) == 0 {
		result.Error = fmt.Errorf("no certificates found")
		result.Valid = false
		return result
	}

	cert := certs[0]
	result.Issuer = cert.Issuer.CommonName
	result.Subject = cert.Subject.CommonName
	result.NotBefore = cert.NotBefore
	result.NotAfter = cert.NotAfter
	result.DaysLeft = int(time.Until(cert.NotAfter).Hours() / 24)
	result.Valid = time.Now().Before(cert.NotAfter) && time.Now().After(cert.NotBefore)

	return result
}
