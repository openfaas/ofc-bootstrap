package validators

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type AuthConfig struct {
	Auth string `json:"auth,omitempty"`
}

type DockerConfigJson struct {
	AuthConfigs map[string]AuthConfig `json:"auths"`
}

func unmarshalRegistryConfig(data []byte) (*DockerConfigJson, error) {
	var registryConfig DockerConfigJson

	err := json.Unmarshal(data, &registryConfig)
	if err != nil {
		return nil, err
	}
	return &registryConfig, nil
}

func ValidateRegistryAuth(registryEndpoint string, configFileBytes []byte) error {

	registryData, unmarshalErr := unmarshalRegistryConfig(configFileBytes)
	if unmarshalErr != nil {
		return unmarshalErr
	}

	noAuthErr := validate(registryData, registryEndpoint)
	if noAuthErr != nil {
		return noAuthErr
	}

	return nil
}

func validate(registryData *DockerConfigJson, endpoint string) error {
	var fileEndpoint string
	if strings.HasPrefix(endpoint, "docker.io") {
		fileEndpoint = "https://index.docker.io/v1/"
	} else {
		fileEndpoint = endpoint
	}
	if endpointConfig, ok := registryData.AuthConfigs[fileEndpoint]; ok {
		if endpointConfig.Auth != "" {
			_, err := base64.StdEncoding.DecodeString(endpointConfig.Auth)
			return err
		} else {
			return errors.New("docker credentials file is not valid (no base64 credentials). Please re-create this file")
		}
	}
	return errors.New(fmt.Sprintf("docker auth file does not contain registry %q that you specified in config. Please use docker login", endpoint))
}
