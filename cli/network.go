package cli

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"

	"github.com/CallMeJaja/zte-cli/router"
)

// networkItems defines all items in the Network module.
var networkItems = []MenuItem{
	{ID: 1, Name: "WAN Connection"},
	{ID: 2, Name: "3G/4G WAN Connection"},
	{ID: 3, Name: "4in6 Tunnel Connection"},
	{ID: 4, Name: "6in4 Tunnel Connection"},
	{ID: 5, Name: "Port Binding"},
	{ID: 6, Name: "DHCP Release First"},
	{ID: 7, Name: "Basic (WLAN)"},
	{ID: 8, Name: "SSID Settings"},
	{ID: 9, Name: "Security"},
	{ID: 10, Name: "Access Control List"},
	{ID: 11, Name: "Associated Devices"},
	{ID: 12, Name: "WMM"},
	{ID: 13, Name: "WiFi Restrictions"},
	{ID: 14, Name: "WPS"},
	{ID: 15, Name: "DHCP Server"},
	{ID: 16, Name: "DHCP Server(IPv6)"},
	{ID: 17, Name: "DHCP Binding"},
	{ID: 18, Name: "DHCP Port Service"},
	{ID: 19, Name: "Prefix Management"},
	{ID: 20, Name: "DHCP Port Service(IPv6)"},
	{ID: 21, Name: "RA Service"},
	{ID: 22, Name: "LOID"},
	{ID: 23, Name: "SN"},
	{ID: 24, Name: "Default Gateway(IPv4)"},
	{ID: 25, Name: "Static Routing(IPv4)"},
	{ID: 26, Name: "Policy Routing(IPv4)"},
	{ID: 27, Name: "Routing Table(IPv4)"},
	{ID: 28, Name: "Default Gateway(IPv6)"},
	{ID: 29, Name: "Static Routing(IPv6)"},
	{ID: 30, Name: "Policy Routing(IPv6)"},
	{ID: 31, Name: "Routing Table(IPv6)"},
	{ID: 32, Name: "Port Locating"},
}

// networkPages maps item IDs to their .gch page names.
var networkPages = map[int]string{
	1:  "IPv46_net_wan2_conf_t.gch",
	2:  "net_ttywan_conf_t.gch",
	3:  "net_dslite_conf_t.gch",
	4:  "net_6in4_conf_t.gch",
	5:  "net_portbind_conf_t.gch",
	6:  "net_dhcpcleanlink_t.gch",
	7:  "net_wlanm_conf1_t.gch",
	8:  "net_wlanm_essid1_t.gch",
	9:  "net_wlanm_secrity1_t.gch",
	10: "net_wlanm_macfilter1_t.gch",
	11: "net_wlanm_assoc1_t.gch",
	12: "net_wlanm_media1_t.gch",
	13: "net_wlanm_off_t.gch",
	14: "net_wlanm_wps1_t.gch",
	15: "net_dhcp_dynamic_t.gch",
	16: "net_v6_dhcp_dynamic_t.gch",
	17: "net_dhcp_static_t.gch",
	18: "net_dhcp_mode_t.gch",
	19: "net_v6_prefix_management_t.gch",
	20: "net_v6_prefix_ban_port_t.gch",
	21: "net_v6_ra_server_t.gch",
	22: "pon_net_ponloid_t.gch",
	23: "pon_net_sn_t.gch",
	24: "net_route_default_t.gch",
	25: "net_route_static_t.gch",
	26: "app_route_policy_t.gch",
	27: "app_route_table_t.gch",
	28: "net_route6_default_t.gch",
	29: "net_route6_static_t.gch",
	30: "app_route6_policy_t.gch",
	31: "app_route6_table_t.gch",
	32: "net_port_locate_t.gch",
}

// RunNetworkModule launches the Network module interactive menu.
func RunNetworkModule(client *router.Client, scanner *bufio.Scanner) {
	for {
		menu := NewMenu("NETWORK", networkItems)
		menu.Render()

		fmt.Print("Pls enter command number: ")
		if !scanner.Scan() {
			return
		}
		choice := strings.TrimSpace(scanner.Text())
		if choice == "" {
			continue
		}
		if choice == "0" {
			return
		}

		id, _ := strconv.Atoi(choice)

		if id == 0 {
			continue
		}

		found := false
		for _, item := range networkItems {
			if item.ID == id {
				found = true
				break
			}
		}

		if !found {
			fmt.Println("\n  [!] Invalid option. Please try again.")
			continue
		}

		displayNetworkItem(client, id, scanner)
	}
}

// displayNetworkItem shows data for a specific network item with Refresh/Edit/Back.
func displayNetworkItem(client *router.Client, id int, scanner *bufio.Scanner) {
	page, ok := networkPages[id]
	if !ok {
		fmt.Println("\n  [!] Page not configured.")
		return
	}

	itemName := ""
	for _, item := range networkItems {
		if item.ID == id {
			itemName = item.Name
			break
		}
	}

	for {
		fmt.Println()
		RenderDataHeader(itemName, 55)

		lines, err := router.DisplayStatusItem(client, page)
		if err != nil {
			errMsg := err.Error()
			if strings.Contains(errMsg, "not found (404)") {
				fmt.Println("  Not Available")
			} else {
				fmt.Printf("  [ERROR] %s\n", errMsg)
			}
		} else if len(lines) == 0 {
			fmt.Println("  No connection or configuration available.")
		} else {
			for _, line := range lines {
				parts := strings.SplitN(line, "||", 2)
				if len(parts) == 2 {
					RenderDataRow(parts[0], parts[1], 22)
				}
			}
		}

		switch id {
		case 8:
			RenderDataFooter(55, "(1) Refresh  |  (2) Change SSID  |  (0) Kembali")
		case 9:
			RenderDataFooter(55, "(1) Refresh  |  (2) Change Password  |  (0) Kembali")
		default:
			RenderDataFooter(55, "(1) Refresh  |  (0) Kembali")
		}

		fmt.Print("Pls enter command number: ")
		if !scanner.Scan() {
			return
		}
		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}
		if input == "0" {
			return
		}
		if input == "1" {
			continue
		}

		if id == 8 && input == "2" {
			changeSSIDPrompt(client, scanner)
			continue
		}
		if id == 9 && input == "2" {
			changeWiFiPasswordFromNetwork(client, scanner)
			continue
		}

		fmt.Println("  [!] Invalid option.")
	}
}

func changeSSIDPrompt(client *router.Client, scanner *bufio.Scanner) {
	fmt.Print("\n  Enter new Wi-Fi SSID name: ")
	if !scanner.Scan() {
		return
	}
	ssid := strings.TrimSpace(scanner.Text())

	if len(ssid) == 0 {
		fmt.Println("  [!] SSID cannot be empty.")
		return
	}

	fmt.Println("  Changing Wi-Fi SSID...")
	ok, err := router.SetWiFiSSID(client, ssid)
	if err != nil {
		fmt.Printf("  [ERROR] %s\n", err)
		return
	}
	if ok {
		fmt.Println("  [OK] Wi-Fi SSID changed successfully.")
	} else {
		fmt.Println("  [!] Failed to change Wi-Fi SSID.")
	}
}

func changeWiFiPasswordFromNetwork(client *router.Client, scanner *bufio.Scanner) {
	fmt.Print("\n  Enter new Wi-Fi password (min 8 characters): ")
	if !scanner.Scan() {
		return
	}
	password := strings.TrimSpace(scanner.Text())

	if len(password) < 8 {
		fmt.Println("  [!] Password must be at least 8 characters.")
		return
	}

	fmt.Println("  Changing Wi-Fi password...")
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
