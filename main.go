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
		// If config not found, offer to create one
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
		// No arguments → interactive mode
		cli.Version = version
		cli.RunInteractive(client)
		return
	}

	// Parse non-interactive commands
	switch args[0] {
	case "status":
		cli.RunStatusCommand(client)
	case "clients":
		cli.RunClientsCommand(client)
	case "wifi":
		handleWiFi(client, args[1:])
	case "reboot":
		handleReboot(client, args[1:])
	case "timeout":
		handleTimeoutCmd(client, args[1:])
	case "bot":
		fmt.Println("  [!] Telegram Bot is not yet implemented (Post-MVP).")
		fmt.Println("      Stay tuned for v1.1.0!")
	default:
		fmt.Fprintf(os.Stderr, "  [ERROR] Unknown command: %s\n\n", args[0])
		printUsage()
		os.Exit(1)
	}
}

func handleWiFi(client *router.Client, args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "  Usage: zte-cli wifi <show|set|qr>")
		fmt.Fprintln(os.Stderr, "    zte-cli wifi show")
		fmt.Fprintln(os.Stderr, "    zte-cli wifi set --ssid <name> --pass <password>")
		fmt.Fprintln(os.Stderr, "    zte-cli wifi qr")
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
	default:
		fmt.Fprintf(os.Stderr, "  [ERROR] Unknown wifi subcommand: %s\n", args[0])
		os.Exit(1)
	}
}

func handleReboot(client *router.Client, args []string) {
	confirm := false
	for _, a := range args {
		if a == "--yes" || a == "-y" {
			confirm = true
		}
	}
	cli.RunRebootCommand(client, confirm)
}

func handleTimeoutCmd(client *router.Client, args []string) {
	if len(args) > 0 {
		cli.RunTimeoutCommand(client, args[0])
	} else {
		cli.RunTimeoutCommand(client, "")
	}
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
    zte-cli                      Interactive menu mode
    zte-cli <command> [options]  Direct command mode

  COMMANDS:
    status                       Show router status summary
    clients                      List connected devices (LAN & Wi-Fi)
    wifi show                    Show current SSID & password
    wifi set --ssid N --pass P   Change Wi-Fi SSID and/or password
    wifi qr                      Show Wi-Fi QR code
    reboot [--yes]               Reboot the router
    timeout [1-30]               Show or set session timeout (minutes)
    config init                  Create config file template
    config path                  Show active config file path
    bot                          Start Telegram bot daemon (Post-MVP)

  OPTIONS:
    --help, -h                   Show this help message
    --version, -v                Show version

  CONFIGURATION:
    Config is searched in this order:
    1. ./config.yaml
    2. ~/.config/zte-cli/config.yaml
    3. /etc/zte-cli/config.yaml

    Run 'zte-cli config init' to generate a template.

`)
}
