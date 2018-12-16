package execute

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"
)

type ExecTask struct {
	Command string
	Shell   bool
}

type ExecResult struct {
	Stdout   string
	Stderr   string
	ExitCode int
}

func (et ExecTask) Execute() (ExecResult, error) {
	fmt.Println(et.Command)

	var cmd *exec.Cmd

	if et.Shell {
		startArgs := strings.Split(et.Command, " ")
		args := []string{"-c"}
		for _, part := range startArgs {
			args = append(args, part)
		}

		cmd = exec.Command("/bin/bash", args...)
	} else {
		cmd = exec.Command(et.Command)
	}

	stdoutPipe, stdoutPipeErr := cmd.StdoutPipe()
	if stdoutPipeErr != nil {
		return ExecResult{}, stdoutPipeErr
	}

	stderrPipe, stderrPipeErr := cmd.StderrPipe()
	if stderrPipeErr != nil {
		return ExecResult{}, stderrPipeErr
	}

	startErr := cmd.Start()

	if startErr != nil {
		return ExecResult{}, startErr
	}

	stdoutBytes, err := ioutil.ReadAll(stdoutPipe)
	if err != nil {
		return ExecResult{}, err
	}

	stderrBytes, err := ioutil.ReadAll(stderrPipe)

	if err != nil {
		return ExecResult{}, err
	}

	fmt.Println("res: " + string(stdoutBytes))

	return ExecResult{
		Stdout: string(stdoutBytes),
		Stderr: string(stderrBytes),
	}, nil
}
