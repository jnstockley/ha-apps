# DNS over HTTPS Home Assistant Apps

[![Lint](https://github.com/jnstockley/ha-apps/actions/workflows/lint.yaml/badge.svg)](https://github.com/jnstockley/ha-apps/actions/workflows/lint.yaml)
[![Builder](https://github.com/jnstockley/ha-apps/actions/workflows/builder.yaml/badge.svg)](https://github.com/jnstockley/ha-apps/actions/workflows/builder.yaml)
[![License: Apache-2.0](https://img.shields.io/badge/license-Apache--2.0-blue.svg)](LICENSE)

Home Assistant app repository for a local DNS-over-HTTPS (DoH) proxy.
The [`doh`](doh/) app accepts ordinary DNS queries over UDP and TCP and
forwards them to configurable DoH upstreams.

## Install

1. In Home Assistant, open **Settings** → **Add-ons** → **Add-on store**.
2. Open the overflow menu, select **Repositories**, and add:
   `https://github.com/jnstockley/ha-apps`
3. Install **DNS over HTTPS**, configure the upstream endpoints, and
   start the app.
4. Configure your DHCP server or individual devices to use the Home
   Assistant host as their DNS server.

The app maps its internal DNS listener to port `53` on the Home
Assistant host by default. See [`doh/DOCS.md`](doh/DOCS.md) for port
mapping and configuration details.

## Repository layout

```text
.
├── doh/                         # DNS-over-HTTPS Home Assistant app
│   ├── config.yaml              # App metadata, ports, and options
│   ├── Dockerfile               # DNSProxy-based runtime image
│   ├── entrypoint.go            # Validates options and starts DNSProxy
│   ├── DOCS.md                  # App documentation
│   └── translations/en.yaml     # Configuration UI text
├── repository.yaml              # Home Assistant repository metadata
└── .github/workflows/           # Build, publish, and release workflows
```

## Supported architectures

The app supports `amd64` and `aarch64`.

## Release automation

Renovate monitors the DNSProxy image tag. Its DNSProxy upgrade PRs also
run [`scripts/bump-doh-version.sh`](scripts/bump-doh-version.sh), which
increments the app manifest patch version. After the PR is merged, the
builder workflow publishes the new image and creates the corresponding
GitHub release.

The Renovate runner must permit `./scripts/bump-doh-version.sh` as a
`postUpgradeTask`; self-hosted Renovate requires this command in its
global `allowedCommands` configuration.

## License

Apache License 2.0 -- see [`LICENSE`](LICENSE).
