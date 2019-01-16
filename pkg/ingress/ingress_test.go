package ingress

import (
	"strings"
	"testing"
)

func Test_applyTemplateWithTLS(t *testing.T) {
	templateValues := IngressTemplate{
		RootDomain: "test.com",
		TLS:        true,
	}

	templateFileName := "../../templates/k8s/ingress.yml"

	generatedValue, err := applyTemplate(templateFileName, templateValues)
	want := "tls"

	if err != nil {
		t.Errorf("expected no error generating template, but got %s", err.Error())
		t.Fail()
		return
	}

	if strings.Contains(string(generatedValue), want) == false {
		t.Errorf("want generated value to contain: %q, generated was: %q", want, string(generatedValue))
		t.Fail()
	}
}
