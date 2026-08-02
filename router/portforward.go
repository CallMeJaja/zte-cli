package router

import (
	"fmt"
	"net/url"
	"strings"
)

// PortForwardRule represents a single port forwarding rule.
type PortForwardRule struct {
	Enabled       bool
	Name          string
	Protocol      string // TCP, UDP, TCP AND UDP
	WANHostStartIP string
	WANHostEndIP   string
	WANConnection  string
	WANStartPort   string
	WANEndPort     string
	MACMapping     bool
	LANHostIP      string
	LANStartPort   string
	LANEndPort     string
}

// FetchPortForwardRules fetches all port forwarding rules.
func FetchPortForwardRules(client *Client) ([]PortForwardRule, error) {
	html, err := client.GetPage(PagePortForward)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch port forwarding page: %w", err)
	}
	return parsePortForwardRules(html), nil
}

// parsePortForwardRules extracts port forwarding rules from HTML.
func parsePortForwardRules(html string) []PortForwardRule {
	rules := make([]PortForwardRule, 0)
	data := extractTransferMeanings(html)

	// The router stores rules in Transfer_meaning calls with indexed keys
	// Pattern: Enable0, Name0, Protocol0, ... Enable1, Name1, Protocol1, ...
	for i := 0; i < 20; i++ {
		suffix := fmt.Sprintf("%d", i)
		name := cleanHexEscapes(data["Name"+suffix])
		if name == "" {
			continue
		}

		rule := PortForwardRule{
			Enabled:        data["Enable"+suffix] == "1",
			Name:           name,
			Protocol:       cleanHexEscapes(data["Protocol"+suffix]),
			WANHostStartIP: cleanHexEscapes(data["RemoteHostStartIP"+suffix]),
			WANHostEndIP:   cleanHexEscapes(data["RemoteHostEndIP"+suffix]),
			WANConnection:  cleanHexEscapes(data["IFName"+suffix]),
			WANStartPort:   cleanHexEscapes(data["ExternalStartPort"+suffix]),
			WANEndPort:     cleanHexEscapes(data["ExternalEndPort"+suffix]),
			LANHostIP:      cleanHexEscapes(data["InternalClient"+suffix]),
			LANStartPort:   cleanHexEscapes(data["InternalStartPort"+suffix]),
			LANEndPort:     cleanHexEscapes(data["InternalEndPort"+suffix]),
		}
		rules = append(rules, rule)
	}

	return rules
}

// AddPortForwardRule adds a new port forwarding rule.
func AddPortForwardRule(client *Client, rule PortForwardRule) (bool, error) {
	// Fetch page to get fresh session_token
	_, err := client.GetPage(PagePortForward)
	if err != nil {
		return false, fmt.Errorf("failed to load port forwarding page: %w", err)
	}

	if client.SessionToken == "" {
		return false, fmt.Errorf("no session token available")
	}

	enableVal := "0"
	if rule.Enabled {
		enableVal = "1"
	}

	macMapping := "0"
	if rule.MACMapping {
		macMapping = "1"
	}

	formData := url.Values{}
	formData.Set("IF_ACTION", "add")
	formData.Set("Enable", enableVal)
	formData.Set("Name", rule.Name)
	formData.Set("Protocol", rule.Protocol)
	formData.Set("RemoteHostStartIP", rule.WANHostStartIP)
	formData.Set("RemoteHostEndIP", rule.WANHostEndIP)
	formData.Set("IFName", rule.WANConnection)
	formData.Set("ExternalStartPort", rule.WANStartPort)
	formData.Set("ExternalEndPort", rule.WANEndPort)
	formData.Set("InternalClient", rule.LANHostIP)
	formData.Set("InternalStartPort", rule.LANStartPort)
	formData.Set("InternalEndPort", rule.LANEndPort)
	formData.Set("MACMapping", macMapping)
	formData.Set("_SESSION_TOKEN", client.SessionToken)

	_, err = client.PostAction(PagePortForward, formData.Encode())
	if err != nil {
		return false, fmt.Errorf("failed to add port forwarding rule: %w", err)
	}

	return true, nil
}

// DeletePortForwardRule deletes a port forwarding rule by name.
func DeletePortForwardRule(client *Client, ruleName string) (bool, error) {
	// Fetch page to get fresh session_token and find rule index
	html, err := client.GetPage(PagePortForward)
	if err != nil {
		return false, fmt.Errorf("failed to load port forwarding page: %w", err)
	}

	rules := parsePortForwardRules(html)
	ruleIdx := -1
	for i, r := range rules {
		if r.Name == ruleName {
			ruleIdx = i
			break
		}
	}

	if ruleIdx == -1 {
		return false, fmt.Errorf("rule '%s' not found", ruleName)
	}

	if client.SessionToken == "" {
		return false, fmt.Errorf("no session token available")
	}

	formData := url.Values{}
	formData.Set("IF_ACTION", "delete")
	formData.Set("IF_INDEX", fmt.Sprintf("%d", ruleIdx))
	formData.Set("_SESSION_TOKEN", client.SessionToken)

	_, err = client.PostAction(PagePortForward, formData.Encode())
	if err != nil {
		return false, fmt.Errorf("failed to delete port forwarding rule: %w", err)
	}

	return true, nil
}

// SetPortForwardEnable enables or disables a port forwarding rule.
func SetPortForwardEnable(client *Client, ruleName string, enable bool) (bool, error) {
	html, err := client.GetPage(PagePortForward)
	if err != nil {
		return false, fmt.Errorf("failed to load port forwarding page: %w", err)
	}

	rules := parsePortForwardRules(html)
	ruleIdx := -1
	for i, r := range rules {
		if r.Name == ruleName {
			ruleIdx = i
			break
		}
	}

	if ruleIdx == -1 {
		return false, fmt.Errorf("rule '%s' not found", ruleName)
	}

	if client.SessionToken == "" {
		return false, fmt.Errorf("no session token available")
	}

	enableVal := "0"
	if enable {
		enableVal = "1"
	}

	formData := url.Values{}
	formData.Set("IF_ACTION", "apply")
	formData.Set("IF_INDEX", fmt.Sprintf("%d", ruleIdx))
	formData.Set("Enable", enableVal)
	formData.Set("_SESSION_TOKEN", client.SessionToken)

	_, err = client.PostAction(PagePortForward, formData.Encode())
	if err != nil {
		return false, fmt.Errorf("failed to update port forwarding rule: %w", err)
	}

	return true, nil
}

// FormatPortForwardRules formats port forwarding rules for display.
func FormatPortForwardRules(rules []PortForwardRule) string {
	if len(rules) == 0 {
		return "  No port forwarding rules configured."
	}

	var sb strings.Builder
	for i, r := range rules {
		status := "OFF"
		if r.Enabled {
			status = "ON"
		}
		sb.WriteString(fmt.Sprintf("  [%d] %s (%s)\n", i+1, r.Name, status))
		sb.WriteString(fmt.Sprintf("      Protocol: %s\n", r.Protocol))
		sb.WriteString(fmt.Sprintf("      WAN: %s:%s-%s\n", r.WANConnection, r.WANStartPort, r.WANEndPort))
		sb.WriteString(fmt.Sprintf("      LAN: %s:%s-%s\n", r.LANHostIP, r.LANStartPort, r.LANEndPort))
		if i < len(rules)-1 {
			sb.WriteString("\n")
		}
	}
	return sb.String()
}
