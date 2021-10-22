package container_entrypoint

import (
	"github.com/paragor/pararun/pkg/cgroups/cgroup_applier"
	"github.com/paragor/pararun/pkg/network"
)

type ContainerSpecification struct {
	Command             string
	Args                []string
	ContainerRootOnHost string
	NetworkConfig       *network.NetworkConfig
	CgroupSpec          *cgroup_applier.CgroupSpec
}

