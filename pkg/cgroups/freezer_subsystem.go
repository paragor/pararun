package cgroups

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
)

type FreezeState string

const (
	FreezerStateFreeze   FreezeState = "FROZEN"
	FreezerStateFreezing FreezeState = "FREEZING"
	FreezerStateThawed   FreezeState = "THAWED"
)

type FreezerSubsystem struct {
	groupName    string
	withChildren bool
}

func NewFreezerSubsystem(groupName string, withChildren bool) (*FreezerSubsystem, error) {
	freezer := &FreezerSubsystem{groupName: groupName, withChildren: withChildren}
	return freezer, freezer.Init()
}

func (s *FreezerSubsystem) Init() error {
	err := os.MkdirAll(s.getPathForCurrentGroup(""), 0755)
	if err != nil {
		return err
	}

	payload := []byte("0")
	if s.withChildren {
		payload = []byte("1")
	}
	err = ioutil.WriteFile(s.getPathForCurrentGroup("cgroup.clone_children"), payload, 0666)
	if err != nil {
		return err
	}
	return nil
}

func (s *FreezerSubsystem) PutProcess(pid int) error {
	return ioutil.WriteFile(s.getPathForCurrentGroup("tasks"), []byte(strconv.Itoa(pid)), 0666)
}

func (s *FreezerSubsystem) ResetProcess(pid int) error {
	return ioutil.WriteFile(path.Join(FreezerDir, "tasks"), []byte(strconv.Itoa(pid)), 0666)
}
func (s *FreezerSubsystem) ResetAll() error {
	pids, err := ioutil.ReadFile(s.getPathForCurrentGroup("tasks"))
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path.Join(FreezerDir, "tasks"), pids, 0666)
	if err != nil {
		return err
	}

	return os.Remove(s.getPathForCurrentGroup(""))
}

func (s *FreezerSubsystem) GetPids() ([]int, error) {
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

func (s *FreezerSubsystem) Freeze() error {
	return ioutil.WriteFile(s.getPathForCurrentGroup("freezer.state"), []byte(FreezerStateFreeze), 0666)
}
func (s *FreezerSubsystem) Thawed() error {
	return ioutil.WriteFile(s.getPathForCurrentGroup("freezer.state"), []byte(FreezerStateThawed), 0666)
}
func (s *FreezerSubsystem) GetState() (FreezeState, error) {
	res, err := ioutil.ReadFile(s.getPathForCurrentGroup("freezer.state"))
	if err != nil {
		return "", err
	}

	return FreezeState(res), nil
}
func (s *FreezerSubsystem) getPathForCurrentGroup(file string) string {
	return path.Join(FreezerDir, s.groupName, file)
}
