package router

import (
	"fmt"
	"regexp"
	"strings"
)

// reRowPattern is the compiled regex for extracting table rows.
var reRowPattern = regexp.MustCompile(`(?s)<td\s+class="tdleft">(.*?)</td>\s*<td[^>]*>(.*?)</td>`)

// StatusPageParser is a function that parses a specific page and returns display lines.
type StatusPageParser func(html string) []string

// statusParsers maps page names to their custom parser functions.
var statusParsers = map[string]StatusPageParser{
	"status_dev_info_t.gch":      parseDeviceInfoPage,
	"pon_status_link_info_t.gch": parsePONStatusPage,
	"status_wlanm_info1_t.gch":   parseWLANStatusPage,
	"pon_status_lan_info_t.gch":  parseLANStatusPage,
	"net_wlanm_essid1_t.gch":     parseSSIDSettingsPage,
	"net_wlanm_secrity1_t.gch":   parseWLANSecurityPage,
	"net_wlanm_assoc1_t.gch":     parseAssociatedDevicesPage,
	"net_dhcp_dynamic_t.gch":     parseDHCPServerPage,
	"pon_net_ponloid_t.gch":      parsePONLOIDPage,
	"IPv46_net_wan2_conf_t.gch":  parseWANConfigPage,
	"pon_net_sn_t.gch":           parsePONSNPage,
	"net_wlanm_conf1_t.gch":      parseWLANBasicPage,
}

// DisplayStatusItem fetches a status page and returns formatted lines.
// Uses custom parsers when available, falls back to generic table parser.
func DisplayStatusItem(client *Client, pageName string) ([]string, error) {
	html, err := client.GetPage(pageName)
	if err != nil {
		return nil, err
	}

	// Check if there's a custom parser for this page
	if parser, ok := statusParsers[pageName]; ok {
		return parser(html), nil
	}

	// Fall back to generic table parser
	return ParseGenericTable(html), nil
}

// ParseGenericTable extracts key-value pairs from a table with tdleft/tdright classes.
func ParseGenericTable(html string) []string {
	lines := make([]string, 0)
	rows := reRowPattern.FindAllStringSubmatch(html, -1)

	for _, row := range rows {
		if len(row) < 3 {
			continue
		}
		label := strings.TrimSpace(stripHTMLTags(row[1]))
		value := strings.TrimSpace(row[2])

		if label == "" {
			continue
		}

		value = extractValueFromCell(value)
		value = decodeHTMLDecimalEntities(value)
		value = cleanHexEscapes(value)
		value = stripHTMLTags(value)
		value = strings.TrimSpace(value)

		if label != "" {
			lines = append(lines, fmt.Sprintf("%s||%s", label, value))
		}
	}
	return lines
}

// parseDeviceInfoPage extracts device info including JS-injected PON Serial Number.
func parseDeviceInfoPage(html string) []string {
	lines := ParseGenericTable(html)

	// PON Serial Number is injected via JS: var sn = "ZTEGcc84ae87"
	// Replace the empty entry from generic parser with the actual value
	if m := rePonSN.FindStringSubmatch(html); len(m) > 1 {
		ponSN := strings.ToUpper(m[1])
		found := false
		for i, line := range lines {
			if strings.HasPrefix(line, "PON Serial Number||") {
				lines[i] = fmt.Sprintf("PON Serial Number||%s", ponSN)
				found = true
				break
			}
		}
		if !found {
			lines = append(lines, fmt.Sprintf("PON Serial Number||%s", ponSN))
		}
	}

	return lines
}

// parsePONStatusPage extracts PON optical data from inline JavaScript variables.
func parsePONStatusPage(html string) []string {
	lines := make([]string, 0)

	// GPON State
	regStates := map[string]string{
		"1": "Initial State (o1)",
		"2": "Standby State (o2)",
		"3": "Serial Number State (o3)",
		"4": "Ranging State (o4)",
		"5": "Operation State (o5)",
		"6": "POPUP State (o6)",
		"7": "Emergency Stop State (o7)",
	}
	if m := rePonRegStatus.FindStringSubmatch(html); len(m) > 1 {
		state := regStates[m[1]]
		if state == "" {
			state = fmt.Sprintf("Unknown State (%s)", m[1])
		}
		lines = append(lines, fmt.Sprintf("GPON State||%s", state))
	}

	// Rx Power
	if m := rePonRxPower.FindStringSubmatch(html); len(m) > 1 {
		var val float64
		fmt.Sscanf(m[1], "%f", &val)
		lines = append(lines, fmt.Sprintf("Rx Power||%.2f dBm", val/10000))
	}

	// Tx Power
	if m := rePonTxPower.FindStringSubmatch(html); len(m) > 1 {
		var val float64
		fmt.Sscanf(m[1], "%f", &val)
		lines = append(lines, fmt.Sprintf("Tx Power||%.2f dBm", val/10000))
	}

	// Voltage, Bias Current, Temperature from HTML fields
	lines = append(lines, fmt.Sprintf("Voltage||%s uV", extractField(html, "Frm_Volt")))
	lines = append(lines, fmt.Sprintf("Bias Current||%s uA", extractField(html, "Frm_Current")))
	lines = append(lines, fmt.Sprintf("Temperature||%s C", extractField(html, "Frm_Temp")))

	return lines
}

// parseSSIDSettingsPage extracts SSID settings from Transfer_meaning calls.
// The page shows one SSID at a time. The current SSID is identified by IF_VIEWID.
func parseSSIDSettingsPage(html string) []string {
	lines := make([]string, 0)
	data := extractTransferMeanings(html)

	// Get current SSID identity
	viewID := data["IF_VIEWID"]
	if viewID != "" {
		viewID = cleanHexEscapes(viewID)
		lines = append(lines, fmt.Sprintf("Current SSID ID||%s", viewID))
	}

	// Get SSID name
	if ssid, ok := data["ESSID"]; ok && ssid != "" {
		lines = append(lines, fmt.Sprintf("SSID Name||%s", cleanHexEscapes(ssid)))
	}

	// Get status
	if enable, ok := data["Enable"]; ok {
		status := "Disabled"
		if enable == "1" {
			status = "Enabled"
		}
		lines = append(lines, fmt.Sprintf("Status||%s", status))
	}

	// Get channel
	if channel, ok := data["Channel"]; ok {
		lines = append(lines, fmt.Sprintf("Channel||%s", cleanHexEscapes(channel)))
	}

	return lines
}

// parseWLANSecurityPage extracts and decrypts the active WLAN security settings.
// Note: It needs the session token from the page to decrypt the password.
func parseWLANSecurityPage(html string) []string {
	lines := make([]string, 0)
	data := extractTransferMeanings(html)

	sessionToken := ""
	if m := reSessionToken.FindStringSubmatch(html); len(m) > 1 {
		sessionToken = m[1]
	}

	password := "N/A"
	if encPass, ok := data["KeyPassphrase"]; ok && encPass != "" {
		encPass = cleanHexEscapes(encPass)
		if dec, err := Decrypt(encPass, sessionToken); err == nil {
			password = dec
		} else {
			password = "[Decryption Error]"
		}
	}

	authType := getAuthType(data["BeaconType"], data["WEPAuthMode"], data["WPAAuthMode"], data["11iAuthMode"])
	encType := getEncryptType(data["BeaconType"], data["WPAEncryptType"], data["11iEncryptType"])

	lines = append(lines, fmt.Sprintf("Auth Type||%s", authType))
	lines = append(lines, fmt.Sprintf("Encryption||%s", encType))
	lines = append(lines, fmt.Sprintf("Password||%s", password))

	return lines
}

// parseAssociatedDevicesPage extracts the list of connected Wi-Fi clients.
func parseAssociatedDevicesPage(html string) []string {
	lines := make([]string, 0)
	// The table lists MAC Addresses. It's usually a standard list.
	// Often represented as: `var AssociatedDeviceMACAddress = "00:11:22:33:44:55,66:77:88:99:AA:BB";` 
	// Or parsed from the table. Let's rely on generic table parser for this since it parses rows.
	
	// Try parsing via generic table first
	genericLines := ParseGenericTable(html)
	if len(genericLines) > 0 {
		return genericLines
	}
	
	// Fallback to searching Transfer_meaning for MACs
	data := extractTransferMeanings(html)
	if macs, ok := data["AssociatedDeviceMACAddress"]; ok && macs != "" {
		macList := strings.Split(cleanHexEscapes(macs), ",")
		for i, mac := range macList {
			lines = append(lines, fmt.Sprintf("Client %d MAC||%s", i+1, mac))
		}
	}

	return lines
}

// parseDHCPServerPage extracts local LAN DHCP parameters.
func parseDHCPServerPage(html string) []string {
	lines := make([]string, 0)
	data := extractTransferMeanings(html)

	if gw, ok := data["IPAddress"]; ok {
		lines = append(lines, fmt.Sprintf("LAN IP Address||%s", cleanHexEscapes(gw)))
	}
	if mask, ok := data["SubnetMask"]; ok {
		lines = append(lines, fmt.Sprintf("Subnet Mask||%s", cleanHexEscapes(mask)))
	}
	if dstart, ok := data["MinAddress"]; ok {
		lines = append(lines, fmt.Sprintf("DHCP Start IP||%s", cleanHexEscapes(dstart)))
	}
	if dend, ok := data["MaxAddress"]; ok {
		lines = append(lines, fmt.Sprintf("DHCP End IP||%s", cleanHexEscapes(dend)))
	}
	if dns, ok := data["DNSServers"]; ok {
		lines = append(lines, fmt.Sprintf("DNS Server||%s", cleanHexEscapes(dns)))
	}
	if lease, ok := data["LeaseTime"]; ok {
		lines = append(lines, fmt.Sprintf("Lease Time (sec)||%s", cleanHexEscapes(lease)))
	}

	return lines
}

// parsePONLOIDPage extracts LOID and Password.
func parsePONLOIDPage(html string) []string {
	lines := make([]string, 0)
	data := extractTransferMeanings(html)

	if loid, ok := data["Loid"]; ok {
		lines = append(lines, fmt.Sprintf("LOID||%s", cleanHexEscapes(loid)))
	}
	if pass, ok := data["Password"]; ok {
		lines = append(lines, fmt.Sprintf("Password||%s", cleanHexEscapes(pass)))
	}

	return lines
}

// parseWANConfigPage extracts WAN connection configuration details.
func parseWANConfigPage(html string) []string {
	lines := make([]string, 0)
	data := extractTransferMeanings(html)

	// Extract WAN interface details
	if name, ok := data["IF_WANNAME0"]; ok && name != "" {
		lines = append(lines, fmt.Sprintf("WAN Name||%s", cleanHexEscapes(name)))
	}
	if ctype, ok := data["IF_WANCALLTYPE0"]; ok && ctype != "" {
		lines = append(lines, fmt.Sprintf("Connection Type||%s", cleanHexEscapes(ctype)))
	}
	if ipmode, ok := data["IF_WANCIPMODE0"]; ok && ipmode != "" {
		lines = append(lines, fmt.Sprintf("IP Mode||%s", cleanHexEscapes(ipmode)))
	}
	if identity, ok := data["IF_WANIDENTITY0"]; ok && identity != "" {
		lines = append(lines, fmt.Sprintf("WAN Identity||%s", cleanHexEscapes(identity)))
	}

	return lines
}

// parsePONSNPage extracts PON SN configuration.
func parsePONSNPage(html string) []string {
	lines := make([]string, 0)
	data := extractTransferMeanings(html)

	if sn, ok := data["SN"]; ok {
		lines = append(lines, fmt.Sprintf("PON SN||%s", cleanHexEscapes(sn)))
	}

	return lines
}

// parseWLANBasicPage extracts basic WLAN configuration.
func parseWLANBasicPage(html string) []string {
	lines := make([]string, 0)
	data := extractTransferMeanings(html)

	if enable, ok := data["Enable"]; ok {
		status := "Disabled"
		if enable == "1" {
			status = "Enabled"
		}
		lines = append(lines, fmt.Sprintf("Wireless Radio||%s", status))
	}
	if channel, ok := data["Channel"]; ok {
		lines = append(lines, fmt.Sprintf("Channel||%s", cleanHexEscapes(channel)))
	}
	if band, ok := data["Band"]; ok {
		lines = append(lines, fmt.Sprintf("Band||%s", cleanHexEscapes(band)))
	}
	if mode, ok := data["11nMode"]; ok {
		lines = append(lines, fmt.Sprintf("802.11 Mode||%s", cleanHexEscapes(mode)))
	}

	return lines
}

// extractValueFromCell extracts value from either input element or text.
func extractValueFromCell(cell string) string {
	vm := reInputVal.FindStringSubmatch(cell)
	if len(vm) > 1 {
		return vm[1]
	}
	return cell
}

// parseWLANStatusPage extracts WLAN interface details from Transfer_meaning calls.
func parseWLANStatusPage(html string) []string {
	lines := make([]string, 0)
	data := extractTransferMeanings(html)

	for i := 0; i < 4; i++ {
		suffix := fmt.Sprintf("%d", i)

		enabled := data["Enable"+suffix] == "1"
		status := "OFF"
		if enabled {
			status = "ON"
		}

		ssid := cleanHexEscapes(data["ESSID"+suffix])
		channel := cleanHexEscapes(data["ChannelInUsed"+suffix])
		mac := cleanHexEscapes(data["Bssid"+suffix])
		rxPkts := cleanHexEscapes(data["TotalPacketsReceived"+suffix])
		txPkts := cleanHexEscapes(data["TotalPacketsSent"+suffix])

		authType := getAuthType(
			data["BeaconType"+suffix],
			data["WEPAuthMode"+suffix],
			data["WPAAuthMode"+suffix],
			data["11iAuthMode"+suffix],
		)
		encType := getEncryptType(
			data["BeaconType"+suffix],
			data["WPAEncryptType"+suffix],
			data["11iEncryptType"+suffix],
		)

		lines = append(lines, fmt.Sprintf("SSID %d Name||%s", i+1, ssid))
		lines = append(lines, fmt.Sprintf("SSID %d Status||%s", i+1, status))
		lines = append(lines, fmt.Sprintf("SSID %d Channel||%s", i+1, channel))
		lines = append(lines, fmt.Sprintf("SSID %d Auth||%s", i+1, authType))
		lines = append(lines, fmt.Sprintf("SSID %d Encryption||%s", i+1, encType))
		lines = append(lines, fmt.Sprintf("SSID %d MAC||%s", i+1, mac))
		lines = append(lines, fmt.Sprintf("SSID %d Packets Rx||%s", i+1, rxPkts))
		lines = append(lines, fmt.Sprintf("SSID %d Packets Tx||%s", i+1, txPkts))
	}

	return lines
}

// parseLANStatusPage extracts Ethernet port details from HTML tables.
func parseLANStatusPage(html string) []string {
	lines := make([]string, 0)
	ports := ParseLANStatus(html)

	for _, port := range ports {
		lines = append(lines, fmt.Sprintf("Port||%s", port.Port))
		lines = append(lines, fmt.Sprintf("Status||%s", port.Status))
		lines = append(lines, fmt.Sprintf("Speed||%s", port.Speed))
		lines = append(lines, fmt.Sprintf("Mode||%s", port.Mode))
		lines = append(lines, fmt.Sprintf("Packets Rx||%s", port.PacketsRx))
		lines = append(lines, fmt.Sprintf("Bytes Rx||%s", port.BytesRx))
		lines = append(lines, fmt.Sprintf("Packets Tx||%s", port.PacketsTx))
		lines = append(lines, fmt.Sprintf("Bytes Tx||%s", port.BytesTx))
		lines = append(lines, fmt.Sprintf("Error Frames||%s", port.ErrorFrames))
	}

	return lines
}
