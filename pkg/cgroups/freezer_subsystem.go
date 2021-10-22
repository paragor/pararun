package cgroups

import (
	"io/ioutil"
)

type FreezeState string

const (
	FreezerStateFreeze   FreezeState = "FROZEN"
	FreezerStateFreezing FreezeState = "FREEZING"
	FreezerStateThawed   FreezeState = "THAWED"
)

type FreezerSubsystem struct {
	abstractSubsystem
}

func NewFreezerSubsystem(groupName string, withChildren bool) (*FreezerSubsystem, error) {
	subsystem := &FreezerSubsystem{
		abstractSubsystem: abstractSubsystem{
			cgroupBasePath: FreezerDir,
			groupName:      groupName,
			withChildren:   withChildren,
		},
	}
	return subsystem, subsystem.init()
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
