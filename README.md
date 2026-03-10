# SyncCast

**Local-First Media Teleportation**

SyncCast is a high-performance, zero-cloud media streaming system that teleports 1080p video from a local source to a Smart TV. It uses the TV's internal browser as a playback target and WebSockets for real-time control, bypassing proprietary casting restrictions entirely. No accounts. No telemetry. No cloud relay. Just your file, your network, your TV.

Licensed under [GPLv3](LICENSE).

---

## Why SyncCast?

Every existing casting solution — Chromecast, AirPlay, Miracast — routes through proprietary protocols, cloud handshakes, or vendor-locked ecosystems. On restricted networks (university residence halls, hospital Wi-Fi, hotel captive portals), these protocols are blocked at the firewall level. SyncCast sidesteps all of this.

It discovers devices on the local subnet, serves the media file over plain HTTP with range-request support, and delivers a browser-based player directly to the TV. The controller is a PWA on your phone. The entire pipeline lives on your LAN.

---

## Architecture

```
┌──────────────┐         ┌──────────────────┐          ┌────────────────┐
│  Mobile PWA  │◄──WS───►│  SyncCast Core   │◄──HTTP──►│   Smart TV     │
│  (Remote)    │         │  (Go binary)     │          │   (Browser)    │
└──────────────┘         └──────────────────┘          └────────────────┘
                                 │
                          ┌──────┴──────┐
                          │ Local Files │
                          └─────────────┘
```

**Core** runs on any machine on the LAN (laptop, Raspberry Pi, NAS). It handles:

- **Discovery** — finds Smart TVs via SSDP multicast and subnet probing
- **Streaming** — serves media over HTTP/1.1 with `Range` header support (HTTP 206) for seeking
- **Control** — WebSocket hub for play/pause/seek/volume commands between the mobile remote and the TV player
- **State** — connection lifecycle management, device tracking, reconnection logic

---

## How It Works

1. **Discover** — SyncCast sends SSDP M-SEARCH packets to `239.255.255.250:1900` targeting DIAL and UPnP MediaRenderer services. If multicast is blocked (common on institutional networks), it falls back to a concurrent /24 subnet scan, probing known device-description endpoints on ports 8008/8009.

2. **Serve** — The selected media file is exposed on a local HTTP server with byte-range support, enabling the TV browser to seek freely without re-downloading the entire file.

3. **Launch** — The TV's browser is directed to a lightweight HTML5 player page served by SyncCast. The player connects back over a WebSocket for bidirectional control.

4. **Control** — Your phone opens the same SyncCast server as a PWA. Play, pause, seek, volume — all commands flow through the WebSocket hub with sub-50ms latency on a typical LAN.

---

## Project Structure

```
SyncCast/
├── cmd/
│   └── synccast/
│       └── main.go              # CLI entry point
├── internal/
│   ├── discovery/
│   │   ├── device.go            # Device struct definition
│   │   ├── discovery.go         # Orchestrator: SSDP-first, subnet-scan fallback
│   │   ├── finder.go            # Local IP resolution
│   │   ├── probe.go             # Concurrent /24 subnet scanner
│   │   └── ssdp.go              # SSDP M-SEARCH implementation
│   ├── control/
│   │   └── soap.go              # SOAP/UPnP device control (planned)
│   ├── state/
│   │   └── manager.go           # Connection state machine (planned)
│   └── streamer/
│       └── server.go            # HTTP media server with Range support (planned)
├── go.mod
├── LICENSE                      # GPLv3
└── README.md
```

---

## Current Status

**What works today:**

- SSDP multicast discovery targeting DIAL and UPnP MediaRenderer devices
- Concurrent subnet scanning as a fallback when multicast is blocked
- Friendly name extraction from device XML descriptors
- Duplicate device deduplication
- Local IP detection

**What's in progress:**

See the roadmap below.

---

## Roadmap

- [x] Network discovery and endpoint mapping (SSDP + subnet scan)
- [ ] Manual IP fallback probe for restricted Wi-Fi environments (hostel/university networks)
- [ ] High-performance media engine with HTTP 206 byte-range serving
- [ ] Real-time WebSocket communication hub
- [ ] Connection state management — device drops, auto-reconnect
- [ ] Interactive TV viewer (HTML5 player served to TV browser)
- [ ] WakeLock integration to prevent TV sleep during playback
- [ ] Mobile remote PWA with home-screen install and offline UI caching
- [ ] CLI toolchain and final integration

---

## Getting Started

### Prerequisites

- Go 1.25+
- A Smart TV with a built-in web browser (Android TV, Tizen, webOS, etc.)
- Both devices on the same local network

### Build

```sh
git clone https://github.com/HRITHIK-SANKAR-R/SyncCast.git
cd SyncCast
go build -o synccast ./cmd/synccast
```

### Run

```sh
./synccast
```

SyncCast will print your local IP, attempt SSDP discovery, and fall back to a subnet scan if needed. Discovered devices are listed with their friendly names and IP addresses.

---

## Restricted Network Mode

On networks that block multicast (802.1X enterprise Wi-Fi, captive portals, AP isolation), SSDP will fail silently. SyncCast handles this automatically by scanning the /24 subnet with 50 concurrent goroutines, probing known Smart TV endpoints. A manual IP entry mode is planned for environments where even scanning is restricted.

---

## Contributing

Contributions are welcome. Please:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/your-feature`)
3. Commit with clear messages
4. Open a pull request against `main`

If you're working on one of the roadmap items, open an issue first to avoid duplicate work.

---

## License

SyncCast is free software released under the [GNU General Public License v3.0](LICENSE).

You are free to use, modify, and redistribute this software under the terms of the GPLv3. See the LICENSE file for the full text.
