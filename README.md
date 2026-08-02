# zte-cli

An open-source CLI tool for managing **ZTE F609** GPON routers directly from your terminal. Built in Go, it produces a single lightweight binary (~6.6MB) with zero runtime dependencies.

> **No need to open the router web panel** — manage everything from the command line.

## Features

### Status & Monitoring
- **Device Information** — model, serial number, hardware/software version
- **WAN Status** — connection name, IP address, gateway, DNS, online duration
- **PON Optical** — GPON state, Rx/Tx power, voltage, bias current, temperature
- **WLAN Status** — 4 SSIDs with auth type, encryption, MAC, packet stats
- **LAN Status** — 4 Ethernet ports with speed, mode, packet stats
- **VoIP Status** — phone registration status
- **Connected Clients** — all LAN & WiFi devices with MAC, IP, hostname

### WiFi Management
- **Show/Change SSID & Password** — with AES-CBC encryption/decryption
- **WiFi Basic Settings** — channel, bandwidth, transmit power, 802.11 mode
- **WiFi Radio Control** — enable/disable WiFi radio
- **QR Code Generation** — scan to connect WiFi

### Network Configuration
- **DHCP Server** — view settings and active leases
- **WAN Connections** — list and manage WAN connections
- **Port Forwarding** — full CRUD (add/delete/enable/disable rules)

### Security
- **Firewall** — view/set firewall level (Off/Low/Middle/High)
- **IP Filter** — block/allow by protocol, IP range, port range
- **MAC Filter** — block/allow by MAC address
- **URL Filter** — block websites by URL

### Applications
- **DDNS** — Dynamic DNS settings (dipc/dyndns/DtDNS/No-IP)
- **DMZ** — Demilitarized Zone host configuration
- **UPnP** — Universal Plug and Play settings

### System Administration
- **Reboot** — restart router remotely
- **Factory Reset** — restore to default settings
- **Session Timeout** — view/set login timeout (1-30 min)
- **System Logs** — view router logs with level filtering
- **Change Password** — update router admin password

### Diagnosis
- **Ping** — ping from router to any host
- **ARP Table** — view ARP entries

### Interface
- **Interactive TUI** — 7 menus with number-based navigation
- **Direct CLI** — 35+ commands for scripting and automation

---

## Installation

### Pre-built Binaries (Recommended)

Download from [GitHub Releases](https://github.com/CallMeJaja/zte-cli/releases):

```bash
# Linux AMD64
wget https://github.com/CallMeJaja/zte-cli/releases/latest/download/zte-cli_linux_amd64.tar.gz
tar -xzf zte-cli_linux_amd64.tar.gz
sudo mv zte-cli /usr/local/bin/

# Linux ARM64 (for STB / Raspberry Pi)
wget https://github.com/CallMeJaja/zte-cli/releases/latest/download/zte-cli_linux_arm64.tar.gz
tar -xzf zte-cli_linux_arm64.tar.gz
sudo mv zte-cli /usr/local/bin/

# macOS
wget https://github.com/CallMeJaja/zte-cli/releases/latest/download/zte-cli_darwin_amd64.tar.gz
tar -xzf zte-cli_darwin_amd64.tar.gz
sudo mv zte-cli /usr/local/bin/
```

### Go Install

```bash
go install github.com/CallMeJaja/zte-cli@latest
```

### Docker

```bash
docker pull ghcr.io/CallMeJaja/zte-cli:latest
docker run -it -v ~/.config/zte-cli:/root/.config/zte-cli ghcr.io/CallMeJaja/zte-cli status
```

### Build from Source

```bash
git clone https://github.com/CallMeJaja/zte-cli.git
cd zte-cli
go build -ldflags="-w -s" -o zte-cli .
sudo mv zte-cli /usr/local/bin/
```

---

## Configuration

Create a config file:

```bash
zte-cli config init
```

Then edit `~/.config/zte-cli/config.yaml`:

```yaml
router_host: "192.168.100.1"
router_username: "admin"
router_password: "admin"
```

### Config Search Priority

1. `./config.yaml` (current directory)
2. `~/.config/zte-cli/config.yaml` (user home, recommended)
3. `/etc/zte-cli/config.yaml` (system-wide)

---

## Usage

### Interactive Mode

```bash
zte-cli
```

```
==========================================================
                 ZTE F609 Router Manager v2.0
==========================================================
(1) Status                         (5) Administration
(2) Network                        (6) Diagnosis
(3) Security                       (7) System
(4) Application                    (0) Exit
==========================================================
```

### Direct Commands

#### Status

```bash
zte-cli status                        # Full status summary
zte-cli clients                       # List connected devices
```

#### WiFi

```bash
zte-cli wifi show                     # Show SSID & password
zte-cli wifi set --ssid N --pass P    # Change SSID/password
zte-cli wifi qr                       # Generate QR code
zte-cli wifi basic                    # Show channel, power, mode
zte-cli wifi set-channel 6            # Set channel (Auto, 1-13)
zte-cli wifi set-power 80             # Set power (20/40/60/80/100)
zte-cli wifi enable                   # Enable WiFi radio
zte-cli wifi disable                  # Disable WiFi radio
```

#### LAN / DHCP

```bash
zte-cli lan dhcp                      # Show DHCP settings
zte-cli lan dhcp-leases               # Show active leases
```

#### WAN

```bash
zte-cli wan list                      # List WAN connections
zte-cli wan delete <name>             # Delete WAN connection
```

#### Port Forwarding

```bash
zte-cli portfw list                   # List all rules
zte-cli portfw add \
  --name "WebServer" \
  --proto TCP \
  --wan-port 80 \
  --lan-ip 192.168.100.10 \
  --lan-port 80
zte-cli portfw delete "WebServer"     # Delete rule
zte-cli portfw enable "WebServer"     # Enable rule
zte-cli portfw disable "WebServer"    # Disable rule
```

#### Firewall & Filters

```bash
zte-cli firewall show                 # Show firewall settings
zte-cli firewall level high           # Set level (off/low/mid/high)
zte-cli firewall ipfilter list        # List IP filter rules
zte-cli firewall macfilter list       # List MAC filter rules
zte-cli firewall urlfilter list       # List URL filter rules
zte-cli firewall urlfilter add "facebook.com"  # Block URL
```

#### Applications

```bash
zte-cli ddns show                     # Show DDNS settings
zte-cli dmz show                      # Show DMZ settings
zte-cli dmz set 192.168.100.10        # Set DMZ host
zte-cli upnp show                     # Show UPnP settings
zte-cli upnp enable                   # Enable UPnP
zte-cli upnp disable                  # Disable UPnP
```

#### Diagnosis

```bash
zte-cli ping 8.8.8.8                  # Ping from router
zte-cli ping google.com               # Ping by hostname
zte-cli arp                           # Show ARP table
```

#### System

```bash
zte-cli reboot                        # Reboot router
zte-cli reboot --yes                  # Reboot without confirmation
zte-cli reset                         # Factory reset (WARNING!)
zte-cli timeout                       # Show session timeout
zte-cli timeout 15                    # Set timeout to 15 min
zte-cli log                           # Show last 50 log entries
zte-cli log 100                       # Show last 100 entries
zte-cli password set \
  --old "admin" \
  --new "newpass" \
  --confirm "newpass"                 # Change password
```

#### Config

```bash
zte-cli config init                   # Create config template
zte-cli config path                   # Show active config path
```

---

## All Commands Reference

| Category | Command | Description |
|----------|---------|-------------|
| Status | `status` | Full status summary |
| Status | `clients` | List connected devices |
| WiFi | `wifi show` | Show SSID & password |
| WiFi | `wifi set --ssid N --pass P` | Change SSID/password |
| WiFi | `wifi qr` | Generate QR code |
| WiFi | `wifi basic` | Show WiFi basic settings |
| WiFi | `wifi set-channel <ch>` | Set channel |
| WiFi | `wifi set-power <20-100>` | Set transmit power |
| WiFi | `wifi enable` | Enable WiFi radio |
| WiFi | `wifi disable` | Disable WiFi radio |
| LAN | `lan dhcp` | Show DHCP settings |
| LAN | `lan dhcp-leases` | Show active leases |
| WAN | `wan list` | List WAN connections |
| WAN | `wan delete <name>` | Delete WAN connection |
| Port FW | `portfw list` | List port forwarding rules |
| Port FW | `portfw add [flags]` | Add port forwarding rule |
| Port FW | `portfw delete <name>` | Delete rule |
| Port FW | `portfw enable <name>` | Enable rule |
| Port FW | `portfw disable <name>` | Disable rule |
| Firewall | `firewall show` | Show firewall settings |
| Firewall | `firewall level <level>` | Set firewall level |
| Firewall | `firewall ipfilter list` | List IP filter rules |
| Firewall | `firewall macfilter list` | List MAC filter rules |
| Firewall | `firewall urlfilter list` | List URL filter rules |
| Firewall | `firewall urlfilter add <url>` | Block URL |
| DDNS | `ddns show` | Show DDNS settings |
| DMZ | `dmz show` | Show DMZ settings |
| DMZ | `dmz set <ip>` | Set DMZ host |
| UPnP | `upnp show` | Show UPnP settings |
| UPnP | `upnp enable` | Enable UPnP |
| UPnP | `upnp disable` | Disable UPnP |
| Diagnosis | `ping <host>` | Ping from router |
| Diagnosis | `arp` | Show ARP table |
| System | `reboot [--yes]` | Reboot router |
| System | `reset [--yes]` | Factory reset |
| System | `timeout [1-30]` | Session timeout |
| System | `log [limit]` | System logs |
| System | `password set` | Change password |
| Config | `config init` | Create config template |
| Config | `config path` | Show config path |

---

## Deployment on STB / Raspberry Pi

```bash
# Download ARM64 binary
wget https://github.com/CallMeJaja/zte-cli/releases/latest/download/zte-cli_linux_arm64.tar.gz
tar -xzf zte-cli_linux_arm64.tar.gz
chmod +x zte-cli
sudo mv zte-cli /usr/local/bin/

# Create config
zte-cli config init
nano ~/.config/zte-cli/config.yaml

# Test
zte-cli status
```

---

## How It Works

### Authentication Flow

1. Fetch login page → extract `Frm_Logintoken` and `Frm_Loginchecktoken`
2. Generate random 8-digit number (`UserRandomNum`)
3. Hash password: `SHA256(password + UserRandomNum)`
4. POST login credentials
5. Fetch `template.gch` → extract `session_token` for subsequent requests

### Wi-Fi Password Encryption

The router stores Wi-Fi passwords encrypted with AES-256-CBC (ZeroPadding). The key and IV are derived from the session token:

- **Key** = `SHA256(session_token)`
- **IV** = `SHA256(session_token[5:15])[:16]`

`zte-cli` automatically decrypts passwords when displaying and encrypts new passwords when changing settings.

### Router Pages

All router settings are accessed via `.gch` pages:

| Category | Pages |
|----------|-------|
| Status | `status_dev_info_t.gch`, `pon_status_link_info_t.gch`, `IPv46_status_wan2_if_t.gch`, `status_wlanm_info1_t.gch`, `pon_status_lan_info_t.gch`, `status_voip_4less_t.gch` |
| Network | `IPv46_net_wan2_conf_t.gch`, `net_wlanm_conf1_t.gch`, `net_wlanm_essid1_t.gch`, `net_wlanm_secrity1_t.gch`, `net_dhcp_dynamic_t.gch` |
| Security | `sec_firewall_conf_t.gch`, `sec_portfilter_conf_t.gch`, `sec_macfilter_conf_t.gch`, `sec_url_filter_t.gch` |
| Application | `app_virtual_conf_t.gch`, `app_ddns_conf_t.gch`, `app_dmz_conf_t.gch`, `app_upnp_conf_t.gch` |
| Administration | `manager_dev_conf_t.gch`, `manager_login_timeout_t.gch`, `manager_log_conf_t.gch`, `manager_aduser_conf_t.gch` |

---

## Project Structure

```
zte-cli/
├── main.go                    # Entry point, CLI routing
├── config/
│   └── config.go              # Config loading
├── router/
│   ├── client.go              # HTTP client, login flow
│   ├── pages.go               # Page constants
│   ├── parser.go              # HTML parsing
│   ├── handlers.go            # Page-specific parsers
│   ├── actions.go             # Core actions
│   ├── crypto.go              # AES encrypt/decrypt
│   ├── portforward.go         # Port forwarding CRUD
│   ├── dhcp.go                # DHCP management
│   ├── wifi.go                # WiFi settings
│   ├── firewall.go            # Firewall & filters
│   ├── applications.go        # DDNS, DMZ, UPnP
│   ├── wan.go                 # WAN management
│   ├── usermgmt.go            # Password management
│   ├── system.go              # Logs, diagnosis, reset
│   └── qr.go                  # QR code generation
├── cli/
│   ├── interactive.go         # TUI menu system
│   ├── menu.go                # Menu rendering
│   ├── commands.go            # Original CLI commands
│   ├── commands_new.go        # New CLI commands
│   ├── display.go             # Display formatting
│   ├── status.go              # Status module
│   └── network.go             # Network module
├── go.mod
├── go.sum
├── Dockerfile
├── .goreleaser.yaml
└── README.md
```

---

## Roadmap

- [x] v1.0.0 — Core router API + Interactive CLI + Direct commands
- [x] v2.0.0 — Full router management (35+ commands, 7 TUI menus)
- [ ] v2.1.0 — WAN connection create/edit
- [ ] v2.2.0 — Telegram Bot integration
- [ ] v2.3.0 — Internet heartbeat monitor & push notifications
- [ ] v2.4.0 — PON signal quality alerts
- [ ] v2.5.0 — Firmware upgrade via CLI

---

## License

[MIT](LICENSE)

## Acknowledgments

- [walangkaji/zte-f609-api](https://github.com/walangkaji/zte-f609-api) — PHP API library for ZTE F609
- [hack-gpon.org](https://hack-gpon.org) — GPON ONT reverse engineering community
