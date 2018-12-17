package utils

import (
	"fmt"
	"github.com/alexellis/ofc-bootstrap/pkg/types"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"testing"
)

func Test_defaultSecretGenerator_defaultSecretDesc(t *testing.T) {
	defaultSecretDesc := newDefaultSecretDesc()
	pass, err := defaultSecretGenerator(nil)

	if err != nil {
		t.Error(err)
	}

	if len(pass) != defaultSecretDesc.length {
		t.Errorf(
			"Length of password is incorrect. Expected: %d. Received: %d",
			defaultSecretDesc.length,
			len(pass),
		)
	}

	re := regexp.MustCompile("[0-9]")
	digitsCount := len(re.FindAllString(pass, -1))

	if digitsCount != defaultSecretDesc.numDigits {
		t.Errorf(
			"Amount of digits is incorrent. Expected: %d. Received: %d",
			defaultSecretDesc.numDigits,
			digitsCount,
		)
	}

	re = regexp.MustCompile("[^a-zA-Z0-9]")
	symbolsCount := len(re.FindAllString(pass, -1))

	if symbolsCount != defaultSecretDesc.numSymbols {
		t.Errorf(
			"Amount of symbols is incorrent. Expected: %d. Received: %d",
			defaultSecretDesc.numSymbols,
			symbolsCount,
		)
	}

	if !defaultSecretDesc.allowRepeat {
		uniqCharsLen := len(RemoveDuplicates(strings.Split(pass, "")))

		if uniqCharsLen != defaultSecretDesc.length {
			t.Errorf(
				"Expected: %d unique characters. Received: %d",
				defaultSecretDesc.length,
				uniqCharsLen,
			)
		}
	}

	if defaultSecretDesc.noUpper {
		re := regexp.MustCompile("[^A-Z]")
		notUpperCharsLen := len(re.FindAllString(pass, -1))

		if notUpperCharsLen != defaultSecretDesc.length {
			t.Errorf(
				"Expected: %d not upper-case characters. Received: %d",
				defaultSecretDesc.length,
				notUpperCharsLen,
			)
		}
	}
}

func TestCreateDockerSecret_successfulSecretGeneration(t *testing.T) {
	expectedName := "foo"
	expectedSecret := "baz_foo_bar"
	expectedCmd := fmt.Sprintf("echo %s | docker secret create %s", expectedSecret, expectedName)
	kvn := types.KeyValueNamespaceTuple{ Name: expectedName }

	cmd := CreateDockerSecret(kvn, func(sd *secretDesc) (string, error) {
		return expectedSecret, nil
	})

	if cmd != expectedCmd {
		t.Errorf("Expected to receive command: `%s`. Received: `%s`", expectedCmd, cmd)
	}
}

func TestCreateDockerSecret_unsuccessfulSecretGeneration(t *testing.T) {
	kvn := types.KeyValueNamespaceTuple{ Name: "foo" }

	if os.Getenv("BE_CRASHING_SECRET_GENERATION") == "1" {
		CreateDockerSecret(kvn, func(sd *secretDesc) (string, error) {
			return "", fmt.Errorf("some error message")
		})

		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestCreateDockerSecret_unsuccessfulSecretGeneration")
	cmd.Env = append(os.Environ(), "BE_CRASHING_SECRET_GENERATION=1")
	err := cmd.Run()

	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}

	t.Fatalf("Process ran with err %v, want os.Exit(1)", err)
}

func TestCreateK8sSecret_successfulSecretGeneration(t *testing.T) {
	expectedName := "foo"
	expectedNamespace := "baz_namespace"
	expectedSecret := "baz_foo_bar"
	expectedCmd := fmt.Sprintf(
		"kubectl create secret generic -n %s %s --from-literal s3-access-key=\"%s\"",
		expectedNamespace,
		expectedName,
		expectedSecret,
	)
	kvn := types.KeyValueNamespaceTuple{ Name: expectedName, Namespace: expectedNamespace }

	cmd := CreateK8sSecret(kvn, func(sd *secretDesc) (string, error) {
		return expectedSecret, nil
	})

	if cmd != expectedCmd {
		t.Errorf("Expected to receive command: `%s`. Received: `%s`", expectedCmd, cmd)
	}
}

func TestCreateK8sSecret_unsuccessfulSecretGeneration(t *testing.T) {
	kvn := types.KeyValueNamespaceTuple{}

	if os.Getenv("BE_CRASHING_SECRET_GENERATION") == "1" {
		CreateK8sSecret(kvn, func(sd *secretDesc) (string, error) {
			return "", fmt.Errorf("some error message")
		})

		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestCreateK8sSecret_unsuccessfulSecretGeneration")
	cmd.Env = append(os.Environ(), "BE_CRASHING_SECRET_GENERATION=1")
	err := cmd.Run()

	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}

	t.Fatalf("Process ran with err %v, want os.Exit(1)", err)
}
