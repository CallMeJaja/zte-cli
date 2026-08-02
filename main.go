package main

import (
	"fmt"
	"os"

	"github.com/CallMeJaja/zte-cli/cli"
	"github.com/CallMeJaja/zte-cli/config"
	"github.com/CallMeJaja/zte-cli/router"
)

var version = "dev"

func main() {
	args := os.Args[1:]

	// Handle help and version before loading config
	if len(args) > 0 {
		switch args[0] {
		case "--help", "-h":
			printUsage()
			return
		case "--version", "-v":
			fmt.Printf("zte-cli v%s\n", version)
			return
		case "config":
			handleConfig(args[1:])
			return
		}
	}

	// Load config
	cfg, _, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "\n  [ERROR] %s\n\n", err)
		fmt.Fprintln(os.Stderr, "  Run 'zte-cli config init' to create a config file.")
		os.Exit(1)
	}

	// Create router client
	client, err := router.NewClient(cfg.RouterHost, cfg.RouterUsername, cfg.RouterPassword)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  [ERROR] Failed to create router client: %s\n", err)
		os.Exit(1)
	}

	// Route to the appropriate handler
	if len(args) == 0 {
		cli.Version = version
		cli.RunInteractive(client)
		return
	}

	// Parse non-interactive commands
	switch args[0] {
	// Status
	case "status":
		cli.RunStatusCommand(client)
	case "clients":
		cli.RunClientsCommand(client)

	// WiFi
	case "wifi":
		handleWiFi(client, args[1:])

	// LAN/DHCP
	case "lan":
		handleLAN(client, args[1:])

	// WAN Management
	case "wan":
		handleWAN(client, args[1:])

	// Port Forwarding
	case "portfw":
		handlePortFwd(client, args[1:])

	// Firewall & Filters
	case "firewall":
		handleFirewall(client, args[1:])

	// DDNS
	case "ddns":
		handleDDNS(client, args[1:])

	// DMZ
	case "dmz":
		handleDMZ(client, args[1:])

	// UPnP
	case "upnp":
		handleUPnP(client, args[1:])

	// Diagnosis
	case "ping":
		host := ""
		if len(args) > 1 {
			host = args[1]
		}
		cli.RunPingCommand(client, host)
	case "arp":
		cli.RunARPTableCommand(client)

	// System Management
	case "reboot":
		handleReboot(client, args[1:])
	case "reset":
		handleReset(client, args[1:])
	case "timeout":
		handleTimeoutCmd(client, args[1:])
	case "log":
		handleLog(client, args[1:])

	// User Management
	case "password":
		handlePassword(client, args[1:])

	case "bot":
		fmt.Println("  [!] Telegram Bot is not yet implemented (Post-MVP).")
		fmt.Println("      Stay tuned for v1.1.0!")

	default:
		fmt.Fprintf(os.Stderr, "  [ERROR] Unknown command: %s\n\n", args[0])
		printUsage()
		os.Exit(1)
	}
}

// ============================================================
// WiFi handler
// ============================================================

func handleWiFi(client *router.Client, args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "  Usage: zte-cli wifi <show|set|qr|basic|set-channel|set-power|enable|disable>")
		os.Exit(1)
	}

	switch args[0] {
	case "show":
		cli.RunWiFiShowCommand(client)
	case "set":
		ssid := getFlagValue(args[1:], "--ssid")
		pass := getFlagValue(args[1:], "--pass")
		cli.RunWiFiSetCommand(client, ssid, pass)
	case "qr":
		cli.RunWiFiQRCommand(client)
	case "basic":
		cli.RunWiFiBasicCommand(client)
	case "set-channel":
		ch := ""
		if len(args) > 1 {
			ch = args[1]
		}
		cli.RunWiFiSetChannelCommand(client, ch)
	case "set-power":
		power := ""
		if len(args) > 1 {
			power = args[1]
		}
		cli.RunWiFiSetPowerCommand(client, power)
	case "enable":
		cli.RunWiFiToggleCommand(client, true)
	case "disable":
		cli.RunWiFiToggleCommand(client, false)
	default:
		fmt.Fprintf(os.Stderr, "  [ERROR] Unknown wifi subcommand: %s\n", args[0])
		os.Exit(1)
	}
}

// ============================================================
// LAN/DHCP handler
// ============================================================

func handleLAN(client *router.Client, args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "  Usage: zte-cli lan <dhcp|dhcp-leases>")
		os.Exit(1)
	}

	switch args[0] {
	case "dhcp":
		cli.RunDHCPShowCommand(client)
	case "dhcp-leases":
		cli.RunDHCPLeasesCommand(client)
	default:
		fmt.Fprintf(os.Stderr, "  [ERROR] Unknown lan subcommand: %s\n", args[0])
		os.Exit(1)
	}
}

// ============================================================
// WAN Management handler
// ============================================================

func handleWAN(client *router.Client, args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "  Usage: zte-cli wan <list|delete>")
		os.Exit(1)
	}

	switch args[0] {
	case "list":
		cli.RunWANListCommand(client)
	case "delete":
		name := ""
		if len(args) > 1 {
			name = args[1]
		}
		cli.RunWANDeleteCommand(client, name)
	default:
		fmt.Fprintf(os.Stderr, "  [ERROR] Unknown wan subcommand: %s\n", args[0])
		os.Exit(1)
	}
}

// ============================================================
// Port Forwarding handler
// ============================================================

func handlePortFwd(client *router.Client, args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "  Usage: zte-cli portfw <list|add|delete|enable|disable>")
		os.Exit(1)
	}

	switch args[0] {
	case "list":
		cli.RunPortFwdListCommand(client)
	case "add":
		name := getFlagValue(args[1:], "--name")
		proto := getFlagValue(args[1:], "--proto")
		wanPort := getFlagValue(args[1:], "--wan-port")
		lanIP := getFlagValue(args[1:], "--lan-ip")
		lanPort := getFlagValue(args[1:], "--lan-port")
		cli.RunPortFwdAddCommand(client, name, proto, wanPort, lanIP, lanPort)
	case "delete":
		name := ""
		if len(args) > 1 {
			name = args[1]
		}
		cli.RunPortFwdDeleteCommand(client, name)
	case "enable":
		name := ""
		if len(args) > 1 {
			name = args[1]
		}
		cli.RunPortFwdEnableCommand(client, name, true)
	case "disable":
		name := ""
		if len(args) > 1 {
			name = args[1]
		}
		cli.RunPortFwdEnableCommand(client, name, false)
	default:
		fmt.Fprintf(os.Stderr, "  [ERROR] Unknown portfw subcommand: %s\n", args[0])
		os.Exit(1)
	}
}

// ============================================================
// Firewall handler
// ============================================================

func handleFirewall(client *router.Client, args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "  Usage: zte-cli firewall <show|level|ipfilter|macfilter|urlfilter>")
		os.Exit(1)
	}

	switch args[0] {
	case "show":
		cli.RunFirewallShowCommand(client)
	case "level":
		level := ""
		if len(args) > 1 {
			level = args[1]
		}
		cli.RunFirewallLevelCommand(client, level)
	case "ipfilter":
		handleIPFilter(client, args[1:])
	case "macfilter":
		handleMACFilter(client, args[1:])
	case "urlfilter":
		handleURLFilter(client, args[1:])
	default:
		fmt.Fprintf(os.Stderr, "  [ERROR] Unknown firewall subcommand: %s\n", args[0])
		os.Exit(1)
	}
}

func handleIPFilter(client *router.Client, args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "  Usage: zte-cli firewall ipfilter <list>")
		os.Exit(1)
	}

	switch args[0] {
	case "list":
		cli.RunIPFilterListCommand(client)
	default:
		fmt.Fprintf(os.Stderr, "  [ERROR] Unknown ipfilter subcommand: %s\n", args[0])
		os.Exit(1)
	}
}

func handleMACFilter(client *router.Client, args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "  Usage: zte-cli firewall macfilter <list>")
		os.Exit(1)
	}

	switch args[0] {
	case "list":
		cli.RunMACFilterListCommand(client)
	default:
		fmt.Fprintf(os.Stderr, "  [ERROR] Unknown macfilter subcommand: %s\n", args[0])
		os.Exit(1)
	}
}

func handleURLFilter(client *router.Client, args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "  Usage: zte-cli firewall urlfilter <list|add>")
		os.Exit(1)
	}

	switch args[0] {
	case "list":
		cli.RunURLFilterListCommand(client)
	case "add":
		url := ""
		if len(args) > 1 {
			url = args[1]
		}
		cli.RunURLFilterAddCommand(client, url)
	default:
		fmt.Fprintf(os.Stderr, "  [ERROR] Unknown urlfilter subcommand: %s\n", args[0])
		os.Exit(1)
	}
}

// ============================================================
// DDNS handler
// ============================================================

func handleDDNS(client *router.Client, args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "  Usage: zte-cli ddns <show>")
		os.Exit(1)
	}

	switch args[0] {
	case "show":
		cli.RunDDNSShowCommand(client)
	default:
		fmt.Fprintf(os.Stderr, "  [ERROR] Unknown ddns subcommand: %s\n", args[0])
		os.Exit(1)
	}
}

// ============================================================
// DMZ handler
// ============================================================

func handleDMZ(client *router.Client, args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "  Usage: zte-cli dmz <show|set>")
		os.Exit(1)
	}

	switch args[0] {
	case "show":
		cli.RunDMZShowCommand(client)
	case "set":
		ip := ""
		if len(args) > 1 {
			ip = args[1]
		}
		cli.RunDMZSetCommand(client, ip)
	default:
		fmt.Fprintf(os.Stderr, "  [ERROR] Unknown dmz subcommand: %s\n", args[0])
		os.Exit(1)
	}
}

// ============================================================
// UPnP handler
// ============================================================

func handleUPnP(client *router.Client, args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "  Usage: zte-cli upnp <show|enable|disable>")
		os.Exit(1)
	}

	switch args[0] {
	case "show":
		cli.RunUPnPShowCommand(client)
	case "enable":
		cli.RunUPnPToggleCommand(client, true)
	case "disable":
		cli.RunUPnPToggleCommand(client, false)
	default:
		fmt.Fprintf(os.Stderr, "  [ERROR] Unknown upnp subcommand: %s\n", args[0])
		os.Exit(1)
	}
}

// ============================================================
// System handlers
// ============================================================

func handleReboot(client *router.Client, args []string) {
	confirm := false
	for _, a := range args {
		if a == "--yes" || a == "-y" {
			confirm = true
		}
	}
	cli.RunRebootCommand(client, confirm)
}

func handleReset(client *router.Client, args []string) {
	confirm := false
	for _, a := range args {
		if a == "--yes" || a == "-y" {
			confirm = true
		}
	}
	cli.RunResetCommand(client, confirm)
}

func handleTimeoutCmd(client *router.Client, args []string) {
	if len(args) > 0 {
		cli.RunTimeoutCommand(client, args[0])
	} else {
		cli.RunTimeoutCommand(client, "")
	}
}

func handleLog(client *router.Client, args []string) {
	if len(args) > 0 {
		cli.RunLogShowCommand(client, args[0])
	} else {
		cli.RunLogShowCommand(client, "")
	}
}

func handlePassword(client *router.Client, args []string) {
	if len(args) == 0 || args[0] != "set" {
		fmt.Fprintln(os.Stderr, "  Usage: zte-cli password set --old <old> --new <new> --confirm <new>")
		os.Exit(1)
	}

	oldPass := getFlagValue(args[1:], "--old")
	newPass := getFlagValue(args[1:], "--new")
	confirmPass := getFlagValue(args[1:], "--confirm")
	cli.RunChangePasswordCommand(client, oldPass, newPass, confirmPass)
}

func handleConfig(args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "  Usage: zte-cli config <init|path>")
		os.Exit(1)
	}

	switch args[0] {
	case "init":
		path, err := config.Init()
		if err != nil {
			fmt.Fprintf(os.Stderr, "  [ERROR] %s\n", err)
			os.Exit(1)
		}
		fmt.Printf("  [OK] Config file created at: %s\n", path)
	case "path":
		fmt.Printf("  %s\n", config.Path())
	default:
		fmt.Fprintf(os.Stderr, "  [ERROR] Unknown config subcommand: %s\n", args[0])
		os.Exit(1)
	}
}

func getFlagValue(args []string, flag string) string {
	for i, a := range args {
		if a == flag && i+1 < len(args) {
			return args[i+1]
		}
	}
	return ""
}

func printUsage() {
	fmt.Print(`
  zte-cli - ZTE F609 Router Manager

  USAGE:
    zte-cli                          Interactive menu mode
    zte-cli <command> [options]      Direct command mode

  STATUS:
    status                           Show full status summary
    clients                          List connected devices (LAN & Wi-Fi)

  WIFI:
    wifi show                        Show current SSID & password
    wifi set --ssid N --pass P       Change Wi-Fi SSID and/or password
    wifi qr                          Show Wi-Fi QR code
    wifi basic                       Show WiFi basic settings (channel, power, mode)
    wifi set-channel <ch>            Set WiFi channel (Auto, 1-13)
    wifi set-power <20-100>          Set transmit power (20/40/60/80/100)
    wifi enable                      Enable WiFi radio
    wifi disable                     Disable WiFi radio

  LAN / DHCP:
    lan dhcp                         Show DHCP server settings
    lan dhcp-leases                  Show active DHCP leases

  WAN:
    wan list                         List WAN connections
    wan delete <name>                Delete WAN connection

  PORT FORWARDING:
    portfw list                      List all port forwarding rules
    portfw add [flags]               Add new rule (--name --proto --wan-port --lan-ip --lan-port)
    portfw delete <name>             Delete rule by name
    portfw enable <name>             Enable rule
    portfw disable <name>            Disable rule

  FIREWALL & FILTERS:
    firewall show                    Show firewall settings
    firewall level <off|low|mid|high>  Set firewall level
    firewall ipfilter list           List IP filter rules
    firewall macfilter list          List MAC filter rules
    firewall urlfilter list          List URL filter rules
    firewall urlfilter add <url>     Add URL to block list

  APPLICATIONS:
    ddns show                        Show DDNS settings
    dmz show                         Show DMZ settings
    dmz set <ip>                     Set DMZ host IP
    upnp show                        Show UPnP settings
    upnp enable                      Enable UPnP
    upnp disable                     Disable UPnP

  DIAGNOSIS:
    ping <host>                      Ping from router
    arp                              Show ARP table

  SYSTEM:
    reboot [--yes]                   Reboot the router
    reset [--yes]                    Factory reset (WARNING!)
    timeout [1-30]                   Show/set session timeout (minutes)
    log [limit]                      Show system logs (default: last 50)
    password set --old O --new N --confirm N  Change router password
    config init                      Create config file template
    config path                      Show active config file path

  OPTIONS:
    --help, -h                       Show this help message
    --version, -v                    Show version

  CONFIGURATION:
    Config is searched in this order:
    1. ./config.yaml
    2. ~/.config/zte-cli/config.yaml
    3. /etc/zte-cli/config.yaml

    Run 'zte-cli config init' to generate a template.

`)
}
