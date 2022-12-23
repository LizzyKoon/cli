// Package cmd gathers the configuration of all commands (names, flags, default
// values etc.) of the Scalingo CLI
package cmd

import (
	"os"

	"github.com/urfave/cli/v2"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/detect"
	"github.com/Scalingo/cli/session"
	"github.com/Scalingo/go-scalingo/v6"
	"github.com/Scalingo/go-scalingo/v6/debug"
)

type AppCommands struct {
	commands []*cli.Command
}

type Command struct {
	*cli.Command
	// Regional flag not available if Global is true
	Global bool
}

func (cmds *AppCommands) addCommand(cmd Command) {
	if !cmd.Global {
		regionFlag := &cli.StringFlag{Name: "region", Value: "", Usage: "Name of the region to use"}
		cmd.Command.Flags = append(cmd.Command.Flags, regionFlag)
		action := cmd.Command.Action
		cmd.Command.Action = func(c *cli.Context) error {
			token := os.Getenv("SCALINGO_API_TOKEN")

			currentUser, err := config.C.CurrentUser()
			if err != nil || currentUser == nil {
				err := session.Login(c.Context, session.LoginOpts{APIToken: token})
				if err != nil {
					errorQuit(err)
				}
			}

			regions, err := config.EnsureRegionsCache(c.Context, config.C, config.GetRegionOpts{
				Token: token,
			})
			if err != nil {
				errorQuit(err)
			}
			currentRegion := regionNameFromFlags(c)

			// Detecting Region from git remote
			if currentRegion == "" {
				currentRegion = detect.GetRegionFromGitRemote(c, &regions)
			}

			if config.C.ScalingoRegion == "" && currentRegion == "" {
				region := getDefaultRegion(regions)
				debug.Printf("[Regions] Use the default region '%s'\n", region.Name)
				currentRegion = region.Name
			}

			if currentRegion != "" {
				config.C.ScalingoRegion = currentRegion
			}

			return action(c)
		}
	}
	cmds.commands = append(cmds.commands, cmd.Command)
}

func getDefaultRegion(regionsCache config.RegionsCache) scalingo.Region {
	defaultRegion := regionsCache.Regions[0]
	for _, region := range regionsCache.Regions {
		if region.Default {
			defaultRegion = region
			break
		}
	}
	return defaultRegion
}

func (cmds *AppCommands) Commands() []*cli.Command {
	return cmds.commands
}

func NewAppCommands() *AppCommands {
	cmds := AppCommands{}
	for _, cmd := range regionalCommands {
		cmds.addCommand(Command{Command: cmd})
	}
	for _, cmd := range globalCommands {
		cmds.addCommand(Command{Global: true, Command: cmd})
	}
	return &cmds
}

var (
	regionalCommands = []*cli.Command{
		// Apps
		&appsCommand,
		&CreateCommand,
		&DestroyCommand,
		&renameCommand,
		&appsInfoCommand,
		&openCommand,
		&dashboardCommand,

		// Apps Actions
		&logsCommand,
		&logsArchivesCommand,
		&runCommand,
		&oneOffStopCommand,

		// Apps Process Actions
		&psCommand,
		&scaleCommand,
		&RestartCommand,

		// Routing Settings
		&forceHTTPSCommand,
		&stickySessionCommand,
		&routerLogsCommand,
		&setCanonicalDomainCommand,
		&unsetCanonicalDomainCommand,

		// Events
		&UserTimelineCommand,
		&TimelineCommand,

		// Environment
		&envCommand,
		&envGetCommand,
		&envSetCommand,
		&envUnsetCommand,

		// Domains
		&DomainsListCommand,
		&DomainsAddCommand,
		&DomainsRemoveCommand,
		&DomainsSSLCommand,

		// Deployments
		&deploymentsListCommand,
		&deploymentLogCommand,
		&deploymentFollowCommand,
		&deploymentDeployCommand,
		&deploymentCacheResetCommand,

		// Collaborators
		&CollaboratorsListCommand,
		&CollaboratorsAddCommand,
		&CollaboratorsRemoveCommand,

		// Stacks
		&stacksListCommand,
		&stacksSetCommand,

		// Addons
		&AddonProvidersListCommand,
		&AddonProvidersPlansCommand,
		&addonsListCommand,
		&addonsAddCommand,
		&addonsRemoveCommand,
		&addonsUpgradeCommand,
		&addonsInfoCommand,

		// Integration Link
		&integrationLinkShowCommand,
		&integrationLinkCreateCommand,
		&integrationLinkUpdateCommand,
		&integrationLinkDeleteCommand,
		&integrationLinkManualDeployCommand,
		&integrationLinkManualReviewAppCommand,

		// Review Apps
		&reviewAppsShowCommand,

		// Notifiers
		&NotifiersListCommand,
		&NotifiersDetailsCommand,
		&NotifiersAddCommand,
		&NotifiersUpdateCommand,
		&NotifiersRemoveCommand,

		// Notification platforms
		&NotificationPlatformListCommand,

		// DB Access
		&DbTunnelCommand,
		&RedisConsoleCommand,
		&MongoConsoleCommand,
		&MySQLConsoleCommand,
		&PgSQLConsoleCommand,
		&InfluxDBConsoleCommand,

		// Databases
		&databaseBackupsConfig,
		&databaseEnableFeature,
		&databaseDisableFeature,

		// Backups
		&backupsListCommand,
		&backupsCreateCommand,
		&backupsDownloadCommand,
		&backupDownloadCommand,

		// Alerts
		&alertsListCommand,
		&alertsAddCommand,
		&alertsUpdateCommand,
		&alertsEnableCommand,
		&alertsDisableCommand,
		&alertsRemoveCommand,

		// Stats
		&StatsCommand,

		// Autoscalers
		&autoscalersListCommand,
		&autoscalersAddCommand,
		&autoscalersRemoveCommand,
		&autoscalersUpdateCommand,
		&autoscalersDisableCommand,
		&autoscalersEnableCommand,

		// Migrations
		&migrationCreateCommand,
		&migrationRunCommand,
		&migrationAbortCommand,
		&migrationListCommand,
		&migrationFollowCommand,

		// Log drains
		&logDrainsAddCommand,
		&logDrainsListCommand,
		&logDrainsRemoveCommand,

		&gitSetup,
		&gitShow,

		// Cron tasks
		&cronTasksListCommand,
	}

	globalCommands = []*cli.Command{
		// SSH keys
		&listSSHKeyCommand,
		&addSSHKeyCommand,
		&removeSSHKeyCommand,

		&integrationsListCommand,
		&integrationsAddCommand,
		&integrationsDeleteCommand,
		&integrationsImportKeysCommand,

		// Sessions
		&LoginCommand,
		&LogoutCommand,
		&RegionsListCommand,
		&ConfigCommand,
		&selfCommand,

		// Version
		&UpdateCommand,

		// Changelog
		&changelogCommand,

		// Help
		&HelpCommand,
	}
)
