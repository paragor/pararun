package cgroup_applier

type CgroupSpec struct {
	Cpuset *CgroupCpusetSpec
}

type CgroupCpusetSpec struct {
	NumaZone uint8
	CpuList  []uint8
}
