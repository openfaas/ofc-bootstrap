package types

import (
	"crypto/sha256"
	"fmt"
	"log"
	"os"

	"github.com/openfaas-incubator/ofc-bootstrap/pkg/execute"
)

func CreateDockerSecret(sec DockerSecret) (string, string) {
	secretCmd := "docker secret create %s %s"

	if len(sec.File.ValueCommand) > 0 {
		task := execute.ExecTask{
			Command: sec.File.ValueCommand,
		}
		_, err := task.Execute()

		if err != nil {
			log.Println(err)
		}
	}

	if len(sec.File.ValueFrom) > 0 {
		return fmt.Sprintf(secretCmd, sec.Name, sec.File.ExpandValueFrom()), ""
	}

	if len(sec.Value) > 0 {
		return fmt.Sprintf(secretCmd, sec.Name, "-"), sec.Value
	}

	val, err := generateSecret()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return fmt.Sprintf(secretCmd, sec.Name, "-"), val
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

		secretCmd = fmt.Sprintf(`%s --from-literal=%s=%s`, secretCmd, key.Name, secretValue)
	}

	for _, file := range kvn.Files {
		if len(file.ValueCommand) > 0 {
			task := execute.ExecTask{
				Command: file.ValueCommand,
			}
			_, err := task.Execute()

			if err != nil {
				log.Println(err)
			}
		}

		secretCmd = fmt.Sprintf("%s --from-file=%s=%s", secretCmd, file.Name, file.ExpandValueFrom())
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
