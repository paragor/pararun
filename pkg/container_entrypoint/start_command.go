package container_entrypoint

import (
	"fmt"
	"github.com/paragor/pararun/pkg/network"
	"github.com/paragor/pararun/pkg/reexec"
	"os"
	"syscall"
)

const (
	ContainerRootDirOnHostEnv = "PARARUN_ROOT"
	NetworkConfigEnv          = "NETWORK_CONFIG_ROOT"
)

func StartContainer(command string, root string, nc *network.NetworkConfig) error {
	configEncoded, err := network.MarshalNetworkConfig(nc)
	if err != nil {
		return fmt.Errorf("cant marshal network nc: %w", err)
	}

	cmd := reexec.Command(command)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = []string{
		"PS1=[pararun] # ", ContainerRootDirOnHostEnv + "=" + root,
		NetworkConfigEnv + "=" + configEncoded,
	}
	cmd.Dir = "/"
	var cloneFlags uintptr = syscall.CLONE_NEWIPC |
		syscall.CLONE_NEWNS |
		syscall.CLONE_NEWPID |
		syscall.CLONE_NEWUSER |
		syscall.CLONE_NEWUTS

	if nc.Type != network.NetworkTypeHost {
		cloneFlags |= syscall.CLONE_NEWNET
	}

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Pdeathsig:    syscall.SIGTERM,
		Cloneflags:   cloneFlags,
		Unshareflags: 0,
		UidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      os.Getuid(),
				Size:        1,
			},
		},
		GidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      os.Getgid(),
				Size:        1,
			},
		},
		GidMappingsEnableSetgroups: false,
		AmbientCaps:                nil,
	}
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
