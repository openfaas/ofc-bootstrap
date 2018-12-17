package types

import (
	"crypto/sha256"
	"fmt"
	"os"
	"strings"

	"github.com/alexellis/ofc-bootstrap/pkg/execute"
)

func CreateDockerSecret(kvn KeyValueNamespaceTuple) string {
	val, err := generateSecret()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return fmt.Sprintf("echo %s | docker secret create %s", val, kvn.Name)
}

func CreateK8sSecret(kvn KeyValueNamespaceTuple) string {
	val, err := generateSecret()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return fmt.Sprintf("kubectl create secret generic -n %s %s --from-literal=%s=%s", kvn.Namespace, kvn.Name, kvn.Name, strings.TrimSpace(val))
}

func generateSecret() (string, error) {
	task := execute.ExecTask{
		Command: "scripts/generate-sha.sh",
		Shell:   false,
	}

	res, err := task.Execute()
	if res.ExitCode != 0 && err != nil {
		err = fmt.Errorf("non-zero exit code")
	}

	h := sha256.New()
	h.Write([]byte(res.Stdout))

	return fmt.Sprintf("%x", h.Sum(nil)), err
}
