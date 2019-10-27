package tls

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"testing"
)

func Test_DigitalOcean_Issuer(t *testing.T) {
	tlsTemplate := TLSTemplate{
		Email:      "sales@openfaas.com",
		IssuerType: "ClusterIssuer",
		DNSService: "digitalocean",
	}

	templatePath := "../../templates/k8s/tls/issuer-prod.yml"
	templateData, err := ioutil.ReadFile(templatePath)
	if err != nil {
		t.Error(err)
		return
	}

	templateRes := template.Must(template.New("prod-issuer").Parse(string(templateData)))
	buf := bytes.Buffer{}

	templateRes.Execute(&buf, &tlsTemplate)

	wantTemplate := `apiVersion: cert-manager.io/v1alpha2
kind: ClusterIssuer
metadata:
  name: letsencrypt-prod
  namespace: openfaas
spec:
  acme:
    email: "sales@openfaas.com"
    server: https://acme-v02.api.letsencrypt.org/directory
    privateKeySecretRef:
      name: letsencrypt-prod
    solvers:
    - dns01:
        digitalocean:
          
            tokenSecretRef:
              name: digitalocean-dns
              key: access-token
          `

	got := string(buf.Bytes())
	if len(got) == 0 {
		t.Errorf("No bytes generated from template")
		t.Fail()
	}

	if debugYAML {
		ioutil.WriteFile("want-"+tlsTemplate.DNSService+".yaml", []byte(wantTemplate), 0700)
		ioutil.WriteFile("got-"+tlsTemplate.DNSService+".yaml", []byte(got), 0700)
	}

	if got != wantTemplate {
		t.Errorf("Want\n`%q`\n, but got\n`%q`", wantTemplate, got)
	}

}

func Test_Route53_Issuer(t *testing.T) {
	tlsTemplate := TLSTemplate{
		Email:       "sales@openfaas.com",
		IssuerType:  "ClusterIssuer",
		DNSService:  "route53",
		Region:      "us-east-1",
		AccessKeyID: "key-id",
	}

	templatePath := "../../templates/k8s/tls/issuer-prod.yml"
	templateData, err := ioutil.ReadFile(templatePath)
	if err != nil {
		t.Error(err)
		return
	}

	templateRes := template.Must(template.New("prod-issuer").Parse(string(templateData)))
	buf := bytes.Buffer{}
	templateRes.Execute(&buf, &tlsTemplate)

	wantTemplate := `apiVersion: cert-manager.io/v1alpha2
kind: ClusterIssuer
metadata:
  name: letsencrypt-prod
  namespace: openfaas
spec:
  acme:
    email: "sales@openfaas.com"
    server: https://acme-v02.api.letsencrypt.org/directory
    privateKeySecretRef:
      name: letsencrypt-prod
    solvers:
    - dns01:
        route53:
          
          region: us-east-1
          # optional if ambient credentials are available; see ambient credentials documentation
          accessKeyID: key-id
          secretAccessKeySecretRef:
            name: "route53-credentials-secret"
            key: secret-access-key
          `

	got := string(buf.Bytes())
	if len(got) == 0 {
		t.Errorf("No bytes generated from template")
		t.Fail()
	}

	if debugYAML {
		ioutil.WriteFile("want-"+tlsTemplate.DNSService+".yaml", []byte(wantTemplate), 0700)
		ioutil.WriteFile("got-"+tlsTemplate.DNSService+".yaml", []byte(got), 0700)
	}

	if got != wantTemplate {
		t.Errorf("Want\n`%q`\n, but got\n`%q`", wantTemplate, got)
	}
}

func Test_GoogleCloudDNS_Issuer(t *testing.T) {
	tlsTemplate := TLSTemplate{
		Email:      "sales@openfaas.com",
		IssuerType: "ClusterIssuer",
		DNSService: "clouddns",
		ProjectID:  "project-1",
	}

	templatePath := "../../templates/k8s/tls/issuer-prod.yml"
	templateData, err := ioutil.ReadFile(templatePath)
	if err != nil {
		t.Error(err)
		return
	}

	templateRes := template.Must(template.New("prod-issuer").Parse(string(templateData)))
	buf := bytes.Buffer{}
	templateRes.Execute(&buf, &tlsTemplate)

	wantTemplate := `apiVersion: cert-manager.io/v1alpha2
kind: ClusterIssuer
metadata:
  name: letsencrypt-prod
  namespace: openfaas
spec:
  acme:
    email: "sales@openfaas.com"
    server: https://acme-v02.api.letsencrypt.org/directory
    privateKeySecretRef:
      name: letsencrypt-prod
    solvers:
    - dns01:
        clouddns:
          
          project: "project-1"
          serviceAccountSecretRef:
            name: "clouddns-service-account"
            key: service-account.json
          `

	got := string(buf.Bytes())
	if len(got) == 0 {
		t.Errorf("No bytes generated from template")
		t.Fail()
	}

	if debugYAML {
		ioutil.WriteFile("want-"+tlsTemplate.DNSService+".yaml", []byte(wantTemplate), 0700)
		ioutil.WriteFile("got-"+tlsTemplate.DNSService+".yaml", []byte(got), 0700)
	}

	if got != wantTemplate {
		t.Errorf("Want\n`%q`\n, but got\n`%q`", wantTemplate, got)
	}
}

var debugYAML bool
