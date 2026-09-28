# Cloudflare DoH

Runs `cloudflare/cloudflared:2025.11.1` as a local DNS server for Home
Assistant. DNS queries received over TCP or UDP are encrypted and
forwarded to Cloudflare's `1.1.1.1` and `1.0.0.1` DNS-over-HTTPS
endpoints. This is the final cloudflared release that supports its local
`proxy-dns` command.

## Use

1. Install and start the app.
2. Ensure no other service on the Home Assistant host uses port `53`.
3. Set your clients' DNS server to the Home Assistant host's IP address.

The app maps its internal listener on port `5053` to host port `53` for
both UDP and TCP. If port `53` is already in use, change the host-side
port mapping in Home Assistant, then configure clients to use that port.

## Configuration

Set multiple **Upstream DoH endpoints** or **Bootstrap DoH endpoints** as
comma-separated HTTPS URLs. `max_upstream_conns` controls the number of
concurrent connections to each upstream; use `0` for no limit. The
metrics listener defaults to `127.0.0.1:39441`, which remains accessible
only inside the add-on container.

## Support

Open an issue on the repository this app was installed from.
