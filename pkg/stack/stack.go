package stack

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"os"

	"github.com/openfaas-incubator/ofc-bootstrap/pkg/types"
)

type gatewayConfig struct {
	Registry     string
	RootDomain   string
	CustomersURL string
	Scheme       string
	S3           types.S3
}

type authConfig struct {
	RootDomain   string
	ClientId     string
	CustomersURL string
	Scheme       string
}

// Apply creates `templates/gateway_config.yml` to be referenced by stack.yml
func Apply(plan types.Plan) error {
	scheme := "http"
	if plan.TLS {
		scheme += "s"
	}

	gwConfigErr := generateTemplate("gateway_config", plan, gatewayConfig{
		Registry:     plan.Registry,
		RootDomain:   plan.RootDomain,
		CustomersURL: plan.CustomersURL,
		Scheme:       scheme,
		S3:           plan.S3,
	})
	if gwConfigErr != nil {
		return gwConfigErr
	}

	githubConfigErr := generateTemplate("github", plan, types.Github{
		AppID:          plan.Github.AppID,
		PrivateKeyFile: plan.Github.PrivateKeyFile,
	})
	if githubConfigErr != nil {
		return githubConfigErr
	}

	dashboardConfigErr := generateTemplate("dashboard_config", plan, gatewayConfig{
		RootDomain: plan.RootDomain, Scheme: scheme,
	})
	if dashboardConfigErr != nil {
		return dashboardConfigErr
	}

	if plan.EnableOAuth {
		ofAuthDepErr := generateTemplate("of-auth-dep", plan, authConfig{
			RootDomain:   plan.RootDomain,
			ClientId:     plan.OAuth.ClientId,
			CustomersURL: plan.CustomersURL,
			Scheme:       scheme,
		})
		if ofAuthDepErr != nil {
			return ofAuthDepErr
		}
	}

	return nil
}

func applyTemplate(templateFileName string, templateType interface{}) ([]byte, error) {
	data, err := ioutil.ReadFile(templateFileName)
	if err != nil {
		return nil, err
	}
	t := template.Must(template.New(templateFileName).Parse(string(data)))

	buffer := new(bytes.Buffer)

	executeErr := t.Execute(buffer, templateType)

	return buffer.Bytes(), executeErr
}

func generateTemplate(fileName string, plan types.Plan, templateType interface{}) error {

	generatedData, err := applyTemplate("templates/"+fileName+".yml", templateType)
	if err != nil {
		return err
	}

	tempFilePath := "tmp/generated-" + fileName + ".yml"
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

	return nil
}
