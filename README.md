# Home Assistant App Template

[![Lint](https://github.com/jnstockley/ha-apps/actions/workflows/lint.yaml/badge.svg)](https://github.com/jnstockley/ha-apps/actions/workflows/lint.yaml)
[![Builder](https://github.com/jnstockley/ha-apps/actions/workflows/builder.yaml/badge.svg)](https://github.com/jnstockley/ha-apps/actions/workflows/builder.yaml)
![Supports amd64](https://img.shields.io/badge/amd64-yes-success.svg)
![Supports aarch64](https://img.shields.io/badge/aarch64-yes-success.svg)
[![License: Apache-2.0](https://img.shields.io/badge/license-Apache--2.0-blue.svg)](LICENSE)

A ready-to-fork template for building a **Home Assistant app** (what used
to be called an "add-on"): a Docker container, managed by the Home
Assistant Supervisor, that exposes its own configurable options and a
web service on a real network port.

It includes:

- A working starter app (`starter_app/`) -- a small Python web service
  demonstrating five different configuration-option types, a health
  endpoint wired up as a watchdog, a JSON API, and an optional Home
  Assistant Core API integration.
- A GitHub Actions pipeline (`.github/workflows/`) that lints, builds
  multi-arch (`amd64` + `aarch64`) images with Docker BuildKit, publishes
  them to the GitHub Container Registry, and cuts a GitHub release.
- Everything needed for Home Assistant to recognize this as an
  installable app repository (`repository.yaml`).
- A VS Code devcontainer for one-command local testing against a real
  Supervisor + Home Assistant Core instance.

## Quickstart

1. **Use this repository as a template** (the green "Use this template"
   button on GitHub), or clone/fork it.
2. Find-and-replace `jnstockley` and `ha-apps`
   throughout the repo (`config.yaml`, `Dockerfile`, `repository.yaml`,
   this file) with your actual GitHub username/org and repo name.
3. Rename `starter_app/` and update its `config.yaml` -- see
   [`starter_app/DOCS.md`](starter_app/DOCS.md#customizing-this-template)
   for the full checklist.
4. Test locally (pick one):
   - **Devcontainer (recommended):** open the repo in VS Code, reopen in
     the container when prompted, then run the **Start Home Assistant**
     task. Your app appears under Settings → Add-ons → Local add-ons.
   - **Plain Docker:** `cd starter_app && docker build -t local/starter-app .`
     then `docker run --rm -p 8090:8090 -v /tmp/starter-data:/data local/starter-app`.
     See [Local app testing](https://developers.home-assistant.io/docs/apps/testing)
     for details.
5. Push to `main`. The `Builder` workflow builds the app on every push
   and pull request; on `main` it also publishes images to
   `ghcr.io/<you>/<slug>` and cuts a GitHub release. Until then, leave
   `image:` commented out in `config.yaml` so Supervisor builds the
   image locally instead of trying to pull one that doesn't exist yet.
6. Once you're ready to distribute it, uncomment `image:` in
   `config.yaml` and share your repository:

   [![Open your Home Assistant instance and show the add app repository dialog with a specific repository URL pre-filled.](https://my.home-assistant.io/badges/supervisor_add_addon_repository.svg)](https://my.home-assistant.io/redirect/supervisor_add_addon_repository/?repository_url=https%3A%2F%2Fgithub.com%2Fjnstockley%2Fha-apps)

   (Update the URL-encoded `repository_url` query parameter to match your
   real repo before sharing this button.)

## Repository layout

```text
.
├── repository.yaml            # Identifies this repo to Home Assistant
├── starter_app/                # One app = one folder. Add more the same way.
│   ├── config.yaml             # Options, ports, arch, schema, etc.
│   ├── Dockerfile
│   ├── server.py               # The actual web service
│   ├── apparmor.txt            # Custom security profile (+1 security point)
│   ├── translations/en.yaml    # UI labels/descriptions for options
│   ├── DOCS.md                 # Shown in the app's "Documentation" tab
│   ├── README.md                # Shown in the app's store listing
│   ├── CHANGELOG.md
│   └── rootfs/etc/services.d/starter_app/
│       ├── run                 # Reads options via bashio, starts the service
│       └── finish              # Correctly halts the app if the service exits
├── .github/workflows/
│   ├── lint.yaml               # hadolint, yamllint, ShellCheck
│   ├── builder.yaml            # Detects which app(s) changed
│   └── build-app.yaml          # Builds, publishes, and releases one app
├── .devcontainer/, .vscode/     # One-command local Supervisor + Core testing
└── .hadolint.yaml, .yamllint.yaml, .github/dependabot.yml
```

## Adding a second app

Copy `starter_app/` to a new folder and repeat the rename steps above.
`builder.yaml` discovers app folders automatically (via
`home-assistant/actions/helpers/find-addons`) and only rebuilds the ones
that actually changed -- no workflow edits needed.

## Supported architectures

Only `amd64` and `aarch64` are wired up, matching what the current
[home-assistant/builder](https://github.com/home-assistant/builder)
composite actions support (32-bit ARM and i386 were dropped in the
`2025.11.0` builder release). If you must support `armv7`/`armhf`, you'll
need an older, unmaintained release of those actions, or a self-hosted
QEMU-based build step -- neither is recommended for a new app in 2026.

## Design choices worth knowing about

- **No `build.yaml`.** Home Assistant retired it in favor of a plain
  `FROM`/`ARG`/`LABEL` Dockerfile -- see the
  ["Migrating app builds to Docker BuildKit"](https://developers.home-assistant.io/blog/2026/04/02/builder-migration)
  post. If you've built Home Assistant add-ons before and are looking
  for it, this is why it's missing.
- **Ports, not Ingress, by default.** Ports are simpler to understand
  and work outside of Home Assistant too. Ingress is more secure (no
  open port on your LAN, Home Assistant handles auth) -- see
  `starter_app/DOCS.md` for how to switch.
- **Zero Python dependencies.** The starter app uses only the standard
  library so the image stays small and there's nothing extra to patch.
- **Native ARM CI runners.** The build pipeline uses GitHub's
  `ubuntu-24.04-arm` runners for the `aarch64` build, so there's no QEMU
  emulation slowing down CI.

## License

Apache License 2.0 -- see [`LICENSE`](LICENSE). Replace the copyright
line with your own name/organization.
