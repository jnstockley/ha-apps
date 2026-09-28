# `proxy-dns` was removed in cloudflared 2026.2.0. Keep this image on the
# final release that supports the local DoH proxy command.
FROM golang:1.25-alpine AS launcher

WORKDIR /src
COPY entrypoint.go .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /cloudflare-doh ./entrypoint.go

FROM cloudflare/cloudflared:2025.11.1

ARG BUILD_ARCH
ARG BUILD_VERSION="0.0.0-dev"

COPY --from=launcher /cloudflare-doh /cloudflare-doh

LABEL \
  io.hass.version="${BUILD_VERSION}" \
  io.hass.type="app" \
  io.hass.arch="${BUILD_ARCH}" \
  org.opencontainers.image.title="Cloudflare DoH" \
  org.opencontainers.image.description="A local DNS server forwarding queries to Cloudflare over HTTPS" \
  org.opencontainers.image.source="https://github.com/jnstockley/ha-apps" \
  org.opencontainers.image.licenses="Apache-2.0" \
  org.opencontainers.image.base.name="docker.io/cloudflare/cloudflared:2025.11.1"

ENTRYPOINT ["/cloudflare-doh"]
