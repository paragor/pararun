package mounts

import (
	"fmt"
	"os"
	"path"
	"syscall"
)

func MountProc(newRoot string) error {
	procfs := path.Join(newRoot, "/proc")
	if err := os.MkdirAll(procfs, 0755); err != nil {
		return fmt.Errorf("cant mkdir new procfs: %w", err)
	}

	if err := syscall.Mount("proc", procfs, "proc", 0, ""); err != nil {
		return fmt.Errorf("cant mount procfs: %w", err)
	}
	return nil
}
