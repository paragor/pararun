package container_entrypoint

import (
	"fmt"
	"github.com/paragor/pararun/pkg/reexec"
	"os"
	"regexp"
	"syscall"
)

const (
	RootEnv = "PARARUN_ROOT"
	HostEnv = "HOST_ROOT"
)

func StartContainer(command string, root string, hostname string) error {
	err := checkHostname(hostname)
	if err != nil {
		return err
	}

	cmd := reexec.Command(command)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = []string{"PS1=[pararun] # ", RootEnv + "=" + root, HostEnv + "=" + hostname}
	cmd.Dir = "/"
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Pdeathsig: syscall.SIGTERM,
		Cloneflags: syscall.CLONE_NEWIPC |
			syscall.CLONE_NEWNET |
			syscall.CLONE_NEWNS |
			syscall.CLONE_NEWPID |
			syscall.CLONE_NEWUSER |
			syscall.CLONE_NEWUTS,
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

var hostnameRe = regexp.MustCompile("^[0-9a-z-A-Z]+$")
func checkHostname(hostname string) error {
	if !hostnameRe.MatchString(hostname) {
		return fmt.Errorf("hostname not pass regexp: %s", hostnameRe.String())
	}
	return nil
}
