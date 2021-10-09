package main

import (
	"github.com/paragor/pararun/pkg/container_entrypoint"
	"github.com/paragor/pararun/pkg/hacks"
	"github.com/paragor/pararun/pkg/image"
	"github.com/paragor/pararun/pkg/reexec"
	"log"
	"os"
)

const imageUrl = "http://dl-cdn.alpinelinux.org/alpine/v3.10/releases/x86_64/alpine-minirootfs-3.10.1-x86_64.tar.gz"
const imageName = "alpine"

const containerRootDir = "/var/lib/pararun/root/"

func main() {
	container_entrypoint.RegisterAllReexecHooks()
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

	err = os.MkdirAll(containerRootDir, 0755)
	if err != nil {
		panic(err)
	}

	if err := imageController.UnpackImage(imageName, containerRootDir); err != nil {
		panic(err)
	}

	log.Println("[PARARUN] start container")
	if err := container_entrypoint.StartContainer("/bin/sh", containerRootDir); err != nil {
		panic(err)
	}
}
