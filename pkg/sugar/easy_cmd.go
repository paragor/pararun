package sugar

import (
	"bytes"
	"fmt"
	"os/exec"
)

func EasyCmd(command string, args ...string) error {
	buffer := bytes.NewBuffer(nil)

	cmd := exec.Command(command, args...)
	cmd.Stdout = buffer
	cmd.Stderr = buffer
	if err := cmd.Run(); err != nil {
		return fmt.Errorf(buffer.String())
	}

	return nil

}
