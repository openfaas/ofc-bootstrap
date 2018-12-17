package types

import (
	"testing"

	"github.com/alexellis/ofc-bootstrap/pkg/execute"
)

// func Test_generateSecret(t *testing.T) {
// 	val, secretErr := generateSecret()
// 	if secretErr != nil {
// 		t.Errorf(secretErr.Error())
// 		t.Fail()
// 	}

// 	if len(val) == 0 {
// 		t.Error("want non-zero-length secret")
// 		t.Fail()
// 	}
// 	log.Println(val)
// }

// func Test_generateSecret_Twice(t *testing.T) {
// 	secret1, _ := generateSecret()
// 	secret2, _ := generateSecret()
// 	if secret1 == secret2 {
// 		t.Errorf("expected random secrets, but they matched")
// 		t.Fail()
// 	}
// }

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
