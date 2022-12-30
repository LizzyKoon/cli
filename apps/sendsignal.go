package apps

import (
	"context"

	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/go-scalingo/v6"
	"github.com/Scalingo/go-utils/errors"
)

func keepUniqueContainersWithNames(containers []scalingo.Container, names []string) map[string]scalingo.Container {
	containersToKill := map[string]scalingo.Container{}

	for _, name := range names {
		hasMatched := false
		for _, container := range containers {
			if container.Label == name {
				containersToKill[container.ID] = container
				hasMatched = true
			}
		}
		if !hasMatched {
			io.Statusf("The name '%v' did not match any container\n", name)
		}
	}

	return containersToKill
}

func SendSignal(ctx context.Context, appName string, signal string, containerNames []string) error {
	if len(containerNames) == 0 {
		return errgo.New("at least one container name should be given")
	}
	if signal == "" {
		return errgo.New("signal must not be empty")
	}

	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client to send signal to application containers")
	}

	containers, err := c.AppsContainersPs(ctx, appName)
	if err != nil {
		return errgo.Notef(err, "fail to list the application containers to get the ID of the container to send the signal")
	}

	containersToKill := keepUniqueContainersWithNames(containers, containerNames)

	for _, container := range containersToKill {
		err := c.ContainersKill(ctx, appName, signal, container.ID)
		if err != nil {
			rootError := errors.RootCause(err)
			io.Statusf("Fail to send signal to container '%v' because of: %v\n", container.Label, rootError)
			continue
		}
		io.Statusf("Sent signal '%v' to '%v' container.\n", signal, container.Label)
	}
	return nil
}
