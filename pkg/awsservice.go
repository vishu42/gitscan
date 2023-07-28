package pkg

import (
	"errors"
	"fmt"
	"os/exec"
)

const (
	AwsBinary = "aws"
)

// VerifyAWSKey verifies the AWS key and secret by running the aws sts get-caller-identity command, returns error if the command fails
func VerifyAWSKey(keyID, keySecret string) (err error) {
	ok, err := BinaryExists(AwsBinary)
	if err != nil {
		return
	}
	if !ok {
		err = errors.New(ErrBinaryNotFound)
		return
	}

	// aws sts get-caller-identity
	cmd := exec.Command(AwsBinary, "sts", "get-caller-identity")
	cmd.Env = []string{
		"AWS_ACCESS_KEY_ID=" + keyID,
		"AWS_SECRET_ACCESS_KEY=" + keySecret,
	}
	_, err = cmd.Output()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return fmt.Errorf("%s", exitError.Stderr)
		}

		return fmt.Errorf("%s: %s", ErrorRunningCmd, err)
	}

	return
}
