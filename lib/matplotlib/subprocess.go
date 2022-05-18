package matplotlib

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
)

const (
	defaultPython       = "python"
	defaultVenvLocation = "/tmp/coaxy-venv.tmp"
)

type Process struct {
	cmd   *exec.Cmd
	err   error
	Stdin io.WriteCloser
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func prepareSubprocess() (string, error) {
	venvDir := defaultVenvLocation

	if !fileExists(venvDir) {
		if err := exec.Command(defaultPython, "-m", "venv", defaultVenvLocation).Run(); err != nil {
			return "", err
		}
	}

	if !fileExists(venvDir) {
		return "", fmt.Errorf("venv directory %q still does not exist", venvDir)
	}

	python := path.Join(venvDir, "bin/python")

	if err := exec.Command(python, "-m", "pip", "install", "matplotlib").Run(); err != nil {
		return "", fmt.Errorf("failed to install matplotlib: %w", err)
	}

	return python, nil
}

func Subprocess() *Process {
	python, err := prepareSubprocess()
	if err != nil {
		return &Process{
			err: err,
		}
	}

	cmd := exec.Command(python)

	inputPipe, err := cmd.StdinPipe()
	if err != nil {
		return &Process{err: err}
	}

	return &Process{
		cmd:   cmd,
		Stdin: inputPipe,
	}
}

func (s *Process) Run() error {
	if s.err != nil {
		return s.err
	}
	return s.cmd.Run()
}

func RunSubprocess(makeInput func(io.Writer) error) error {
	p := Subprocess()
	if p.err != nil {
		return p.err
	}

	var scriptErr error
	go func() {
		scriptErr = makeInput(p.Stdin)
		p.Stdin.Close()
	}()

	err := p.Run()

	if scriptErr != nil {
		return scriptErr
	}

	return err
}
