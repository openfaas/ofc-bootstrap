package types

import (
	"os"
	"strings"
)

const (
	// DefaultFeature filter is for the features which are mandatory
	DefaultFeature = "default"
	// GitHubFeature filter enables secrets created with the scm_github filter
	GitHubFeature = "scm_github"
	// GitLabFeature filter is the feature which enables secret creation for GitLab
	GitLabFeature = "scm_gitlab"
	// Auth filter enables OAuth secret creation
	Auth = "auth"
	// GCPDNS filter enables the creation of secrets for Google Cloud Platform DNS when TLS is enabled
	GCPDNS = "gcp_dns01"
	// DODNS filter enables the creation of secrets for Digital Ocean DNS when TLS is enabled
	DODNS = "do_dns01"
	// Route53DNS filter enables the creation of secrets for Amazon Route53 DNS when TLS is enabled
	Route53DNS = "route53_dns01"

	// CloudDNS is the dns_service field in yaml file for Google Cloud Platform
	CloudDNS = "clouddns"
	// DigitalOcean is the dns_service field in yaml file for Digital Ocean
	DigitalOcean = "digitalocean"
	// Route53 is the dns_service field in yaml file for Amazon
	Route53 = "route53"
	// GitLabManager repository manager name as displayed in the init.yaml file
	GitLabSCM = "gitlab"
	// GitHubManager repository manager name as displayed in the init.yaml file
	GitHubSCM = "github"
)

type Plan struct {
	Features             []string                 `yaml:"features"`
	Orchestration        string                   `yaml:"orchestration"`
	Secrets              []KeyValueNamespaceTuple `yaml:"secrets"`
	RootDomain           string                   `yaml:"root_domain"`
	Registry             string                   `yaml:"registry"`
	CustomersURL         string                   `yaml:"customers_url"`
	SCM                  string                   `yaml:"scm"`
	Github               Github                   `yaml:"github"`
	Gitlab               Gitlab                   `yaml:"gitlab"`
	TLS                  bool                     `yaml:"tls"`
	OAuth                OAuth                    `yaml:"oauth"`
	S3                   S3                       `yaml:"s3"`
	EnableOAuth          bool                     `yaml:"enable_oauth"`
	TLSConfig            TLSConfig                `yaml:"tls_config"`
	Slack                Slack                    `yaml:"slack"`
	Ingress              string                   `yaml:"ingress"`
	Deployment           Deployment               `yaml:"deployment"`
	EnableDockerfileLang bool                     `yaml:"enable_dockerfile_lang"`
	ScaleToZero          bool                     `yaml:"scale_to_zero"`
	OpenFaaSCloudVersion string                   `yaml:"openfaas_cloud_version"`
	NetworkPolicies      bool                     `yaml:"network_policies"`
}

// Deployment is the deployment section of YAML concerning
// functions as deployed
type Deployment struct {
	CustomTemplate []string `yaml:"custom_templates"`
}

// FormatCustomTemplates are formatted in a CSV format with a space after each comma
func (d Deployment) FormatCustomTemplates() string {
	val := ""
	for _, templateURL := range d.CustomTemplate {
		val = val + templateURL + ", "
	}

	return strings.TrimRight(val, " ,")
}

type KeyValueTuple struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

type FileSecret struct {
	Name      string `yaml:"name"`
	ValueFrom string `yaml:"value_from"`

	// ValueCommand is a command to execute to generate
	// a secret file specified in ValueFrom
	ValueCommand string `yaml:"value_command"`
}

// ExpandValueFrom expands ~ to the home directory of the current user
// kept in the HOME env-var.
func (fs FileSecret) ExpandValueFrom() string {
	value := fs.ValueFrom
	value = strings.Replace(value, "~", os.Getenv("HOME"), -1)
	return value
}

type KeyValueNamespaceTuple struct {
	Name      string          `yaml:"name"`
	Literals  []KeyValueTuple `yaml:"literals"`
	Namespace string          `yaml:"namespace"`
	Files     []FileSecret    `yaml:"files"`
	Type      string          `yaml:"type"`
	Filters   []string        `yaml:"filters"`
}

type Github struct {
	AppID          string `yaml:"app_id"`
	PrivateKeyFile string `yaml:"private_key_filename"`
}

type Gitlab struct {
	GitLabInstance string `yaml:"gitlab_instance"`
}

type Slack struct {
	URL string `yaml:"url"`
}

type OAuth struct {
	ClientId             string `yaml:"client_id"`
	OAuthProviderBaseURL string `yaml:"oauth_provider_base_url"`
}

type S3 struct {
	Url    string `yaml:"s3_url"`
	Region string `yaml:"s3_region"`
	TLS    bool   `yaml:"s3_tls"`
	Bucket string `yaml:"s3_bucket"`
}

type TLSConfig struct {
	Email       string `yaml:"email"`
	DNSService  string `yaml:"dns_service"`
	ProjectID   string `yaml:"project_id"`
	IssuerType  string `yaml:"issuer_type"`
	Region      string `yaml:"region"`
	AccessKeyID string `yaml:"access_key_id"`
}
