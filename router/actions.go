package router

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// FetchStatus fetches and consolidates full system status.
func FetchStatus(client *Client) (*SystemStatus, error) {
	status := &SystemStatus{}

	// Device Info
	html, err := client.GetPage("status_dev_info_t.gch")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch device info: %w", err)
	}
	status.Device = ParseDeviceInfo(html)

	// PON Status
	html, err = client.GetPage("pon_status_link_info_t.gch")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch PON status: %w", err)
	}
	status.PON = ParsePONStatus(html)

	// WAN Status
	html, err = client.GetPage("IPv46_status_wan2_if_t.gch")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch WAN status: %w", err)
	}
	status.WAN = ParseWANStatus(html)

	// WLAN Status
	html, err = client.GetPage("status_wlanm_info1_t.gch")
	if err == nil {
		status.WLANs = ParseWLANStatus(html)
	}

	// LAN Status
	html, err = client.GetPage("pon_status_lan_info_t.gch")
	if err == nil {
		status.LANs = ParseLANStatus(html)
	}

	return status, nil
}

// FetchWLANStatus fetches Wi-Fi interface details.
func FetchWLANStatus(client *Client) ([]WLANInterface, error) {
	html, err := client.GetPage("status_wlanm_info1_t.gch")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch WLAN status: %w", err)
	}
	return ParseWLANStatus(html), nil
}

// FetchLANStatus fetches Ethernet port details.
func FetchLANStatus(client *Client) ([]LANPort, error) {
	html, err := client.GetPage("pon_status_lan_info_t.gch")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch LAN status: %w", err)
	}
	return ParseLANStatus(html), nil
}

// FetchWLANSecurity fetches and decrypts Wi-Fi security settings.
func FetchWLANSecurity(client *Client) (*WLANSecurity, error) {
	sec := &WLANSecurity{}

	// Get security page (has encrypted password)
	html, err := client.GetPage("net_wlanm_secrity1_t.gch")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch WLAN security page: %w", err)
	}

	// Save the session_token BEFORE fetching other pages (each page generates a new token)
	securityToken := client.SessionToken

	data := extractTransferMeanings(html)

	// Get SSID name from the SSID settings page
	htmlSSID, err := client.GetPage("net_wlanm_essid1_t.gch")
	if err == nil {
		ssidData := extractTransferMeanings(htmlSSID)
		if v, ok := ssidData["ESSID"]; ok {
			sec.SSID = cleanHexEscapes(v)
		}
	}

	// Decrypt password using the session_token from the security page
	if encPass, ok := data["KeyPassphrase"]; ok && encPass != "" {
		encPass = cleanHexEscapes(encPass)
		decrypted, err := Decrypt(encPass, securityToken)
		if err != nil {
			sec.Password = "[Decryption Error]"
		} else {
			sec.Password = decrypted
		}
	}

	// Auth & encryption types
	sec.AuthType = getAuthType(
		data["BeaconType"],
		data["WEPAuthMode"],
		data["WPAAuthMode"],
		data["11iAuthMode"],
	)
	sec.EncryptType = getEncryptType(
		data["BeaconType"],
		data["WPAEncryptType"],
		data["11iEncryptType"],
	)
	sec.BeaconType = data["BeaconType"]

	return sec, nil
}

// SetWiFiPassword changes the Wi-Fi password (WPA KeyPassphrase).
func SetWiFiPassword(client *Client, password string) (bool, error) {
	// Load security page to get a fresh session_token
	_, err := client.GetPage("net_wlanm_secrity1_t.gch")
	if err != nil {
		return false, fmt.Errorf("failed to load security page: %w", err)
	}

	if client.SessionToken == "" {
		return false, fmt.Errorf("no session token available")
	}

	// Encrypt password using AES-CBC with session token
	encryptedPass, err := Encrypt(password, client.SessionToken)
	if err != nil {
		return false, fmt.Errorf("failed to encrypt password: %w", err)
	}

	formData := url.Values{}
	formData.Set("IF_ACTION", "apply")
	formData.Set("IF_INDEX", "0")
	formData.Set("KeyPassphrase", encryptedPass)
	formData.Set("_SESSION_TOKEN", client.SessionToken)

	_, err = client.PostAction("net_wlanm_secrity1_t.gch", formData.Encode())
	if err != nil {
		return false, fmt.Errorf("failed to submit password change: %w", err)
	}

	return true, nil
}

// SetWiFiSSID changes the Wi-Fi SSID name.
func SetWiFiSSID(client *Client, ssid string) (bool, error) {
	_, err := client.GetPage("net_wlanm_essid1_t.gch")
	if err != nil {
		return false, fmt.Errorf("failed to load SSID page: %w", err)
	}

	if client.SessionToken == "" {
		return false, fmt.Errorf("no session token available")
	}

	formData := url.Values{}
	formData.Set("IF_ACTION", "apply")
	formData.Set("IF_INDEX", "0")
	formData.Set("ESSID", ssid)
	formData.Set("_SESSION_TOKEN", client.SessionToken)

	_, err = client.PostAction("net_wlanm_essid1_t.gch", formData.Encode())
	if err != nil {
		return false, fmt.Errorf("failed to submit SSID change: %w", err)
	}

	return true, nil
}

// Reboot sends a reboot command to the router.
func Reboot(client *Client) (bool, error) {
	// Load reboot page to get fresh session_token
	_, err := client.GetPage("manager_dev_conf_t.gch")
	if err != nil {
		// Fallback: try fetching any page to get session_token
		_, _ = client.GetPage("status_dev_info_t.gch")
	}

	if client.SessionToken == "" {
		return false, fmt.Errorf("no session token available")
	}

	formData := fmt.Sprintf(
		"IF_ACTION=devrestart&IF_ERRORSTR=SUCC&IF_ERRORPARAM=SUCC&IF_ERRORTYPE=-1&flag=1&_SESSION_TOKEN=%s",
		client.SessionToken,
	)

	body, err := client.PostAction("manager_dev_conf_t.gch", formData)
	if err != nil {
		return false, fmt.Errorf("failed to send reboot command: %w", err)
	}

	if strings.Contains(body, "flag") {
		return true, nil
	}

	return false, nil
}

// FetchTimeout retrieves the current session timeout value in minutes.
func FetchTimeout(client *Client) (int, error) {
	html, err := client.GetPage("manager_login_timeout_t.gch")
	if err != nil {
		return 0, fmt.Errorf("failed to fetch timeout page: %w", err)
	}

	data := extractTransferMeanings(html)
	if v, ok := data["Timeout"]; ok {
		timeout := cleanHexEscapes(v)
		// Parse integer
		val, err := strconv.Atoi(timeout)
		if err != nil {
			return 0, fmt.Errorf("invalid timeout value: %s", timeout)
		}
		return val, nil
	}

	return 0, fmt.Errorf("could not parse timeout value")
}

// SetTimeout changes the session timeout (1-30 minutes).
func SetTimeout(client *Client, minutes int) (bool, error) {
	// Fetch page to get fresh session_token
	_, err := client.GetPage("manager_login_timeout_t.gch")
	if err != nil {
		return false, fmt.Errorf("failed to load timeout page: %w", err)
	}

	if client.SessionToken == "" {
		return false, fmt.Errorf("no session token available")
	}

	formData := fmt.Sprintf(
		"IF_ACTION=apply&Timeout=%d&_SESSION_TOKEN=%s",
		minutes, client.SessionToken,
	)

	_, err = client.PostAction("manager_login_timeout_t.gch", formData)
	if err != nil {
		return false, fmt.Errorf("failed to submit timeout change: %w", err)
	}

	return true, nil
}
