package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/alexellis/ofc-bootstrap/pkg/execute"
	"github.com/alexellis/ofc-bootstrap/pkg/ingress"
	"github.com/alexellis/ofc-bootstrap/pkg/stack"
	"github.com/alexellis/ofc-bootstrap/pkg/tls"
	"github.com/alexellis/ofc-bootstrap/pkg/types"
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

func taskGivesStdout(tool string) error {
	task := execute.ExecTask{Command: tool}
	res, err := task.Execute()
	if err != nil {
		return fmt.Errorf("could not run: '%s', error: %s", tool, err)
	}
	if len(res.Stdout) == 0 {
		return fmt.Errorf("error executing '%s', no output was given - tool is available in PATH")
	}
	return nil
}

func validateTools(tools []string) error {

	for _, tool := range tools {
		err := taskGivesStdout(tool)
		if err != nil {
			return err
		}
	}

	return nil

}

func validatePlan(plan types.Plan) error {
	for _, secret := range plan.Secrets {
		if len(secret.Files) > 0 {
			for _, file := range secret.Files {
				if len(file.ValueCommand) == 0 {
					if _, err := os.Stat(file.ExpandValueFrom()); err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

func main() {

	vars := Vars{}
	flag.StringVar(&vars.YamlFile, "yaml", "init.yaml", "YAML file for bootstrap")
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

	plan := types.Plan{}
	unmarshalErr := yaml.Unmarshal(yamlBytes, &plan)
	if unmarshalErr != nil {
		fmt.Fprintf(os.Stderr, "-yaml file gave error: %s\n", unmarshalErr.Error())
		os.Exit(1)
	}

	log.Println("Validating tools available in PATH")

	var tools []string
	if plan.Orchestration == OrchestrationK8s {
		tools = []string{"kubectl version --client", "openssl version", "helm version -c", "faas-cli version"}
	}

	validateToolsErr := validateTools(tools)

	if validateToolsErr != nil {
		panic(validateToolsErr)
	}

	validateErr := validatePlan(plan)
	if validateErr != nil {
		panic(validateErr)

	}

	fmt.Fprintf(os.Stdout, "Plan loaded from: %s\n", vars.YamlFile)

	os.Mkdir("tmp", 0700)

	start := time.Now()
	err := process(plan)
	done := time.Since(start)

	if err != nil {
		fmt.Fprintf(os.Stdout, "Plan failed after %f seconds\nError: %s", done.Seconds(), err.Error())

		os.Exit(1)
	}

	fmt.Fprintf(os.Stdout, "Plan completed in %f seconds\n", done.Seconds())
}

func process(plan types.Plan) error {

	fmt.Println(plan)
	if plan.Orchestration == OrchestrationK8s {
		fmt.Println("Orchestration: Kubernetes")
	} else if plan.Orchestration == OrchestrationSwarm {
		fmt.Println("Orchestration: Swarm")
	}

	if plan.Orchestration == OrchestrationK8s {

		nsErr := createNamespaces()
		if nsErr != nil {
			log.Println(nsErr)
		}

		fmt.Println("Building Ingress")
		tillerErr := installTiller()
		if tillerErr != nil {
			log.Println(tillerErr)
		}

		retries := 260
		for i := 0; i < retries; i++ {
			log.Printf("Is tiller ready? %d/%d\n", i+1, retries)
			ready := tillerReady()
			if ready {
				break
			}
			time.Sleep(time.Second * 2)
		}

		installIngressErr := installIngressController()
		if installIngressErr != nil {
			log.Println(installIngressErr.Error())
		}

		createSecrets(plan)

		minioErr := installMinio()
		if minioErr != nil {
			log.Println(minioErr)
		}

		cmErr := installCertmanager()
		if cmErr != nil {
			log.Println(cmErr)
		}

		functionAuthErr := createFunctionsAuth()
		if functionAuthErr != nil {
			log.Println(functionAuthErr.Error())
		}

		ofErr := installOpenfaas()
		if ofErr != nil {
			log.Println(ofErr)
		}

		ingressErr := ingress.Apply(plan)
		if ingressErr != nil {
			log.Println(ingressErr)
		}

		if plan.TLS {
			tlsErr := tls.Apply(plan)
			if tlsErr != nil {
				log.Println(tlsErr)
			}
		}

		fmt.Println("Creating stack.yml")

		planErr := stack.Apply(plan)
		if planErr != nil {
			log.Println(planErr)
		}

		sealedSecretsErr := installSealedSecrets()
		if sealedSecretsErr != nil {
			log.Println(sealedSecretsErr)
		}

		for i := 0; i < retries; i++ {
			log.Printf("Are SealedSecrets ready? %d/%d\n", i+1, retries)
			ready := sealedSecretsReady()
			if ready {
				break
			}
			time.Sleep(time.Second * 2)
		}

		pubCert := exportSealedSecretPubCert()
		writeErr := ioutil.WriteFile("tmp/pubcert.pem", []byte(pubCert), 0700)
		if writeErr != nil {
			log.Println(writeErr)
		}

		cloneErr := cloneCloudComponents()
		if (cloneErr) != nil {
			return cloneErr
		}

		deployErr := deployCloudComponents(plan)
		if (deployErr) != nil {
			return deployErr
		}
	}

	return nil
}

func installTiller() error {
	log.Println("Creating Tiller")

	task1 := execute.ExecTask{
		Command: "scripts/create-tiller-sa.sh",
		Shell:   true,
	}

	res1, err1 := task1.Execute()

	if err1 != nil {
		return err1
	}

	log.Println(res1.Stdout)
	log.Println(res1.Stderr)

	task2 := execute.ExecTask{
		Command: "scripts/create-tiller.sh",
		Shell:   true,
	}

	res2, err2 := task2.Execute()

	if err2 != nil {
		return err2
	}

	log.Println(res2.Stdout)
	log.Println(res2.Stderr)

	return nil
}

func createFunctionsAuth() error {
	log.Println("Creating secrets for functions to consume")

	task := execute.ExecTask{
		Command: "scripts/create-functions-auth.sh",
		Shell:   true,
	}

	res, err := task.Execute()

	if err != nil {
		return err
	}

	log.Println(res.Stdout)
	log.Println(res.Stderr)

	return nil
}

func installIngressController() error {
	log.Println("Creating Ingress Controller")

	task := execute.ExecTask{
		Command: "scripts/install-nginx.sh",
		Shell:   true,
	}

	res, err := task.Execute()

	if err != nil {
		return err
	}

	log.Println(res.Stdout)
	log.Println(res.Stderr)

	return nil
}

func installSealedSecrets() error {
	log.Println("Creating SealedSecrets")

	task := execute.ExecTask{
		Command: "scripts/install-sealedsecrets.sh",
		Shell:   true,
	}

	res, err := task.Execute()

	if err != nil {
		return err
	}

	log.Println(res.Stdout)
	log.Println(res.Stderr)

	return nil
}

func installOpenfaas() error {
	log.Println("Creating OpenFaaS")

	task := execute.ExecTask{
		Command: "scripts/install-openfaas.sh",
		Shell:   true,
	}

	res, err := task.Execute()

	if err != nil {
		return err
	}

	log.Println(res.Stdout)
	log.Println(res.Stderr)

	return nil
}

func installMinio() error {
	log.Println("Creating Minio")

	task := execute.ExecTask{
		Command: "scripts/install-minio.sh",
		Shell:   true,
	}

	res, err := task.Execute()

	if err != nil {
		return err
	}

	log.Println(res.Stdout)
	log.Println(res.Stderr)

	return nil
}

func installCertmanager() error {
	log.Println("Creating Cert-Manager")

	task := execute.ExecTask{
		Command: "scripts/install-cert-manager.sh",
		Shell:   true,
	}

	res, err := task.Execute()

	if err != nil {
		return err
	}

	log.Println(res.Stdout)
	log.Println(res.Stderr)

	return nil
}

func createNamespaces() error {
	log.Println("Creating namespaces")

	task := execute.ExecTask{
		Command: "scripts/create-namespaces.sh",
		Shell:   true,
	}

	res, err := task.Execute()

	if err != nil {
		return err
	}

	log.Println(res.Stdout)
	log.Println(res.Stderr)

	return nil
}

func createSecrets(plan types.Plan) error {

	fmt.Println(plan.Secrets)

	for _, secret := range plan.Secrets {

		var command execute.ExecTask
		if plan.Orchestration == OrchestrationK8s {
			command = execute.ExecTask{
				Command: types.CreateK8sSecret(secret),
				Shell:   false,
			}
		} else if plan.Orchestration == OrchestrationSwarm {
			command = execute.ExecTask{Command: types.CreateDockerSecret(secret)}
		}

		res, err := command.Execute()

		if err != nil {
			log.Println(err)
		}

		fmt.Println(res)
	}

	return nil
}

func sealedSecretsReady() bool {

	task := execute.ExecTask{
		Command: "./scripts/get-sealedsecretscontroller.sh",
		Shell:   true,
	}

	res, err := task.Execute()
	fmt.Println("sealedsecretscontroller", res.ExitCode, res.Stdout, res.Stderr, err)
	return res.Stdout == "1"
}

func exportSealedSecretPubCert() string {

	task := execute.ExecTask{
		Command: "./scripts/export-sealed-secret-pubcert.sh",
		Shell:   true,
	}

	res, err := task.Execute()
	fmt.Println("secrets cert", res.ExitCode, res.Stdout, res.Stderr, err)
	return res.Stdout
}

func tillerReady() bool {

	task := execute.ExecTask{
		Command: "./scripts/get-tiller.sh",
		Shell:   true,
	}

	res, err := task.Execute()
	fmt.Println("tiller", res.ExitCode, res.Stdout, res.Stderr, err)
	return res.Stdout == "1"
}

func cloneCloudComponents() error {
	task := execute.ExecTask{
		Command: "./scripts/clone-cloud-components.sh",
		Shell:   true,
	}

	res, err := task.Execute()
	if err != nil {
		return err
	}

	fmt.Println(res)

	return nil
}

func deployCloudComponents(plan types.Plan) error {

	env := ""
	if plan.EnableOAuth {
		env = "ENABLE_OAUTH=true"
	}
	task := execute.ExecTask{
		Command: "./scripts/deploy-cloud-components.sh",
		Shell:   true,
		Env:     []string{env},
	}

	res, err := task.Execute()
	if err != nil {
		return err
	}

	fmt.Println(res)

	return nil
}
