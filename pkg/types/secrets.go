package types

import (
	"crypto/sha256"
	"fmt"
	"os"

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
	secretCmd := fmt.Sprintf("kubectl create secret generic -n %s %s", kvn.Namespace, kvn.Name)

	for _, key := range kvn.Literals {
		secretValue := key.Value
		if len(secretValue) == 0 {
			val, err := generateSecret()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			secretValue = val
		}

		secretCmd = fmt.Sprintf("%s --from-literal=%s=%s", secretCmd, key.Name, secretValue)
	}
	return secretCmd
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
