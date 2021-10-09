package container_entrypoint

import (
	"github.com/paragor/pararun/pkg/mounts"
	"github.com/paragor/pararun/pkg/reexec"
	"os"
)

// RegisterAllReexecHooks pre start хуки в неймспейсе контейнера
func RegisterAllReexecHooks() {
	reexec.Register(func() {
		err := mounts.MountProc(os.Getenv(RootEnv))
		if err != nil {
			panic(err)
		}
	})
	reexec.Register(func() {
		err := mounts.PivotRoot(os.Getenv(RootEnv))
		if err != nil {
			panic(err)
		}
	})

	reexec.Register(func() {
		os.Unsetenv(RootEnv)
	})

}
