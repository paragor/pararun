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

type ContainerSpecification struct {
	Command             string
	Args                []string
	ContainerRootOnHost string
	NetworkConfig       *network.NetworkConfig
}

func StartContainer(containerSpec *ContainerSpecification) error {
	configEncoded, err := network.MarshalNetworkConfig(containerSpec.NetworkConfig)
	if err != nil {
		return fmt.Errorf("cant marshal network nc: %w", err)
	}

	cmd := reexec.Command(insideContainerEntrypoint, append([]string{containerSpec.Command}, containerSpec.Args...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = []string{
		"PS1=[pararun] # ", ContainerRootDirOnHostEnv + "=" + containerSpec.ContainerRootOnHost,
		NetworkConfigEnv + "=" + configEncoded,
	}
	cmd.Dir = "/"
	var cloneFlags uintptr = syscall.CLONE_NEWIPC |
		syscall.CLONE_NEWNS |
		syscall.CLONE_NEWPID |
		syscall.CLONE_NEWUSER |
		syscall.CLONE_NEWUTS

	if containerSpec.NetworkConfig.Type != network.NetworkTypeHost {
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
	if err := cmd.Start(); err != nil {
		return err
	}

	if containerSpec.NetworkConfig.Type == network.NetworkTypeBridge {
		if err := network.SetupBridgeNetwork(containerSpec.NetworkConfig, cmd.Process.Pid); err != nil {
			killErr := syscall.Kill(cmd.Process.Pid, syscall.SIGINT)
			if killErr != nil {
				return fmt.Errorf("cant kill porcess on error upLoopback: %w; origin err: %w", killErr, err)
			}
			return fmt.Errorf("cant setup bridge: %w", err)
		}
	}

	if err := cmd.Wait(); err != nil {
		return err
	}
	return nil
}
