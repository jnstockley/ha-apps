package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

const optionsPath = "/data/options.json"

type options struct {
	Upstreams        string `json:"upstreams"`
	Bootstrap        string `json:"bootstrap"`
	MaxUpstreamConns int    `json:"max_upstream_conns"`
	Metrics          string `json:"metrics"`
}

var defaultOptions = options{
	Upstreams:        "https://1.1.1.1/dns-query,https://1.0.0.1/dns-query",
	Bootstrap:        "https://162.159.36.1/dns-query,https://162.159.46.1/dns-query",
	MaxUpstreamConns: 5,
	Metrics:          "127.0.0.1:39441",
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
	bootstrap, err := parseEndpoints("bootstrap", config.Bootstrap)
	if err != nil {
		fatal(err)
	}
	if config.MaxUpstreamConns < 0 {
		fatal(errors.New("max_upstream_conns must be zero or greater"))
	}
	if _, _, err := net.SplitHostPort(config.Metrics); err != nil {
		fatal(fmt.Errorf("metrics must be a host:port address: %w", err))
	}

	cloudflared, err := exec.LookPath("cloudflared")
	if err != nil {
		fatal(fmt.Errorf("find cloudflared: %w", err))
	}

	args := []string{
		cloudflared,
		"--no-autoupdate",
		"proxy-dns",
		"--address", "0.0.0.0",
		"--port", "5053",
		"--metrics", config.Metrics,
		"--max-upstream-conns", strconv.Itoa(config.MaxUpstreamConns),
	}
	for _, endpoint := range upstreams {
		args = append(args, "--upstream", endpoint)
	}
	for _, endpoint := range bootstrap {
		args = append(args, "--bootstrap", endpoint)
	}

	if err := syscall.Exec(cloudflared, args, os.Environ()); err != nil {
		fatal(fmt.Errorf("start cloudflared: %w", err))
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

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "cloudflare-doh: %v\n", err)
	os.Exit(1)
}
