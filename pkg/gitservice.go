package pkg

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

const (
	GitBinary         = "git"
	ErrBinaryNotFound = "binary not found"
	ErrorRunningCmd   = "error running command"
)

// CloneRepo clones a github repository
func CloneRepo(repo, workdir string) (err error) {
	ok, err := BinaryExists(GitBinary)
	if err != nil {
		return
	}
	if !ok {
		err = errors.New(ErrBinaryNotFound)
		return
	}

	// git clone
	cmd := exec.Command(GitBinary, "clone", repo)
	cmd.Dir = workdir
	_, err = cmd.Output()
	if err != nil {
		// if clone fails, the error will be of type *exec.ExitError
		if exitError, ok := err.(*exec.ExitError); ok {
			return fmt.Errorf("%s", exitError.Stderr)
		}

		return fmt.Errorf("%s: %s", ErrorRunningCmd, err)
	}

	return
}

func GetAllCommits(workdir string) (commits []string, err error) {
	ok, err := BinaryExists(GitBinary)
	if err != nil {
		return
	}
	if !ok {
		err = errors.New(ErrBinaryNotFound)
		return
	}

	// git log
	cmd := exec.Command(GitBinary, "log", "--all", "--pretty=format:%H")
	cmd.Dir = workdir
	output, err := cmd.Output()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("%s", exitError.Stderr)
		}

		return nil, fmt.Errorf("%s: %s", ErrorRunningCmd, err)
	}

	s := bufio.NewScanner(strings.NewReader(string(output)))
	for s.Scan() {
		commits = append(commits, s.Text())
	}

	return
}

func GetCommitFiles(workdir, commit string) (files []string, err error) {
	ok, err := BinaryExists(GitBinary)
	if err != nil {
		return
	}
	if !ok {
		err = errors.New(ErrBinaryNotFound)
		return
	}

	// git ls-tree
	cmd := exec.Command(GitBinary, "ls-tree", "--name-only", "-r", commit)
	cmd.Dir = workdir
	output, err := cmd.Output()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("%s", exitError.Stderr)
		}

		return nil, fmt.Errorf("%s: %s", ErrorRunningCmd, err)
	}

	s := bufio.NewScanner(strings.NewReader(string(output)))
	for s.Scan() {
		files = append(files, s.Text())
	}

	return
}

func GetFileContent(workdir, commit, file string) (r io.Reader, err error) {
	ok, err := BinaryExists(GitBinary)
	if err != nil {
		return
	}
	if !ok {
		err = errors.New(ErrBinaryNotFound)
		return
	}

	// git show
	cmd := exec.Command(GitBinary, "show", fmt.Sprintf("%s:%s", commit, file))
	cmd.Dir = workdir
	output, err := cmd.Output()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("%s", exitError.Stderr)
		}

		return nil, fmt.Errorf("%s: %s", ErrorRunningCmd, err)
	}

	return strings.NewReader(string(output)), nil
}
