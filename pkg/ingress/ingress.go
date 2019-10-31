package ingress

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"

	"github.com/alexellis/go-execute"
	"github.com/openfaas-incubator/ofc-bootstrap/pkg/types"
)

type IngressTemplate struct {
	RootDomain string
	TLS        bool
	IssuerType string
}

// Apply templates and applies any ingress records required
// for the OpenFaaS Cloud ingress configuration
func Apply(plan types.Plan) error {

	err := apply("ingress-wildcard.yml", "ingress-wildcard", IngressTemplate{
		RootDomain: plan.RootDomain,
		TLS:        plan.TLS,
		IssuerType: plan.TLSConfig.IssuerType,
	}, plan.DryRun)

	if err != nil {
		return err
	}

	err1 := apply("ingress-auth.yml", "ingress-auth", IngressTemplate{
		RootDomain: plan.RootDomain,
		TLS:        plan.TLS,
		IssuerType: plan.TLSConfig.IssuerType,
	}, plan.DryRun)

	if err1 != nil {
		return err1
	}

	return nil
}

func apply(source string, name string, ingress IngressTemplate, dryRun bool) error {

	generatedData, err := applyTemplate("templates/k8s/"+source, ingress)
	if err != nil {
		return fmt.Errorf("unable to read template %s (%s), error: %q", name, "templates/k8s/"+source, err)
	}

	tempFilePath := "tmp/generated-ingress-" + name + ".yaml"
	file, fileErr := os.Create(tempFilePath)
	if fileErr != nil {
		return fileErr
	}
	defer file.Close()

	_, writeErr := file.Write(generatedData)
	file.Close()

	if writeErr != nil {
		return writeErr
	}

	execTask := execute.ExecTask{
		Command: "kubectl",
		Args:    []string{"apply", "-f", tempFilePath},
		Shell:   false,
	}

	if dryRun {
		execTask.Args = append(execTask.Args, "--dry-run=true")
	}

	execRes, execErr := execTask.Execute()
	if execErr != nil {
		return execErr
	}

	log.Println(execRes.ExitCode, execRes.Stdout, execRes.Stderr)

	return nil
}

func applyTemplate(templateFileName string, templateValues IngressTemplate) ([]byte, error) {
	data, err := ioutil.ReadFile(templateFileName)
	if err != nil {
		return nil, err
	}
	t := template.Must(template.New(templateFileName).Parse(string(data)))

	buffer := new(bytes.Buffer)

	executeErr := t.Execute(buffer, templateValues)

	return buffer.Bytes(), executeErr
}
