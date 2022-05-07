package gnuplot

import (
	"io"
	"os/exec"
)

type Process struct {
	cmd   *exec.Cmd
	err   error
	Stdin io.WriteCloser
}

func Subprocess() *Process {
	cmd := exec.Command("gnuplot")

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
