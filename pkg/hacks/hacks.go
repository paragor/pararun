package hacks

import (
	"github.com/paragor/pararun/pkg/hacks/hacks_centos7"
	"os"
)

func RunHacks() error {
	if isCentos() {
		return hacks_centos7.EnableUserNamespaces()
	}

	return nil
}

func isCentos() bool {
	_, err := os.Stat("/etc/redhat-release")
	return err == nil
}
