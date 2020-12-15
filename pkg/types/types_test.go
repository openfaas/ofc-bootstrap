package types

import (
	"os"
	"testing"
)

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
