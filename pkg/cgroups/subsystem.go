package cgroups

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
)

type CgroupSubsystem interface {
	PutProcess(pid int) error
	ResetProcess(pid int) error
	ResetAll() error
	GetPids() ([]int, error)
}

type abstractSubsystem struct {
	cgroupBasePath string
	groupName      string
	withChildren bool
}

func (s *abstractSubsystem) init() error {
	cgroupSubsystemPath := s.getPathForCurrentGroup("")

	if _, err := os.Stat(cgroupSubsystemPath); err == nil {
		err = os.Remove(cgroupSubsystemPath)
		if err != nil {
			return err
		}
	} else if !os.IsNotExist(err) {
		return err
	}

	err := os.MkdirAll(cgroupSubsystemPath, 0755)
	if err != nil {
		return err
	}

	payload := []byte("0")
	if s.withChildren {
		payload = []byte("1")
	}
	err = ioutil.WriteFile(path.Join(cgroupSubsystemPath, "cgroup.clone_children"), payload, 0666)
	if err != nil {
		return err
	}
	return nil
}

func (s *abstractSubsystem) PutProcess(pid int) error {
	return ioutil.WriteFile(s.getPathForCurrentGroup("tasks"), []byte(strconv.Itoa(pid)), 0666)
}

func (s *abstractSubsystem) ResetProcess(pid int) error {
	return ioutil.WriteFile(path.Join(FreezerDir, "tasks"), []byte(strconv.Itoa(pid)), 0666)
}
func (s *abstractSubsystem) ResetAll() error {
	return os.Remove(s.getPathForCurrentGroup(""))
}

func (s *abstractSubsystem) GetPids() ([]int, error) {
	res, err := ioutil.ReadFile(s.getPathForCurrentGroup("tasks"))
	if err != nil {
		return nil, err
	}

	pidsStrs := strings.Split(string(res), "\n")
	pids := make([]int, 0, len(pidsStrs))
	for _, pidStr := range pidsStrs {
		pid, err := strconv.Atoi(pidStr)
		if err != nil {
			return nil, fmt.Errorf("cant parse tasks id for %s:%w", pidStr, err)
		}
		pids = append(pids, pid)
	}
	return pids, nil
}

func (s *abstractSubsystem) getPathForCurrentGroup(file string) string {
	return path.Join(s.cgroupBasePath, s.groupName, file)
}
