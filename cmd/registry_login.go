package cmd

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var registryLoginCommand = &cobra.Command{
	Use:          "registry-login",
	Short:        "Generate and save the registry authentication file",
	SilenceUsage: true,
	RunE:         generateRegistryAuthFile,
}

func init() {
	rootCommand.AddCommand(registryLoginCommand)

	registryLoginCommand.Flags().String("server", "https://index.docker.io/v1/", "The server URL, it is defaulted to the docker registry")
	registryLoginCommand.Flags().String("username", "", "The Registry Username")
	registryLoginCommand.Flags().String("password", "", "The registry password")
	registryLoginCommand.Flags().BoolP("password-stdin", "s", false, "Reads the gateway password from stdin")

	registryLoginCommand.Flags().Bool("ecr", false, "If we are using ECR we need a different set of flags, so if this is set, we need to set --username and --password")
	registryLoginCommand.Flags().String("account-id", "", "Your AWS Account id")
	registryLoginCommand.Flags().String("region", "", "Your AWS region")
}

func generateRegistryAuthFile(command *cobra.Command, _ []string) error {
	ecrEnabled, _ := command.Flags().GetBool("ecr")
	accountID, _ := command.Flags().GetString("account-id")
	region, _ := command.Flags().GetString("region")
	username, _ := command.Flags().GetString("username")
	password, _ := command.Flags().GetString("password")
	server, _ := command.Flags().GetString("server")
	passStdIn, _ := command.Flags().GetBool("password-stdin")

	if ecrEnabled {
		return generateECRFile(accountID, region)
	} else {
		if passStdIn {
			passwordStdin, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				return err
			}
			return generateFile(username, strings.TrimSpace(string(passwordStdin)), server)
		} else {
			return generateFile(username, password, server)
		}
	}
}

func generateFile(username string, password string, server string) error {

	fileBytes, err := generateRegistryAuth(server, username, password)
	if err != nil {
		return err
	}

	writeErr := writeFileToOFCTmp(fileBytes)

	if writeErr != nil {
		fmt.Print("Registry credentials saved to ./credentials/config.json")
	}

	return writeErr
}

func generateECRFile(accountID string, region string) error {

	fileBytes, err := generateECRRegistryAuth(accountID, region)
	if err != nil {
		return err
	}

	writeErr := writeFileToOFCTmp(fileBytes)

	if writeErr != nil {
		fmt.Print("Registry credentials saved to ./credentials/config.json")
	}
	return writeErr
}

func generateRegistryAuth(server, username, password string) ([]byte, error) {
	if len(username) == 0 || len(password) == 0 || len(server) == 0 {
		return nil, errors.New("both --username and --password must be used, and provided, for us to generate a valid file")
	}

	encodedString := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", username, password)))
	data := RegistryAuth{
		AuthConfigs: map[string]Auth{
			server: {Base64AuthString: encodedString},
		},
	}

	registryBytes, err := json.MarshalIndent(data, "", " ")

	return registryBytes, err
}

func generateECRRegistryAuth(accountID, region string) ([]byte, error) {
	if len(accountID) == 0 || len(region) == 0 {
		return nil, errors.New("you must provide an --account-id and --region when using --ecr")
	}

	data := ECRRegistryAuth{
		CredsStore: "ecr-login",
		CredHelpers: map[string]string{
			fmt.Sprintf("%s.dkr.ecr.%s.amazonaws.com", accountID, region): "ecr-login",
		},
	}

	registryBytes, err := json.MarshalIndent(data, "", " ")

	return registryBytes, err
}

func writeFileToOFCTmp(fileBytes []byte) error {
	path := "./credentials"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.Mkdir(path, 0744)
		if err != nil {
			return err
		}
	}

	writeErr := ioutil.WriteFile(filepath.Join(path, "config.json"), fileBytes, 0744)

	return writeErr

}

type Auth struct {
	Base64AuthString string `json:"auth"`
}

type RegistryAuth struct {
	AuthConfigs map[string]Auth `json:"auths"`
}

type ECRRegistryAuth struct {
	CredsStore  string            `json:"credsStore"`
	CredHelpers map[string]string `json:"credHelpers"`
}
