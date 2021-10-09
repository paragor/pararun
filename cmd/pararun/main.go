package main

import (
	"fmt"
	"github.com/paragor/pararun/pkg/cgroups"
)

func main() {
	freezer, err := cgroups.NewFreezerSubsystem("test", true)
	if err != nil {
		panic(err)
	}

	for {

		fmt.Printf("Pid to freeze: ")
		var pid int
		_, err = fmt.Scanf("%d", &pid)
		if err != nil {
			panic(err)
		}
		err = freezer.PutProcess(pid)
		if err != nil {
			panic(err)
		}

		err = freezer.Freeze()
		if err != nil {
			panic(err)
		}

		fmt.Println("Press [ENTER] to unfreeze process.")
		fmt.Scanf("a")

		err = freezer.ResetAll()
		if err != nil {
			panic(err)
		}
	}
}
