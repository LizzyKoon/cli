package cmd

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/urfave/cli"

	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/scmintegrations"
	"github.com/Scalingo/go-scalingo/v4"
)

var (
	integrationsListCommand = cli.Command{
		Name:     "integrations",
		Category: "Integrations",
		Usage:    "List your integrations",
		Description: `List all the integrations associated with your account:

	$ scalingo integrations

	# See also commands 'integrations-add', 'integrations-delete', 'integrations-import-keys'`,

		Action: func(c *cli.Context) {
			err := scmintegrations.List()
			if err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "integrations")
		},
	}

	integrationsAddCommand = cli.Command{
		Name:     "integrations-add",
		Category: "Integrations",
		Flags: []cli.Flag{appFlag,
			cli.StringFlag{Name: "url", Usage: "URL of the integration"},
			cli.StringFlag{Name: "token", Usage: "Token of the integration"},
		},
		Usage: "Link your Scalingo account with your account on a tool such as github.com",
		Description: `Link your Scalingo account with your account on a tool. After creating the link, you will be able to automatically deploy when pushing to your repository, or create Review Apps for all pull requests created.

	For github.com:
	$ scalingo integrations-add github

	For gitlab.com:
	$ scalingo integrations-add gitlab

	For GitHub Enterprise:
	$ scalingo integrations-add github-enterprise --url https://ghe.example.com --token personal-access-token

	For GitLab Self-hosted:
	$ scalingo integrations-add gitlab-self-hosted --url https://gitlab.example.com --token personal-access-token

	# See also commands 'integrations', 'integrations-delete', 'integrations-import-keys'`,

		Action: func(c *cli.Context) {
			if c.NArg() != 1 {
				_ = cli.ShowCommandHelp(c, "integrations-add")
				return
			}

			integrationURL := c.String("url")
			token := c.String("token")
			scmType := scalingo.SCMType(c.Args()[0])

			switch scmType {
			case scalingo.SCMGithubType, scalingo.SCMGitlabType:
				break
			case scalingo.SCMGithubEnterpriseType, scalingo.SCMGitlabSelfHostedType:
				if integrationURL == "" || token == "" {
					errorQuit(errors.New("both --url and --token must be set"))
				}

				u, err := url.Parse(integrationURL)
				if err != nil || u.Scheme == "" || u.Host == "" {
					errorQuit(fmt.Errorf("'%s' is not a valid URL", integrationURL))
				}
			default:
				errorQuit(errors.New(
					"unknown integration. Available integrations: github, github-enterprise, gitlab, gitlab-self-hosted",
				))
			}

			args := scmintegrations.CreateArgs{
				SCMType: scmType,
				URL:     integrationURL,
				Token:   token,
			}

			err := scmintegrations.Create(args)
			if err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "integrations-add")
		},
	}

	integrationsDeleteCommand = cli.Command{
		Name:     "integrations-delete",
		Category: "Integrations",
		Usage:    "Unlink your Scalingo account and your account on a tool such as github.com",
		Description: `Unlink your Scalingo account and your account on a tool:

	$ scalingo integrations-delete integration-type
	OR
	$ scalingo integrations-delete integration-uuid

	Examples:
	$ scalingo integrations-delete github-enterprise
	$ scalingo integrations-delete gitlab

	# See also commands 'integrations', 'integrations-add', 'integrations-import-keys'`,

		Action: func(c *cli.Context) {
			if c.NArg() != 1 {
				_ = cli.ShowCommandHelp(c, "integrations-delete")
				return
			}

			err := scmintegrations.Delete(c.Args()[0])
			if err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "integrations-delete")
		},
	}

	integrationsImportKeysCommand = cli.Command{
		Name:     "integrations-import-keys",
		Category: "Integrations",
		Usage:    "Import public SSH keys from integration account",
		Description: `Import public SSH keys from integration account:

	$ scalingo integrations-import-keys integration-type
	OR
	$ scalingo integrations-import-keys integration-uuid

	Examples:
	$ scalingo integrations-import-keys github
	$ scalingo integrations-import-keys gitlab-self-hosted

	# See also commands 'integrations', 'integrations-add', 'integrations-delete'`,

		Action: func(c *cli.Context) {
			if c.NArg() != 1 {
				_ = cli.ShowCommandHelp(c, "integrations-import-keys")
				return
			}

			err := scmintegrations.ImportKeys(c.Args()[0])
			if err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "integrations-import-keys")
		},
	}
)
