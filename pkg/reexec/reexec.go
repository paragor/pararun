package reexec

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

var registeredInitializers = make(map[string]func(), 0)

func Register(name string, initializer func()) {
	if _, exists := registeredInitializers[name]; exists {
		panic(fmt.Errorf("hook %s already exists", name))
	}
	registeredInitializers[name] = initializer
}

func Init() bool {
	if initializer, ok := registeredInitializers[os.Args[0]]; ok {
		initializer()
		return true
	}
	return false
}

func Command(hook string, args ...string) *exec.Cmd {
	return &exec.Cmd{
		Path: "/proc/self/exe",
		Args: append([]string{hook}, args...),
		SysProcAttr: &syscall.SysProcAttr{
			Pdeathsig: syscall.SIGTERM,
		},
	}
}
