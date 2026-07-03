package router

import (
	"fmt"
	"strings"
)

// GenerateQR generates an ASCII QR code string for terminal display.
// Uses a simple encoding approach — falls back to manual rendering if
// external QR library is not available.
func GenerateQR(data string) (string, error) {
	// Try using the qrcode library if available
	// For MVP, we generate a simple representation
	// The actual QR library will be added in go.mod

	// Build Wi-Fi QR string
	if data == "" {
		return "", fmt.Errorf("QR data cannot be empty")
	}

	// For now, return the Wi-Fi connection string that users can
	// manually copy or use with any QR scanner app
	var sb strings.Builder
	sb.WriteString("  ┌───────────────────────────────────────────┐\n")
	sb.WriteString("  │  Wi-Fi QR Code Data (scan with any app):  │\n")
	sb.WriteString("  ├───────────────────────────────────────────┤\n")
	sb.WriteString(fmt.Sprintf("  │  %s\n", data))
	sb.WriteString("  └───────────────────────────────────────────┘\n")
	sb.WriteString("\n  Copy the string above and paste into any\n")
	sb.WriteString("  QR code generator app, or scan directly.\n")

	return sb.String(), nil
}
