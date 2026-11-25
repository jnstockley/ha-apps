FROM golang:1.27-alpine AS launcher

WORKDIR /src
COPY entrypoint.go .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /doh ./entrypoint.go

FROM adguard/dnsproxy:v0.86.0

ARG BUILD_ARCH
ARG BUILD_VERSION="0.0.0-dev"

COPY --from=launcher /doh /doh

LABEL \
  io.hass.version="${BUILD_VERSION}" \
  io.hass.type="app" \
  io.hass.arch="${BUILD_ARCH}" \
  org.opencontainers.image.title="DNS over HTTPS" \
  org.opencontainers.image.description="A local DNS server forwarding queries over HTTPS" \
  org.opencontainers.image.source="https://github.com/jnstockley/ha-apps" \
  org.opencontainers.image.licenses="Apache-2.0" \
  org.opencontainers.image.base.name="docker.io/adguard/dnsproxy:v0.85.0"

ENTRYPOINT ["/doh"]
