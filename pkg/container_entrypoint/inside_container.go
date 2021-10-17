package container_entrypoint

import (
	"fmt"
	"github.com/paragor/pararun/pkg/mounts"
	"github.com/paragor/pararun/pkg/network"
	"github.com/paragor/pararun/pkg/reexec"
	"io/ioutil"
	"os"
	"syscall"
	"time"
)

const insideContainerEntrypoint = "inside_container_hook"

// RegisterInsideContainerEntrypoint pre start хуки в неймспейсе контейнера
func RegisterInsideContainerEntrypoint() {
	reexec.Register(insideContainerEntrypoint, func() {
		resolvConfContent, err := ioutil.ReadFile("/etc/resolv.conf")
		if err != nil {
			panic(fmt.Errorf("cant read /etc/resolv.conf: %w", err))
		}

		if err := mounts.MountProc(os.Getenv(ContainerRootDirOnHostEnv)); err != nil {
			panic(err)
		}
		if err := mounts.PivotRoot(os.Getenv(ContainerRootDirOnHostEnv)); err != nil {
			panic(err)
		}

		nc, err := network.UnmarshalNetworkConfig(os.Getenv(NetworkConfigEnv))
		if err != nil {
			panic(err)
		}

		if err := network.ContainerApplyConfig(nc, string(resolvConfContent)); err != nil {
			panic(err)
		}

		if err := network.WaitNetwork(time.Second * 5); err != nil {
			panic(err)
		}

		if err := syscall.Exec(os.Args[1], os.Args[1:], os.Environ()); err != nil {
			panic(err)
		}

		_ = os.Unsetenv(ContainerRootDirOnHostEnv)
	})

}
