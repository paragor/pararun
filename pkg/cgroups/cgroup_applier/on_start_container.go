package cgroup_applier

import (
	"fmt"
	"github.com/paragor/pararun/pkg/cgroups"
	"strconv"
)

func megaCloser(closers []func() error) func() error {
	return func() error {
		var resultErr error
		for _, fn := range closers {
			if err := fn(); err != nil {
				resultErr = fmt.Errorf("prev err: %w | %w", resultErr, err)
			}
		}

		return resultErr
	}

}

func ApplyCgroupConfig(spec *CgroupSpec, pid int) (close func() error, resultErr error) {
	var closers []func() error
	defer func() {
		if resultErr != nil {
			for _, fn := range closers {
				_ = fn()
			}
		}
	}()

	if spec == nil {
		return nil, fmt.Errorf("spec is nil")
	}

	if spec.Cpuset != nil {
		cpuset, err := cgroups.NewCpusetSubsystem(strconv.Itoa(pid), true)
		if err != nil {
			return nil, err
		}
		if err := cpuset.SetNumaZone(spec.Cpuset.NumaZone); err != nil {
			return nil, fmt.Errorf("on set numa zone: %w", err)
		}

		if err := cpuset.SetCpuset(spec.Cpuset.CpuList); err != nil {
			return nil, fmt.Errorf("on set cpuset: %w", err)
		}

		if err := cpuset.PutProcess(pid); err != nil {
			return nil, fmt.Errorf("on put task: %w", err)
		}

		closers = append(closers, func() error {
			fmt.Println(cpuset.ResetAll())
			return nil
		})
	}

	return megaCloser(closers), nil
}
