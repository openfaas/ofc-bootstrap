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
	// CloudflareDNS filter enables the creation of secrets for Cloudflare DNS when TLS is enabled
	CloudflareDNS = "cloudflare_dns01"

	// CloudDNS is the dns_service field in yaml file for Google Cloud Platform
	CloudDNS = "clouddns"
	// DigitalOcean is the dns_service field in yaml file for Digital Ocean
	DigitalOcean = "digitalocean"
	// Route53 is the dns_service field in yaml file for Amazon
	Route53 = "route53"
	// Cloudflare for dns_service
	Cloudflare = "cloudflare"

	// GitLabSCM repository manager name as displayed in the init.yaml file
	GitLabSCM = "gitlab"
	// GitHubSCM repository manager name as displayed in the init.yaml file
	GitHubSCM = "github"

	// ECRFeature enable ECR
	ECRFeature = "ecr"
)

type Plan struct {
	Features             []string                 `yaml:"features,omitempty"`
	Orchestration        string                   `yaml:"orchestration,omitempty"`
	Secrets              []KeyValueNamespaceTuple `yaml:"secrets,omitempty"`
	RootDomain           string                   `yaml:"root_domain,omitempty"`
	Registry             string                   `yaml:"registry,omitempty"`
	CustomersURL         string                   `yaml:"customers_url,omitempty"`
	SCM                  string                   `yaml:"scm,omitempty"`
	Github               Github                   `yaml:"github,omitempty"`
	Gitlab               Gitlab                   `yaml:"gitlab,omitempty"`
	TLS                  bool                     `yaml:"tls,omitempty"`
	OAuth                OAuth                    `yaml:"oauth,omitempty"`
	S3                   S3                       `yaml:"s3,omitempty"`
	EnableOAuth          bool                     `yaml:"enable_oauth,omitempty"`
	TLSConfig            TLSConfig                `yaml:"tls_config,omitempty"`
	Slack                Slack                    `yaml:"slack,omitempty"`
	Ingress              string                   `yaml:"ingress,omitempty"`
	Deployment           Deployment               `yaml:"deployment,omitempty"`
	EnableDockerfileLang bool                     `yaml:"enable_dockerfile_lang,omitempty"`
	ScaleToZero          bool                     `yaml:"scale_to_zero,omitempty"`
	OpenFaaSCloudVersion string                   `yaml:"openfaas_cloud_version,omitempty"`
	NetworkPolicies      bool                     `yaml:"network_policies,omitempty"`
	BuildBranch          string                   `yaml:"build_branch,omitempty"`
	EnableECR            bool                     `yaml:"enable_ecr,omitempty"`
	ECRConfig            ECRConfig                `yaml:"ecr_config,omitempty"`
	CustomersSecret      bool                     `yaml:"customers_secret,omitempty"`
	IngressOperator      bool                     `yaml:"enable_ingress_operator,omitempty"`
	IngressConfig        IngressConfig            `yaml:"ingress_config,omitempty"`
}

// IngressConfig is for configuring Ingress type settings, like max connections
type IngressConfig struct {
	MaxConnections    string `yaml:"max_connections,omitempty"`
	RequestsPerMinute string `yaml:"requests_per_minute,omitempty"`
}

// Deployment is the deployment section of YAML concerning
// functions as deployed
type Deployment struct {
	CustomTemplate []string `yaml:"custom_templates,omitempty"`
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
	Name  string `yaml:"name,omitempty"`
	Value string `yaml:"value,omitempty"`
}

type FileSecret struct {
	Name      string `yaml:"name,omitempty"`
	ValueFrom string `yaml:"value_from,omitempty"`

	// ValueCommand is a command to execute to generate
	// a secret file specified in ValueFrom
	ValueCommand string `yaml:"value_command,omitempty"`
}

// ExpandValueFrom expands ~ to the home directory of the current user
// kept in the HOME env-var.
func (fs FileSecret) ExpandValueFrom() string {
	value := fs.ValueFrom
	value = strings.Replace(value, "~", os.Getenv("HOME"), -1)
	return value
}

type KeyValueNamespaceTuple struct {
	Name      string          `yaml:"name,omitempty"`
	Literals  []KeyValueTuple `yaml:"literals,omitempty"`
	Namespace string          `yaml:"namespace,omitempty"`
	Files     []FileSecret    `yaml:"files,omitempty"`
	Type      string          `yaml:"type,omitempty"`
	Filters   []string        `yaml:"filters,omitempty"`
}

type Github struct {
	AppID          string `yaml:"app_id,omitempty"`
	PrivateKeyFile string `yaml:"private_key_filename,omitempty"`
}

type Gitlab struct {
	GitLabInstance string `yaml:"gitlab_instance,omitempty"`
}

type Slack struct {
	URL string `yaml:"url,omitempty"`
}

type OAuth struct {
	ClientId             string `yaml:"client_id,omitempty"`
	OAuthProviderBaseURL string `yaml:"oauth_provider_base_url,omitempty"`
}

type S3 struct {
	Url    string `yaml:"s3_url,omitempty"`
	Region string `yaml:"s3_region,omitempty"`
	TLS    bool   `yaml:"s3_tls,omitempty"`
	Bucket string `yaml:"s3_bucket,omitempty"`
}

type TLSConfig struct {
	Email       string `yaml:"email,omitempty"`
	DNSService  string `yaml:"dns_service,omitempty"`
	ProjectID   string `yaml:"project_id,omitempty"`
	IssuerType  string `yaml:"issuer_type,omitempty"`
	Region      string `yaml:"region,omitempty"`
	AccessKeyID string `yaml:"access_key_id,omitempty"`
}

type ECRConfig struct {
	ECRRegion string `yaml:"ecr_region,omitempty"`
}

type OFCValues struct {
	EdgeAuth        EdgeAuth        `yaml:"edgeAuth,omitempty"`
	TLS             OFCTLS          `yaml:"tls,omitempty"`
	Ingress         Ingress         `yaml:"ingress,omitempty"`
	Global          Global          `yaml:"global,omitempty"`
	NetworkPolicies NetworkPolicies `yaml:"networkPolicies,omitempty"`
	Customers       Customers       `yaml:"customers,omitempty"`
}

type NetworkPolicies struct {
	Enabled bool `yaml:"enabled,omitempty"`
}

type EdgeAuth struct {
	EnableOauth2         bool   `yaml:"enableOauth2,omitempty"`
	OauthProvider        string `yaml:"oauthProvider,omitempty"`
	OauthProviderBaseURL string `yaml:"oauthProviderBaseURL,omitempty"`
	ClientID             string `yaml:"clientId,omitempty"`
}

type OFCTLS struct {
	Enabled    bool           `yaml:"enabled"`
	Email      string         `yaml:"email,omitempty"`
	IssuerType string         `yaml:"issuerType,omitempty"`
	DNSService string         `yaml:"dnsService,omitempty"`
	Route53    OFCRoute53     `yaml:"route53,omitempty"`
	Cloudflare OFCCloudflare  `yaml:"cloudflare,omitempty"`
	CloudDNS   CloudDNSConfig `yaml:"clouddns,omitempty"`
}

type CloudDNSConfig struct {
	ProjectID string `yaml:"projectID,omitempty"`
}

type OFCRoute53 struct {
	Region      string `yaml:"region,omitempty"`
	AccessKeyID string `yaml:"accessKeyID,omitempty"`
}

type OFCCloudflare struct {
	Email     string `yaml:"email,omitempty"`
	ProjectID string `yaml:"projectID,omitempty"`
}

type Ingress struct {
	MaxConnections    string `yaml:"maxConnections,omitempty"`
	RequestsPerMinute string `yaml:"requestsPerMinute,omitempty"`
}

type Customers struct {
	URL             string `yaml:"url,omitempty"`
	CustomersSecret bool   `yaml:"customersSecret,omitempty"`
}

type Global struct {
	RootDomain string `yaml:"rootDomain,omitempty"`
	EnableECR  bool   `yaml:"enableECR,omitempty"`
}
