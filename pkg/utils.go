package pkg

import (
	"fmt"
	"os/exec"

	"github.com/go-cmd/cmd"
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

// HandleStatus returns an error if the command failed to execute or there is a go error in status object
func HandleStatus(s cmd.Status) error {
	switch {
	case s.Error != nil:
		return s.Error
	case s.Exit != 0:
		err := fmt.Errorf("error while running %s\n%q", s.Cmd, s.Stderr)
		return err
	default:
		for _, line := range s.Stdout {
			fmt.Println(line)
		}
	}

	return nil
}

// RmDir removes a directory
func RmDir(dir string) (err error) {
	rmdir := cmd.NewCmd("rm", "-rf", dir)
	status := <-rmdir.Start()
	err = HandleStatus(status)
	if err != nil {
		return err
	}
	return nil
}

// MkDir creates a directory
func MkDir(dir string) (err error) {
	mkdir := cmd.NewCmd("mkdir", dir)
	status := <-mkdir.Start()
	err = HandleStatus(status)
	if err != nil {
		return err
	}
	return nil
}
