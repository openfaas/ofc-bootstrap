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

	gwConfigErr := generateTemplate("gateway_config", plan, gatewayConfig{
		Registry:     plan.Registry,
		RootDomain:   plan.RootDomain,
		CustomersURL: plan.CustomersURL,
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
		RootDomain: plan.RootDomain,
	})
	if dashboardConfigErr != nil {
		return dashboardConfigErr
	}

	return nil
}

func generateTemplate(fileName string, plan types.Plan, templateType interface{}) error {

	data, err := ioutil.ReadFile("templates/" + fileName + ".yml")
	if err != nil {
		return err
	}

	t := template.Must(template.New(fileName).Parse(string(data)))
	tempFilePath := "tmp/generated-" + fileName + ".yml"
	file, fileErr := os.Create(tempFilePath)
	if fileErr != nil {
		return fileErr
	}
	defer file.Close()

	executeErr := t.Execute(file, templateType)
	file.Close()

	if executeErr != nil {
		return executeErr
	}

	return nil
}
