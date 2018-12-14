package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"

	yaml "gopkg.in/yaml.v2"
)

type Vars struct {
	YamlFile string
	Verbose  bool
}

const (
	OrchestrationK8s   = "kubernetes"
	OrchestrationSwarm = "swarm"
)

type Plan struct {
	Orchestration string                   `yaml:"orchestration"`
	Secrets       []KeyValueNamespaceTuple `yaml:"secrets"`
}

type KeyValueNamespaceTuple struct {
	Name      string `yaml:"name"`
	Value     string `yaml:"value"`
	Namespace string `yaml:"namespace"`
}

func main() {

	vars := Vars{}
	flag.StringVar(&vars.YamlFile, "yaml", "", "YAML file for bootstrap")
	flag.BoolVar(&vars.Verbose, "verbose", false, "control verbosity")
	flag.Parse()

	if len(vars.YamlFile) == 0 {
		fmt.Fprintf(os.Stderr, "No -yaml flag given\n")
		os.Exit(1)
	}

	yamlBytes, yamlErr := ioutil.ReadFile(vars.YamlFile)
	if yamlErr != nil {
		fmt.Fprintf(os.Stderr, "-yaml file gave error: %s\n", yamlErr.Error())
		os.Exit(1)
	}

	plan := Plan{}
	unmarshalErr := yaml.Unmarshal(yamlBytes, &plan)
	if unmarshalErr != nil {
		fmt.Fprintf(os.Stderr, "-yaml file gave error: %s\n", unmarshalErr.Error())
		os.Exit(1)
	}

	fmt.Fprintf(os.Stdout, "Plan loaded from: %s\n", vars.YamlFile)

	start := time.Now()
	err := process(plan)
	done := time.Since(start)

	if err != nil {
		fmt.Fprintf(os.Stdout, "Plan failed after %f seconds\nError: %s", done.Seconds(), err.Error())

		os.Exit(1)
	}

	fmt.Fprintf(os.Stdout, "Plan completed in %f seconds\n", done.Seconds())
}

func process(plan Plan) error {

	fmt.Println(plan)
	if plan.Orchestration == OrchestrationK8s {
		fmt.Println("Orchestration: Kubernetes")
	} else if plan.Orchestration == OrchestrationSwarm {
		fmt.Println("Orchestration: Swarm")
	}

	createSecrets(plan)

	return nil
}

func createSecrets(plan Plan) error {

	fmt.Println(plan.Secrets)

	var command ExecTask
	if plan.Orchestration == OrchestrationK8s {
		command = ExecTask{Command: createK8sSecret(plan.Secrets[0])}
	} else if plan.Orchestration == OrchestrationSwarm {
		command = ExecTask{Command: createDockerSecret(plan.Secrets[0])}
	}

	res, err := command.Execute()
	fmt.Println(res)

	return err
}

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
		args := []string{"-c", "\""}
		for _, part := range startArgs {
			args = append(args, part)
		}
		args = append(args, "\"")

		fmt.Println(args)

		cmd = exec.Command("/bin/bash", args...)
	} else {
		cmd = exec.Command(et.Command)
	}

	stdoutPipe, stdoutPipeErr := cmd.StdoutPipe()
	if stdoutPipeErr != nil {
		return ExecResult{}, stdoutPipeErr
	}

	startErr := cmd.Start()
	if startErr != nil {
		return ExecResult{}, startErr
	}

	stdoutBytes, err := ioutil.ReadAll(stdoutPipe)
	if err != nil {
		return ExecResult{}, err
	}
	fmt.Println("res: " + string(stdoutBytes))

	return ExecResult{
		Stdout: string(stdoutBytes),
	}, nil
}

func generateSecret() (string, error) {
	task := ExecTask{
		Command: "head -c 16 /dev/urandom | shasum",
		Shell:   true,
	}

	res, err := task.Execute()
	if res.ExitCode != 0 && err != nil {
		err = fmt.Errorf("non-zero exit code")
	}

	return res.Stdout, err
}

func createDockerSecret(kvn KeyValueNamespaceTuple) string {
	val, err := generateSecret()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return fmt.Sprintf("echo %s | docker secret create %s", val, kvn.Name)
}

func createK8sSecret(kvn KeyValueNamespaceTuple) string {
	val, err := generateSecret()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return fmt.Sprintf("kubectl create secret generic -n %s %s --from-literal s3-access-key=\"%s\"", kvn.Namespace, kvn.Name, val)
}
