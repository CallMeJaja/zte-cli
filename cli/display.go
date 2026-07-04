package cli

import (
	"fmt"

	"github.com/CallMeJaja/zte-cli/router"
)

// displayStatus prints the full router status summary (used by direct commands).
func displayStatus(client *router.Client) {
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

// displayClients prints connected devices (used by direct commands).
func displayClients(client *router.Client) {
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

// displayWiFiDetails prints Wi-Fi settings (used by direct commands).
func displayWiFiDetails(client *router.Client) {
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

// displayWiFiQR prints Wi-Fi QR code data (used by direct commands).
func displayWiFiQR(client *router.Client) {
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

	qr, err := router.GenerateQR(qrData)
	if err != nil {
		fmt.Printf("  [ERROR] Could not generate QR code: %s\n", err)
		fmt.Printf("  Wi-Fi string: %s\n", qrData)
		return
	}
	fmt.Println(qr)
}

// truncate shortens a string to maxLen characters.
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-2] + ".."
}
