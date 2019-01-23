package stack

import (
	"strings"
	"testing"
)

func Test_applyTemplateWithAuth(t *testing.T) {

	templateValues := authConfig{
		ClientId:     "7gbfgsbh9gbgbg786gs7bs",
		CustomersURL: "https://raw.githubusercontent.com/test/path/CUSTOMERS",
		Scheme:       "http",
	}

	templateFileName := "../../templates/of-auth-dep.yml"

	generatedValue, err := applyTemplate(templateFileName, templateValues)

	if err != nil {
		t.Errorf("expected no error generating template, but got %s", err.Error())
		t.Fail()
		return
	}

	values := []string{"7gbfgsbh9gbgbg786gs7bs", "https://raw.githubusercontent.com/test/path/CUSTOMERS"}
	for _, want := range values {
		if strings.Contains(string(generatedValue), want) == false {
			t.Errorf("want generated value to contain: %q, generated was: %q", want, string(generatedValue))
			t.Fail()
		}
	}
}
