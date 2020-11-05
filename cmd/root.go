// Copyright (c) OpenFaaS Author(s) 2019. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package cmd

import (
	"fmt"

	"github.com/morikuni/aec"
	"github.com/openfaas/ofc-bootstrap/version"
	"github.com/spf13/cobra"
)

var (
	// Version as per git repo
	Version string

	// GitCommit as per git repo
	GitCommit string
)

// WelcomeMessage to introduce ofc-bootstrap
const WelcomeMessage = "Welcome to ofc-bootstrap! Find out more at https://github.com/openfaas/ofc-bootstrap"

func init() {
	rootCommand.AddCommand(versionCmd)
	rootCommand.Flags().StringArrayP("yaml", "f", []string{""}, "The init.yaml plan file")
}

var rootCommand = &cobra.Command{
	Use:   "ofc-bootstrap",
	Short: "Bootstrap OpenFaaS Cloud.",
	Long: `
Bootstrap OpenFaaS Cloud
`,
	RunE:         runRootCommand,
	SilenceUsage: true,
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display version information.",
	Run:   parseBaseCommand,
}

func getVersion() string {
	if len(Version) != 0 {
		return Version
	}
	return "dev"
}

func parseBaseCommand(_ *cobra.Command, _ []string) {
	printLogo()

	fmt.Printf(
		`ofc-bootstrap
Bootstrap your own OpenFaaS Cloud within 100 seconds

Commit: %s
Version: %s

`, version.GitCommit, version.GetVersion())
}

func Execute(version, gitCommit string) error {

	// Get Version and GitCommit values from main.go.
	Version = version
	GitCommit = gitCommit

	if err := rootCommand.Execute(); err != nil {
		return err
	}
	return nil
}

func runRootCommand(cmd *cobra.Command, args []string) error {

	if cmd.Flags().Changed("yaml") {
		return fmt.Errorf("a breaking change was introduced, you now need to use ofc-bootstrap apply --file init.yaml")
	}

	printLogo()
	cmd.Help()

	return nil
}

func printLogo() {
	logoText := aec.WhiteF.Apply(version.Logo)
	fmt.Println(logoText)
}
