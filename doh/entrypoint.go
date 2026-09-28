package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const optionsPath = "/data/options.json"
const dnsproxyPath = "/opt/dnsproxy/dnsproxy"

type options struct {
	Upstreams    string `json:"upstreams"`
	Bootstrap    string `json:"bootstrap"`
	UpstreamMode string `json:"upstream_mode"`
	Timeout      string `json:"timeout"`
	Cache        bool   `json:"cache"`
	CacheSize    int    `json:"cache_size"`
	Verbose      bool   `json:"verbose"`
}

var defaultOptions = options{
	Upstreams:    "https://dns.google/dns-query,https://dns.quad9.net/dns-query",
	Bootstrap:    "8.8.8.8:53,9.9.9.9:53",
	UpstreamMode: "load_balance",
	Timeout:      "10s",
	CacheSize:    65536,
}

func main() {
	config, err := loadOptions()
	if err != nil {
		fatal(err)
	}

	upstreams, err := parseEndpoints("upstreams", config.Upstreams)
	if err != nil {
		fatal(err)
	}
	bootstrap, err := parseBootstrap(config.Bootstrap)
	if err != nil {
		fatal(err)
	}
	if !validUpstreamMode(config.UpstreamMode) {
		fatal(fmt.Errorf("upstream_mode must be load_balance, parallel, or fastest_addr, got %q", config.UpstreamMode))
	}
	timeout, err := time.ParseDuration(config.Timeout)
	if err != nil {
		fatal(fmt.Errorf("timeout must be a Go duration such as 10s: %w", err))
	}
	if timeout <= 0 {
		fatal(errors.New("timeout must be greater than zero"))
	}
	if config.CacheSize < 0 {
		fatal(errors.New("cache_size must be zero or greater"))
	}

	args := []string{
		dnsproxyPath,
		"--listen", "0.0.0.0",
		"--port", "5053",
		"--upstream-mode", config.UpstreamMode,
		"--timeout", config.Timeout,
	}
	for _, endpoint := range upstreams {
		args = append(args, "--upstream", endpoint)
	}
	for _, endpoint := range bootstrap {
		args = append(args, "--bootstrap", endpoint)
	}
	if config.Cache {
		args = append(args, "--cache", "--cache-size", strconv.Itoa(config.CacheSize))
	}
	if config.Verbose {
		args = append(args, "--verbose")
	}

	if err := syscall.Exec(dnsproxyPath, args, os.Environ()); err != nil {
		fatal(fmt.Errorf("start DNSProxy: %w", err))
	}
}

func loadOptions() (options, error) {
	config := defaultOptions
	file, err := os.Open(optionsPath)
	if errors.Is(err, os.ErrNotExist) {
		return config, nil
	}
	if err != nil {
		return options{}, fmt.Errorf("open %s: %w", optionsPath, err)
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return options{}, fmt.Errorf("decode %s: %w", optionsPath, err)
	}
	return config, nil
}

func parseEndpoints(name, value string) ([]string, error) {
	var endpoints []string
	for _, item := range strings.Split(value, ",") {
		endpoint := strings.TrimSpace(item)
		parsed, err := url.ParseRequestURI(endpoint)
		if err != nil || parsed.Scheme != "https" || parsed.Host == "" {
			return nil, fmt.Errorf("%s contains an invalid HTTPS URL: %q", name, endpoint)
		}
		endpoints = append(endpoints, endpoint)
	}
	if len(endpoints) == 0 {
		return nil, fmt.Errorf("%s must contain at least one HTTPS URL", name)
	}
	return endpoints, nil
}

func parseBootstrap(value string) ([]string, error) {
	var resolvers []string
	for _, item := range strings.Split(value, ",") {
		resolver := strings.TrimSpace(item)
		host, port, err := net.SplitHostPort(resolver)
		if err != nil || host == "" || port == "" {
			return nil, fmt.Errorf("bootstrap contains an invalid DNS server address: %q", resolver)
		}
		resolvers = append(resolvers, resolver)
	}
	if len(resolvers) == 0 {
		return nil, errors.New("bootstrap must contain at least one DNS server address")
	}
	return resolvers, nil
}

func validUpstreamMode(mode string) bool {
	return mode == "load_balance" || mode == "parallel" || mode == "fastest_addr"
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "doh: %v\n", err)
	os.Exit(1)
}
