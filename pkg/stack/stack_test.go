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

	templateFileName := "../../templates/edge-auth-dep.yml"

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

func Test_gitlabTemplates(t *testing.T) {
	gitLabInstance := "https://gitlab.test.o6s.io/"

	gitlabTemplateFileName := "../../templates/gitlab.yml"

	generatedValue, err := applyTemplate(gitlabTemplateFileName, gitlabConfig{
		GitLabInstance:      gitLabInstance,
		CustomersSecretPath: "",
	})

	if err != nil {
		t.Errorf("expected no error generating template, but got %s", err.Error())
		t.Fail()
		return
	}

	want := gitLabInstance
	if strings.Contains(string(generatedValue), want) == false {
		t.Errorf("want generated value to contain: %q, generated was: %q", want, string(generatedValue))
		t.Fail()
	}
}
