package checker

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/sync/semaphore"

	"github.com/TheRemyyy/gopunch/internal/health"
)

// Options configures the checking behavior
type Options struct {
	Timeout         time.Duration
	Method          string
	Headers         map[string]string
	Body            string
	Insecure        bool
	FollowRedirects bool
	ExpectedCodes   []int
	Retries         int
}

// Result represents the outcome of a check
type Result struct {
	URL        string
	StatusCode int    // HTTP status code
	Status     string // HTTP status string
	Info       string // Extra info for non-HTTP checks (e.g. "Open", "Valid")
	Duration   time.Duration
	Size       int64
	Headers    http.Header
	Success    bool
	Error      error
	Retries    int
}

// CheckURLs performs concurrent health checks
func CheckURLs(urls []string, opts Options, concurrency int) []Result {
	results := make([]Result, len(urls))
	sem := semaphore.NewWeighted(int64(concurrency))
	var wg sync.WaitGroup

	for i, u := range urls {
		wg.Add(1)
		go func(idx int, target string) {
			defer wg.Done()
			_ = sem.Acquire(context.Background(), 1)
			defer sem.Release(1)

			// Detect scheme
			if strings.HasPrefix(target, "tcp://") {
				results[idx] = checkTCP(target, opts)
			} else if strings.HasPrefix(target, "dns://") {
				results[idx] = checkDNS(target, opts)
			} else if strings.HasPrefix(target, "ssl://") {
				results[idx] = checkSSL(target, opts)
			} else {
				// Default to HTTP
				if !strings.HasPrefix(target, "http") {
					target = "https://" + target
				}
				results[idx] = checkHTTP(target, opts)
			}
		}(i, u)
	}

	wg.Wait()
	return results
}

func checkTCP(target string, opts Options) Result {
	u, err := url.Parse(target)
	if err != nil {
		return Result{URL: target, Error: err}
	}

	host := u.Hostname()
	portStr := u.Port()
	port := 80
	if portStr != "" {
		port, _ = strconv.Atoi(portStr)
	}

	res := health.CheckTCP(host, port, opts.Timeout)

	result := Result{
		URL:      target,
		Duration: res.Duration,
		Success:  res.Open,
		Error:    res.Error,
	}

	if res.Open {
		result.Info = "Open"
	}
	return result
}

func checkDNS(target string, opts Options) Result {
	u, err := url.Parse(target)
	if err != nil {
		// handle "dns://hostname" or regular hostname
		target = strings.TrimPrefix(target, "dns://")
	} else {
		target = u.Hostname()
		if target == "" {
			target = u.Path // Handle dns://hostname case if simple
		}
	}

	// If parse failed somewhat or target still has scheme
	target = strings.TrimPrefix(target, "dns://")

	res := health.CheckDNS(target, opts.Timeout)

	result := Result{
		URL:      "dns://" + target,
		Duration: res.Duration,
		Success:  res.Resolved,
		Error:    res.Error,
	}

	if res.Resolved {
		result.Info = fmt.Sprintf("%d IPs", len(res.IPs))
	}
	return result
}

func checkSSL(target string, opts Options) Result {
	u, err := url.Parse(target)
	if err != nil {
		return Result{URL: target, Error: err}
	}

	host := u.Hostname()
	portStr := u.Port()
	port := 443
	if portStr != "" {
		port, _ = strconv.Atoi(portStr)
	}

	res := health.CheckSSL(host, port, opts.Timeout)

	result := Result{
		URL:      target,
		Duration: res.Duration,
		Success:  res.Valid,
		Error:    res.Error,
	}

	if res.Valid {
		result.Info = fmt.Sprintf("Valid (%d days)", res.DaysLeft)
	}
	return result
}

func checkHTTP(url string, opts Options) Result {
	var result Result
	result.URL = url

	client := createClient(opts)

	for attempt := 0; attempt <= opts.Retries; attempt++ {
		result = doHTTPCheck(url, opts, client)
		result.Retries = attempt

		if result.Error == nil && result.Success {
			break
		}

		if attempt < opts.Retries {
			time.Sleep(time.Duration(100*(1<<attempt)) * time.Millisecond)
		}
	}

	return result
}

func createClient(opts Options) *http.Client {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: opts.Insecure,
		},
		DisableKeepAlives:     true,
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   10,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	client := &http.Client{
		Timeout:   opts.Timeout,
		Transport: transport,
	}

	if !opts.FollowRedirects {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}

	return client
}

func doHTTPCheck(url string, opts Options, client *http.Client) Result {
	result := Result{URL: url}
	start := time.Now()

	var body io.Reader
	if opts.Body != "" {
		body = bytes.NewBufferString(opts.Body)
	}

	method := opts.Method
	if method == "" {
		method = "GET"
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		result.Error = fmt.Errorf("request creation failed: %w", err)
		result.Duration = time.Since(start)
		return result
	}

	req.Header.Set("User-Agent", "GoPunch/2.0")
	if opts.Body != "" {
		req.Header.Set("Content-Type", "application/json")
	}

	for key, value := range opts.Headers {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	result.Duration = time.Since(start)

	if err != nil {
		result.Error = fmt.Errorf("request failed: %w", err)
		return result
	}
	defer resp.Body.Close()

	result.StatusCode = resp.StatusCode
	result.Status = resp.Status
	result.Headers = resp.Header

	bodyBytes, _ := io.ReadAll(resp.Body)
	result.Size = int64(len(bodyBytes))

	result.Success = isSuccessCode(resp.StatusCode, opts.ExpectedCodes)

	return result
}

func isSuccessCode(code int, expected []int) bool {
	if len(expected) > 0 {
		for _, e := range expected {
			if code == e {
				return true
			}
		}
		return false
	}
	return code >= 200 && code < 400
}
