package router

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

// Regex patterns for parsing .gch pages
var (
	reTransferMeaning = regexp.MustCompile(`Transfer_meaning\('([^']*)',\s*'([^']*)'\);`)
	rePonRxPower      = regexp.MustCompile(`RxPower = "([^"]+)"`)
	rePonTxPower      = regexp.MustCompile(`TxPower = "([^"]+)"`)
	rePonRegStatus    = regexp.MustCompile(`RegStatus = "([^"]+)"`)
	rePonSN           = regexp.MustCompile(`var sn = "([^"]+)";`)
	reHTMLDecimal     = regexp.MustCompile(`&#([0-9]+);`)
)

// DeviceInfo holds basic router identification.
type DeviceInfo struct {
	Model           string
	SerialNumber    string
	HardwareVersion string
	SoftwareVersion string
	BootLoaderVersion string
	PONSerialNumber string
	BatchNumber     string
}

// PONStatus holds optical module readings.
type PONStatus struct {
	GPONState  string
	RxPower    string // dBm
	TxPower    string // dBm
	Voltage    string // uV
	BiasCurrent string // uA
	Temperature string // C
}

// WANStatus holds WAN connection details.
type WANStatus struct {
	ConnectionName   string
	IPVersion        string
	NAT              string
	IPAddress        string
	SubnetMask       string
	Gateway          string
	DNS              string
	ConnectionStatus string
	OnlineDuration   string
	DisconnectReason string
	WANMAC           string
}

// WLANInterface holds Wi-Fi interface details.
type WLANInterface struct {
	Index         int
	Enabled       bool
	SSIDName      string
	AuthType      string
	EncryptType   string
	MACAddress    string
	PacketsRx     string
	PacketsTx     string
	BytesRx       string
	BytesTx       string
	ErrorsRx      string
	ErrorsTx      string
	DiscardedRx   string
	DiscardedTx   string
	Channel       string
	RadioStatus   bool
}

// LANPort holds Ethernet port details.
type LANPort struct {
	Port        string
	Status      string
	Speed       string
	Mode        string
	PacketsRx   string
	BytesRx     string
	PacketsTx   string
	BytesTx     string
	ErrorFrames string
}

// WLANSecurity holds Wi-Fi security configuration.
type WLANSecurity struct {
	SSID          string
	Password      string
	AuthType      string
	EncryptType   string
	BeaconType    string
}

// SystemStatus is the consolidated status snapshot.
type SystemStatus struct {
	Device DeviceInfo
	PON    PONStatus
	WAN    WANStatus
	WLANs  []WLANInterface
	LANs   []LANPort
}

// ParseDeviceInfo extracts model, serial number, hardware/software version, etc.
func ParseDeviceInfo(html string) DeviceInfo {
	info := DeviceInfo{}

	info.Model = extractField(html, "Frm_ModelName")
	info.SerialNumber = extractField(html, "Frm_SerialNumber")
	info.HardwareVersion = extractField(html, "Frm_HardwareVer")
	info.SoftwareVersion = extractField(html, "Frm_SoftwareVer")
	info.BootLoaderVersion = extractField(html, "Frm_BootVer")
	info.BatchNumber = extractField(html, "Frm_SoftwareVerExtent")

	// PON Serial Number is injected via JS: var sn = "ZTEGcc84ae87";
	if m := rePonSN.FindStringSubmatch(html); len(m) > 1 {
		info.PONSerialNumber = strings.ToUpper(m[1])
	}

	return info
}

// ParsePONStatus extracts optical power, state, and module details.
func ParsePONStatus(html string) PONStatus {
	status := PONStatus{}

	// GPON Registration State
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
		status.GPONState = regStates[m[1]]
		if status.GPONState == "" {
			status.GPONState = fmt.Sprintf("Unknown State (%s)", m[1])
		}
	}

	// Rx Power
	if m := rePonRxPower.FindStringSubmatch(html); len(m) > 1 {
		if v, err := strconv.ParseFloat(m[1], 64); err == nil {
			status.RxPower = fmt.Sprintf("%.2f dBm", v/10000)
		}
	}

	// Tx Power
	if m := rePonTxPower.FindStringSubmatch(html); len(m) > 1 {
		if v, err := strconv.ParseFloat(m[1], 64); err == nil {
			status.TxPower = fmt.Sprintf("%.2f dBm", v/10000)
		}
	}

	// Voltage, Bias Current, Temperature from HTML fields
	status.Voltage = extractField(html, "Frm_Volt")
	status.BiasCurrent = extractField(html, "Frm_Current")
	status.Temperature = extractField(html, "Frm_Temp")

	return status
}

// ParseWANStatus extracts WAN connection details from the status page.
// The WAN page uses <input id="..." value="..."> elements rather than Transfer_meaning.
func ParseWANStatus(html string) WANStatus {
	status := WANStatus{}

	// Extract values from input elements by ID
	status.ConnectionName = extractInputValue(html, "TextWANCName0")
	if status.ConnectionName == "" {
		status.ConnectionName = extractInputValue(html, "TextETHCName0")
	}

	status.IPVersion = extractInputValue(html, "TextPPPIpMode0")
	if status.IPVersion == "" {
		status.IPVersion = extractInputValue(html, "TextETHIpMode0")
	}

	status.NAT = extractInputValue(html, "TextPPPIsNAT0")
	if status.NAT == "" {
		status.NAT = extractInputValue(html, "TextETHIsNAT0")
	}

	status.IPAddress = extractInputValue(html, "TextPPPIPAddress0")
	if status.IPAddress == "" {
		status.IPAddress = extractInputValue(html, "TextETHIPAddress0")
	}

	// Subnet and Gateway don't always have IDs, extract from the table pattern
	status.SubnetMask = extractFieldValue(html, "Subnet Mask")
	status.Gateway = extractFieldValue(html, "Gateway")
	status.DNS = extractInputValue(html, "TextPPPDNS0")
	if status.DNS == "" {
		status.DNS = extractInputValue(html, "TextETHDNS0")
	}

	status.ConnectionStatus = extractInputValue(html, "TextPPPConStatus0")
	if status.ConnectionStatus == "" {
		status.ConnectionStatus = extractInputValue(html, "TextETHConStatus0")
	}

	status.OnlineDuration = extractInputValue(html, "TextPPPUpTime0")
	if status.OnlineDuration == "" {
		status.OnlineDuration = extractInputValue(html, "TextETHUpTime0")
	}

	status.DisconnectReason = extractFieldValue(html, "Disconnect Reason")
	status.WANMAC = extractInputValue(html, "TextPPPWorkIFMac0")
	if status.WANMAC == "" {
		status.WANMAC = extractInputValue(html, "TextETHWorkIFMac0")
	}

	return status
}

// ParseWLANStatus extracts Wi-Fi interface details for all 4 SSIDs.
func ParseWLANStatus(html string) []WLANInterface {
	data := extractTransferMeanings(html)
	wlans := make([]WLANInterface, 0, 4)

	for i := 0; i < 4; i++ {
		suffix := strconv.Itoa(i)
		wlan := WLANInterface{
			Index: i,
		}

		if v, ok := data["Enable"+suffix]; ok {
			wlan.Enabled = v == "1"
		}
		if v, ok := data["ESSID"+suffix]; ok {
			wlan.SSIDName = cleanHexEscapes(v)
		}
		if v, ok := data["BeaconType"+suffix]; ok {
			wlan.AuthType = getAuthType(
				v,
				data["WEPAuthMode"+suffix],
				data["WPAAuthMode"+suffix],
				data["11iAuthMode"+suffix],
			)
			wlan.EncryptType = getEncryptType(
				v,
				data["WPAEncryptType"+suffix],
				data["11iEncryptType"+suffix],
			)
		}
		if v, ok := data["Bssid"+suffix]; ok {
			wlan.MACAddress = cleanHexEscapes(v)
		}
		if v, ok := data["TotalPacketsReceived"+suffix]; ok {
			wlan.PacketsRx = cleanHexEscapes(v)
		}
		if v, ok := data["TotalPacketsSent"+suffix]; ok {
			wlan.PacketsTx = cleanHexEscapes(v)
		}
		if v, ok := data["TotalBytesReceived"+suffix]; ok {
			wlan.BytesRx = cleanHexEscapes(v)
		}
		if v, ok := data["TotalBytesSent"+suffix]; ok {
			wlan.BytesTx = cleanHexEscapes(v)
		}
		if v, ok := data["ErrorsReceived"+suffix]; ok {
			wlan.ErrorsRx = cleanHexEscapes(v)
		}
		if v, ok := data["ErrorsSent"+suffix]; ok {
			wlan.ErrorsTx = cleanHexEscapes(v)
		}
		if v, ok := data["DiscardPacketsReceived"+suffix]; ok {
			wlan.DiscardedRx = cleanHexEscapes(v)
		}
		if v, ok := data["DiscardPacketsSent"+suffix]; ok {
			wlan.DiscardedTx = cleanHexEscapes(v)
		}
		if v, ok := data["ChannelInUsed"+suffix]; ok {
			wlan.Channel = cleanHexEscapes(v)
		}
		if v, ok := data["RadioStatus"+suffix]; ok {
			wlan.RadioStatus = v == "1"
		}

		wlans = append(wlans, wlan)
	}

	return wlans
}

// ParseLANStatus extracts Ethernet port status from HTML tables.
func ParseLANStatus(html string) []LANPort {
	ports := make([]LANPort, 0, 4)

	// The page has 4 tables, each for one LAN port.
	// Each table has sequential rows: Ethernet Port, Status, Speed, Mode, Packets Rx/Bytes Rx, Packets Tx/Bytes Tx, Error Frames
	// Row format: <td class="tdleft">LABEL</td> <td class="tdright">VALUE</td>

	tablePattern := regexp.MustCompile(`(?s)<table[^>]*id="TestContent\d?"[^>]*>(.*?)</table>`)
	rowPattern := regexp.MustCompile(`(?s)<td\s+class="tdleft">(.*?)</td>\s*<td[^>]*>(.*?)</td>`)

	tables := tablePattern.FindAllStringSubmatch(html, -1)

	for _, table := range tables {
		if len(table) < 2 {
			continue
		}
		tableContent := table[1]

		rows := rowPattern.FindAllStringSubmatch(tableContent, -1)
		port := LANPort{}
		fieldIdx := 0

		for _, row := range rows {
			if len(row) < 3 {
				continue
			}
			value := stripHTMLTags(row[2])
			value = strings.TrimSpace(value)

			switch fieldIdx {
			case 0:
				port.Port = value
			case 1:
				port.Status = value
			case 2:
				port.Speed = value
			case 3:
				port.Mode = value
			case 4:
				parts := strings.SplitN(value, "/", 2)
				if len(parts) == 2 {
					port.PacketsRx = strings.TrimSpace(parts[0])
					port.BytesRx = strings.TrimSpace(parts[1])
				}
			case 5:
				parts := strings.SplitN(value, "/", 2)
				if len(parts) == 2 {
					port.PacketsTx = strings.TrimSpace(parts[0])
					port.BytesTx = strings.TrimSpace(parts[1])
				}
			case 6:
				// Error Frames may contain inline JavaScript, extract just the numeric value
				port.ErrorFrames = extractScriptValue(value)
			}
			fieldIdx++
		}

		if port.Port != "" {
			ports = append(ports, port)
		}
	}

	return ports
}

// extractInputValue extracts the value from an <input id="..." value="..."> element.
var reInputValue = regexp.MustCompile(`(?s)id="([^"]*)"[^>]*value="?([^">\s]*)"?\s`)

func extractInputValue(html string, id string) string {
	pattern := regexp.MustCompile(fmt.Sprintf(`(?s)id="%s"[^>]*value="?([^">\s]*)"?\s`, id))
	m := pattern.FindStringSubmatch(html)
	if len(m) > 1 {
		raw := strings.TrimSpace(m[1])
		raw = decodeHTMLDecimalEntities(raw)
		raw = cleanHexEscapes(raw)
		return raw
	}
	return ""
}

// extractFieldValue extracts a value from the row following a label.
// Looks for pattern: <td class="tdleft">LABEL</td>\n<td class="tdright">VALUE</td>
// Handles both plain text values and <input value="..."> elements.
func extractFieldValue(html string, label string) string {
	pattern := regexp.MustCompile(fmt.Sprintf(`(?s)class="tdleft">%s</td>\s*<td[^>]*>(.*?)</td>`, regexp.QuoteMeta(label)))
	m := pattern.FindStringSubmatch(html)
	if len(m) > 1 {
		cellContent := m[1]
		// Try to extract value from <input> element first
		inputVal := regexp.MustCompile(`(?i)value="?([^">\s]*)"?\s`)
		vm := inputVal.FindStringSubmatch(cellContent)
		if len(vm) > 1 {
			raw := strings.TrimSpace(vm[1])
			raw = decodeHTMLDecimalEntities(raw)
			raw = cleanHexEscapes(raw)
			return raw
		}
		// Fallback: strip HTML tags and return text content
		raw := stripHTMLTags(cellContent)
		raw = strings.TrimSpace(raw)
		raw = decodeHTMLDecimalEntities(raw)
		raw = cleanHexEscapes(raw)
		return raw
	}
	return ""
}

// extractField extracts a value from an HTML element by its ID.
// Looks for patterns like: <td id="Frm_ModelName" ...>VALUE</td>
// Uses (?s) flag so dot matches newlines (values may span multiple lines).
func extractField(html string, id string) string {
	pattern := regexp.MustCompile(fmt.Sprintf(`(?s)id="%s"[^>]*>(.*?)</td>`, id))
	m := pattern.FindStringSubmatch(html)
	if len(m) > 1 {
		raw := m[1]
		raw = strings.TrimSpace(raw)
		raw = decodeHTMLDecimalEntities(raw)
		return raw
	}
	return ""
}

// extractTransferMeanings extracts all Transfer_meaning('key','value') pairs from the HTML.
func extractTransferMeanings(html string) map[string]string {
	data := make(map[string]string)
	matches := reTransferMeaning.FindAllStringSubmatch(html, -1)
	for _, m := range matches {
		if len(m) > 2 {
			data[m[1]] = m[2]
		}
	}
	return data
}

// cleanHexEscapes converts \x2d to -, \x2e to ., etc.
var reHexEscape = regexp.MustCompile(`\\x([0-9a-fA-F]{2})`)

func cleanHexEscapes(s string) string {
	return reHexEscape.ReplaceAllStringFunc(s, func(match string) string {
		m := reHexEscape.FindStringSubmatch(match)
		if len(m) < 2 {
			return match
		}
		val, err := strconv.ParseInt(m[1], 16, 32)
		if err != nil {
			return match
		}
		return string(rune(val))
	})
}

// decodeHTMLDecimalEntities converts &#70;&#54;&#48;&#57; to readable text.
func decodeHTMLDecimalEntities(s string) string {
	return reHTMLDecimal.ReplaceAllStringFunc(s, func(match string) string {
		numStr := match[2 : len(match)-1]
		val, err := strconv.Atoi(numStr)
		if err != nil {
			return match
		}
		return string(rune(val))
	})
}

// stripHTMLTags removes HTML tags from a string.
func stripHTMLTags(s string) string {
	re := regexp.MustCompile(`<[^>]*>`)
	s = re.ReplaceAllString(s, "")
	s = strings.ReplaceAll(s, "&nbsp;", " ")
	s = strings.ReplaceAll(s, "&amp;", "&")
	s = strings.ReplaceAll(s, "&lt;", "<")
	s = strings.ReplaceAll(s, "&gt;", ">")
	return s
}

// extractScriptValue extracts a numeric value from inline JavaScript blocks.
// Handles patterns like: var str_InError = 0;
func extractScriptValue(s string) string {
	// If it's a simple value (no script), return as-is
	if !strings.Contains(s, "var ") && !strings.Contains(s, "<SCRIPT>") {
		return s
	}
	// Try to extract the first numeric assignment
	re := regexp.MustCompile(`var\s+\w+\s*=\s*(\d+)`)
	m := re.FindStringSubmatch(s)
	if len(m) > 1 {
		return m[1]
	}
	return "0"
}

// getAuthType determines the authentication type from beacon/mode values.
func getAuthType(beaconType, wepAuth, wpaAuth, i11Auth string) string {
	switch {
	case beaconType == "None":
		return "Open System"
	case beaconType == "Basic" && wepAuth == "SharedAuthentication":
		return "Shared Key"
	case beaconType == "WPA" && wpaAuth == "PSKAuthentication":
		return "WPA-PSK"
	case beaconType == "11i" && i11Auth == "PSKAuthentication":
		return "WPA2-PSK"
	case beaconType == "WPAand11i" && wpaAuth == "PSKAuthentication" && i11Auth == "PSKAuthentication":
		return "WPA/WPA2-PSK"
	case beaconType == "WPA" && wpaAuth == "EAPAuthentication":
		return "WPA-EAP"
	case beaconType == "11i" && i11Auth == "EAPAuthentication":
		return "WPA2-EAP"
	default:
		return beaconType
	}
}

// getEncryptType determines the encryption type from beacon/type values.
func getEncryptType(beaconType, wpaEnc, i11Enc string) string {
	switch {
	case beaconType == "None":
		return "None"
	case beaconType == "Basic":
		return "WEP"
	case wpaEnc == "TKIPEncryption" || i11Enc == "TKIPEncryption":
		return "TKIP"
	case wpaEnc == "AESEncryption" || i11Enc == "AESEncryption":
		return "AES"
	case wpaEnc == "TKIPandAESEncryption" || i11Enc == "TKIPandAESEncryption":
		return "TKIP+AES"
	default:
		return wpaEnc
	}
}

// ReadResponseBody reads the full response body and returns it as a string.
func ReadResponseBody(resp *http.Response) string {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	// Handle non-UTF8 characters
	if !utf8.Valid(body) {
		body = []byte(strings.ToValidUTF8(string(body), "?"))
	}
	return string(body)
}
