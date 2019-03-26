package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/openfaas-incubator/ofc-bootstrap/pkg/execute"
	"github.com/openfaas-incubator/ofc-bootstrap/pkg/ingress"
	"github.com/openfaas-incubator/ofc-bootstrap/pkg/stack"
	"github.com/openfaas-incubator/ofc-bootstrap/pkg/tls"
	"github.com/openfaas-incubator/ofc-bootstrap/pkg/types"
	yaml "gopkg.in/yaml.v2"
)

// Vars are variables parsed from flags
type Vars struct {
	YamlFile string
	Verbose  bool
}

const (
	// OrchestrationK8s uses Kubernetes
	OrchestrationK8s = "kubernetes"

	// OrchestrationSwarm uses Docker Swarm
	OrchestrationSwarm = "swarm"
)

func taskGivesStdout(tool string) error {
	task := execute.ExecTask{Command: tool}
	res, err := task.Execute()
	if err != nil {
		return fmt.Errorf("could not run: '%s', error: %s", tool, err)
	}
	if len(res.Stdout) == 0 {
		return fmt.Errorf("error executing '%s', no output was given - tool is available in PATH", task.Command)
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
		if featureEnabled(plan.Features, secret.Filters) {
			err := filesExists(secret.Files)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func filesExists(files []types.FileSecret) error {
	if len(files) > 0 {
		for _, file := range files {
			if len(file.ValueCommand) == 0 {
				if _, err := os.Stat(file.ExpandValueFrom()); err != nil {
					return err
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

	var featuresErr error
	plan, featuresErr = filterFeatures(plan)
	if featuresErr != nil {
		fmt.Fprintf(os.Stderr, "error while retreiving features: %s\n", featuresErr.Error())
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

		log.Printf("Is tiller rolled out?\n")
		tillerReady()

		if err := helmRepoUpdate(); err != nil {
			log.Println(err.Error())
		}

		installIngressErr := installIngressController(plan.Ingress)
		if installIngressErr != nil {
			log.Println(installIngressErr.Error())
		}

		createSecrets(plan)

		saErr := patchFnServiceaccount()
		if saErr != nil {
			log.Println(saErr)
		}

		minioErr := installMinio()
		if minioErr != nil {
			log.Println(minioErr)
		}

		if plan.TLS {
			cmErr := installCertmanager()
			if cmErr != nil {
				log.Println(cmErr)
			}
		}

		functionAuthErr := createFunctionsAuth()
		if functionAuthErr != nil {
			log.Println(functionAuthErr.Error())
		}

		ofErr := installOpenfaas()
		if ofErr != nil {
			log.Println(ofErr)
		}

		if plan.TLS {
			log.Printf("Is cert-manager rolled out?\n")
			certManagerReady()

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

		log.Printf("Are SealedSecrets rolled out?\n")
		sealedSecretsReady()

		pubCert := exportSealedSecretPubCert()
		writeErr := ioutil.WriteFile("tmp/pubcert.pem", []byte(pubCert), 0700)
		if writeErr != nil {
			log.Println(writeErr)
		}

		cloneErr := cloneCloudComponents()
		if (cloneErr) != nil {
			return cloneErr
		}

		openfaasGatewayReady()

		deployErr := deployCloudComponents(plan)
		if (deployErr) != nil {
			return deployErr
		}
	}

	return nil
}

func openfaasGatewayReady() {
	task := execute.ExecTask{
		Command: "/scripts/get-openfaas-gateway.sh",
		Shell:   true,
	}

	res, err := task.Execute()
	fmt.Println("cert-manager", res.ExitCode, res.Stdout, res.Stderr, err)
}

func helmRepoUpdate() error {
	log.Println("Updating helm repos")

	task := execute.ExecTask{
		Command: "helm repo update",
	}

	taskRes, taskErr := task.Execute()

	if taskErr != nil {
		return taskErr
	}

	log.Println(taskRes.Stdout)
	log.Println(taskRes.Stderr)

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

func installIngressController(ingress string) error {
	log.Println("Creating Ingress Controller")

	env := os.Environ()

	if ingress == "host" {
		env = append(os.Environ(), "ADDITIONAL_SET=,controller.hostNetwork=true,controller.daemonset.useHostPort=true,dnsPolicy=ClusterFirstWithHostNet,controller.kind=DaemonSet")
	}

	task := execute.ExecTask{
		Command: "scripts/install-nginx.sh",
		Shell:   true,
		Env:     env,
	}

	res, err := task.Execute()

	if err != nil {
		return err
	}

	log.Println(res.ExitCode, res.Stdout, res.Stderr)

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

	log.Println(res.ExitCode, res.Stdout, res.Stderr)

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

func patchFnServiceaccount() error {
	log.Println("Patching openfaas-fn serviceaccount for pull secrets")

	task := execute.ExecTask{
		Command: "scripts/patch-fn-serviceaccount.sh",
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

	log.Println(res.ExitCode, res.Stdout, res.Stderr)

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

	log.Println(res.ExitCode, res.Stdout, res.Stderr)

	return nil
}

func createSecrets(plan types.Plan) error {

	for _, secret := range plan.Secrets {
		if featureEnabled(plan.Features, secret.Filters) {
			fmt.Printf("Creating secret: %s\n", secret.Name)

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
	}

	return nil
}

func sealedSecretsReady() {

	task := execute.ExecTask{
		Command: "./scripts/get-sealedsecretscontroller.sh",
		Shell:   true,
	}

	res, err := task.Execute()
	fmt.Println("sealedsecretscontroller", res.ExitCode, res.Stdout, res.Stderr, err)

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

func certManagerReady() {
	task := execute.ExecTask{
		Command: "./scripts/get-cert-manager.sh",
		Shell:   true,
	}

	res, err := task.Execute()
	fmt.Println("cert-manager", res.ExitCode, res.Stdout, res.Stderr, err)
}

func tillerReady() {

	task := execute.ExecTask{
		Command: "./scripts/get-tiller.sh",
		Shell:   true,
	}

	res, err := task.Execute()
	fmt.Println("tiller", res.ExitCode, res.Stdout, res.Stderr, err)
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

	authEnv := ""
	if plan.EnableOAuth {
		authEnv = "ENABLE_OAUTH=true"
	}
	gitlabEnv := ""
	if plan.SCM == "gitlab" {
		gitlabEnv = "GITLAB=true"
	}

	task := execute.ExecTask{
		Command: "./scripts/deploy-cloud-components.sh",
		Shell:   true,
		Env:     []string{authEnv, gitlabEnv},
	}

	res, err := task.Execute()
	if err != nil {
		return err
	}

	fmt.Println(res)

	return nil
}

func featureEnabled(features []string, secretFeatures []string) bool {
	for _, feature := range features {
		for _, secretFeature := range secretFeatures {
			if feature == secretFeature {
				return true
			}
		}
	}
	return false
}

func filterFeatures(plan types.Plan) (types.Plan, error) {
	var err error

	plan.Features = append(plan.Features, types.DefaultFeature)

	plan, err = filterGitRepositoryManager(plan)
	if err != nil {
		return plan, fmt.Errorf("Error while filtering features: %s", err.Error())
	}

	if plan.TLS == true {
		plan, err = filterDNSFeature(plan)
		if err != nil {
			return plan, fmt.Errorf("Error while filtering features: %s", err.Error())
		}
	}

	if plan.EnableOAuth == true {
		plan.Features = append(plan.Features, types.Auth)
	}

	return plan, err
}

func filterDNSFeature(plan types.Plan) (types.Plan, error) {
	if plan.TLSConfig.DNSService == types.DigitalOcean {
		plan.Features = append(plan.Features, types.DODNS)
	} else if plan.TLSConfig.DNSService == types.CloudDNS {
		plan.Features = append(plan.Features, types.GCPDNS)
	} else if plan.TLSConfig.DNSService == types.Route53 {
		plan.Features = append(plan.Features, types.Route53DNS)
	} else {
		return plan, fmt.Errorf("Error unavailable DNS service provider: %s", plan.TLSConfig.DNSService)
	}
	return plan, nil
}

func filterGitRepositoryManager(plan types.Plan) (types.Plan, error) {
	if plan.SCM == types.GitLabSCM {
		plan.Features = append(plan.Features, types.GitLabFeature)
	} else if plan.SCM == types.GitHubSCM {
		plan.Features = append(plan.Features, types.GitHubFeature)
	} else {
		return plan, fmt.Errorf("Error unsupported Git repository manager: %s", plan.SCM)
	}
	return plan, nil
}
