package stack

import (
	"strings"
	"testing"
)

func Test_applyTemplateWithAuth(t *testing.T) {

	clientID := "test_oauth_app_client_id"
	customersURL := "https://raw.githubusercontent.com/test/path/CUSTOMERS"

	templateValues := authConfig{
		ClientId:     clientID,
		CustomersURL: customersURL,
		Scheme:       "http",
	}

	templateFileName := "../../templates/k8s/of-auth-dep.yml"

	generatedValue, err := applyTemplate(templateFileName, templateValues)

	if err != nil {
		t.Errorf("expected no error generating template, but got %s", err.Error())
		t.Fail()
		return
	}

	values := []string{clientID, customersURL}
	for _, want := range values {
		if strings.Contains(string(generatedValue), want) == false {
			t.Errorf("want generated value to contain: %q, generated was: %q", want, string(generatedValue))
			t.Fail()
		}
	}
}
