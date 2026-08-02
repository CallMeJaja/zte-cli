package router

import (
	"fmt"
	"net/url"
	"strings"
)

// FirewallSettings holds firewall configuration.
type FirewallSettings struct {
	AntiHacking bool
	Level       string // Off, Low, Middle, High
}

// IPFilterRule represents an IP filter rule.
type IPFilterRule struct {
	Enabled    bool
	Name       string
	Protocol   string // TCP, UDP, TCP AND UDP, ICMP, ANY
	SrcIPStart string
	SrcIPEnd   string
	DstIPStart string
	DstIPEnd   string
	SrcPortStart string
	SrcPortEnd   string
	DstPortStart string
	DstPortEnd   string
	Ingress    string // LAN, WAN connection name
	Egress     string // LAN, WAN connection name
	Mode       string // Discard, Permit
}

// MACFilterRule represents a MAC filter rule.
type MACFilterRule struct {
	Enabled  bool
	Mode     string // Discard
	Type     string // Bridge, Route, Bridge+Route
	Protocol string // IP, ARP, RARP, PPPoE, ALL
	SrcMAC   string
	DstMAC   string
}

// URLFilterRule represents a URL filter rule.
type URLFilterRule struct {
	Enabled bool
	Mode    string // Discard, Permit
	URL     string
}

// ============================================================
// Firewall
// ============================================================

// FetchFirewallSettings fetches firewall configuration.
func FetchFirewallSettings(client *Client) (*FirewallSettings, error) {
	html, err := client.GetPage(PageFirewall)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch firewall settings: %w", err)
	}
	return parseFirewallSettings(html), nil
}

func parseFirewallSettings(html string) *FirewallSettings {
	data := extractTransferMeanings(html)
	level := "Off"
	switch data["FirewallLevel"] {
	case "1":
		level = "Low"
	case "2":
		level = "Middle"
	case "3":
		level = "High"
	}
	return &FirewallSettings{
		AntiHacking: data["AntiHacking"] == "1",
		Level:       level,
	}
}

// SetFirewallSettings updates firewall configuration.
func SetFirewallSettings(client *Client, settings FirewallSettings) (bool, error) {
	_, err := client.GetPage(PageFirewall)
	if err != nil {
		return false, fmt.Errorf("failed to load firewall page: %w", err)
	}

	if client.SessionToken == "" {
		return false, fmt.Errorf("no session token available")
	}

	antiHacking := "0"
	if settings.AntiHacking {
		antiHacking = "1"
	}

	level := "0"
	switch settings.Level {
	case "Low":
		level = "1"
	case "Middle":
		level = "2"
	case "High":
		level = "3"
	}

	formData := url.Values{}
	formData.Set("IF_ACTION", "apply")
	formData.Set("AntiHacking", antiHacking)
	formData.Set("FirewallLevel", level)
	formData.Set("_SESSION_TOKEN", client.SessionToken)

	_, err = client.PostAction(PageFirewall, formData.Encode())
	if err != nil {
		return false, fmt.Errorf("failed to update firewall settings: %w", err)
	}

	return true, nil
}

// ============================================================
// IP Filter
// ============================================================

// FetchIPFilterRules fetches all IP filter rules.
func FetchIPFilterRules(client *Client) ([]IPFilterRule, error) {
	html, err := client.GetPage(PageIPFilter)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch IP filter rules: %w", err)
	}
	return parseIPFilterRules(html), nil
}

func parseIPFilterRules(html string) []IPFilterRule {
	rules := make([]IPFilterRule, 0)
	data := extractTransferMeanings(html)

	for i := 0; i < 20; i++ {
		suffix := fmt.Sprintf("%d", i)
		name := cleanHexEscapes(data["Name"+suffix])
		if name == "" {
			continue
		}

		rule := IPFilterRule{
			Enabled:      data["Enable"+suffix] == "1",
			Name:         name,
			Protocol:     cleanHexEscapes(data["Protocol"+suffix]),
			SrcIPStart:   cleanHexEscapes(data["SrcIPStart"+suffix]),
			SrcIPEnd:     cleanHexEscapes(data["SrcIPEnd"+suffix]),
			DstIPStart:   cleanHexEscapes(data["DstIPStart"+suffix]),
			DstIPEnd:     cleanHexEscapes(data["DstIPEnd"+suffix]),
			SrcPortStart: cleanHexEscapes(data["SrcPortStart"+suffix]),
			SrcPortEnd:   cleanHexEscapes(data["SrcPortEnd"+suffix]),
			DstPortStart: cleanHexEscapes(data["DstPortStart"+suffix]),
			DstPortEnd:   cleanHexEscapes(data["DstPortEnd"+suffix]),
			Ingress:      cleanHexEscapes(data["Ingress"+suffix]),
			Egress:       cleanHexEscapes(data["Egress"+suffix]),
			Mode:         cleanHexEscapes(data["Action"+suffix]),
		}
		rules = append(rules, rule)
	}

	return rules
}

// AddIPFilterRule adds a new IP filter rule.
func AddIPFilterRule(client *Client, rule IPFilterRule) (bool, error) {
	_, err := client.GetPage(PageIPFilter)
	if err != nil {
		return false, fmt.Errorf("failed to load IP filter page: %w", err)
	}

	if client.SessionToken == "" {
		return false, fmt.Errorf("no session token available")
	}

	enableVal := "0"
	if rule.Enabled {
		enableVal = "1"
	}

	formData := url.Values{}
	formData.Set("IF_ACTION", "add")
	formData.Set("Enable", enableVal)
	formData.Set("Name", rule.Name)
	formData.Set("Protocol", rule.Protocol)
	formData.Set("SrcIPStart", rule.SrcIPStart)
	formData.Set("SrcIPEnd", rule.SrcIPEnd)
	formData.Set("DstIPStart", rule.DstIPStart)
	formData.Set("DstIPEnd", rule.DstIPEnd)
	formData.Set("SrcPortStart", rule.SrcPortStart)
	formData.Set("SrcPortEnd", rule.SrcPortEnd)
	formData.Set("DstPortStart", rule.DstPortStart)
	formData.Set("DstPortEnd", rule.DstPortEnd)
	formData.Set("Ingress", rule.Ingress)
	formData.Set("Egress", rule.Egress)
	formData.Set("Action", rule.Mode)
	formData.Set("_SESSION_TOKEN", client.SessionToken)

	_, err = client.PostAction(PageIPFilter, formData.Encode())
	if err != nil {
		return false, fmt.Errorf("failed to add IP filter rule: %w", err)
	}

	return true, nil
}

// DeleteIPFilterRule deletes an IP filter rule by name.
func DeleteIPFilterRule(client *Client, ruleName string) (bool, error) {
	html, err := client.GetPage(PageIPFilter)
	if err != nil {
		return false, fmt.Errorf("failed to load IP filter page: %w", err)
	}

	rules := parseIPFilterRules(html)
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

	_, err = client.PostAction(PageIPFilter, formData.Encode())
	if err != nil {
		return false, fmt.Errorf("failed to delete IP filter rule: %w", err)
	}

	return true, nil
}

// ============================================================
// MAC Filter
// ============================================================

// FetchMACFilterRules fetches all MAC filter rules.
func FetchMACFilterRules(client *Client) ([]MACFilterRule, error) {
	html, err := client.GetPage(PageMACFilter)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch MAC filter rules: %w", err)
	}
	return parseMACFilterRules(html), nil
}

func parseMACFilterRules(html string) []MACFilterRule {
	rules := make([]MACFilterRule, 0)
	data := extractTransferMeanings(html)

	for i := 0; i < 20; i++ {
		suffix := fmt.Sprintf("%d", i)
		srcMAC := cleanHexEscapes(data["SrcMAC"+suffix])
		if srcMAC == "" {
			continue
		}

		rule := MACFilterRule{
			Enabled:  data["Enable"+suffix] == "1",
			Mode:     cleanHexEscapes(data["Action"+suffix]),
			Type:     cleanHexEscapes(data["Type"+suffix]),
			Protocol: cleanHexEscapes(data["Protocol"+suffix]),
			SrcMAC:   srcMAC,
			DstMAC:   cleanHexEscapes(data["DstMAC"+suffix]),
		}
		rules = append(rules, rule)
	}

	return rules
}

// AddMACFilterRule adds a new MAC filter rule.
func AddMACFilterRule(client *Client, rule MACFilterRule) (bool, error) {
	_, err := client.GetPage(PageMACFilter)
	if err != nil {
		return false, fmt.Errorf("failed to load MAC filter page: %w", err)
	}

	if client.SessionToken == "" {
		return false, fmt.Errorf("no session token available")
	}

	enableVal := "0"
	if rule.Enabled {
		enableVal = "1"
	}

	formData := url.Values{}
	formData.Set("IF_ACTION", "add")
	formData.Set("Enable", enableVal)
	formData.Set("Action", rule.Mode)
	formData.Set("Type", rule.Type)
	formData.Set("Protocol", rule.Protocol)
	formData.Set("SrcMAC", rule.SrcMAC)
	formData.Set("DstMAC", rule.DstMAC)
	formData.Set("_SESSION_TOKEN", client.SessionToken)

	_, err = client.PostAction(PageMACFilter, formData.Encode())
	if err != nil {
		return false, fmt.Errorf("failed to add MAC filter rule: %w", err)
	}

	return true, nil
}

// DeleteMACFilterRule deletes a MAC filter rule by source MAC.
func DeleteMACFilterRule(client *Client, srcMAC string) (bool, error) {
	html, err := client.GetPage(PageMACFilter)
	if err != nil {
		return false, fmt.Errorf("failed to load MAC filter page: %w", err)
	}

	rules := parseMACFilterRules(html)
	ruleIdx := -1
	for i, r := range rules {
		if r.SrcMAC == srcMAC {
			ruleIdx = i
			break
		}
	}

	if ruleIdx == -1 {
		return false, fmt.Errorf("rule with source MAC '%s' not found", srcMAC)
	}

	if client.SessionToken == "" {
		return false, fmt.Errorf("no session token available")
	}

	formData := url.Values{}
	formData.Set("IF_ACTION", "delete")
	formData.Set("IF_INDEX", fmt.Sprintf("%d", ruleIdx))
	formData.Set("_SESSION_TOKEN", client.SessionToken)

	_, err = client.PostAction(PageMACFilter, formData.Encode())
	if err != nil {
		return false, fmt.Errorf("failed to delete MAC filter rule: %w", err)
	}

	return true, nil
}

// ============================================================
// URL Filter
// ============================================================

// FetchURLFilterRules fetches all URL filter rules.
func FetchURLFilterRules(client *Client) ([]URLFilterRule, error) {
	html, err := client.GetPage(PageURLFilter)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch URL filter rules: %w", err)
	}
	return parseURLFilterRules(html), nil
}

func parseURLFilterRules(html string) []URLFilterRule {
	rules := make([]URLFilterRule, 0)
	data := extractTransferMeanings(html)

	// URL filter stores URLs in a list
	for i := 0; i < 20; i++ {
		suffix := fmt.Sprintf("%d", i)
		urlVal := cleanHexEscapes(data["URL"+suffix])
		if urlVal == "" {
			continue
		}

		rule := URLFilterRule{
			Enabled: data["Enable"+suffix] == "1",
			Mode:    cleanHexEscapes(data["Action"+suffix]),
			URL:     urlVal,
		}
		rules = append(rules, rule)
	}

	return rules
}

// AddURLFilterRule adds a new URL filter rule.
func AddURLFilterRule(client *Client, rule URLFilterRule) (bool, error) {
	_, err := client.GetPage(PageURLFilter)
	if err != nil {
		return false, fmt.Errorf("failed to load URL filter page: %w", err)
	}

	if client.SessionToken == "" {
		return false, fmt.Errorf("no session token available")
	}

	enableVal := "0"
	if rule.Enabled {
		enableVal = "1"
	}

	formData := url.Values{}
	formData.Set("IF_ACTION", "add")
	formData.Set("Enable", enableVal)
	formData.Set("Action", rule.Mode)
	formData.Set("URL", rule.URL)
	formData.Set("_SESSION_TOKEN", client.SessionToken)

	_, err = client.PostAction(PageURLFilter, formData.Encode())
	if err != nil {
		return false, fmt.Errorf("failed to add URL filter rule: %w", err)
	}

	return true, nil
}

// DeleteURLFilterRule deletes a URL filter rule by URL.
func DeleteURLFilterRule(client *Client, urlStr string) (bool, error) {
	html, err := client.GetPage(PageURLFilter)
	if err != nil {
		return false, fmt.Errorf("failed to load URL filter page: %w", err)
	}

	rules := parseURLFilterRules(html)
	ruleIdx := -1
	for i, r := range rules {
		if r.URL == urlStr {
			ruleIdx = i
			break
		}
	}

	if ruleIdx == -1 {
		return false, fmt.Errorf("URL filter rule '%s' not found", urlStr)
	}

	if client.SessionToken == "" {
		return false, fmt.Errorf("no session token available")
	}

	formData := url.Values{}
	formData.Set("IF_ACTION", "delete")
	formData.Set("IF_INDEX", fmt.Sprintf("%d", ruleIdx))
	formData.Set("_SESSION_TOKEN", client.SessionToken)

	_, err = client.PostAction(PageURLFilter, formData.Encode())
	if err != nil {
		return false, fmt.Errorf("failed to delete URL filter rule: %w", err)
	}

	return true, nil
}

// ============================================================
// Formatters
// ============================================================

// FormatFirewallSettings formats firewall settings for display.
func FormatFirewallSettings(s *FirewallSettings) string {
	antiHacking := "Disabled"
	if s.AntiHacking {
		antiHacking = "Enabled"
	}
	return fmt.Sprintf("  Anti-Hacking: %s\n  Firewall Level: %s", antiHacking, s.Level)
}

// FormatIPFilterRules formats IP filter rules for display.
func FormatIPFilterRules(rules []IPFilterRule) string {
	if len(rules) == 0 {
		return "  No IP filter rules configured."
	}

	var sb strings.Builder
	for i, r := range rules {
		status := "OFF"
		if r.Enabled {
			status = "ON"
		}
		sb.WriteString(fmt.Sprintf("  [%d] %s (%s) - %s\n", i+1, r.Name, status, r.Mode))
		sb.WriteString(fmt.Sprintf("      Protocol: %s\n", r.Protocol))
		sb.WriteString(fmt.Sprintf("      Src: %s:%s → %s:%s\n", r.SrcIPStart, r.SrcPortStart, r.SrcIPEnd, r.SrcPortEnd))
		sb.WriteString(fmt.Sprintf("      Dst: %s:%s → %s:%s\n", r.DstIPStart, r.DstPortStart, r.DstIPEnd, r.DstPortEnd))
		sb.WriteString(fmt.Sprintf("      Ingress: %s, Egress: %s\n", r.Ingress, r.Egress))
		if i < len(rules)-1 {
			sb.WriteString("\n")
		}
	}
	return sb.String()
}

// FormatMACFilterRules formats MAC filter rules for display.
func FormatMACFilterRules(rules []MACFilterRule) string {
	if len(rules) == 0 {
		return "  No MAC filter rules configured."
	}

	var sb strings.Builder
	for i, r := range rules {
		status := "OFF"
		if r.Enabled {
			status = "ON"
		}
		sb.WriteString(fmt.Sprintf("  [%d] %s - %s (%s)\n", i+1, r.SrcMAC, r.Mode, status))
		sb.WriteString(fmt.Sprintf("      Type: %s, Protocol: %s\n", r.Type, r.Protocol))
		if r.DstMAC != "" {
			sb.WriteString(fmt.Sprintf("      Dst MAC: %s\n", r.DstMAC))
		}
		if i < len(rules)-1 {
			sb.WriteString("\n")
		}
	}
	return sb.String()
}

// FormatURLFilterRules formats URL filter rules for display.
func FormatURLFilterRules(rules []URLFilterRule) string {
	if len(rules) == 0 {
		return "  No URL filter rules configured."
	}

	var sb strings.Builder
	for i, r := range rules {
		status := "OFF"
		if r.Enabled {
			status = "ON"
		}
		sb.WriteString(fmt.Sprintf("  [%d] %s (%s) - %s\n", i+1, r.URL, status, r.Mode))
		if i < len(rules)-1 {
			sb.WriteString("\n")
		}
	}
	return sb.String()
}
