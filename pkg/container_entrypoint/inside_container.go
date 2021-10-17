package container_entrypoint

import (
	"fmt"
	"github.com/paragor/pararun/pkg/mounts"
	"github.com/paragor/pararun/pkg/network"
	"github.com/paragor/pararun/pkg/reexec"
	"io/ioutil"
	"os"
	"time"
)

// RegisterAllReexecHooks pre start хуки в неймспейсе контейнера
func RegisterAllReexecHooks() {
	reexec.Register(func() {
		resolvConfContent, err := ioutil.ReadFile("/etc/resolv.conf")
		if err != nil {
			panic(fmt.Errorf("cant read /etc/resolv.conf: %w", err))
		}

		err = mounts.MountProc(os.Getenv(ContainerRootDirOnHostEnv))
		if err != nil {
			panic(err)
		}
		err = mounts.PivotRoot(os.Getenv(ContainerRootDirOnHostEnv))
		if err != nil {
			panic(err)
		}

		nc, err := network.UnmarshalNetworkConfig(os.Getenv(NetworkConfigEnv))
		if err != nil {
			panic(err)
		}

		err = network.ContainerApplyConfig(nc, string(resolvConfContent))
		if err != nil {
			panic(err)
		}

		err = network.WaitNetwork(time.Second * 5)
		if err != nil {
			panic(err)
		}
	})

	reexec.Register(func() {
		_ = os.Unsetenv(ContainerRootDirOnHostEnv)
	})

}
