package tls

import (
	"html/template"
	"io/ioutil"
	"log"
	"os"

	execute "github.com/alexellis/go-execute/pkg/v1"
	"github.com/openfaas/ofc-bootstrap/pkg/types"
)

// TLSTemplate TLS configuration
type TLSTemplate struct {
	RootDomain  string
	Email       string
	DNSService  string
	ProjectID   string
	IssuerType  string
	Region      string
	AccessKeyID string
}

// Apply executes the plan
func Apply(plan types.Plan) error {

	tlsTemplatesList, _ := listTLSTemplates()
	tlsTemplate := TLSTemplate{
		RootDomain:  plan.RootDomain,
		Email:       plan.TLSConfig.Email,
		DNSService:  plan.TLSConfig.DNSService,
		ProjectID:   plan.TLSConfig.ProjectID,
		IssuerType:  plan.TLSConfig.IssuerType,
		Region:      plan.TLSConfig.Region,
		AccessKeyID: plan.TLSConfig.AccessKeyID,
	}

	for _, template := range tlsTemplatesList {
		tempFilePath, tlsTemplateErr := generateTemplate(template, tlsTemplate)
		if tlsTemplateErr != nil {
			return tlsTemplateErr
		}

		applyErr := applyTemplate(tempFilePath)
		if applyErr != nil {
			return applyErr
		}
	}

	return nil
}

func listTLSTemplates() ([]string, error) {

	return []string{
		"issuer-prod.yml",
		"issuer-staging.yml",
		"wildcard-domain-cert.yml",
		"auth-domain-cert.yml",
	}, nil
}

func generateTemplate(fileName string, tlsTemplate TLSTemplate) (string, error) {
	tlsTemplatesPath := "templates/k8s/tls/"

	data, err := ioutil.ReadFile(tlsTemplatesPath + fileName)
	if err != nil {
		return "", err
	}

	t := template.Must(template.New(fileName).Parse(string(data)))
	tempFilePath := "tmp/generated-tls-" + fileName
	file, fileErr := os.Create(tempFilePath)
	if fileErr != nil {
		return "", fileErr
	}
	defer file.Close()

	executeErr := t.Execute(file, tlsTemplate)

	if executeErr != nil {
		return "", executeErr
	}

	return tempFilePath, nil
}

func applyTemplate(tempFilePath string) error {

	execTask := execute.ExecTask{
		Command:     "kubectl apply -f " + tempFilePath,
		Shell:       false,
		StreamStdio: false,
	}

	execRes, execErr := execTask.Execute()
	if execErr != nil {
		return execErr
	}

	log.Println(execRes.ExitCode, execRes.Stdout, execRes.Stderr)

	return nil
}
