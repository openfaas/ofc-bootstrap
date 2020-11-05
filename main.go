package main

import (
	"os"

	"github.com/openfaas/ofc-bootstrap/cmd"

	"github.com/openfaas/ofc-bootstrap/version"
)

func main() {

	if err := cmd.Execute(version.Version, version.GitCommit); err != nil {
		// fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}
	return
}
