package types

import (
	"fmt"
	"log"
	"os"

	execute "github.com/alexellis/go-execute/pkg/v1"
	"github.com/sethvargo/go-password/password"
)

func BuildSecretTask(kvn KeyValueNamespaceTuple) execute.ExecTask {
	task := execute.ExecTask{
		Command:     "kubectl",
		Args:        []string{"create", "secret", "generic", "-n=" + kvn.Namespace, kvn.Name},
		StreamStdio: false,
	}

	if len(kvn.Type) > 0 {
		task.Args = append(task.Args, "--type="+kvn.Type)
	}

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
		task.Args = append(task.Args, fmt.Sprintf("--from-literal=%s=%s", key.Name, secretValue))
	}

	for _, file := range kvn.Files {
		filePath := file.ExpandValueFrom()
		if len(file.ValueCommand) > 0 {
			if _, err := os.Stat(filePath); err != nil {

				valueTask := execute.ExecTask{
					Command:     file.ValueCommand,
					StreamStdio: false,
				}
				res, err := valueTask.Execute()
				if err != nil {
					log.Fatal(fmt.Errorf("error executing value_command: %s", file.ValueCommand))
				}

				if res.ExitCode != 0 {
					log.Fatal(fmt.Errorf("error running value_command: %s, stderr: %s", file.ValueCommand, res.Stderr))
				}
			} else {
				fmt.Printf("%s exists, not running value_command\n", filePath)
			}
		}

		task.Args = append(task.Args, fmt.Sprintf("--from-file=%s=%s", file.Name, file.ExpandValueFrom()))

	}

	return task
}

func generateSecret() (string, error) {
	pass, err := password.Generate(25, 10, 0, false, true)
	if err != nil {
		return "", err
	}
	return pass, nil
}
