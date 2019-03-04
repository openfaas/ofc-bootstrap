package types

import (
	"os"
	"strings"
)

const (
	// InternalTrust filter enables creation of payload secret
	InternalTrust = "internal_trust"
	// BasicAuth enables creation of basic-auth secret for OF gateway
	BasicAuth = "basic_auth"
	// GitHub filter enables secrets created with the scm_github filter
	GitHub = "scm_github"
	// GitLab filter is the feature
	GitLab = "scm_gitlab"
	// Auth filter enables OAuth secret creation
	Auth = "auth"
	// GCPDNS filter enables the creation of secrets for Google Cloud Platform DNS when TLS is enabled
	GCPDNS = "gcp_dns01"
	// DODNS filter enables the creation of secrets for Digital Ocean DNS when TLS is enabled
	DODNS = "do_dns01"
	// DODNS filter enables the creation of secrets for Amazon Route53 DNS when TLS is enabled
	Route53DNS = "route53_dns01"
	// S3Bucket enables creation of secrets for S3 buckets or minio case you would like to see logs
	S3Bucket = "s3"
	// Registry filter enables creation of registry secret
	Registry = "registry"

	// CloudDns is the dns_service field in yaml file for Google Cloud Platform
	CloudDns = "clouddns"
	// CloudDns is the dns_service field in yaml file for Digital Ocean
	DigitalOcean = "digitalocean"
	// Route53 is the dns_service field in yaml file for Amazon
	Route53 = "route53"
)

type Plan struct {
	Features      []string                 `yaml:"features"`
	Orchestration string                   `yaml:"orchestration"`
	Secrets       []KeyValueNamespaceTuple `yaml:"secrets"`
	RootDomain    string                   `yaml:"root_domain"`
	Registry      string                   `yaml:"registry"`
	CustomersURL  string                   `yaml:"customers_url"`
	Github        Github                   `yaml:"github"`
	TLS           bool                     `yaml:"tls"`
	OAuth         OAuth                    `yaml:"oauth"`
	S3            S3                       `yaml:"s3"`
	EnableOAuth   bool                     `yaml:"enable_oauth"`
	TLSConfig     TLSConfig                `yaml:"tls_config"`
	Slack         Slack                    `yaml:"slack"`
	Ingress       string                   `yaml:"ingress"`
	BasicAuth     bool                     `yaml:"basic_auth"`
	InternalTrust bool                     `yaml:"internal_trust"`
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

type Slack struct {
	URL string `yaml:"url"`
}

type OAuth struct {
	ClientId string `yaml:"client_id"`
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
