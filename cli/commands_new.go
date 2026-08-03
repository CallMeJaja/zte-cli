package cli

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/CallMeJaja/zte-cli/router"
)

// ============================================================
// Port Forwarding Commands
// ============================================================

// RunPortFwdListCommand handles: zte-cli portfw list
func RunPortFwdListCommand(client *router.Client) {
	ok, err := client.Login()
	if err != nil || !ok {
		fmt.Fprintf(os.Stderr, "  [ERROR] Login failed: %v\n", err)
		os.Exit(1)
	}

	rules, err := router.FetchPortForwardRules(client)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  [ERROR] %s\n", err)
		os.Exit(1)
	}

	fmt.Println("\n  Port Forwarding Rules:")
	fmt.Println("  " + strings.Repeat("─", 50))
	fmt.Println(router.FormatPortForwardRules(rules))
	fmt.Println()
}

// RunPortFwdAddCommand handles: zte-cli portfw add [flags]
func RunPortFwdAddCommand(client *router.Client, name, proto, wanPort, lanIP, lanPort string) {
	ok, err := client.Login()
	if err != nil || !ok {
		fmt.Fprintf(os.Stderr, "  [ERROR] Login failed: %v\n", err)
		os.Exit(1)
	}

	if name == "" || proto == "" || wanPort == "" || lanIP == "" || lanPort == "" {
		fmt.Fprintln(os.Stderr, "  [ERROR] All fields required: --name, --proto, --wan-port, --lan-ip, --lan-port")
		os.Exit(1)
	}

	rule := router.PortForwardRule{
		Enabled:      true,
		Name:         name,
		Protocol:     proto,
		WANStartPort: wanPort,
		WANEndPort:   wanPort,
		LANHostIP:    lanIP,
		LANStartPort: lanPort,
		LANEndPort:   lanPort,
	}

	ok, err = router.AddPortForwardRule(client, rule)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  [ERROR] %s\n", err)
		os.Exit(1)
	}

	if ok {
		fmt.Printf("  [OK] Port forwarding rule '%s' added.\n", name)
	}
}

// RunPortFwdDeleteCommand handles: zte-cli portfw delete <name>
func RunPortFwdDeleteCommand(client *router.Client, name string) {
	ok, err := client.Login()
	if err != nil || !ok {
		fmt.Fprintf(os.Stderr, "  [ERROR] Login failed: %v\n", err)
		os.Exit(1)
	}

	if name == "" {
		fmt.Fprintln(os.Stderr, "  [ERROR] Rule name required.")
		os.Exit(1)
	}

	ok, err = router.DeletePortForwardRule(client, name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  [ERROR] %s\n", err)
		os.Exit(1)
	}

	if ok {
		fmt.Printf("  [OK] Port forwarding rule '%s' deleted.\n", name)
	}
}

// RunPortFwdEnableCommand handles: zte-cli portfw enable/disable <name>
func RunPortFwdEnableCommand(client *router.Client, name string, enable bool) {
	ok, err := client.Login()
	if err != nil || !ok {
		fmt.Fprintf(os.Stderr, "  [ERROR] Login failed: %v\n", err)
		os.Exit(1)
	}

	if name == "" {
		fmt.Fprintln(os.Stderr, "  [ERROR] Rule name required.")
		os.Exit(1)
	}

	ok, err = router.SetPortForwardEnable(client, name, enable)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  [ERROR] %s\n", err)
		os.Exit(1)
	}

	action := "disabled"
	if enable {
		action = "enabled"
	}
	if ok {
		fmt.Printf("  [OK] Port forwarding rule '%s' %s.\n", name, action)
	}
}

// ============================================================
// DHCP/LAN Commands
// ============================================================

// RunDHCPShowCommand handles: zte-cli lan dhcp show
func RunDHCPShowCommand(client *router.Client) {
	ok, err := client.Login()
	if err != nil || !ok {
		fmt.Fprintf(os.Stderr, "  [ERROR] Login failed: %v\n", err)
		os.Exit(1)
	}

	settings, err := router.FetchDHCPSettings(client)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  [ERROR] %s\n", err)
		os.Exit(1)
	}

	fmt.Println("\n  DHCP Server Settings:")
	fmt.Println("  " + strings.Repeat("─", 50))
	fmt.Println(router.FormatDHCPSettings(settings))
	fmt.Println()
}

// RunDHCPLeasesCommand handles: zte-cli lan dhcp leases
func RunDHCPLeasesCommand(client *router.Client) {
	ok, err := client.Login()
	if err != nil || !ok {
		fmt.Fprintf(os.Stderr, "  [ERROR] Login failed: %v\n", err)
		os.Exit(1)
	}

	leases, err := router.FetchDHCPLeases(client)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  [ERROR] %s\n", err)
		os.Exit(1)
	}

	fmt.Println("\n  DHCP Leases:")
	fmt.Println("  " + strings.Repeat("─", 60))
	fmt.Println(router.FormatDHCPLeases(leases))
	fmt.Println()
}

// ============================================================
// WiFi Advanced Commands
// ============================================================

// RunWiFiBasicCommand handles: zte-cli wifi basic
func RunWiFiBasicCommand(client *router.Client) {
	ok, err := client.Login()
	if err != nil || !ok {
		fmt.Fprintf(os.Stderr, "  [ERROR] Login failed: %v\n", err)
		os.Exit(1)
	}

	settings, err := router.FetchWiFiBasicSettings(client)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  [ERROR] %s\n", err)
		os.Exit(1)
	}

	fmt.Println("\n  WiFi Basic Settings:")
	fmt.Println("  " + strings.Repeat("─", 50))
	fmt.Println(router.FormatWiFiBasicSettings(settings))
	fmt.Println()
}

// RunWiFiSetChannelCommand handles: zte-cli wifi set-channel <ch>
func RunWiFiSetChannelCommand(client *router.Client, channel string) {
	ok, err := client.Login()
	if err != nil || !ok {
		fmt.Fprintf(os.Stderr, "  [ERROR] Login failed: %v\n", err)
		os.Exit(1)
	}

	ok, err = router.SetWiFiChannel(client, channel)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  [ERROR] %s\n", err)
		os.Exit(1)
	}

	if ok {
		fmt.Printf("  [OK] WiFi channel set to %s.\n", channel)
	}
}

// RunWiFiSetPowerCommand handles: zte-cli wifi set-power <20-100>
func RunWiFiSetPowerCommand(client *router.Client, power string) {
	ok, err := client.Login()
	if err != nil || !ok {
		fmt.Fprintf(os.Stderr, "  [ERROR] Login failed: %v\n", err)
		os.Exit(1)
	}

	// Validate and map power value
	powerMap := map[string]string{
		"20": "20%", "40": "40%", "60": "60%", "80": "80%", "100": "100%",
	}
	mapped, ok := powerMap[power]
	if !ok {
		fmt.Fprintln(os.Stderr, "  [ERROR] Power must be 20, 40, 60, 80, or 100.")
		os.Exit(1)
	}

	ok, err = router.SetWiFiTransmitPower(client, mapped)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  [ERROR] %s\n", err)
		os.Exit(1)
	}

	if ok {
		fmt.Printf("  [OK] WiFi transmit power set to %s.\n", mapped)
	}
}

// RunWiFiToggleCommand handles: zte-cli wifi enable/disable
func RunWiFiToggleCommand(client *router.Client, enable bool) {
	ok, err := client.Login()
	if err != nil || !ok {
		fmt.Fprintf(os.Stderr, "  [ERROR] Login failed: %v\n", err)
		os.Exit(1)
	}

	ok, err = router.SetWiFiRFMode(client, enable)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  [ERROR] %s\n", err)
		os.Exit(1)
	}

	action := "disabled"
	if enable {
		action = "enabled"
	}
	if ok {
		fmt.Printf("  [OK] WiFi radio %s.\n", action)
	}
}

// ============================================================
// Firewall Commands
// ============================================================

// RunFirewallShowCommand handles: zte-cli firewall show
func RunFirewallShowCommand(client *router.Client) {
	ok, err := client.Login()
	if err != nil || !ok {
		fmt.Fprintf(os.Stderr, "  [ERROR] Login failed: %v\n", err)
		os.Exit(1)
	}

	settings, err := router.FetchFirewallSettings(client)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  [ERROR] %s\n", err)
		os.Exit(1)
	}

	fmt.Println("\n  Firewall Settings:")
	fmt.Println("  " + strings.Repeat("─", 50))
	fmt.Println(router.FormatFirewallSettings(settings))
	fmt.Println()
}

// RunFirewallLevelCommand handles: zte-cli firewall level <level>
func RunFirewallLevelCommand(client *router.Client, level string) {
	ok, err := client.Login()
	if err != nil || !ok {
		fmt.Fprintf(os.Stderr, "  [ERROR] Login failed: %v\n", err)
		os.Exit(1)
	}

	validLevels := map[string]bool{"off": true, "low": true, "middle": true, "mid": true, "high": true}
	if !validLevels[strings.ToLower(level)] {
		fmt.Fprintln(os.Stderr, "  [ERROR] Level must be: off, low, middle, high")
		os.Exit(1)
	}

	// Capitalize first letter
	level = strings.ToUpper(level[:1]) + strings.ToLower(level[1:])
	if level == "Mid" {
		level = "Middle"
	}

	ok, err = router.SetFirewallSettings(client, router.FirewallSettings{
		AntiHacking: true,
		Level:       level,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "  [ERROR] %s\n", err)
		os.Exit(1)
	}

	if ok {
		fmt.Printf("  [OK] Firewall level set to %s.\n", level)
	}
}

// ============================================================
// IP Filter Commands
// ============================================================

// RunIPFilterListCommand handles: zte-cli firewall ipfilter list
func RunIPFilterListCommand(client *router.Client) {
	ok, err := client.Login()
	if err != nil || !ok {
		fmt.Fprintf(os.Stderr, "  [ERROR] Login failed: %v\n", err)
		os.Exit(1)
	}

	rules, err := router.FetchIPFilterRules(client)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  [ERROR] %s\n", err)
		os.Exit(1)
	}

	fmt.Println("\n  IP Filter Rules:")
	fmt.Println("  " + strings.Repeat("─", 50))
	fmt.Println(router.FormatIPFilterRules(rules))
	fmt.Println()
}

// ============================================================
// MAC Filter Commands
// ============================================================

// RunMACFilterListCommand handles: zte-cli firewall macfilter list
func RunMACFilterListCommand(client *router.Client) {
	ok, err := client.Login()
	if err != nil || !ok {
		fmt.Fprintf(os.Stderr, "  [ERROR] Login failed: %v\n", err)
		os.Exit(1)
	}

	rules, err := router.FetchMACFilterRules(client)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  [ERROR] %s\n", err)
		os.Exit(1)
	}

	fmt.Println("\n  MAC Filter Rules:")
	fmt.Println("  " + strings.Repeat("─", 50))
	fmt.Println(router.FormatMACFilterRules(rules))
	fmt.Println()
}

// ============================================================
// URL Filter Commands
// ============================================================

// RunURLFilterListCommand handles: zte-cli firewall urlfilter list
func RunURLFilterListCommand(client *router.Client) {
	ok, err := client.Login()
	if err != nil || !ok {
		fmt.Fprintf(os.Stderr, "  [ERROR] Login failed: %v\n", err)
		os.Exit(1)
	}

	rules, err := router.FetchURLFilterRules(client)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  [ERROR] %s\n", err)
		os.Exit(1)
	}

	fmt.Println("\n  URL Filter Rules:")
	fmt.Println("  " + strings.Repeat("─", 50))
	fmt.Println(router.FormatURLFilterRules(rules))
	fmt.Println()
}

// RunURLFilterAddCommand handles: zte-cli firewall urlfilter add <url>
func RunURLFilterAddCommand(client *router.Client, urlStr string) {
	ok, err := client.Login()
	if err != nil || !ok {
		fmt.Fprintf(os.Stderr, "  [ERROR] Login failed: %v\n", err)
		os.Exit(1)
	}

	if urlStr == "" {
		fmt.Fprintln(os.Stderr, "  [ERROR] URL required.")
		os.Exit(1)
	}

	ok, err = router.AddURLFilterRule(client, router.URLFilterRule{
		Enabled: true,
		Mode:    "Discard",
		URL:     urlStr,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "  [ERROR] %s\n", err)
		os.Exit(1)
	}

	if ok {
		fmt.Printf("  [OK] URL filter rule added: %s\n", urlStr)
	}
}

// ============================================================
// DDNS Commands
// ============================================================

// RunDDNSShowCommand handles: zte-cli ddns show
func RunDDNSShowCommand(client *router.Client) {
	ok, err := client.Login()
	if err != nil || !ok {
		fmt.Fprintf(os.Stderr, "  [ERROR] Login failed: %v\n", err)
		os.Exit(1)
	}

	settings, err := router.FetchDDNSSettings(client)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  [ERROR] %s\n", err)
		os.Exit(1)
	}

	fmt.Println("\n  DDNS Settings:")
	fmt.Println("  " + strings.Repeat("─", 50))
	fmt.Println(router.FormatDDNSSettings(settings))
	fmt.Println()
}

// ============================================================
// DMZ Commands
// ============================================================

// RunDMZShowCommand handles: zte-cli dmz show
func RunDMZShowCommand(client *router.Client) {
	ok, err := client.Login()
	if err != nil || !ok {
		fmt.Fprintf(os.Stderr, "  [ERROR] Login failed: %v\n", err)
		os.Exit(1)
	}

	settings, err := router.FetchDMZSettings(client)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  [ERROR] %s\n", err)
		os.Exit(1)
	}

	fmt.Println("\n  DMZ Settings:")
	fmt.Println("  " + strings.Repeat("─", 50))
	fmt.Println(router.FormatDMZSettings(settings))
	fmt.Println()
}

// RunDMZSetCommand handles: zte-cli dmz set <ip>
func RunDMZSetCommand(client *router.Client, ip string) {
	ok, err := client.Login()
	if err != nil || !ok {
		fmt.Fprintf(os.Stderr, "  [ERROR] Login failed: %v\n", err)
		os.Exit(1)
	}

	if ip == "" {
		fmt.Fprintln(os.Stderr, "  [ERROR] IP address required.")
		os.Exit(1)
	}

	ok, err = router.SetDMZSettings(client, router.DMZSettings{
		Enabled: true,
		HostIP:  ip,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "  [ERROR] %s\n", err)
		os.Exit(1)
	}

	if ok {
		fmt.Printf("  [OK] DMZ host set to %s.\n", ip)
	}
}

// ============================================================
// UPnP Commands
// ============================================================

// RunUPnPShowCommand handles: zte-cli upnp show
func RunUPnPShowCommand(client *router.Client) {
	ok, err := client.Login()
	if err != nil || !ok {
		fmt.Fprintf(os.Stderr, "  [ERROR] Login failed: %v\n", err)
		os.Exit(1)
	}

	settings, err := router.FetchUPnPSettings(client)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  [ERROR] %s\n", err)
		os.Exit(1)
	}

	fmt.Println("\n  UPnP Settings:")
	fmt.Println("  " + strings.Repeat("─", 50))
	fmt.Println(router.FormatUPnPSettings(settings))
	fmt.Println()
}

// RunUPnPToggleCommand handles: zte-cli upnp enable/disable
func RunUPnPToggleCommand(client *router.Client, enable bool) {
	ok, err := client.Login()
	if err != nil || !ok {
		fmt.Fprintf(os.Stderr, "  [ERROR] Login failed: %v\n", err)
		os.Exit(1)
	}

	settings, err := router.FetchUPnPSettings(client)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  [ERROR] %s\n", err)
		os.Exit(1)
	}

	settings.Enabled = enable
	ok, err = router.SetUPnPSettings(client, *settings)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  [ERROR] %s\n", err)
		os.Exit(1)
	}

	action := "disabled"
	if enable {
		action = "enabled"
	}
	if ok {
		fmt.Printf("  [OK] UPnP %s.\n", action)
	}
}

// ============================================================
// Diagnosis Commands
// ============================================================

// RunPingCommand handles: zte-cli ping <host>
func RunPingCommand(client *router.Client, host string) {
	ok, err := client.Login()
	if err != nil || !ok {
		fmt.Fprintf(os.Stderr, "  [ERROR] Login failed: %v\n", err)
		os.Exit(1)
	}

	if host == "" {
		fmt.Fprintln(os.Stderr, "  [ERROR] Host required.")
		os.Exit(1)
	}

	fmt.Printf("  Pinging %s from router...\n\n", host)
	result, err := router.Ping(client, host)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  [ERROR] %s\n", err)
		os.Exit(1)
	}

	fmt.Println(result)
}

// RunARPTableCommand handles: zte-cli arp
func RunARPTableCommand(client *router.Client) {
	ok, err := client.Login()
	if err != nil || !ok {
		fmt.Fprintf(os.Stderr, "  [ERROR] Login failed: %v\n", err)
		os.Exit(1)
	}

	result, err := router.FetchARPTable(client)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  [ERROR] %s\n", err)
		os.Exit(1)
	}

	fmt.Println("\n  ARP Table:")
	fmt.Println("  " + strings.Repeat("─", 50))
	fmt.Println(result)
	fmt.Println()
}

// ============================================================
// System Management Commands
// ============================================================

// RunLogShowCommand handles: zte-cli log show [limit]
func RunLogShowCommand(client *router.Client, limitStr string) {
	ok, err := client.Login()
	if err != nil || !ok {
		fmt.Fprintf(os.Stderr, "  [ERROR] Login failed: %v\n", err)
		os.Exit(1)
	}

	limit := 50 // default
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	logs, err := router.FetchLogs(client)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  [ERROR] %s\n", err)
		os.Exit(1)
	}

	fmt.Println("\n  System Logs (last", limit, "entries):")
	fmt.Println("  " + strings.Repeat("─", 60))
	fmt.Println(router.FormatLogs(logs, limit))
	fmt.Println()
}

// RunResetCommand handles: zte-cli reset
func RunResetCommand(client *router.Client, confirm bool) {
	ok, err := client.Login()
	if err != nil || !ok {
		fmt.Fprintf(os.Stderr, "  [ERROR] Login failed: %v\n", err)
		os.Exit(1)
	}

	if !confirm {
		fmt.Print("  ⚠️  WARNING: This will reset ALL settings to factory defaults! [y/N]: ")
		var answer string
		_, _ = fmt.Scanln(&answer)
		if answer != "y" && answer != "yes" {
			fmt.Println("  Factory reset cancelled.")
			return
		}
	}

	fmt.Println("  Sending factory reset command...")
	ok, err = router.FactoryReset(client)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  [ERROR] %s\n", err)
		os.Exit(1)
	}
	if ok {
		fmt.Println("  [OK] Router is resetting to factory defaults...")
	}
}

// ============================================================
// WAN Management Commands
// ============================================================

// RunWANListCommand handles: zte-cli wan list
func RunWANListCommand(client *router.Client) {
	ok, err := client.Login()
	if err != nil || !ok {
		fmt.Fprintf(os.Stderr, "  [ERROR] Login failed: %v\n", err)
		os.Exit(1)
	}

	conns, err := router.FetchWANConnections(client)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  [ERROR] %s\n", err)
		os.Exit(1)
	}

	fmt.Println("\n  WAN Connections:")
	fmt.Println("  " + strings.Repeat("─", 50))
	fmt.Println(router.FormatWANConnections(conns))
	fmt.Println()
}

// RunWANDeleteCommand handles: zte-cli wan delete <name>
func RunWANDeleteCommand(client *router.Client, name string) {
	ok, err := client.Login()
	if err != nil || !ok {
		fmt.Fprintf(os.Stderr, "  [ERROR] Login failed: %v\n", err)
		os.Exit(1)
	}

	if name == "" {
		fmt.Fprintln(os.Stderr, "  [ERROR] Connection name required.")
		os.Exit(1)
	}

	ok, err = router.DeleteWANConnection(client, name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  [ERROR] %s\n", err)
		os.Exit(1)
	}

	if ok {
		fmt.Printf("  [OK] WAN connection '%s' deleted.\n", name)
	}
}

// ============================================================
// User Management Commands
// ============================================================

// RunChangePasswordCommand handles: zte-cli password set
func RunChangePasswordCommand(client *router.Client, oldPass, newPass, confirmPass string) {
	ok, err := client.Login()
	if err != nil || !ok {
		fmt.Fprintf(os.Stderr, "  [ERROR] Login failed: %v\n", err)
		os.Exit(1)
	}

	if oldPass == "" || newPass == "" || confirmPass == "" {
		fmt.Fprintln(os.Stderr, "  [ERROR] All fields required: --old, --new, --confirm")
		os.Exit(1)
	}

	ok, err = router.ChangePassword(client, oldPass, newPass, confirmPass)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  [ERROR] %s\n", err)
		os.Exit(1)
	}

	if ok {
		fmt.Println("  [OK] Password changed successfully.")
	}
}
