package ingress

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"log"
	"os"

	"github.com/openfaas-incubator/ofc-bootstrap/pkg/execute"
	"github.com/openfaas-incubator/ofc-bootstrap/pkg/types"
)

type IngressTemplate struct {
	RootDomain string
	TLS        bool
	IssuerType string
	DNSService string
}

func Apply(plan types.Plan) error {

	err := apply("ingress-wildcard.yml", "ingress-wildcard", IngressTemplate{
		RootDomain: plan.RootDomain,
		TLS:        plan.TLS,
		IssuerType: plan.TLSConfig.IssuerType,
		DNSService: plan.TLSConfig.DNSService,
	})

	if err != nil {
		return err
	}

	err1 := apply("ingress.yml", "ingress", IngressTemplate{
		RootDomain: plan.RootDomain,
		TLS:        plan.TLS,
		IssuerType: plan.TLSConfig.IssuerType,
		DNSService: plan.TLSConfig.DNSService,
	})

	if err1 != nil {
		return err1
	}

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

func apply(source string, name string, ingress IngressTemplate) error {

	generatedData, err := applyTemplate("templates/k8s/"+source, ingress)
	if err != nil {
		return err
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
		Command: "kubectl apply -f " + tempFilePath,
		Shell:   false,
	}

	execRes, execErr := execTask.Execute()
	if execErr != nil {
		return execErr
	}

	log.Println(execRes.Stdout)

	return nil
}
