package main

import (
	"context"
	"fmt"
	"github.com/google/go-github/github"
	"github.com/shmileee/github-to-terraform/pkg/repository"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"text/template"
)

// version is the published version of the utility
var version string

const (
	CollaboratorType string = "collaborator-type"
	Organization     string = "org"
	RepositoryName   string = "repo-name"
	RepositoryType   string = "repo-type"
	TemplatePath     string = "template"
	VerboseFlag      string = "verbose"
)

// initFlags is where command line flags are instantiated
func initFlags(flag *pflag.FlagSet) {

	flag.BoolP(VerboseFlag, "v", false, "log messages at the debug level.")
	flag.String(TemplatePath, "templates/collaborators.tpl", "Template to render collaborators from")
	flag.String(Organization, "Appsilon", "GitHub Organisation")
	flag.String(CollaboratorType, "outside", "Collaborator time: outside, direct, all")
	flag.String(RepositoryType, "private", "Limit by repo type (public, private) if one repo not specified")
	flag.String(RepositoryName, "", "Collaborators for specific repo")

	flag.SortFlags = false
}

// checkConfig is how the input to command line flags are checked
func checkConfig(v *viper.Viper) error {

	return nil
}

func main() {
	root := cobra.Command{
		Use:   "github-to-terraform [flags]",
		Short: "Get current GitHub configuration and create Terraform code for it",
		Long:  "Get current GitHub configuration and create Terraform code for it",
	}

	completionCommand := &cobra.Command{
		Use:   "completion",
		Short: "Generates bash completion scripts",
		Long:  "To install completion scripts run:\ngithub-to-terraform completion > /usr/local/etc/bash_completion.d/github-to-terraform",
		RunE: func(cmd *cobra.Command, args []string) error {
			return root.GenBashCompletion(os.Stdout)
		},
	}
	root.AddCommand(completionCommand)

	repositoryCollaboratorCommand := &cobra.Command{
		Use:                   "collaborators [flags]",
		DisableFlagsInUseLine: true,
		Short:                 "Retrieve repository collaborators from GitHub and save them as Terraform resources",
		Long:                  "Retrieve repository collaborators from GitHub and save them as Terraform resources",
		RunE:                  repositoryCollaboratorFunction,
	}
	initFlags(repositoryCollaboratorCommand.Flags())
	root.AddCommand(repositoryCollaboratorCommand)

	versionCommand := &cobra.Command{
		Use:                   "version",
		DisableFlagsInUseLine: true,
		Short:                 "Print the version",
		Long:                  "Print the version",
		RunE:                  versionFunction,
	}
	root.AddCommand(versionCommand)

	if err := root.Execute(); err != nil {
		panic(err)
	}
}

func versionFunction(cmd *cobra.Command, args []string) error {
	if len(version) == 0 {
		fmt.Println("development")
		return nil
	}
	fmt.Println(version)
	return nil
}

type Output struct {
	Repository   *github.Repository
	Collaborator *github.User
}

const (
	perPage = 50
)

func repositoryCollaboratorFunction(cmd *cobra.Command, args []string) error {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()

	err := cmd.ParseFlags(args)
	if err != nil {
		return err
	}

	flag := cmd.Flags()

	v := viper.New()
	bindErr := v.BindPFlags(flag)
	if bindErr != nil {
		return bindErr
	}
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	// Create the logger
	// Remove the prefix and any datetime data
	logger := log.New(os.Stdout, "", log.LstdFlags)

	verbose := v.GetBool(VerboseFlag)
	if !verbose {
		// Disable any logging that isn't attached to the logger unless using the verbose flag
		log.SetOutput(ioutil.Discard)
		log.SetFlags(0)

		// Remove the flags for the logger
		logger.SetFlags(0)
	}

	// Check the config and exit with usage details if there is a problem
	checkConfigErr := checkConfig(v)
	if checkConfigErr != nil {
		return checkConfigErr
	}

	tPath := v.GetString(TemplatePath)
	org := v.GetString(Organization)
	collaboratorType := v.GetString(CollaboratorType)
	repoType := v.GetString(RepositoryType)
	repoName := v.GetString(RepositoryName)

	out := new(Output)

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	repos := make(chan *github.Repository, perPage)

	repoConfig := repository.ListConfig{
		ListOptions:  github.ListOptions{PerPage: perPage},
		Type:         repoType,
		Organization: org,
		Fast:         true,
		RepoName:     repoName,
	}

	go func() {
		err := repository.List(client, repos, repoConfig)
		if err != nil {
			panic(err)
		}
		close(repos)
	}()

	// list collaborators by type
	opt := &github.ListCollaboratorsOptions{
		Affiliation: collaboratorType,
		ListOptions: github.ListOptions{PerPage: 50},
	}
	t := template.Must(template.ParseFiles(tPath))
	for repo := range repos {
		out.Repository = repo
	Collaborators:
		for {
			collabs, resp, err := client.Repositories.ListCollaborators(context.Background(), org, *out.Repository.Name, opt)
			if err != nil {
				fmt.Println(err)
			}
			for i := range collabs {
				out.Collaborator = collabs[i]
				err = t.Execute(os.Stdout, out)
				if err != nil {
					panic(err)
				}
			}

			if resp.NextPage == 0 {
				break Collaborators
			}
			opt.Page = resp.NextPage
		}
	}

	return nil
}
