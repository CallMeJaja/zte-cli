package router

// ============================================================
// All .gch page names for ZTE F609 router
// ============================================================

// Status pages
const (
	PageDeviceInfo    = "status_dev_info_t.gch"
	PagePONStatus     = "pon_status_link_info_t.gch"
	PageWANStatus     = "IPv46_status_wan2_if_t.gch"
	PageWLANStatus    = "status_wlanm_info1_t.gch"
	PageLANStatus     = "pon_status_lan_info_t.gch"
	PageVoIPStatus    = "status_voip_4less_t.gch"
)

// Network - WAN pages
const (
	PageWANConfig     = "IPv46_net_wan2_conf_t.gch"
	PagePortBinding   = "net_portbinding_t.gch"
)

// Network - WLAN pages
const (
	PageWLANBasic     = "net_wlanm_conf1_t.gch"
	PageWLANSSID      = "net_wlanm_essid1_t.gch"
	PageWLANSecurity  = "net_wlanm_secrity1_t.gch"
	PageWLANACL       = "net_wlanm_acl1_t.gch"
	PageWLANAssocDev  = "net_wlanm_assoc1_t.gch"
	PageWLANWMM       = "net_wlanm_wmm1_t.gch"
	PageWLANWPS       = "net_wlanm_wps1_t.gch"
)

// Network - LAN pages
const (
	PageDHCPServer    = "net_dhcp_dynamic_t.gch"
	PageDHCPBinding   = "net_dhcp_bind_t.gch"
	PageDHCPPortSvc   = "net_dhcp_portservice_t.gch"
)

// Network - PON pages
const (
	PagePONLOID       = "pon_net_ponloid_t.gch"
	PagePONSN         = "pon_net_sn_t.gch"
)

// Network - Routing pages
const (
	PageRouteDefault  = "net_route_default_t.gch"
	PageStaticRoute   = "net_route_static_t.gch"
	PagePolicyRoute   = "net_route_policy_t.gch"
	PageRouteTable    = "net_route_table_t.gch"
)

// Security pages
const (
	PageFirewall      = "sec_firewall_conf_t.gch"
	PageIPFilter      = "sec_portfilter_conf_t.gch"
	PageMACFilter     = "sec_macfilter_conf_t.gch"
	PageURLFilter     = "sec_url_filter_t.gch"
	PageServiceCtrl   = "sec_service_ctrl_t.gch"
	PageALG           = "sec_alg_conf_t.gch"
)

// Application pages
const (
	PageDDNS          = "app_ddns_conf_t.gch"
	PageDMZ           = "app_dmz_conf_t.gch"
	PageUPnP          = "app_upnp_conf_t.gch"
	PageUPnPMapping   = "app_upnp_portmap_t.gch"
	PagePortForward   = "app_virtual_conf_t.gch"
	PagePortTrigger   = "app_port_trigger_t.gch"
	PageDNSService    = "app_dns_conf_t.gch"
	PageSNTP          = "app_sntp_conf_t.gch"
	PageUSBStorage    = "app_usb_storage_t.gch"
	PageSamba         = "app_samba_conf_t.gch"
	PageFTP           = "app_ftp_conf_t.gch"
)

// Administration pages
const (
	PageTR069         = "net_tr069_basic_t.gch"
	PageUserMgmt      = "manager_aduser_conf_t.gch"
	PageLoginTimeout  = "manager_login_timeout_t.gch"
	PageReboot        = "manager_dev_conf_t.gch"
	PageSWUpgrade     = "manager_dev_version_t.gch"
	PageConfigMgmt    = "manager_dev_conf_backup_t.gch"
	PageLogMgmt       = "manager_log_conf_t.gch"
	PagePing          = "manager_dev_ping_t.gch"
	PageTraceroute    = "manager_dev_traceroute_t.gch"
	PageARPMacTable   = "manager_arp_table_t.gch"
)
