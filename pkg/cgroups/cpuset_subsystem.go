package cgroups

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

type CpusetSubsystem struct {
	abstractSubsystem
}

func NewCpusetSubsystem(groupName string, withChildren bool) (*CpusetSubsystem, error) {
	subsystem := &CpusetSubsystem{
		abstractSubsystem: abstractSubsystem{
			cgroupBasePath: CpusetDir,
			groupName:      groupName,
			withChildren:   withChildren,
		},
	}
	return subsystem, subsystem.init()
}

func (s *CpusetSubsystem) SetNumaZone(zone uint8) error {
	return ioutil.WriteFile(s.getPathForCurrentGroup("cpuset.mems"), []byte(strconv.Itoa(int(zone))), 0666)
}
func (s *CpusetSubsystem) GetNumaZone() (uint8, error) {
	res, err := ioutil.ReadFile(s.getPathForCurrentGroup("cpuset.mems"))
	if err != nil {
		return 0, err
	}

	numa, err := strconv.Atoi(strings.TrimSpace(string(res)))
	if err != nil {
		return 0, err
	}
	return uint8(numa), nil
}
func (s *CpusetSubsystem) SetCpuset(cpus []uint8) error {

	sort.Sort(sortUint8(cpus))

	var cpusStrs []string
	for _, cpu := range cpus {
		cpusStrs = append(cpusStrs, strconv.Itoa(int(cpu)))
	}

	payload := strings.Join(cpusStrs, ",")
	err := ioutil.WriteFile(s.getPathForCurrentGroup("cpuset.cpus"), []byte(payload), 0666)
	if err != nil {
		return err
	}

	resultCpus, err := s.GetCpuset()
	if err != nil {
		return err
	}

	if !reflect.DeepEqual(resultCpus, cpus) {
		return fmt.Errorf("not applied")
	}

	return nil
}
func (s *CpusetSubsystem) GetCpuset() ([]uint8, error) {
	res, err := ioutil.ReadFile(s.getPathForCurrentGroup("cpuset.cpus"))
	if err != nil {
		return nil, err
	}
	cpusRawList := strings.Split(strings.TrimSpace(string(res)), ",")
	var cpus []uint8
	for _, cpuRaw := range cpusRawList {
		if cpuRaw == "" {
			continue
		}

		if strings.Contains(cpuRaw, "-") {
			parts := strings.Split(cpuRaw, "-")
			part0, err := strconv.Atoi(parts[0])
			if err != nil {
				return nil, err
			}
			part1, err := strconv.Atoi(parts[1])
			if err != nil {
				return nil, err
			}
			for i := part0; i <= part1; i++ {
				cpus = append(cpus, uint8(i))
			}
			continue
		}

		atoi, err := strconv.Atoi(cpuRaw)
		if err != nil {
			return nil, err
		}
		cpus = append(cpus, uint8(atoi))
	}

	sort.Sort(sortUint8(cpus))
	return cpus, nil

}


type sortUint8 []uint8

func (x sortUint8) Len() int           { return len(x) }
func (x sortUint8) Less(i, j int) bool { return x[i] < x[j] }
func (x sortUint8) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

