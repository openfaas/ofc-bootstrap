package ingress

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"

	execute "github.com/alexellis/go-execute/pkg/v1"
	"github.com/openfaas/ofc-bootstrap/pkg/types"
)

type IngressTemplate struct {
	RootDomain string
	TLS        bool
	IssuerType string
}

// Apply templates and applies any ingress records required
// for the OpenFaaS Cloud ingress configuration
func Apply(plan types.Plan) error {

	if err := apply("ingress-wildcard.yml", "ingress-wildcard", IngressTemplate{
		RootDomain: plan.RootDomain,
		TLS:        plan.TLS,
		IssuerType: plan.TLSConfig.IssuerType,
	}); err != nil {
		return err
	}

	if err := apply("ingress-auth.yml", "ingress-auth", IngressTemplate{
		RootDomain: plan.RootDomain,
		TLS:        plan.TLS,
		IssuerType: plan.TLSConfig.IssuerType,
	}); err != nil {
		return err
	}

	return nil
}

func apply(source string, name string, ingress IngressTemplate) error {

	generatedData, err := applyTemplate("templates/k8s/"+source, ingress)
	if err != nil {
		return fmt.Errorf("unable to read template %s (%s), error: %q", name, "templates/k8s/"+source, err)
	}

	tempFilePath := "tmp/generated-ingress-" + name + ".yaml"
	file, err := os.Create(tempFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(generatedData)
	file.Close()

	if err != nil {
		return err
	}

	execTask := execute.ExecTask{
		Command:     "kubectl",
		Args:        []string{"apply", "-f", tempFilePath},
		Shell:       false,
		StreamStdio: false,
	}

	execRes, err := execTask.Execute()
	if err != nil {
		return err
	}

	log.Println(execRes.Stdout, execRes.Stderr)

	return nil
}

func applyTemplate(templateFileName string, templateValues IngressTemplate) ([]byte, error) {
	data, err := ioutil.ReadFile(templateFileName)
	if err != nil {
		return nil, err
	}

	t := template.Must(template.New(templateFileName).Parse(string(data)))
	buffer := new(bytes.Buffer)
	if err := t.Execute(buffer, templateValues); err != nil {
		return []byte{}, err
	}

	return buffer.Bytes(), nil
}
