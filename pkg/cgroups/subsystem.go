package cgroups

type CgroupSubsystem interface {
	Init() error
	PutProcess(pid int) error
	ResetProcess(pid int) error
	ResetAll() error
	GetPids() ([]int, error)
}
