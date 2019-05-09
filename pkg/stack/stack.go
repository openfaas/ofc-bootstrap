package stack

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"os"

	"github.com/openfaas-incubator/ofc-bootstrap/pkg/types"
)

type gatewayConfig struct {
	Registry             string
	RootDomain           string
	CustomersURL         string
	Scheme               string
	S3                   types.S3
	CustomTemplates      string
	EnableDockerfileLang bool
}

type authConfig struct {
	RootDomain           string
	ClientId             string
	CustomersURL         string
	Scheme               string
	OAuthProvider        string
	OAuthProviderBaseURL string
}

// Apply creates `templates/gateway_config.yml` to be referenced by stack.yml
func Apply(plan types.Plan) error {
	scheme := "http"
	if plan.TLS {
		scheme += "s"
	}

	gwConfigErr := generateTemplate("gateway_config", plan, gatewayConfig{
		Registry:             plan.Registry,
		RootDomain:           plan.RootDomain,
		CustomersURL:         plan.CustomersURL,
		Scheme:               scheme,
		S3:                   plan.S3,
		CustomTemplates:      plan.Deployment.FormatCustomTemplates(),
		EnableDockerfileLang: plan.EnableDockerfileLang,
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

	if slackConfigErr := generateTemplate("slack", plan, types.Slack{
		URL: plan.Slack.URL,
	}); slackConfigErr != nil {
		return slackConfigErr
	}

	if plan.SCM == "gitlab" {
		gitlabConfigErr := generateTemplate("gitlab", plan, types.Gitlab{
			GitLabInstance: plan.Gitlab.GitLabInstance,
		})
		if gitlabConfigErr != nil {
			return gitlabConfigErr
		}
	}

	dashboardConfigErr := generateTemplate("dashboard_config", plan, gatewayConfig{
		RootDomain: plan.RootDomain, Scheme: scheme,
	})
	if dashboardConfigErr != nil {
		return dashboardConfigErr
	}

	if plan.EnableOAuth {
		ofAuthDepErr := generateTemplate("edge-auth-dep", plan, authConfig{
			RootDomain:           plan.RootDomain,
			ClientId:             plan.OAuth.ClientId,
			CustomersURL:         plan.CustomersURL,
			Scheme:               scheme,
			OAuthProvider:        plan.SCM,
			OAuthProviderBaseURL: plan.OAuth.OAuthProviderBaseURL,
		})
		if ofAuthDepErr != nil {
			return ofAuthDepErr
		}
	}

	isGitHub := plan.SCM == "github"

	stackErr := generateTemplate("stack", plan, stackConfig{
		GitHub: isGitHub,
	})

	if stackErr != nil {
		return stackErr
	}

	return nil
}

type stackConfig struct {
	GitHub bool
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
