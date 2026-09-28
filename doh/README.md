# DNS over HTTPS

Runs the maintained `adguard/dnsproxy` image as a local DNS server for
Home Assistant. DNS queries received over TCP or UDP are encrypted and
forwarded to configurable DNS-over-HTTPS endpoints.

## Use

1. Install and start the app.
2. Ensure no other service on the Home Assistant host uses port `53`.
3. Set your clients' DNS server to the Home Assistant host's IP address.

The app maps its internal listener on port `5053` to host port `53` for
both UDP and TCP. If port `53` is already in use, change the host-side
port mapping in Home Assistant, then configure clients to use that port.

## Configuration

Set multiple **Upstream DoH endpoints** as comma-separated HTTPS URLs.
**Bootstrap DNS servers** are comma-separated DNS server addresses
with ports, such as `1.1.1.1:53`, that resolve upstream hostnames before
the proxy is available. DNSProxy also exposes upstream selection,
timeout, caching, and verbose-logging options.

## Support

Open an issue on the repository this app was installed from.
