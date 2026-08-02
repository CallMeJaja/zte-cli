package router

import (
	"fmt"
	"net/url"
	"strings"
)

// WANConnection holds WAN connection details.
type WANConnection struct {
	Name          string
	Type          string // Route, Bridge
	ServiceList   string // INTERNET, TR069, VoIP, etc.
	MTU           string
	LinkType      string // PPP, IP
	IPVersion     string // IPv4, IPv6, IPv4/v6
	PPPTransType  string // PPPoE
	Username      string
	Password      string
	AuthType      string // Auto, PAP, CHAP
	ConnTrigger   string // Always On, On Demand, Manual
	NATEnabled    bool
	VLANEnabled   bool
	VLANID        string
	Status        string // Connected, Disconnected
	IPAddress     string
}

// FetchWANConnections fetches list of WAN connections from the status page.
func FetchWANConnections(client *Client) ([]WANConnection, error) {
	html, err := client.GetPage(PageWANStatus)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch WAN status: %w", err)
	}
	return parseWANConnections(html), nil
}

// parseWANConnections extracts WAN connection info from HTML.
func parseWANConnections(html string) []WANConnection {
	connections := make([]WANConnection, 0)
	data := extractTransferMeanings(html)

	// Try to get connection names from the page
	// The WAN status page may have multiple connections
	for i := 0; i < 8; i++ {
		suffix := fmt.Sprintf("%d", i)
		name := cleanHexEscapes(data["ConnectionName"+suffix])
		if name == "" {
			name = cleanHexEscapes(data["IF_WANNAME"+suffix])
		}
		if name == "" {
			continue
		}

		conn := WANConnection{
			Name:        name,
			Status:      cleanHexEscapes(data["ConnectionStatus"+suffix]),
			IPAddress:   cleanHexEscapes(data["IPAddress"+suffix]),
			Type:        cleanHexEscapes(data["ConnectionType"+suffix]),
			ServiceList: cleanHexEscapes(data["ServiceList"+suffix]),
		}
		connections = append(connections, conn)
	}

	// If no indexed connections found, try single connection
	if len(connections) == 0 {
		name := cleanHexEscapes(data["ConnectionName"])
		if name == "" {
			name = cleanHexEscapes(data["IF_WANNAME0"])
		}
		if name != "" {
			connections = append(connections, WANConnection{
				Name:      name,
				Status:    cleanHexEscapes(data["ConnectionStatus"]),
				IPAddress: cleanHexEscapes(data["IPAddress"]),
			})
		}
	}

	return connections
}

// FetchWANConfig fetches WAN connection configuration page.
func FetchWANConfig(client *Client) (string, error) {
	html, err := client.GetPage(PageWANConfig)
	if err != nil {
		return "", fmt.Errorf("failed to fetch WAN config: %w", err)
	}
	return html, nil
}

// DeleteWANConnection deletes a WAN connection by name.
func DeleteWANConnection(client *Client, connName string) (bool, error) {
	_, err := client.GetPage(PageWANConfig)
	if err != nil {
		return false, fmt.Errorf("failed to load WAN config page: %w", err)
	}

	if client.SessionToken == "" {
		return false, fmt.Errorf("no session token available")
	}

	formData := url.Values{}
	formData.Set("IF_ACTION", "delete")
	formData.Set("IF_NAME", connName)
	formData.Set("_SESSION_TOKEN", client.SessionToken)

	_, err = client.PostAction(PageWANConfig, formData.Encode())
	if err != nil {
		return false, fmt.Errorf("failed to delete WAN connection: %w", err)
	}

	return true, nil
}

// CreateWANConnection creates a new WAN connection.
func CreateWANConnection(client *Client, conn WANConnection) (bool, error) {
	_, err := client.GetPage(PageWANConfig)
	if err != nil {
		return false, fmt.Errorf("failed to load WAN config page: %w", err)
	}

	if client.SessionToken == "" {
		return false, fmt.Errorf("no session token available")
	}

	vlanEnable := "0"
	if conn.VLANEnabled {
		vlanEnable = "1"
	}

	natEnable := "0"
	if conn.NATEnabled {
		natEnable = "1"
	}

	formData := url.Values{}
	formData.Set("IF_ACTION", "add")
	formData.Set("IF_NAME", conn.Name)
	formData.Set("ConnectionType", conn.Type)
	formData.Set("ServiceList", conn.ServiceList)
	formData.Set("MTU", conn.MTU)
	formData.Set("LinkType", conn.LinkType)
	formData.Set("IPVersion", conn.IPVersion)
	formData.Set("PPPTransType", conn.PPPTransType)
	formData.Set("Username", conn.Username)
	formData.Set("Password", conn.Password)
	formData.Set("AuthType", conn.AuthType)
	formData.Set("ConnectionTrigger", conn.ConnTrigger)
	formData.Set("NATEnabled", natEnable)
	formData.Set("VLANEnable", vlanEnable)
	formData.Set("VLANID", conn.VLANID)
	formData.Set("_SESSION_TOKEN", client.SessionToken)

	_, err = client.PostAction(PageWANConfig, formData.Encode())
	if err != nil {
		return false, fmt.Errorf("failed to create WAN connection: %w", err)
	}

	return true, nil
}

// FormatWANConnections formats WAN connections for display.
func FormatWANConnections(conns []WANConnection) string {
	if len(conns) == 0 {
		return "  No WAN connections found."
	}

	var sb strings.Builder
	for i, c := range conns {
		sb.WriteString(fmt.Sprintf("  [%d] %s\n", i+1, c.Name))
		if c.Type != "" {
			sb.WriteString(fmt.Sprintf("      Type: %s\n", c.Type))
		}
		if c.ServiceList != "" {
			sb.WriteString(fmt.Sprintf("      Service: %s\n", c.ServiceList))
		}
		if c.IPVersion != "" {
			sb.WriteString(fmt.Sprintf("      IP Version: %s\n", c.IPVersion))
		}
		if c.Status != "" {
			sb.WriteString(fmt.Sprintf("      Status: %s\n", c.Status))
		}
		if c.IPAddress != "" {
			sb.WriteString(fmt.Sprintf("      IP: %s\n", c.IPAddress))
		}
		if i < len(conns)-1 {
			sb.WriteString("\n")
		}
	}
	return sb.String()
}
