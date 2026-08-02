package router

import (
	"fmt"
	"net/url"
	"strings"
)

// DDNSSettings holds DDNS configuration.
type DDNSSettings struct {
	Enabled      bool
	ServiceType  string // dipc, dyndns, DtDNS, No-IP
	Server       string
	Username     string
	Password     string
	WANConnection string
	Domain       string
}

// DMZSettings holds DMZ configuration.
type DMZSettings struct {
	Enabled      bool
	WANConnection string
	MACMapping   bool
	HostIP       string
}

// UPnPSettings holds UPnP configuration.
type UPnPSettings struct {
	Enabled           bool
	IPv4WANConnection string
	AdvPeriod         string // minutes
	AdvTTL            string // hops
}

// ============================================================
// DDNS
// ============================================================

// FetchDDNSSettings fetches DDNS configuration.
func FetchDDNSSettings(client *Client) (*DDNSSettings, error) {
	html, err := client.GetPage(PageDDNS)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch DDNS settings: %w", err)
	}
	return parseDDNSSettings(html), nil
}

func parseDDNSSettings(html string) *DDNSSettings {
	data := extractTransferMeanings(html)
	return &DDNSSettings{
		Enabled:       data["Enable"] == "1",
		ServiceType:   cleanHexEscapes(data["ServiceType"]),
		Server:        cleanHexEscapes(data["Server"]),
		Username:      cleanHexEscapes(data["Username"]),
		WANConnection: cleanHexEscapes(data["IFName"]),
		Domain:        cleanHexEscapes(data["Domain"]),
	}
}

// SetDDNSSettings updates DDNS configuration.
func SetDDNSSettings(client *Client, settings DDNSSettings) (bool, error) {
	_, err := client.GetPage(PageDDNS)
	if err != nil {
		return false, fmt.Errorf("failed to load DDNS page: %w", err)
	}

	if client.SessionToken == "" {
		return false, fmt.Errorf("no session token available")
	}

	enableVal := "0"
	if settings.Enabled {
		enableVal = "1"
	}

	formData := url.Values{}
	formData.Set("IF_ACTION", "apply")
	formData.Set("Enable", enableVal)
	formData.Set("ServiceType", settings.ServiceType)
	formData.Set("Server", settings.Server)
	formData.Set("Username", settings.Username)
	formData.Set("Password", settings.Password)
	formData.Set("IFName", settings.WANConnection)
	formData.Set("Domain", settings.Domain)
	formData.Set("_SESSION_TOKEN", client.SessionToken)

	_, err = client.PostAction(PageDDNS, formData.Encode())
	if err != nil {
		return false, fmt.Errorf("failed to update DDNS settings: %w", err)
	}

	return true, nil
}

// ============================================================
// DMZ
// ============================================================

// FetchDMZSettings fetches DMZ configuration.
func FetchDMZSettings(client *Client) (*DMZSettings, error) {
	html, err := client.GetPage(PageDMZ)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch DMZ settings: %w", err)
	}
	return parseDMZSettings(html), nil
}

func parseDMZSettings(html string) *DMZSettings {
	data := extractTransferMeanings(html)
	return &DMZSettings{
		Enabled:       data["Enable"] == "1",
		WANConnection: cleanHexEscapes(data["IFName"]),
		MACMapping:    data["MACMapping"] == "1",
		HostIP:        cleanHexEscapes(data["DMZIPAddress"]),
	}
}

// SetDMZSettings updates DMZ configuration.
func SetDMZSettings(client *Client, settings DMZSettings) (bool, error) {
	_, err := client.GetPage(PageDMZ)
	if err != nil {
		return false, fmt.Errorf("failed to load DMZ page: %w", err)
	}

	if client.SessionToken == "" {
		return false, fmt.Errorf("no session token available")
	}

	enableVal := "0"
	if settings.Enabled {
		enableVal = "1"
	}

	macMapping := "0"
	if settings.MACMapping {
		macMapping = "1"
	}

	formData := url.Values{}
	formData.Set("IF_ACTION", "apply")
	formData.Set("Enable", enableVal)
	formData.Set("IFName", settings.WANConnection)
	formData.Set("MACMapping", macMapping)
	formData.Set("DMZIPAddress", settings.HostIP)
	formData.Set("_SESSION_TOKEN", client.SessionToken)

	_, err = client.PostAction(PageDMZ, formData.Encode())
	if err != nil {
		return false, fmt.Errorf("failed to update DMZ settings: %w", err)
	}

	return true, nil
}

// ============================================================
// UPnP
// ============================================================

// FetchUPnPSettings fetches UPnP configuration.
func FetchUPnPSettings(client *Client) (*UPnPSettings, error) {
	html, err := client.GetPage(PageUPnP)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch UPnP settings: %w", err)
	}
	return parseUPnPSettings(html), nil
}

func parseUPnPSettings(html string) *UPnPSettings {
	data := extractTransferMeanings(html)
	return &UPnPSettings{
		Enabled:           data["Enable"] == "1",
		IPv4WANConnection: cleanHexEscapes(data["IFName"]),
		AdvPeriod:         cleanHexEscapes(data["AdvertisementPeriod"]),
		AdvTTL:            cleanHexEscapes(data["AdvertisementTTL"]),
	}
}

// SetUPnPSettings updates UPnP configuration.
func SetUPnPSettings(client *Client, settings UPnPSettings) (bool, error) {
	_, err := client.GetPage(PageUPnP)
	if err != nil {
		return false, fmt.Errorf("failed to load UPnP page: %w", err)
	}

	if client.SessionToken == "" {
		return false, fmt.Errorf("no session token available")
	}

	enableVal := "0"
	if settings.Enabled {
		enableVal = "1"
	}

	formData := url.Values{}
	formData.Set("IF_ACTION", "apply")
	formData.Set("Enable", enableVal)
	formData.Set("IFName", settings.IPv4WANConnection)
	formData.Set("AdvertisementPeriod", settings.AdvPeriod)
	formData.Set("AdvertisementTTL", settings.AdvTTL)
	formData.Set("_SESSION_TOKEN", client.SessionToken)

	_, err = client.PostAction(PageUPnP, formData.Encode())
	if err != nil {
		return false, fmt.Errorf("failed to update UPnP settings: %w", err)
	}

	return true, nil
}

// ============================================================
// Formatters
// ============================================================

// FormatDDNSSettings formats DDNS settings for display.
func FormatDDNSSettings(s *DDNSSettings) string {
	status := "Disabled"
	if s.Enabled {
		status = "Enabled"
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("  Status:         %s\n", status))
	sb.WriteString(fmt.Sprintf("  Service Type:   %s\n", s.ServiceType))
	sb.WriteString(fmt.Sprintf("  Server:         %s\n", s.Server))
	sb.WriteString(fmt.Sprintf("  Username:       %s\n", s.Username))
	sb.WriteString(fmt.Sprintf("  WAN Connection: %s\n", s.WANConnection))
	sb.WriteString(fmt.Sprintf("  Domain:         %s\n", s.Domain))
	return sb.String()
}

// FormatDMZSettings formats DMZ settings for display.
func FormatDMZSettings(s *DMZSettings) string {
	status := "Disabled"
	if s.Enabled {
		status = "Enabled"
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("  Status:         %s\n", status))
	sb.WriteString(fmt.Sprintf("  WAN Connection: %s\n", s.WANConnection))
	sb.WriteString(fmt.Sprintf("  Host IP:        %s\n", s.HostIP))
	return sb.String()
}

// FormatUPnPSettings formats UPnP settings for display.
func FormatUPnPSettings(s *UPnPSettings) string {
	status := "Disabled"
	if s.Enabled {
		status = "Enabled"
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("  Status:           %s\n", status))
	sb.WriteString(fmt.Sprintf("  WAN Connection:   %s\n", s.IPv4WANConnection))
	sb.WriteString(fmt.Sprintf("  Adv Period:       %s min\n", s.AdvPeriod))
	sb.WriteString(fmt.Sprintf("  Adv TTL:          %s hops\n", s.AdvTTL))
	return sb.String()
}
