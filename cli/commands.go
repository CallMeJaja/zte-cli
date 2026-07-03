package cli

import (
	"fmt"
	"os"
	"strconv"

	"github.com/CallMeJaja/zte-cli/router"
)

// RunStatusCommand handles: zte-cli status
func RunStatusCommand(client *router.Client) {
	ok, err := client.Login()
	if err != nil || !ok {
		fmt.Fprintf(os.Stderr, "  [ERROR] Login failed: %v\n", err)
		os.Exit(1)
	}
	handleStatus(client)
}

// RunClientsCommand handles: zte-cli clients
func RunClientsCommand(client *router.Client) {
	ok, err := client.Login()
	if err != nil || !ok {
		fmt.Fprintf(os.Stderr, "  [ERROR] Login failed: %v\n", err)
		os.Exit(1)
	}
	handleClients(client)
}

// RunWiFiShowCommand handles: zte-cli wifi show
func RunWiFiShowCommand(client *router.Client) {
	ok, err := client.Login()
	if err != nil || !ok {
		fmt.Fprintf(os.Stderr, "  [ERROR] Login failed: %v\n", err)
		os.Exit(1)
	}
	showWiFiDetails(client)
}

// RunWiFiSetCommand handles: zte-cli wifi set --ssid <name> --pass <password>
func RunWiFiSetCommand(client *router.Client, ssid string, password string) {
	ok, err := client.Login()
	if err != nil || !ok {
		fmt.Fprintf(os.Stderr, "  [ERROR] Login failed: %v\n", err)
		os.Exit(1)
	}

	if ssid != "" {
		fmt.Printf("  Setting SSID to '%s'...\n", ssid)
		ok, err := router.SetWiFiSSID(client, ssid)
		if err != nil {
			fmt.Fprintf(os.Stderr, "  [ERROR] %s\n", err)
			os.Exit(1)
		}
		if ok {
			fmt.Println("  [OK] SSID updated.")
		}
	}

	if password != "" {
		if len(password) < 8 {
			fmt.Fprintf(os.Stderr, "  [ERROR] Password must be at least 8 characters.\n")
			os.Exit(1)
		}
		fmt.Printf("  Setting password to '%s'...\n", password)
		ok, err := router.SetWiFiPassword(client, password)
		if err != nil {
			fmt.Fprintf(os.Stderr, "  [ERROR] %s\n", err)
			os.Exit(1)
		}
		if ok {
			fmt.Println("  [OK] Password updated.")
		}
	}

	if ssid == "" && password == "" {
		fmt.Fprintln(os.Stderr, "  [ERROR] Please specify --ssid and/or --pass.")
		os.Exit(1)
	}
}

// RunWiFiQRCommand handles: zte-cli wifi qr
func RunWiFiQRCommand(client *router.Client) {
	ok, err := client.Login()
	if err != nil || !ok {
		fmt.Fprintf(os.Stderr, "  [ERROR] Login failed: %v\n", err)
		os.Exit(1)
	}
	showWiFiQR(client)
}

// RunRebootCommand handles: zte-cli reboot [--yes]
func RunRebootCommand(client *router.Client, confirm bool) {
	ok, err := client.Login()
	if err != nil || !ok {
		fmt.Fprintf(os.Stderr, "  [ERROR] Login failed: %v\n", err)
		os.Exit(1)
	}

	if !confirm {
		fmt.Print("  Are you sure you want to reboot the router? [y/N]: ")
		var answer string
		fmt.Scanln(&answer)
		if answer != "y" && answer != "yes" {
			fmt.Println("  Reboot cancelled.")
			return
		}
	}

	fmt.Println("  Sending reboot command...")
	ok, err = router.Reboot(client)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  [ERROR] %s\n", err)
		os.Exit(1)
	}
	if ok {
		fmt.Println("  [OK] Router is rebooting...")
	} else {
		fmt.Println("  [!] Reboot command may have failed.")
	}
}

// RunTimeoutCommand handles: zte-cli timeout [1-30]
func RunTimeoutCommand(client *router.Client, value string) {
	ok, err := client.Login()
	if err != nil || !ok {
		fmt.Fprintf(os.Stderr, "  [ERROR] Login failed: %v\n", err)
		os.Exit(1)
	}

	if value == "" {
		// Just show current timeout
		timeout, err := router.FetchTimeout(client)
		if err != nil {
			fmt.Fprintf(os.Stderr, "  [ERROR] %s\n", err)
			os.Exit(1)
		}
		fmt.Printf("  Current session timeout: %d minutes\n", timeout)
		return
	}

	newTimeout, err := strconv.Atoi(value)
	if err != nil || newTimeout < 1 || newTimeout > 30 {
		fmt.Fprintf(os.Stderr, "  [ERROR] Timeout must be between 1 and 30 minutes.\n")
		os.Exit(1)
	}

	fmt.Printf("  Setting timeout to %d minutes...\n", newTimeout)
	ok, err = router.SetTimeout(client, newTimeout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  [ERROR] %s\n", err)
		os.Exit(1)
	}
	if ok {
		fmt.Printf("  [OK] Session timeout set to %d minutes.\n", newTimeout)
	} else {
		fmt.Println("  [!] Failed to set timeout.")
	}
}
