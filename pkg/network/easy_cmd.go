package network

import (
	"bytes"
	"fmt"
	"os/exec"
)

func easyCmd(command string, args ...string) error {
	buffer := bytes.NewBuffer(nil)

	cmd := exec.Command(command, args...)
	cmd.Stdout = buffer
	cmd.Stderr = buffer
	if err := cmd.Run(); err != nil {
		return fmt.Errorf(buffer.String())
	}

	return nil

}
