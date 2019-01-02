package types

import (
	"os"
	"strings"
)

type Plan struct {
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
}

type Github struct {
	AppID          string `yaml:"app_id"`
	PrivateKeyFile string `yaml:"private_key_filename"`
}

type OAuth struct {
	ClientId     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
}

type S3 struct {
	Url    string `yaml:"s3_url"`
	Region string `yaml:"s3_region"`
	TLS    bool   `yaml:"s3_tls"`
	Bucket string `yaml:"s3_bucket"`
}

type TLSConfig struct {
	Email             string `yaml:"email"`
	DNSService        string `yaml:"dns_service"`
	ProjectID         string `yaml:"project_id"`
	IssuerType        string `yaml:"issuer_type"`
	LetsencryptServer string `yaml:"letsencrypt_server"`
	Region            string `yaml:"region"`
	AccessKeyID       string `yaml:"access_key_id"`
}
