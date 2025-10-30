package models

type NetworkInterface struct {
	Name    string   `json:"name"`
	MAC     string   `json:"mac"`
	IPs     []string `json:"ips"`
	Netmask string   `json:"netmask"`
	Gateway string   `json:"gateway"`
	IsUp    bool     `json:"is_up"`
}
