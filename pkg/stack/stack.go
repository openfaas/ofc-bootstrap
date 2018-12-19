package stack

import (
	"html/template"
	"io/ioutil"
	"os"

	"github.com/alexellis/ofc-bootstrap/pkg/types"
)

type gatewayConfig struct {
	Registry     string
	RootDomain   string
	CustomersURL string
}

// Apply creates `templates/gateway_config.yml` to be referenced by stack.yml
func Apply(plan types.Plan) error {

	gwConfigErr := generateGatewayConfig(plan)
	if gwConfigErr != nil {
		return gwConfigErr
	}

	githubConfigErr := generateGithubConfig(plan)
	if githubConfigErr != nil {
		return githubConfigErr
	}

	return nil
}

func generateGatewayConfig(plan types.Plan) error {

	name := "gateway_config"
	data, err := ioutil.ReadFile("templates/gateway_config.yml")
	if err != nil {
		return err
	}

	t := template.Must(template.New(name).Parse(string(data)))
	tempFilePath := "tmp/generated-" + name + ".yml"
	file, fileErr := os.Create(tempFilePath)
	if fileErr != nil {
		return fileErr
	}
	defer file.Close()

	executeErr := t.Execute(file, gatewayConfig{
		Registry:     plan.Registry,
		RootDomain:   plan.RootDomain,
		CustomersURL: plan.CustomersURL,
	})
	file.Close()

	if executeErr != nil {
		return executeErr
	}

	return nil
}

func generateGithubConfig(plan types.Plan) error {

	name := "github"
	data, err := ioutil.ReadFile("templates/github.yml")
	if err != nil {
		return err
	}

	t := template.Must(template.New(name).Parse(string(data)))
	tempFilePath := "tmp/generated-" + name + ".yml"
	file, fileErr := os.Create(tempFilePath)
	if fileErr != nil {
		return fileErr
	}
	defer file.Close()

	executeErr := t.Execute(file, types.Github{
		AppID:          plan.Github.AppID,
		PrivateKeyFile: plan.Github.PrivateKeyFile,
	})
	file.Close()

	if executeErr != nil {
		return executeErr
	}

	return nil
}
