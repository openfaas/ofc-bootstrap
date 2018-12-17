package types

type Plan struct {
	Orchestration string                   `yaml:"orchestration"`
	Secrets       []KeyValueNamespaceTuple `yaml:"secrets"`
	RootDomain    string                   `yaml:"root_domain"`
	Registry      string                   `yaml:"registry"`
	CustomersURL  string                   `yaml:"customers_url"`
	FunctionStack string                   `yaml:"function_stack"`
}

type KeyValueNamespaceTuple struct {
	Name      string `yaml:"name"`
	Value     string `yaml:"value"`
	Namespace string `yaml:"namespace"`
}
