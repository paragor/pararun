package network

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"syscall"
	"time"
)

func HostApplyConfig(nc *NetworkConfig) error {
	return nil
}

func ContainerApplyConfig(nc *NetworkConfig, hostResolveConf string ) error {
	err := syscall.Sethostname([]byte(nc.Hostname))
	if err != nil {
		panic(err)
	}

	err = setupResolvConf(nc, hostResolveConf)
	if err != nil {
		return err
	}

	return nil
}

func setupResolvConf(nc *NetworkConfig, hostResolveConf string) error {
	err := os.MkdirAll("/etc", 0755)
	if err != nil {
		return fmt.Errorf("cant mkdir /etc: %w", err)
	}

	if len(nc.Nameservers) > 0 {

		resolvConfContent := ""
		for _, ns := range nc.Nameservers {
			resolvConfContent += fmt.Sprintf("nameserver %s\n", ns.String())
		}
		err = ioutil.WriteFile("/etc/resolv.conf", []byte(resolvConfContent), 0644)
		if err != nil {
			return fmt.Errorf("cant write /etc/resolv.conf: %w", err)
		}
	} else if nc.Type == NetworkTypeHost {
		err = ioutil.WriteFile("/etc/resolv.conf", []byte(hostResolveConf), 0644)
		if err != nil {
			return fmt.Errorf("cant write /etc/resolv.conf: %w", err)
		}
	}
	return nil
}

func WaitNetwork(duration time.Duration) error {
	deadline := time.Now().Add(duration)
	for {
		if time.Now().After(deadline) {
			return fmt.Errorf("network set up deadline")
		}

		ifaces, err := net.Interfaces()
		if err != nil {
			return fmt.Errorf("cant get ifaces: %w", err)
		}

		if len(ifaces) > 1 {
			return nil
		}

		time.Sleep(time.Millisecond * 100)
	}

}
