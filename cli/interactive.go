package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/CallMeJaja/zte-cli/router"
)

// Version is set by main before launching interactive mode.
var Version = "dev"

var banner string

func init() {
	banner = fmt.Sprintf(`
==========================================================
                 ZTE F609 Router Manager v%s
==========================================================
(1) Status                         (5) Administration
(2) Network                        (6) Diagnosis
(3) Security                       (7) System
(4) Application                    (0) Exit
==========================================================
`, Version)
}

// RunInteractive launches the interactive TUI menu loop.
func RunInteractive(client *router.Client) {
	scanner := bufio.NewScanner(os.Stdin)

	// Login first
	fmt.Println("\n  Connecting to router...")
	ok, err := client.Login()
	if err != nil {
		fmt.Printf("\n  [ERROR] %s\n", err)
		return
	}
	if !ok {
		fmt.Println("\n  [ERROR] Login failed.")
		return
	}
	fmt.Println("  [OK] Logged in successfully.")

	for {
		fmt.Print(banner)
		fmt.Print("  Enter command number: ")

		if !scanner.Scan() {
			break
		}
		choice := strings.TrimSpace(scanner.Text())
		if choice == "" {
			continue
		}

		switch choice {
		case "1":
			runStatusMenu(client, scanner)
		case "2":
			runNetworkMenu(client, scanner)
		case "3":
			runSecurityMenu(client, scanner)
		case "4":
			runApplicationMenu(client, scanner)
		case "5":
			runAdminMenu(client, scanner)
		case "6":
			runDiagnosisMenu(client, scanner)
		case "7":
			runSystemMenu(client, scanner)
		case "0":
			fmt.Println("\n  Goodbye!")
			return
		default:
			fmt.Println("\n  [!] Invalid option. Please try again.")
		}
	}
}

// ============================================================
// Status Menu
// ============================================================

func runStatusMenu(client *router.Client, scanner *bufio.Scanner) {
	for {
		fmt.Println(`
==========================================================
  STATUS
==========================================================
  (1) Device Information
  (2) WAN Status
  (3) PON Optical Status
  (4) WLAN Status
  (5) LAN Status
  (6) VoIP Status
  (7) Connected Clients
  (0) Back
==========================================================`)
		fmt.Print("  Enter choice: ")
		if !scanner.Scan() {
			return
		}
		choice := strings.TrimSpace(scanner.Text())

		switch choice {
		case "1":
			displayStatus(client)
		case "2":
			status, err := router.FetchStatus(client)
			if err != nil {
				fmt.Printf("  [ERROR] %s\n", err)
			} else {
				fmt.Printf("\n  WAN: %s | IP: %s | Status: %s\n\n",
					status.WAN.ConnectionName, status.WAN.IPAddress, status.WAN.ConnectionStatus)
			}
		case "3":
			status, err := router.FetchStatus(client)
			if err != nil {
				fmt.Printf("  [ERROR] %s\n", err)
			} else {
				fmt.Printf("\n  GPON: %s | Rx: %s | Tx: %s\n\n",
					status.PON.GPONState, status.PON.RxPower, status.PON.TxPower)
			}
		case "4":
			displayClients(client)
		case "5":
			lans, err := router.FetchLANStatus(client)
			if err != nil {
				fmt.Printf("  [ERROR] %s\n", err)
			} else {
				for _, l := range lans {
					fmt.Printf("  %s: %s | Speed: %s\n", l.Port, l.Status, l.Speed)
				}
				fmt.Println()
			}
		case "6":
			html, err := client.GetPage("status_voip_4less_t.gch")
			if err != nil {
				fmt.Printf("  [ERROR] %s\n", err)
			} else {
				lines := router.ParseGenericTable(html)
				for _, line := range lines {
					fmt.Println("  " + line)
				}
			}
		case "7":
			displayClients(client)
		case "0":
			return
		default:
			fmt.Println("  [!] Invalid option.")
		}
	}
}

// ============================================================
// Network Menu
// ============================================================

func runNetworkMenu(client *router.Client, scanner *bufio.Scanner) {
	for {
		fmt.Println(`
==========================================================
  NETWORK
==========================================================
  (1) WAN Connections
  (2) WiFi Settings
  (3) WiFi Security (SSID/Password)
  (4) DHCP Settings
  (5) DHCP Leases
  (6) Port Forwarding
  (0) Back
==========================================================`)
		fmt.Print("  Enter choice: ")
		if !scanner.Scan() {
			return
		}
		choice := strings.TrimSpace(scanner.Text())

		switch choice {
		case "1":
			conns, err := router.FetchWANConnections(client)
			if err != nil {
				fmt.Printf("  [ERROR] %s\n", err)
			} else {
				fmt.Println("\n  WAN Connections:")
				fmt.Println(router.FormatWANConnections(conns))
			}
		case "2":
			settings, err := router.FetchWiFiBasicSettings(client)
			if err != nil {
				fmt.Printf("  [ERROR] %s\n", err)
			} else {
				fmt.Println("\n  WiFi Basic Settings:")
				fmt.Println(router.FormatWiFiBasicSettings(settings))
			}
		case "3":
			displayWiFiDetails(client)
		case "4":
			settings, err := router.FetchDHCPSettings(client)
			if err != nil {
				fmt.Printf("  [ERROR] %s\n", err)
			} else {
				fmt.Println("\n  DHCP Settings:")
				fmt.Println(router.FormatDHCPSettings(settings))
			}
		case "5":
			leases, err := router.FetchDHCPLeases(client)
			if err != nil {
				fmt.Printf("  [ERROR] %s\n", err)
			} else {
				fmt.Println("\n  DHCP Leases:")
				fmt.Println(router.FormatDHCPLeases(leases))
			}
		case "6":
			rules, err := router.FetchPortForwardRules(client)
			if err != nil {
				fmt.Printf("  [ERROR] %s\n", err)
			} else {
				fmt.Println("\n  Port Forwarding:")
				fmt.Println(router.FormatPortForwardRules(rules))
			}
		case "0":
			return
		default:
			fmt.Println("  [!] Invalid option.")
		}
	}
}

// ============================================================
// Security Menu
// ============================================================

func runSecurityMenu(client *router.Client, scanner *bufio.Scanner) {
	for {
		fmt.Println(`
==========================================================
  SECURITY
==========================================================
  (1) Firewall Settings
  (2) IP Filter Rules
  (3) MAC Filter Rules
  (4) URL Filter Rules
  (0) Back
==========================================================`)
		fmt.Print("  Enter choice: ")
		if !scanner.Scan() {
			return
		}
		choice := strings.TrimSpace(scanner.Text())

		switch choice {
		case "1":
			settings, err := router.FetchFirewallSettings(client)
			if err != nil {
				fmt.Printf("  [ERROR] %s\n", err)
			} else {
				fmt.Println("\n  Firewall Settings:")
				fmt.Println(router.FormatFirewallSettings(settings))
			}
		case "2":
			rules, err := router.FetchIPFilterRules(client)
			if err != nil {
				fmt.Printf("  [ERROR] %s\n", err)
			} else {
				fmt.Println("\n  IP Filter Rules:")
				fmt.Println(router.FormatIPFilterRules(rules))
			}
		case "3":
			rules, err := router.FetchMACFilterRules(client)
			if err != nil {
				fmt.Printf("  [ERROR] %s\n", err)
			} else {
				fmt.Println("\n  MAC Filter Rules:")
				fmt.Println(router.FormatMACFilterRules(rules))
			}
		case "4":
			rules, err := router.FetchURLFilterRules(client)
			if err != nil {
				fmt.Printf("  [ERROR] %s\n", err)
			} else {
				fmt.Println("\n  URL Filter Rules:")
				fmt.Println(router.FormatURLFilterRules(rules))
			}
		case "0":
			return
		default:
			fmt.Println("  [!] Invalid option.")
		}
	}
}

// ============================================================
// Application Menu
// ============================================================

func runApplicationMenu(client *router.Client, scanner *bufio.Scanner) {
	for {
		fmt.Println(`
==========================================================
  APPLICATION
==========================================================
  (1) DDNS Settings
  (2) DMZ Settings
  (3) UPnP Settings
  (4) WiFi QR Code
  (0) Back
==========================================================`)
		fmt.Print("  Enter choice: ")
		if !scanner.Scan() {
			return
		}
		choice := strings.TrimSpace(scanner.Text())

		switch choice {
		case "1":
			settings, err := router.FetchDDNSSettings(client)
			if err != nil {
				fmt.Printf("  [ERROR] %s\n", err)
			} else {
				fmt.Println("\n  DDNS Settings:")
				fmt.Println(router.FormatDDNSSettings(settings))
			}
		case "2":
			settings, err := router.FetchDMZSettings(client)
			if err != nil {
				fmt.Printf("  [ERROR] %s\n", err)
			} else {
				fmt.Println("\n  DMZ Settings:")
				fmt.Println(router.FormatDMZSettings(settings))
			}
		case "3":
			settings, err := router.FetchUPnPSettings(client)
			if err != nil {
				fmt.Printf("  [ERROR] %s\n", err)
			} else {
				fmt.Println("\n  UPnP Settings:")
				fmt.Println(router.FormatUPnPSettings(settings))
			}
		case "4":
			displayWiFiQR(client)
		case "0":
			return
		default:
			fmt.Println("  [!] Invalid option.")
		}
	}
}

// ============================================================
// Administration Menu
// ============================================================

func runAdminMenu(client *router.Client, scanner *bufio.Scanner) {
	for {
		fmt.Println(`
==========================================================
  ADMINISTRATION
==========================================================
  (1) Change Password
  (2) Session Timeout
  (0) Back
==========================================================`)
		fmt.Print("  Enter choice: ")
		if !scanner.Scan() {
			return
		}
		choice := strings.TrimSpace(scanner.Text())

		switch choice {
		case "1":
			fmt.Print("  Old Password: ")
			scanner.Scan()
			oldPass := strings.TrimSpace(scanner.Text())
			fmt.Print("  New Password: ")
			scanner.Scan()
			newPass := strings.TrimSpace(scanner.Text())
			fmt.Print("  Confirm Password: ")
			scanner.Scan()
			confirmPass := strings.TrimSpace(scanner.Text())

			ok, err := router.ChangePassword(client, oldPass, newPass, confirmPass)
			if err != nil {
				fmt.Printf("  [ERROR] %s\n", err)
			} else if ok {
				fmt.Println("  [OK] Password changed successfully.")
			}
		case "2":
			timeout, err := router.FetchTimeout(client)
			if err != nil {
				fmt.Printf("  [ERROR] %s\n", err)
			} else {
				fmt.Printf("  Current session timeout: %d minutes\n", timeout)
			}
		case "0":
			return
		default:
			fmt.Println("  [!] Invalid option.")
		}
	}
}

// ============================================================
// Diagnosis Menu
// ============================================================

func runDiagnosisMenu(client *router.Client, scanner *bufio.Scanner) {
	for {
		fmt.Println(`
==========================================================
  DIAGNOSIS
==========================================================
  (1) Ping
  (2) ARP Table
  (0) Back
==========================================================`)
		fmt.Print("  Enter choice: ")
		if !scanner.Scan() {
			return
		}
		choice := strings.TrimSpace(scanner.Text())

		switch choice {
		case "1":
			fmt.Print("  Host: ")
			scanner.Scan()
			host := strings.TrimSpace(scanner.Text())
			if host == "" {
				fmt.Println("  [!] Host required.")
				continue
			}
			fmt.Printf("  Pinging %s...\n", host)
			result, err := router.Ping(client, host)
			if err != nil {
				fmt.Printf("  [ERROR] %s\n", err)
			} else {
				fmt.Println(result)
			}
		case "2":
			result, err := router.FetchARPTable(client)
			if err != nil {
				fmt.Printf("  [ERROR] %s\n", err)
			} else {
				fmt.Println("\n  ARP Table:")
				fmt.Println(result)
			}
		case "0":
			return
		default:
			fmt.Println("  [!] Invalid option.")
		}
	}
}

// ============================================================
// System Menu
// ============================================================

func runSystemMenu(client *router.Client, scanner *bufio.Scanner) {
	for {
		fmt.Println(`
==========================================================
  SYSTEM
==========================================================
  (1) Reboot Router
  (2) Factory Reset
  (3) System Logs
  (0) Back
==========================================================`)
		fmt.Print("  Enter choice: ")
		if !scanner.Scan() {
			return
		}
		choice := strings.TrimSpace(scanner.Text())

		switch choice {
		case "1":
			fmt.Print("  Are you sure? [y/N]: ")
			scanner.Scan()
			answer := strings.TrimSpace(strings.ToLower(scanner.Text()))
			if answer == "y" || answer == "yes" {
				ok, err := router.Reboot(client)
				if err != nil {
					fmt.Printf("  [ERROR] %s\n", err)
				} else if ok {
					fmt.Println("  [OK] Router is rebooting...")
				}
			} else {
				fmt.Println("  Reboot cancelled.")
			}
		case "2":
			fmt.Print("  ⚠️  WARNING: This will reset ALL settings! [y/N]: ")
			scanner.Scan()
			answer := strings.TrimSpace(strings.ToLower(scanner.Text()))
			if answer == "y" || answer == "yes" {
				ok, err := router.FactoryReset(client)
				if err != nil {
					fmt.Printf("  [ERROR] %s\n", err)
				} else if ok {
					fmt.Println("  [OK] Router is resetting...")
				}
			} else {
				fmt.Println("  Reset cancelled.")
			}
		case "3":
			logs, err := router.FetchLogs(client)
			if err != nil {
				fmt.Printf("  [ERROR] %s\n", err)
			} else {
				fmt.Println("\n  System Logs:")
				fmt.Println(router.FormatLogs(logs, 30))
			}
		case "0":
			return
		default:
			fmt.Println("  [!] Invalid option.")
		}
	}
}
