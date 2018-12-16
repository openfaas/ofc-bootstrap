package ingress

import (
	"io/ioutil"
	"log"
	"os"
	"text/template"

	"github.com/alexellis/ofc-bootstrap/pkg/execute"
	"github.com/alexellis/ofc-bootstrap/pkg/types"
)

type IngressTemplate struct {
	RootDomain string
}

func Apply(plan types.Plan) error {

	err := apply("ingress-wildcard.yml", "ingress-wildcard", IngressTemplate{
		RootDomain: plan.RootDomain,
	})

	if err != nil {
		return err
	}

	err1 := apply("ingress.yml", "ingress", IngressTemplate{
		RootDomain: plan.RootDomain,
	})

	if err1 != nil {
		return err1
	}

	return nil
}

func apply(source string, name string, ingress IngressTemplate) error {

	data, err := ioutil.ReadFile("templates/k8s/" + source)
	if err != nil {
		return err
	}

	t := template.Must(template.New(name).Parse(string(data)))
	tempFilePath := "tmp/generated-ingress-" + name + ".yaml"
	file, fileErr := os.Create(tempFilePath)
	if fileErr != nil {
		return fileErr
	}
	defer file.Close()

	executeErr := t.Execute(file, IngressTemplate{
		RootDomain: ingress.RootDomain,
	})
	file.Close()

	if executeErr != nil {
		return executeErr
	}

	execTask := execute.ExecTask{
		Command: "kubectl apply -f " + tempFilePath,
		Shell:   false,
	}

	execRes, execErr := execTask.Execute()
	if err != nil {
		return execErr
	}

	log.Println(execRes.Stdout)

	return nil
}
