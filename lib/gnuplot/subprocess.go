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
