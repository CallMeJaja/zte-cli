package cli

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/CallMeJaja/zte-cli/router"
)

const banner = `
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
`

const wifiMenu = `
 --------------------------------------------------
                 Wi-Fi Settings
 --------------------------------------------------
 1. Show Current SSID & Password
 2. Change Wi-Fi Password
 3. Show Wi-Fi QR Code
 0. Back
 --------------------------------------------------
`

const timeoutMenu = `
 --------------------------------------------------
             Session Timeout
 --------------------------------------------------
 Current timeout is displayed above.
 Enter new timeout (1-30 minutes), or 0 to cancel:
 --------------------------------------------------
`

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
	fmt.Println("  [OK] Logged in successfully.\n")

	for {
		fmt.Print(banner)
		fmt.Print("  Choose option [0-5]: ")

		if !scanner.Scan() {
			break
		}
		choice := strings.TrimSpace(scanner.Text())

		switch choice {
		case "1":
			handleStatus(client)
		case "2":
			handleClients(client)
		case "3":
			handleWiFiMenu(client, scanner)
		case "4":
			handleReboot(client, scanner)
		case "5":
			handleTimeout(client, scanner)
		case "0":
			fmt.Println("\n  Goodbye!")
			return
		default:
			fmt.Println("\n  [!] Invalid option. Please try again.")
		}
	}
}

func handleStatus(client *router.Client) {
	fmt.Println("\n  Fetching router status...")

	status, err := router.FetchStatus(client)
	if err != nil {
		fmt.Printf("  [ERROR] %s\n", err)
		return
	}

	fmt.Println()
	fmt.Println("  ┌─────────────────────────────────────────────────┐")
	fmt.Println("  │             DEVICE INFORMATION                  │")
	fmt.Println("  ├─────────────────────────────────────────────────┤")
	fmt.Printf("  │  %-18s │ %-28s │\n", "Model", status.Device.Model)
	fmt.Printf("  │  %-18s │ %-28s │\n", "Serial Number", status.Device.SerialNumber)
	fmt.Printf("  │  %-18s │ %-28s │\n", "HW Version", status.Device.HardwareVersion)
	fmt.Printf("  │  %-18s │ %-28s │\n", "SW Version", status.Device.SoftwareVersion)
	fmt.Printf("  │  %-18s │ %-28s │\n", "PON Serial", status.Device.PONSerialNumber)
	fmt.Printf("  │  %-18s │ %-28s │\n", "Batch Number", status.Device.BatchNumber)
	fmt.Println("  ├─────────────────────────────────────────────────┤")
	fmt.Println("  │             PON OPTICAL STATUS                  │")
	fmt.Println("  ├─────────────────────────────────────────────────┤")
	fmt.Printf("  │  %-18s │ %-28s │\n", "GPON State", status.PON.GPONState)
	fmt.Printf("  │  %-18s │ %-28s │\n", "Rx Power", status.PON.RxPower)
	fmt.Printf("  │  %-18s │ %-28s │\n", "Tx Power", status.PON.TxPower)
	fmt.Printf("  │  %-18s │ %-28s │\n", "Voltage", status.PON.Voltage+" uV")
	fmt.Printf("  │  %-18s │ %-28s │\n", "Bias Current", status.PON.BiasCurrent+" uA")
	fmt.Printf("  │  %-18s │ %-28s │\n", "Temperature", status.PON.Temperature+" C")
	fmt.Println("  ├─────────────────────────────────────────────────┤")
	fmt.Println("  │             WAN CONNECTION                      │")
	fmt.Println("  ├─────────────────────────────────────────────────┤")
	fmt.Printf("  │  %-18s │ %-28s │\n", "Connection", status.WAN.ConnectionName)
	fmt.Printf("  │  %-18s │ %-28s │\n", "Status", status.WAN.ConnectionStatus)
	fmt.Printf("  │  %-18s │ %-28s │\n", "IP Address", status.WAN.IPAddress)
	fmt.Printf("  │  %-18s │ %-28s │\n", "Gateway", status.WAN.Gateway)
	fmt.Printf("  │  %-18s │ %-28s │\n", "DNS", status.WAN.DNS)
	fmt.Printf("  │  %-18s │ %-28s │\n", "Online Duration", status.WAN.OnlineDuration)
	if status.WAN.DisconnectReason != "" {
		fmt.Printf("  │  %-18s │ %-28s │\n", "Disconnect Reason", status.WAN.DisconnectReason)
	}
	fmt.Println("  └─────────────────────────────────────────────────┘")
	fmt.Println()
}

func handleClients(client *router.Client) {
	fmt.Println("\n  Fetching connected devices...")

	wlans, err := router.FetchWLANStatus(client)
	if err != nil {
		fmt.Printf("  [ERROR] %s\n", err)
		return
	}

	lans, err := router.FetchLANStatus(client)
	if err != nil {
		fmt.Printf("  [ERROR] %s\n", err)
		return
	}

	fmt.Println()
	fmt.Println("  ┌─────────────────────────────────────────────────┐")
	fmt.Println("  │             WI-FI INTERFACES                    │")
	fmt.Println("  ├─────────────────────────────────────────────────┤")
	for _, w := range wlans {
		status := "OFF"
		if w.Enabled {
			status = "ON"
		}
		fmt.Printf("  │  SSID %-2d: %-12s  Status: %-4s  Ch: %-4s │\n",
			w.Index+1, truncate(w.SSIDName, 12), status, w.Channel)
		fmt.Printf("  │    Auth: %-12s  Enc: %-10s          │\n", w.AuthType, w.EncryptType)
		fmt.Printf("  │    MAC:  %-20s  Rx/Tx: %s/%s │\n", w.MACAddress, w.PacketsRx, w.PacketsTx)
		if w.Index < len(wlans)-1 {
			fmt.Println("  │                                                 │")
		}
	}
	fmt.Println("  ├─────────────────────────────────────────────────┤")
	fmt.Println("  │             ETHERNET PORTS                      │")
	fmt.Println("  ├─────────────────────────────────────────────────┤")
	for _, l := range lans {
		fmt.Printf("  │  %-10s │ Status: %-10s │ Speed: %-10s │\n", l.Port, l.Status, l.Speed)
		fmt.Printf("  │  Mode: %-8s  Errors: %-8s               │\n", l.Mode, l.ErrorFrames)
	}
	fmt.Println("  └─────────────────────────────────────────────────┘")
	fmt.Println()
}

func handleWiFiMenu(client *router.Client, scanner *bufio.Scanner) {
	for {
		fmt.Print(wifiMenu)
		fmt.Print("  Choose option [0-3]: ")

		if !scanner.Scan() {
			return
		}
		choice := strings.TrimSpace(scanner.Text())

		switch choice {
		case "1":
			showWiFiDetails(client)
		case "2":
			changeWiFiPassword(client, scanner)
		case "3":
			showWiFiQR(client)
		case "0":
			return
		default:
			fmt.Println("\n  [!] Invalid option.")
		}
	}
}

func showWiFiDetails(client *router.Client) {
	fmt.Println("\n  Fetching Wi-Fi settings...")

	sec, err := router.FetchWLANSecurity(client)
	if err != nil {
		fmt.Printf("  [ERROR] %s\n", err)
		return
	}

	fmt.Println()
	fmt.Println("  ┌─────────────────────────────────────────────────┐")
	fmt.Println("  │             WI-FI SETTINGS                      │")
	fmt.Println("  ├─────────────────────────────────────────────────┤")
	fmt.Printf("  │  %-18s │ %-28s │\n", "SSID Name", sec.SSID)
	fmt.Printf("  │  %-18s │ %-28s │\n", "Password", sec.Password)
	fmt.Printf("  │  %-18s │ %-28s │\n", "Auth Type", sec.AuthType)
	fmt.Printf("  │  %-18s │ %-28s │\n", "Encryption", sec.EncryptType)
	fmt.Println("  └─────────────────────────────────────────────────┘")
	fmt.Println()
}

func changeWiFiPassword(client *router.Client, scanner *bufio.Scanner) {
	fmt.Print("\n  Enter new Wi-Fi password (min 8 characters): ")
	if !scanner.Scan() {
		return
	}
	password := strings.TrimSpace(scanner.Text())

	if len(password) < 8 {
		fmt.Println("  [!] Password must be at least 8 characters.")
		return
	}

	fmt.Printf("  Changing Wi-Fi password to '%s'...\n", password)
	ok, err := router.SetWiFiPassword(client, password)
	if err != nil {
		fmt.Printf("  [ERROR] %s\n", err)
		return
	}
	if ok {
		fmt.Println("  [OK] Wi-Fi password changed successfully.")
	} else {
		fmt.Println("  [!] Failed to change Wi-Fi password.")
	}
}

func showWiFiQR(client *router.Client) {
	sec, err := router.FetchWLANSecurity(client)
	if err != nil {
		fmt.Printf("  [ERROR] %s\n", err)
		return
	}

	auth := "WPA"
	if sec.AuthType == "WPA2-PSK" || sec.AuthType == "WPA/WPA2-PSK" {
		auth = "WPA"
	}

	qrData := fmt.Sprintf("WIFI:T:%s;S:%s;P:%s;;", auth, sec.SSID, sec.Password)

	fmt.Println()
	fmt.Println("  Scan this QR code with your phone to connect:")
	fmt.Println()

	// Use a simple QR code renderer
	qr, err := router.GenerateQR(qrData)
	if err != nil {
		fmt.Printf("  [ERROR] Could not generate QR code: %s\n", err)
		fmt.Printf("  Wi-Fi string: %s\n", qrData)
		return
	}
	fmt.Println(qr)
}

func handleReboot(client *router.Client, scanner *bufio.Scanner) {
	fmt.Print("\n  Are you sure you want to reboot the router? [y/N]: ")
	if !scanner.Scan() {
		return
	}
	answer := strings.TrimSpace(strings.ToLower(scanner.Text()))
	if answer != "y" && answer != "yes" {
		fmt.Println("  Reboot cancelled.")
		return
	}

	fmt.Println("  Sending reboot command...")
	ok, err := router.Reboot(client)
	if err != nil {
		fmt.Printf("  [ERROR] %s\n", err)
		return
	}
	if ok {
		fmt.Println("  [OK] Router is rebooting...")
	} else {
		fmt.Println("  [!] Reboot command may have failed.")
	}
}

func handleTimeout(client *router.Client, scanner *bufio.Scanner) {
	timeout, err := router.FetchTimeout(client)
	if err != nil {
		fmt.Printf("  [ERROR] %s\n", err)
		return
	}

	fmt.Printf("\n  Current session timeout: %d minutes\n", timeout)

	fmt.Print(timeoutMenu)
	fmt.Print("  Enter new timeout (1-30) or 0 to cancel: ")

	if !scanner.Scan() {
		return
	}
	input := strings.TrimSpace(scanner.Text())
	newTimeout, err := strconv.Atoi(input)
	if err != nil || newTimeout < 0 || newTimeout > 30 {
		fmt.Println("  [!] Invalid input. Must be between 1 and 30.")
		return
	}
	if newTimeout == 0 {
		fmt.Println("  Cancelled.")
		return
	}

	fmt.Printf("  Setting timeout to %d minutes...\n", newTimeout)
	ok, err := router.SetTimeout(client, newTimeout)
	if err != nil {
		fmt.Printf("  [ERROR] %s\n", err)
		return
	}
	if ok {
		fmt.Printf("  [OK] Session timeout set to %d minutes.\n", newTimeout)
	} else {
		fmt.Println("  [!] Failed to set timeout.")
	}
}

// truncate shortens a string to maxLen characters.
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-2] + ".."
}
