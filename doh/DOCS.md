# DNS over HTTPS Documentation

This app uses the official `adguard/dnsproxy:v0.85.0` image. It accepts
standard DNS queries and forwards them over HTTPS to configurable DoH
upstreams. It uses a maintained DNS proxy with native
DNS-over-HTTPS-upstream support.

## Network configuration

DNSProxy listens on container port `5053` for both UDP and TCP. By
default, Home Assistant publishes that listener as port `53` on the
Home Assistant host.

Configure DHCP or individual devices to use the Home Assistant host's
IP address as their DNS server. If the host port `53` is unavailable,
edit the port mapping in the app's **Network** settings. Devices must
then support and be configured for the selected non-standard DNS port.

Both transport protocols are published: UDP handles normal DNS traffic,
while TCP is required for larger DNS responses and fallback queries.

## Upstream resolvers

The app forwards DNS queries to the comma-separated **Upstream DoH
endpoints** configured in the add-on. The defaults are:

- `https://dns.google/dns-query`
- `https://dns.quad9.net/dns-query`

The comma-separated **Bootstrap DNS servers** resolve upstream hostnames
before the proxy is available. Each address must include a port. They
default to:

- `8.8.8.8:53`
- `9.9.9.9:53`

## DNSProxy settings

**Upstream mode** maps to DNSProxy's `--upstream-mode`: `load_balance`
uses one resolver per query, `parallel` queries every resolver, and
`fastest_addr` selects the fastest responder. **Upstream timeout** maps
to `--timeout` and accepts Go duration syntax, such as `10s` or `500ms`.

Enable **DNS cache** to cache responses; **DNS cache size** sets the
maximum cache size in bytes. Enable **Verbose logging** only when
troubleshooting. No queries are served if DNSProxy cannot reach an
upstream endpoint.
