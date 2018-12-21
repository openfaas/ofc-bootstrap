package types

import (
	"os"
	"testing"
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

func TestFileSecret_ExpandValueFrom(t *testing.T) {
	os.Setenv("HOME", "/home/user")
	fs := FileSecret{
		ValueFrom: "~/.docker/config.json",
	}
	want := "/home/user/.docker/config.json"
	got := fs.ExpandValueFrom()
	if got != want {
		t.Errorf("error want: %s, got %s", want, got)
		t.Fail()
	}
}
