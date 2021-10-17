package network

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

const namespacePrefix = "pararun."

func SetupBridgeNetwork(nc *NetworkConfig, containerPid int) error {
	namespaceName, err := setupNetworkNamespace(containerPid)
	if err != nil {
		return fmt.Errorf("cant setup veth on host: %w", err)
	}
	defer unSetupNetworkNamespace(namespaceName)

	if err := upLoopback(containerPid); err != nil {
		return fmt.Errorf("cant up loopback: %w", err)
	}

	if _, err := net.InterfaceByName(nc.BridgeConfig.BridgeName); err != nil {
		err := setupBridge(nc)
		if err != nil {
			return fmt.Errorf("cant setup bridge on host: %w", err)
		}
	}

	if err := setupVeth(nc, namespaceName); err != nil {
		return fmt.Errorf("cant setup veth on host: %w", err)
	}

	return nil
}
func setupNetworkNamespace(pid int) (string, error) {
	ns := getNamespaceName(pid)
	err := os.Symlink(fmt.Sprintf("/proc/%d/ns/net", pid), fmt.Sprintf("/var/run/netns/%s", ns))
	if err != nil {
		return "", fmt.Errorf("cant create symlink for namespace: %w", err)
	}

	return ns, nil
}
func unSetupNetworkNamespace(ns string) error {
	err := os.Remove(fmt.Sprintf("/var/run/netns/%s", ns))
	if err != nil {
		return fmt.Errorf("cant remove symlink: %w", err)
	}
	return nil
}

func setupBridge(nc *NetworkConfig) error {
	return nil
}
func setupVeth(nc *NetworkConfig, ns string) error {
	hostVeth := nc.BridgeConfig.VethName + ".h"
	containerVeth := nc.BridgeConfig.VethName + ".c"
	fullNsPath := fmt.Sprintf("/var/run/netns/%s", ns)
	commands := [][]string{
		{"iptables", "-t", "nat", "-I", "POSTROUTING", "-s", nc.BridgeConfig.BridgeNet.String(), "!", "-o", hostVeth, "-j", "MASQUERADE"},
		{"ip", "link", "add", hostVeth, "type", "veth", "peer", "name", containerVeth},
		{"ip", "addr", "add", nc.BridgeConfig.BridgeNet.String(), "dev", hostVeth},
		{"ip", "link", "set", containerVeth, "netns", ns},
		{"ip", "link", "set", hostVeth, "up"},
		{"nsenter", "--net=" + fullNsPath, "ip", "link", "sh"},
		{"nsenter", "--net=" + fullNsPath, "ip", "addr", "sh"},
		{"nsenter", "--net=" + fullNsPath, "ip", "link", "set", "dev", containerVeth, "name", "eth0"},
		{"nsenter", "--net=" + fullNsPath, "ip", "link", "set", "eth0", "up"},
		{"nsenter", "--net=" + fullNsPath, "ip", "addr", "add", nc.BridgeConfig.ContainerNet.String(), "dev", "eth0"},
		{"nsenter", "--net=" + fullNsPath, "ip", "route", "replace", "default", "via", nc.BridgeConfig.BridgeNet.IP.String(), "dev", "eth0"},
	}

	for _, args := range commands {
		if err := easyCmd(args[0], args[1:]...); err != nil {
			return fmt.Errorf("cant exec '%s': %w", strings.Join(args, " "), err)
		}
	}
	return nil
}

func getNamespaceName(pid int) string {
	return namespacePrefix + strconv.Itoa(pid)
}

func upLoopback(pid int) error {
	if err := easyCmd("nsenter", "--target", strconv.Itoa(pid), "-n", "ip", "link", "set", "lo", "up"); err != nil {
		return err
	}
	return nil
}
