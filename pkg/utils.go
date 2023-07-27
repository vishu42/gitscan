package pkg

import (
	"fmt"
	"os/exec"
)

// BinaryExists throws an error if the binary does not exist
func BinaryExists(binary string) (bool, error) {
	cmd := exec.Command("which", binary)
	_, err := cmd.Output()
	if err != nil {
		// if the binary does not exist, the error will be of type *exec.ExitError
		if exitError, ok := err.(*exec.ExitError); ok {
			return false, fmt.Errorf("%s", exitError.Stderr)
		}

		return false, fmt.Errorf("error running command: %s", err)
	}

	return true, nil
}
