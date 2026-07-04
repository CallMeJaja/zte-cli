package cli

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/CallMeJaja/zte-cli/router"
)

// statusItems defines all items in the Status module.
var statusItems = []MenuItem{
	{ID: 1, Name: "Device Information"},
	{ID: 2, Name: "WAN Connection"},
	{ID: 3, Name: "3G/4G WAN Connection"},
	{ID: 4, Name: "4in6 Tunnel Connection"},
	{ID: 5, Name: "6in4 Tunnel Connection"},
	{ID: 6, Name: "PON Information"},
	{ID: 7, Name: "Mobile Network"},
	{ID: 8, Name: "WLAN Status"},
	{ID: 9, Name: "Ethernet Status"},
	{ID: 10, Name: "USB Status"},
	{ID: 11, Name: "VoIP Status"},
}

// statusPages maps item IDs to their .gch page names.
var statusPages = map[int]string{
	1:  "status_dev_info_t.gch",
	2:  "IPv46_status_wan2_if_t.gch",
	3:  "status_ttywan_info_t.gch",
	4:  "status_dslite_if_t.gch",
	5:  "status_6in4_info_t.gch",
	6:  "pon_status_link_info_t.gch",
	7:  "status_mobnet_info_t.gch",
	8:  "status_wlanm_info1_t.gch",
	9:  "pon_status_lan_info_t.gch",
	10: "status_usb_info_t.gch",
	11: "status_voip_4less_t.gch",
}


// RunStatusModule launches the Status module interactive menu.
func RunStatusModule(client *router.Client, scanner *bufio.Scanner) {
	for {
		menu := NewMenu("STATUS", statusItems)
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

		// Parse choice
		var id int
		fmt.Sscanf(choice, "%d", &id)

		if id == 0 {
			continue
		}

		// Find the item
		found := false
		for _, item := range statusItems {
			if item.ID == id {
				found = true
				break
			}
		}

		if !found {
			fmt.Println("\n  [!] Invalid option. Please try again.")
			continue
		}

		// Display data view with refresh loop
		displayStatusItem(client, id, scanner)
	}
}

// displayStatusItem shows data for a specific status item with Refresh/Back.
func displayStatusItem(client *router.Client, id int, scanner *bufio.Scanner) {
	page, ok := statusPages[id]
	if !ok {
		fmt.Println("\n  [!] Page not configured.")
		return
	}

	// Find the item name
	itemName := ""
	for _, item := range statusItems {
		if item.ID == id {
			itemName = item.Name
			break
		}
	}

	for {
		fmt.Println()
		RenderDataHeader(itemName, 55)

		// Fetch and display data
		lines, err := router.DisplayStatusItem(client, page)
		if err != nil {
			errMsg := err.Error()
			if strings.Contains(errMsg, "not found (404)") {
				fmt.Println("  Not Available")
			} else {
				fmt.Printf("  [ERROR] %s\n", errMsg)
			}
		} else if len(lines) == 0 {
			fmt.Println("  No data available")
		} else {
			for _, line := range lines {
				parts := strings.SplitN(line, "||", 2)
				if len(parts) == 2 {
					RenderDataRow(parts[0], parts[1], 22)
				}
			}
		}

		RenderDataFooter(55, "(1) Refresh  |  (0) Kembali")

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
			continue // Refresh
		}
		fmt.Println("  [!] Invalid option.")
	}
}
