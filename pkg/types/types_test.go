package types

import (
	"testing"

	"github.com/alexellis/ofc-bootstrap/pkg/execute"
)

func TestExec_WithShell(t *testing.T) {
	task := execute.ExecTask{Command: "/bin/ls /", Shell: true}
	res, err := task.Execute()
	if err != nil {
		t.Errorf(err.Error())
		t.Fail()
	}

	if len(res.Stdout) == 0 {
		t.Errorf("want data, but got empty")
		t.Fail()
	}

	if len(res.Stderr) != 0 {
		t.Errorf("want empty, but got: %s", res.Stderr)
		t.Fail()
	}
}
