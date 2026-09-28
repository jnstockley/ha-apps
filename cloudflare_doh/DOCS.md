# Cloudflare DoH Documentation

This app uses the official `cloudflare/cloudflared:2025.11.1` image in
`proxy-dns` mode. It accepts standard DNS queries and forwards them over
HTTPS to Cloudflare DNS. cloudflared removed `proxy-dns` in version
`2026.2.0`, so the image is deliberately pinned to the final supported
release.

## Network configuration

cloudflared listens on container port `5053` for both UDP and TCP. By
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

- `https://1.1.1.1/dns-query`
- `https://1.0.0.1/dns-query`

The comma-separated **Bootstrap DoH endpoints** resolve upstream
hostnames before the proxy is available. They default to:

- `https://162.159.36.1/dns-query`
- `https://162.159.46.1/dns-query`

**Maximum upstream connections** maps to cloudflared's
`--max-upstream-conns` option. Set it to `0` to allow unlimited
connections.

**Metrics address** maps to `--metrics`. It must be a `host:port`
address. The default (`127.0.0.1:39441`) does not expose metrics outside
the add-on. No queries are served if cloudflared cannot reach an
upstream endpoint.
