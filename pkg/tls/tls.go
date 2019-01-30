package tls

import (
	"html/template"
	"io/ioutil"
	"log"
	"os"

	"github.com/openfaas-incubator/ofc-bootstrap/pkg/execute"
	"github.com/openfaas-incubator/ofc-bootstrap/pkg/types"
)

// TLSTemplate TLS configuration
type TLSTemplate struct {
	RootDomain              string
	Email                   string
	DNSService              string
	ProjectID               string
	IssuerType              string
	Region                  string
	AccessKeyID             string
	DigitalOceanAccessToken string
}

var tlsTemplatesPath = "templates/k8s/tls/"

// Apply executes the plan
func Apply(plan types.Plan) error {

	tlsTemplatesList, _ := listTLSTemplates()
	tlsTemplate := TLSTemplate{
		RootDomain:              plan.RootDomain,
		Email:                   plan.TLSConfig.Email,
		DNSService:              plan.TLSConfig.DNSService,
		ProjectID:               plan.TLSConfig.ProjectID,
		IssuerType:              plan.TLSConfig.IssuerType,
		Region:                  plan.TLSConfig.Region,
		AccessKeyID:             plan.TLSConfig.AccessKeyID,
		DigitalOceanAccessToken: plan.TLSConfig.DigitalOceanAccessToken,
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
	file, err := os.Open(tlsTemplatesPath)

	if err != nil {
		log.Fatalf("failed opening directory: %s, %s", tlsTemplatesPath, err)
		return nil, err
	}
	defer file.Close()

	list, _ := file.Readdirnames(0)
	if err != nil {
		log.Fatalf("failed reading filenames in directory %s, %s", tlsTemplatesPath, err)
		return nil, err
	}
	return list, nil
}

func generateTemplate(fileName string, tlsTemplate TLSTemplate) (string, error) {

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
