package router

import (
	"fmt"
	"net/url"
	"strings"
)

// WiFiBasicSettings holds WiFi basic configuration.
type WiFiBasicSettings struct {
	RFMode        string // Enabled, Disabled, Scheduled
	Mode          string // 802.11 mode
	Country       string
	BandWidth     string // 20MHz, 40MHz, Auto
	Channel       string // Auto, 1-13
	SGIEnabled    bool
	BeaconInterval string
	TransmitPower string // 100%, 80%, 60%, 40%, 20%
	QoSType       string // Disabled, WMM, SSID
	RTSThreshold  string
	DTIMInterval  string
}

// FetchWiFiBasicSettings fetches WiFi basic configuration.
func FetchWiFiBasicSettings(client *Client) (*WiFiBasicSettings, error) {
	html, err := client.GetPage(PageWLANBasic)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch WiFi basic settings: %w", err)
	}
	return parseWiFiBasicSettings(html), nil
}

// parseWiFiBasicSettings extracts WiFi basic settings from HTML.
func parseWiFiBasicSettings(html string) *WiFiBasicSettings {
	data := extractTransferMeanings(html)
	settings := &WiFiBasicSettings{
		RFMode:         cleanHexEscapes(data["RadioStatus"]),
		Mode:           cleanHexEscapes(data["11nMode"]),
		Country:        cleanHexEscapes(data["Country"]),
		BandWidth:      cleanHexEscapes(data["BandWidth"]),
		Channel:        cleanHexEscapes(data["Channel"]),
		SGIEnabled:     data["SGI"] == "1",
		BeaconInterval: cleanHexEscapes(data["BeaconInterval"]),
		TransmitPower:  cleanHexEscapes(data["TransmitPower"]),
		QoSType:        cleanHexEscapes(data["QoSType"]),
		RTSThreshold:   cleanHexEscapes(data["RTSThreshold"]),
		DTIMInterval:   cleanHexEscapes(data["DTIMInterval"]),
	}

	// Map RF mode value
	switch settings.RFMode {
	case "1":
		settings.RFMode = "Enabled"
	case "0":
		settings.RFMode = "Disabled"
	case "2":
		settings.RFMode = "Scheduled"
	}

	// Map transmit power
	switch settings.TransmitPower {
	case "1":
		settings.TransmitPower = "100%"
	case "2":
		settings.TransmitPower = "80%"
	case "3":
		settings.TransmitPower = "60%"
	case "4":
		settings.TransmitPower = "40%"
	case "5":
		settings.TransmitPower = "20%"
	}

	return settings
}

// SetWiFiBasicSettings updates WiFi basic configuration.
func SetWiFiBasicSettings(client *Client, settings WiFiBasicSettings) (bool, error) {
	_, err := client.GetPage(PageWLANBasic)
	if err != nil {
		return false, fmt.Errorf("failed to load WiFi basic page: %w", err)
	}

	if client.SessionToken == "" {
		return false, fmt.Errorf("no session token available")
	}

	// Map RF mode
	rfMode := "1" // Enabled
	switch settings.RFMode {
	case "Disabled":
		rfMode = "0"
	case "Scheduled":
		rfMode = "2"
	}

	// Map transmit power
	power := "1" // 100%
	switch settings.TransmitPower {
	case "80%":
		power = "2"
	case "60%":
		power = "3"
	case "40%":
		power = "4"
	case "20%":
		power = "5"
	}

	sgi := "0"
	if settings.SGIEnabled {
		sgi = "1"
	}

	formData := url.Values{}
	formData.Set("IF_ACTION", "apply")
	formData.Set("RadioStatus", rfMode)
	formData.Set("11nMode", settings.Mode)
	formData.Set("Country", settings.Country)
	formData.Set("BandWidth", settings.BandWidth)
	formData.Set("Channel", settings.Channel)
	formData.Set("SGI", sgi)
	formData.Set("BeaconInterval", settings.BeaconInterval)
	formData.Set("TransmitPower", power)
	formData.Set("QoSType", settings.QoSType)
	formData.Set("RTSThreshold", settings.RTSThreshold)
	formData.Set("DTIMInterval", settings.DTIMInterval)
	formData.Set("_SESSION_TOKEN", client.SessionToken)

	_, err = client.PostAction(PageWLANBasic, formData.Encode())
	if err != nil {
		return false, fmt.Errorf("failed to update WiFi basic settings: %w", err)
	}

	return true, nil
}

// SetWiFiChannel sets the WiFi channel.
func SetWiFiChannel(client *Client, channel string) (bool, error) {
	settings, err := FetchWiFiBasicSettings(client)
	if err != nil {
		return false, err
	}
	settings.Channel = channel
	return SetWiFiBasicSettings(client, *settings)
}

// SetWiFiTransmitPower sets the WiFi transmit power.
func SetWiFiTransmitPower(client *Client, power string) (bool, error) {
	settings, err := FetchWiFiBasicSettings(client)
	if err != nil {
		return false, err
	}
	settings.TransmitPower = power
	return SetWiFiBasicSettings(client, *settings)
}

// SetWiFiRFMode enables or disables the WiFi radio.
func SetWiFiRFMode(client *Client, enabled bool) (bool, error) {
	settings, err := FetchWiFiBasicSettings(client)
	if err != nil {
		return false, err
	}
	if enabled {
		settings.RFMode = "Enabled"
	} else {
		settings.RFMode = "Disabled"
	}
	return SetWiFiBasicSettings(client, *settings)
}

// FormatWiFiBasicSettings formats WiFi basic settings for display.
func FormatWiFiBasicSettings(s *WiFiBasicSettings) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("  RF Mode:          %s\n", s.RFMode))
	sb.WriteString(fmt.Sprintf("  Mode:             %s\n", s.Mode))
	sb.WriteString(fmt.Sprintf("  Country:          %s\n", s.Country))
	sb.WriteString(fmt.Sprintf("  Bandwidth:        %s\n", s.BandWidth))
	sb.WriteString(fmt.Sprintf("  Channel:          %s\n", s.Channel))
	sgi := "Disabled"
	if s.SGIEnabled {
		sgi = "Enabled"
	}
	sb.WriteString(fmt.Sprintf("  SGI:              %s\n", sgi))
	sb.WriteString(fmt.Sprintf("  Beacon Interval:  %s ms\n", s.BeaconInterval))
	sb.WriteString(fmt.Sprintf("  Transmit Power:   %s\n", s.TransmitPower))
	sb.WriteString(fmt.Sprintf("  QoS Type:         %s\n", s.QoSType))
	sb.WriteString(fmt.Sprintf("  RTS Threshold:    %s\n", s.RTSThreshold))
	sb.WriteString(fmt.Sprintf("  DTIM Interval:    %s\n", s.DTIMInterval))
	return sb.String()
}
