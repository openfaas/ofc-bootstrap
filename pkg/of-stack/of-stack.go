package of_stack

import (
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/alexellis/ofc-bootstrap/pkg/execute"
)

func isValidUrl(s string) bool {
	if _, err := url.ParseRequestURI(s); err != nil {
		return false
	}

	return true
}

func checkFileExists(s string) bool {
	fi, err := os.Stat(s)

	return err == nil && !fi.IsDir()
}

func Deploy(path string) error {
	if path == "" {
		return fmt.Errorf("provided path can't be empty")
	}

	if !isValidUrl(path) && !checkFileExists(path) {
		return fmt.Errorf("file %s doesn't exist or is a directory", path)
	}

	execTask := execute.ExecTask{
		Command: "faas-cli deploy -f " + path,
		Shell:   false,
	}

	execRes, execErr := execTask.Execute()

	if execErr != nil {
		return execErr
	}

	log.Println(execRes.Stdout)

	return nil
}
