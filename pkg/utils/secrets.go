package utils

import (
	"fmt"
	"os"

	"github.com/alexellis/ofc-bootstrap/pkg/types"

	"github.com/sethvargo/go-password/password"
)

type secretDesc struct {
	length      int
	numDigits   int
	numSymbols  int
	noUpper     bool
	allowRepeat bool
}

func newDefaultSecretDesc() *secretDesc {
	return &secretDesc{
		length:      16,
		numDigits:   4,
		numSymbols:  4,
		allowRepeat: true,
	}
}

func defaultSecretGenerator(sd *secretDesc) (string, error) {
	secretParams := sd

	if secretParams == nil {
		secretParams = newDefaultSecretDesc()
	}

	return password.Generate(
		secretParams.length,
		secretParams.numDigits,
		secretParams.numSymbols,
		secretParams.noUpper,
		secretParams.allowRepeat,
	)
}

func CreateDockerSecret(kvn types.KeyValueNamespaceTuple, secretGenerator func(*secretDesc) (string, error)) string {
	if secretGenerator == nil {
		secretGenerator = defaultSecretGenerator
	}

	val, err := secretGenerator(nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)

		os.Exit(1)
	}

	return fmt.Sprintf("echo %s | docker secret create %s", val, kvn.Name)
}

func CreateK8sSecret(kvn types.KeyValueNamespaceTuple, secretGenerator func(*secretDesc) (string, error)) string {
	if secretGenerator == nil {
		secretGenerator = defaultSecretGenerator
	}

	val, err := secretGenerator(nil)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return fmt.Sprintf(
		"kubectl create secret generic -n %s %s --from-literal s3-access-key=\"%s\"",
		kvn.Namespace,
		kvn.Name,
		val,
	)
}
