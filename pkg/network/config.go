package network

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
	"regexp"
)

type NetworkType string

const (
	NetworkTypeHost   NetworkType = "host"
	NetworkTypeBridge NetworkType = "bridge"
)

type NetworkConfig struct {
	Nameservers []net.IP `json:"nameservers"`
	Hostname    string   `json:"hostname"`

	BridgeConfig  *BridgeConfig `json:"bridge_config"`
	Type          NetworkType   `json:"type"`
}

type BridgeConfig struct {
	BridgeName string    `json:"bridge_name"`
	BridgeNet  net.IPNet `json:"bridge_net"`

	VethName     string    `json:"veth_name"`
	ContainerNet net.IPNet `json:"container_net"`
}

func ValidateConfig(nc *NetworkConfig) error {
	if err := checkHostname(nc.Hostname); err != nil {
		return err
	}

	if nc.BridgeConfig != nil {
		if err := checkIface(nc.BridgeConfig.VethName); err != nil {
			return err
		}
		if err := checkIface(nc.BridgeConfig.BridgeName); err != nil {
			return err
		}
	}
	return nil
}

var ifaceRe = regexp.MustCompile("^[0-9a-z.]+$")

func checkIface(iface string) error {
	if !ifaceRe.MatchString(iface) {
		return fmt.Errorf("iface '%w' not pass regexp: %s", iface, ifaceRe.String())
	}
	return nil
}

var hostnameRe = regexp.MustCompile("^[0-9a-z-A-Z]+$")

func checkHostname(hostname string) error {
	if !hostnameRe.MatchString(hostname) {
		return fmt.Errorf("hostname not pass regexp: %s", hostnameRe.String())
	}
	return nil
}

func MarshalNetworkConfig(nc *NetworkConfig) (string, error) {
	jsonContent, err := json.Marshal(nc)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(jsonContent), nil
}

func UnmarshalNetworkConfig(config string) (*NetworkConfig, error) {
	jsonContent, err := base64.StdEncoding.DecodeString(config)
	if err != nil {
		return nil, err
	}
	var nc NetworkConfig
	err = json.Unmarshal(jsonContent, &nc)
	if err != nil {
		return nil, err
	}
	return &nc, nil
}
