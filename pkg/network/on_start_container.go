package network

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"syscall"
	"time"
)

func SetupHostname(nc *NetworkConfig) error {
	err := syscall.Sethostname([]byte(nc.Hostname))
	if err != nil {
		panic(err)
	}

	file, err := os.OpenFile("/etc/hosts", syscall.O_RDWR | syscall.O_CREAT |syscall.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("cant open /etc/hosts: %w", err)
	}

	if _, err := file.Seek(0, io.SeekEnd); err != nil {
		return fmt.Errorf("cant seek file /etc/hosts: %w", err)
	}

	if _, err := file.WriteString(fmt.Sprintf("\n127.0.0.1 %s\n", nc.Hostname)); err != nil {
		return fmt.Errorf("cant append file /etc/hosts: %w", err)
	}

	return nil
}

func SetupResolvConf(nc *NetworkConfig, hostResolveConf string) error {
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
