package main

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/paragor/pararun/pkg/cgroups/cgroup_applier"
	"github.com/paragor/pararun/pkg/container_entrypoint"
	"github.com/paragor/pararun/pkg/hacks"
	"github.com/paragor/pararun/pkg/image"
	"github.com/paragor/pararun/pkg/network"
	"github.com/paragor/pararun/pkg/reexec"
	"log"
	"math/rand"
	"net"
	"os"
	"path"
	"time"
)

const imageUrl = "http://dl-cdn.alpinelinux.org/alpine/v3.10/releases/x86_64/alpine-minirootfs-3.10.1-x86_64.tar.gz"
const imageName = "alpine"

const containerRootDir = "/var/lib/pararun/root/"

func main() {
	rand.Seed(time.Now().UnixNano())
	container_entrypoint.RegisterInsideContainerEntrypoint()
	if reexec.Init() {
		os.Exit(0)
	}

	err := hacks.RunHacks()
	if err != nil {
		panic(err)
	}

	imageController, err := image.NewImageController("/srv/images/")
	if err != nil {
		panic(err)
	}

	log.Println("[PARARUN] download image")
	err = imageController.DownloadImage(imageUrl, imageName, false)
	if err != nil {
		panic(err)
	}

	var containerName string
	fmt.Printf("container name: ")
	fmt.Scanf("%s", &containerName)

	containerFs := path.Join(containerRootDir, containerName)

	err = os.MkdirAll(containerFs, 0755)
	if err != nil {
		panic(err)
	}

	if err := imageController.UnpackImage(imageName, containerFs); err != nil {
		panic(err)
	}

	networkConfig := &network.NetworkConfig{
		BridgeConfig: &network.BridgeConfig{
			BridgeName: "pararun.br0",
			BridgeNet: net.IPNet{
				IP:   net.IPv4(192, 168, 123, 1),
				Mask: net.IPv4Mask(255, 255, 255, 0),
			},
			ContainerNet: net.IPNet{
				IP:   net.IPv4(192, 168, 123, byte(rand.Intn(250)+1)),
				Mask: net.IPv4Mask(255, 255, 255, 0),
			},
		},
		Nameservers: []net.IP{
			net.IPv4(1, 1, 1, 1),
			net.IPv4(8, 8, 8, 8),
		},
		Hostname: uuid.New().String(),
		Type:     network.NetworkTypeBridge,
	}
	if err := network.ValidateConfig(networkConfig); err != nil {
		panic(err)
	}
	containerSpec := &container_entrypoint.ContainerSpecification{
		Command:             "/bin/sh",
		Args:                []string{},
		ContainerRootOnHost: containerFs,
		NetworkConfig:       networkConfig,
		CgroupSpec: &cgroup_applier.CgroupSpec{Cpuset: &cgroup_applier.CgroupCpusetSpec{
			NumaZone: 0,
			CpuList:  []uint8{0},
		}},
	}

	log.Println("[PARARUN] start container")
	if err := container_entrypoint.StartContainer(containerSpec); err != nil {
		panic(err)
	}
}
