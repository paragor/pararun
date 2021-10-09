package mounts

import (
	"fmt"
	"os"
	"path"
	"syscall"
)

const PivotOldRootDirName = ".pivot_old_root"

func PivotRoot(newRoot string) error {
	oldRoot := path.Join(newRoot, PivotOldRootDirName)

	if err := syscall.Mount(newRoot, newRoot, "", syscall.MS_BIND|syscall.MS_REC, ""); err != nil {
		return fmt.Errorf("mount new root: %w", err)
	}

	if err := os.MkdirAll(oldRoot, 0700); err != nil {
		return fmt.Errorf("mkdir oldroot: %w", err)
	}

	if err := syscall.PivotRoot(newRoot, oldRoot); err != nil {
		return fmt.Errorf("pivot syscall: %w", err)
	}

	if err := os.Chdir("/"); err != nil {
		return fmt.Errorf("chdir: %w", err)
	}
	oldRootInNewFs := "/" + PivotOldRootDirName
	if err := syscall.Unmount(
		oldRootInNewFs,
		syscall.MNT_DETACH,
	); err != nil {
	return fmt.Errorf("unmount old root: %w", err)
	}

	if err := os.RemoveAll(oldRootInNewFs); err != nil {
	return fmt.Errorf("remove old root dir: %w", err)
	}

	return nil
}
