# zte-cli

An open-source interactive CLI tool for managing **ZTE F609** GPON routers. Built in Go, it produces a single lightweight binary (~6MB) with zero runtime dependencies.

## Features

- **Interactive TUI Menu** — navigate router settings with simple number inputs
- **Direct CLI Commands** — for scripting and cron job automation
- **Router Status** — device info, PON optical power, WAN IP, connection status
- **Connected Devices** — list LAN ports and Wi-Fi clients
- **Wi-Fi Management** — show/change SSID & password, generate QR code
- **Remote Reboot** — reboot router from terminal
- **Session Timeout** — view and change router login timeout
- **AES-CBC Decryption** — automatically decrypts Wi-Fi password from router's encrypted storage

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

## Configuration

Create a config file at `~/.config/zte-cli/config.yaml`:

```bash
zte-cli config init
```

Then edit the file:

```yaml
router_host: "192.168.100.1"
router_username: "admin"
router_password: "admin"
```

### Config Search Priority

1. `./config.yaml` (current directory)
2. `~/.config/zte-cli/config.yaml` (user home, recommended)
3. `/etc/zte-cli/config.yaml` (system-wide)

## Usage

### Interactive Mode

```bash
zte-cli
```

```
 ==================================================
          ZTE F609 Router Manager v0.1.0
 ==================================================
 1. Show Router Status
 2. Connected Devices (LAN & Wi-Fi)
 3. Wi-Fi Settings
 4. Reboot Router
 5. Session Timeout Settings
 0. Exit
 ==================================================
 Choose option [0-5]: _
```

### Direct Commands

```bash
# Show router status summary
zte-cli status

# List connected devices
zte-cli clients

# Show current Wi-Fi SSID & password
zte-cli wifi show

# Change Wi-Fi password
zte-cli wifi set --pass newpassword123

# Change Wi-Fi SSID and password
zte-cli wifi set --ssid MyNewWiFi --pass newpassword123

# Show Wi-Fi QR code
zte-cli wifi qr

# Reboot router
zte-cli reboot --yes

# Show current session timeout
zte-cli timeout

# Set session timeout to 30 minutes
zte-cli timeout 30

# Show config file location
zte-cli config path
```

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

# (Optional) Run as system service
sudo tee /etc/systemd/system/zte-cli.service > /dev/null <<EOF
[Unit]
Description=ZTE F609 CLI Manager
After=network.target

[Service]
Type=simple
ExecStart=/usr/local/bin/zte-cli bot
Restart=on-failure
RestartSec=10

[Install]
WantedBy=multi-user.target
EOF

sudo systemctl enable zte-cli
sudo systemctl start zte-cli
```

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

## Roadmap

- [x] MVP v1.0.0 — Core router API + Interactive CLI + Direct commands
- [ ] v1.1.0 — Telegram Bot command handler
- [ ] v1.2.0 — Internet heartbeat monitor & push notifications
- [ ] v1.3.0 — PON signal quality alerts

## License

[MIT](LICENSE)

## Acknowledgments

- [walangkaji/zte-f609-api](https://github.com/walangkaji/zte-f609-api) — PHP API library for ZTE F609
- [hack-gpon.org](https://hack-gpon.org) — GPON ONT reverse engineering community
