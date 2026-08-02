package router

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

// DHCPSettings holds DHCP server configuration.
type DHCPSettings struct {
	Enabled       bool
	LANIP         string
	SubnetMask    string
	StartIP       string
	EndIP         string
	DNSServer1    string
	DNSServer2    string
	DNSServer3    string
	DefaultGW     string
	LeaseTime     string
	AssignISP     bool
}

// DHCPLease represents an active DHCP lease.
type DHCPLease struct {
	MACAddress    string
	IPAddress     string
	LeaseTime     string
	HostName      string
	Port          string
}

// FetchDHCPSettings fetches DHCP server configuration.
func FetchDHCPSettings(client *Client) (*DHCPSettings, error) {
	html, err := client.GetPage(PageDHCPServer)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch DHCP settings: %w", err)
	}
	return parseDHCPSettings(html), nil
}

// parseDHCPSettings extracts DHCP settings from HTML.
func parseDHCPSettings(html string) *DHCPSettings {
	data := extractTransferMeanings(html)
	settings := &DHCPSettings{
		Enabled:    data["DHCPEnable"] == "1",
		LANIP:      cleanHexEscapes(data["IPAddress"]),
		SubnetMask: cleanHexEscapes(data["SubnetMask"]),
		StartIP:    cleanHexEscapes(data["MinAddress"]),
		EndIP:      cleanHexEscapes(data["MaxAddress"]),
		DNSServer1: cleanHexEscapes(data["DNSServers"]),
		DefaultGW:  cleanHexEscapes(data["DefaultGateway"]),
		LeaseTime:  cleanHexEscapes(data["LeaseTime"]),
	}

	// Parse DNS servers (may be comma-separated)
	dns := cleanHexEscapes(data["DNSServers"])
	if dns != "" {
		parts := strings.Split(dns, ",")
		if len(parts) > 0 {
			settings.DNSServer1 = parts[0]
		}
		if len(parts) > 1 {
			settings.DNSServer2 = parts[1]
		}
		if len(parts) > 2 {
			settings.DNSServer3 = parts[2]
		}
	}

	return settings
}

// FetchDHCPLeases fetches active DHCP leases from the DHCP page.
func FetchDHCPLeases(client *Client) ([]DHCPLease, error) {
	html, err := client.GetPage(PageDHCPServer)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch DHCP leases: %w", err)
	}
	return parseDHCPLeases(html), nil
}

// parseDHCPLeases extracts DHCP lease entries from HTML tables.
func parseDHCPLeases(html string) []DHCPLease {
	leases := make([]DHCPLease, 0)

	// Look for the allocated address table
	// Pattern: rows with MAC, IP, Lease Time, Hostname, Port
	reLeaseRow := `(?s)<td[^>]*>([0-9a-fA-F:]+)</td>\s*<td[^>]*>([0-9.]+)</td>\s*<td[^>]*>([^<]*)</td>`
	pattern := compileRegex(reLeaseRow)
	matches := pattern.FindAllStringSubmatch(html, -1)

	for _, m := range matches {
		if len(m) >= 4 {
			lease := DHCPLease{
				MACAddress: strings.TrimSpace(m[1]),
				IPAddress:  strings.TrimSpace(m[2]),
				LeaseTime:  strings.TrimSpace(m[3]),
			}
			leases = append(leases, lease)
		}
	}

	return leases
}

// SetDHCPSettings updates DHCP server configuration.
func SetDHCPSettings(client *Client, settings DHCPSettings) (bool, error) {
	_, err := client.GetPage(PageDHCPServer)
	if err != nil {
		return false, fmt.Errorf("failed to load DHCP page: %w", err)
	}

	if client.SessionToken == "" {
		return false, fmt.Errorf("no session token available")
	}

	enableVal := "0"
	if settings.Enabled {
		enableVal = "1"
	}

	assignISP := "0"
	if settings.AssignISP {
		assignISP = "1"
	}

	formData := url.Values{}
	formData.Set("IF_ACTION", "apply")
	formData.Set("DHCPEnable", enableVal)
	formData.Set("IPAddress", settings.LANIP)
	formData.Set("SubnetMask", settings.SubnetMask)
	formData.Set("MinAddress", settings.StartIP)
	formData.Set("MaxAddress", settings.EndIP)
	formData.Set("DefaultGateway", settings.DefaultGW)
	formData.Set("LeaseTime", settings.LeaseTime)
	formData.Set("AssignISP", assignISP)

	if settings.DNSServer1 != "" {
		formData.Set("DNSServers", settings.DNSServer1)
	}
	if settings.DNSServer2 != "" {
		formData.Set("DNSserver2", settings.DNSServer2)
	}
	if settings.DNSServer3 != "" {
		formData.Set("DNSserver3", settings.DNSServer3)
	}

	formData.Set("_SESSION_TOKEN", client.SessionToken)

	_, err = client.PostAction(PageDHCPServer, formData.Encode())
	if err != nil {
		return false, fmt.Errorf("failed to update DHCP settings: %w", err)
	}

	return true, nil
}

// FormatDHCPSettings formats DHCP settings for display.
func FormatDHCPSettings(s *DHCPSettings) string {
	var sb strings.Builder
	status := "Disabled"
	if s.Enabled {
		status = "Enabled"
	}

	sb.WriteString(fmt.Sprintf("  DHCP Server:    %s\n", status))
	sb.WriteString(fmt.Sprintf("  LAN IP:         %s\n", s.LANIP))
	sb.WriteString(fmt.Sprintf("  Subnet Mask:    %s\n", s.SubnetMask))
	sb.WriteString(fmt.Sprintf("  Start IP:       %s\n", s.StartIP))
	sb.WriteString(fmt.Sprintf("  End IP:         %s\n", s.EndIP))
	sb.WriteString(fmt.Sprintf("  Default GW:     %s\n", s.DefaultGW))
	sb.WriteString(fmt.Sprintf("  DNS Server 1:   %s\n", s.DNSServer1))
	if s.DNSServer2 != "" {
		sb.WriteString(fmt.Sprintf("  DNS Server 2:   %s\n", s.DNSServer2))
	}
	if s.DNSServer3 != "" {
		sb.WriteString(fmt.Sprintf("  DNS Server 3:   %s\n", s.DNSServer3))
	}
	sb.WriteString(fmt.Sprintf("  Lease Time:     %s sec\n", s.LeaseTime))

	return sb.String()
}

// FormatDHCPLeases formats DHCP leases for display.
func FormatDHCPLeases(leases []DHCPLease) string {
	if len(leases) == 0 {
		return "  No active DHCP leases."
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("  %-18s %-16s %-12s %s\n", "MAC Address", "IP Address", "Lease", "Hostname"))
	sb.WriteString(fmt.Sprintf("  %-18s %-16s %-12s %s\n", "------------------", "----------------", "------------", "--------"))

	for _, l := range leases {
		hostname := l.HostName
		if hostname == "" {
			hostname = "-"
		}
		sb.WriteString(fmt.Sprintf("  %-18s %-16s %-12s %s\n", l.MACAddress, l.IPAddress, l.LeaseTime, hostname))
	}

	return sb.String()
}

// compileRegex is a helper to compile regex.
func compileRegex(pattern string) *regexp.Regexp {
	return regexp.MustCompile(pattern)
}
