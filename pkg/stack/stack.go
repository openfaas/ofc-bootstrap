package stack

import (
	"html/template"
	"io/ioutil"
	"os"

	"github.com/alexellis/ofc-bootstrap/pkg/types"
)

type Stack struct {
	Registry     string
	RootDomain   string
	CustomersURL string
}

func Apply(plan types.Plan) error {

	name := "stack"
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

	executeErr := t.Execute(file, Stack{
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
