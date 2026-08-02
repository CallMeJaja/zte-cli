package router

import (
	"fmt"
	"strings"
)

// LogEntry represents a system log entry.
type LogEntry struct {
	Timestamp string
	Level     string
	Message   string
}

// FetchLogs fetches system logs from the router.
func FetchLogs(client *Client) ([]LogEntry, error) {
	html, err := client.GetPage(PageLogMgmt)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch logs: %w", err)
	}
	return parseLogs(html), nil
}

// parseLogs extracts log entries from HTML.
func parseLogs(html string) []LogEntry {
	logs := make([]LogEntry, 0)

	// Logs are in the textarea or in Transfer_meaning calls
	// Pattern: P0000-00-00T00:00:00 [Level] Message
	reLogLine := `P(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2})\s+\[(\w+)\]\s+(.+?)(?:\n|$)`
	pattern := compileRegex(reLogLine)
	matches := pattern.FindAllStringSubmatch(html, -1)

	for _, m := range matches {
		if len(m) >= 4 {
			logs = append(logs, LogEntry{
				Timestamp: m[1],
				Level:     m[2],
				Message:   strings.TrimSpace(m[3]),
			})
		}
	}

	return logs
}

// Ping runs a ping test from the router.
func Ping(client *Client, host string) (string, error) {
	_, err := client.GetPage(PagePing)
	if err != nil {
		return "", fmt.Errorf("failed to load ping page: %w", err)
	}

	if client.SessionToken == "" {
		return "", fmt.Errorf("no session token available")
	}

	formData := fmt.Sprintf(
		"IF_ACTION=apply&IPAddress=%s&IF_INDEX=0&_SESSION_TOKEN=%s",
		host, client.SessionToken,
	)

	body, err := client.PostAction(PagePing, formData)
	if err != nil {
		return "", fmt.Errorf("ping request failed: %w", err)
	}

	return parsePingResult(body), nil
}

// parsePingResult extracts ping output from HTML.
func parsePingResult(html string) string {
	// Look for the result textarea or div
	result := extractField(html, "PingResult")
	if result == "" {
		// Try Transfer_meaning
		data := extractTransferMeanings(html)
		if v, ok := data["Result"]; ok {
			result = cleanHexEscapes(v)
		}
	}
	if result == "" {
		result = "No ping result available."
	}
	return result
}

// FetchARPTable fetches the ARP table from the router.
func FetchARPTable(client *Client) (string, error) {
	html, err := client.GetPage(PageARPMacTable)
	if err != nil {
		return "", fmt.Errorf("failed to fetch ARP table: %w", err)
	}

	// Parse ARP table entries
	lines := ParseGenericTable(html)
	if len(lines) == 0 {
		return "  No ARP entries found.", nil
	}

	var sb strings.Builder
	for _, line := range lines {
		sb.WriteString("  " + line + "\n")
	}
	return sb.String(), nil
}

// RebootDevice sends a reboot command.
func RebootDevice(client *Client) (bool, error) {
	return Reboot(client) // Reuse existing function
}

// FactoryReset sends a factory reset command.
func FactoryReset(client *Client) (bool, error) {
	_, err := client.GetPage(PageReboot)
	if err != nil {
		_, _ = client.GetPage(PageDeviceInfo) // fallback to get token
	}

	if client.SessionToken == "" {
		return false, fmt.Errorf("no session token available")
	}

	formData := fmt.Sprintf(
		"IF_ACTION=restoredefault&IF_ERRORSTR=SUCC&IF_ERRORPARAM=SUCC&IF_ERRORTYPE=-1&_SESSION_TOKEN=%s",
		client.SessionToken,
	)

	_, err = client.PostAction(PageReboot, formData)
	if err != nil {
		return false, fmt.Errorf("failed to send factory reset command: %w", err)
	}

	return true, nil
}

// FormatLogs formats log entries for display.
func FormatLogs(logs []LogEntry, limit int) string {
	if len(logs) == 0 {
		return "  No log entries found."
	}

	start := 0
	if limit > 0 && limit < len(logs) {
		start = len(logs) - limit
	}

	var sb strings.Builder
	for _, log := range logs[start:] {
		sb.WriteString(fmt.Sprintf("  [%s] %-12s %s\n", log.Timestamp, log.Level, log.Message))
	}

	return sb.String()
}
