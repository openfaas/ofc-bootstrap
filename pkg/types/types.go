package types

type Plan struct {
	Orchestration string                   `yaml:"orchestration"`
	Secrets       []KeyValueNamespaceTuple `yaml:"secrets"`
	RootDomain    string                   `yaml:"root_domain"`
	Registry      string                   `yaml:"registry"`
	CustomersURL  string                   `yaml:"customers_url"`
}

type KeyValueTuple struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

type KeyValueNamespaceTuple struct {
	Name      string          `yaml:"name"`
	Literals  []KeyValueTuple `yaml:"literals"`
	Namespace string          `yaml:"namespace"`
}
