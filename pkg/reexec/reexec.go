package reexec

import (
	"os"
	"os/exec"
	"syscall"
)

const ExecCommand = "reexec"

var registeredInitializers = make([]func(), 0)

func Register(initializer func()) {
	registeredInitializers = append(registeredInitializers, initializer)
}

func Init() bool {
	if os.Args[0] == ExecCommand {
		for _, initializer := range registeredInitializers {
			initializer()
		}

		if err := syscall.Exec(os.Args[1], os.Args[1:], os.Environ()); err != nil {
			panic(err)
		}
		return true
	}
	return false
}

func Command(args ...string) *exec.Cmd {
	return &exec.Cmd{
		Path: "/proc/self/exe",
		Args: append([]string{ExecCommand}, args...),
		SysProcAttr: &syscall.SysProcAttr{
			Pdeathsig: syscall.SIGTERM,
		},
	}
}
