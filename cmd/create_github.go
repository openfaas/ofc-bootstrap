package cmd

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/inlets/inletsctl/pkg/names"
	"github.com/openfaas-incubator/ofc-bootstrap/pkg/github"
	"github.com/openfaas-incubator/ofc-bootstrap/pkg/types"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
)

var createGitHubAppCommand = &cobra.Command{
	Use:          "create-github-app",
	Short:        "Create a GitHub App",
	SilenceUsage: true,
	RunE:         createGitHubAppE,
}

func init() {
	rootCommand.AddCommand(createGitHubAppCommand)

	createGitHubAppCommand.Flags().String("name", "", "The name of your GitHub App for OpenFaaS Cloud, leave blank to generate")
	createGitHubAppCommand.Flags().String("root-domain", "", "The root domain for your app i.e. ofc.example.com")
	createGitHubAppCommand.Flags().Bool("insecure", false, "Use HTTP instead of HTTPS for webhooks")
	createGitHubAppCommand.Flags().Int("port", 30010, "HTTP port to listen to for GitHub App automation")
}

func createGitHubAppE(command *cobra.Command, _ []string) error {

	name := ""
	if nameFlagVal, _ := command.Flags().GetString("name"); len(nameFlagVal) > 0 {
		name = nameFlagVal

	} else {
		name = "OFC " + strings.Replace(names.GetRandomName(10), "_", " ", -1)
	}
	var rootDomain string
	if rootDomain, _ = command.Flags().GetString("root-domain"); len(rootDomain) == 0 {
		return fmt.Errorf("give a value for --root-domain")
	}

	port, portErr := command.Flags().GetInt("port")
	if portErr != nil {
		return fmt.Errorf("--port error: %s", portErr.Error())
	}

	scheme := "https"
	if insecure, _ := command.Flags().GetBool("insecure"); insecure {
		scheme = "http"
	}

	inputMap := map[string]string{
		"AppName":     name,
		"GitHubEvent": fmt.Sprintf("%s://system.%s/github-event", scheme, rootDomain),
	}

	if _, err := os.Stat(path.Join("./pkg/github", "index.html")); err != nil {
		return fmt.Errorf(`cannot find template "index.html", run this command from the ofc-bootstrap repository`)
	}

	fmt.Printf("Name: %s\tRoot domain: %s\tScheme: %v\n", name, rootDomain, scheme)

	launchBrowser := true

	context, cancel := context.WithCancel(context.TODO())
	defer cancel()
	resCh := make(chan github.AppResult)
	go func() {
		appRes := <-resCh
		fmt.Printf("GitHub App result received.\n")
		printResult(rootDomain, appRes)

		cancel()
	}()

	err := receiveGitHubApp(context, inputMap, resCh, launchBrowser, port)

	if err != nil {
		return err
	}

	return nil
}

func receiveGitHubApp(ctx context.Context, inputMap map[string]string, resCh chan github.AppResult, launchBrowser bool, listenPort int) error {

	server := &http.Server{
		Addr:           fmt.Sprintf(":%d", listenPort),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20, // Max header of 1MB
		Handler:        http.HandlerFunc(github.MakeHandler(inputMap, resCh)),
	}

	go func() {
		fmt.Printf("Starting local HTTP server on port %d\n", listenPort)
		if err := server.ListenAndServe(); err != nil {
			panic(err)
		}
	}()

	defer server.Shutdown(ctx)

	localURL, err := url.Parse("http://" + "127.0.0.1" + server.Addr)

	if err != nil {
		return err
	}

	fmt.Printf("Launching browser: %s\n", localURL.String())
	if launchBrowser {
		err := launchURL(localURL.String())
		if err != nil {
			return errors.Wrap(err, "unable to launch browser")
		}
	}

	fmt.Printf("Please complete the workflow in the browser to create your GitHub App.\n")

	<-ctx.Done()
	return nil
}

func printResult(rootDomain string, appRes github.AppResult) {
	p := types.Plan{
		RootDomain: rootDomain,
		Github: types.Github{
			AppID: fmt.Sprintf("%d", appRes.ID),
		},
		Secrets: []types.KeyValueNamespaceTuple{
			types.KeyValueNamespaceTuple{
				Name: "github-webhook-secret",
				Literals: []types.KeyValueTuple{
					types.KeyValueTuple{
						Name:  "github-webhook-secret",
						Value: appRes.WebhookSecret,
					},
				},
				Filters:   []string{"scm_github"},
				Namespace: "openfaas-fn",
			},
			types.KeyValueNamespaceTuple{
				Name: "private-key",
				Literals: []types.KeyValueTuple{
					types.KeyValueTuple{
						Name:  "private-key",
						Value: appRes.PEM,
					},
				},
				Filters:   []string{"scm_github"},
				Namespace: "openfaas-fn",
			},
		},
	}
	res, _ := yaml.Marshal(p)

	fmt.Printf("App: %s\tURL: %s\nYAML file\n\n", appRes.Name, appRes.URL)

	fmt.Printf("%s\n", string(res))

}

// launchURL opens a URL with the default browser for Linux, MacOS or Windows.
func launchURL(serverURL string) error {
	ctx := context.Background()
	var command *exec.Cmd
	switch runtime.GOOS {
	case "linux":
		command = exec.CommandContext(ctx, "sh", "-c", fmt.Sprintf(`xdg-open "%s"`, serverURL))
	case "darwin":
		command = exec.CommandContext(ctx, "sh", "-c", fmt.Sprintf(`open "%s"`, serverURL))
	case "windows":
		escaped := strings.Replace(serverURL, "&", "^&", -1)
		command = exec.CommandContext(ctx, "cmd", "/c", fmt.Sprintf(`start %s`, escaped))
	}
	command.Stdout = os.Stdout
	command.Stdin = os.Stdin
	command.Stderr = os.Stderr
	return command.Run()
}
