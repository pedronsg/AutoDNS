# AutoDNS

[![License: LGPL v3](https://img.shields.io/badge/License-LGPL_v3-blue.svg)](https://www.gnu.org/licenses/lgpl-3.0)

Dynamic DNS updater for OrbitOS. Keeps your public IP in sync with No-IP, DuckDNS, DynDNS, and Cloudflare records.

The app runs a local web UI (default port **9033**) and polls your public IP on a configurable interval. When the IP changes, it updates every enabled DNS entry.

---

## Requirements

- Go 1.21+
- [Orbit Studio VS Code extension](https://marketplace.visualstudio.com/items?itemName=orbit-os.orbit-studio)
- Python 3 (used by the build script)
- An OrbitOS device reachable on the network

---

## Running locally

```bash
go run ./cmd/AutoDNS
```

Config is stored at `cmd/AutoDNS/orb/data/config.json`. You can override the device host at startup:

```bash
go run ./cmd/AutoDNS -host 192.168.5.226
```

Open `http://127.0.0.1:9033` to access the web UI.

---

## Build

Produces an `.orb` package under `.orbit/`.

### VS Code (Orbit Studio extension)

Run the default build task:

```
Ctrl+Shift+B  →  Orbit: Build ORB
```

Or via the sidebar: open the **Orbit** panel and click **Build**.

### CLI

```bash
./orbit-os-sdk-go/scripts/build_orb.sh -path cmd/AutoDNS
```

Defaults to `linux/arm64`. To target a different architecture:

```bash
./orbit-os-sdk-go/scripts/build_orb.sh -path cmd/AutoDNS -arch amd64
```

Output: `.orbit/autodns_v<version>.orb`

---

## Deploy

The device host and package path are configured in [`orbit.project.json`](orbit.project.json).

### VS Code (Orbit Studio extension)

Open the Command Palette and run:

```
Tasks: Run Task  →  Orbit: Deploy to device
```

Or via the sidebar: open the **Orbit** panel and click **Deploy**.

### CLI

```bash
go run ./orbit-os-sdk-go/cmd/orbit-deploy -root .
```

Installs the `.orb` from `.orbit/` onto the device defined in `orbit.project.json`.

---

## Project structure

```
cmd/AutoDNS/
├── main.go              # Entry point, device connection
├── metadata.json        # App manifest (name, version, permissions)
├── config/              # Config struct and persistence
├── providers/           # DNS provider implementations
├── updater/             # IP polling and update loop
├── web/                 # HTTP server and embedded UI
└── orb/
    ├── icon.svg         # App icon
    └── data/            # Runtime data — config.json is written here
```

---

## Configuration

Settings are saved to `orb/data/config.json` on first run and editable through the web UI.

| Field | Default | Description |
|---|---|---|
| `check_interval_min` | `5` | How often to poll for IP changes (minutes) |
| `ip_detector_url` | `https://api.ipify.org` | Service used to detect the public IP |
| `web_port` | `9033` | Local web UI port |
| `device_host` | `192.168.5.226` | OrbitOS device IP (overridable with `-host`) |
