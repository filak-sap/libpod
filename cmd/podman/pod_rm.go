package main

import (
	"fmt"
	"github.com/containers/libpod/libpod/adapter"

	"github.com/containers/libpod/cmd/podman/cliconfig"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	podRmCommand     cliconfig.PodRmValues
	podRmDescription = fmt.Sprintf(`
podman rm will remove one or more pods from the host. The pod name or ID can
be used.  A pod with containers will not be removed without --force.
If --force is specified, all containers will be stopped, then removed.
`)
	_podRmCommand = &cobra.Command{
		Use:   "rm",
		Short: "Remove one or more pods",
		Long:  podRmDescription,
		RunE: func(cmd *cobra.Command, args []string) error {
			podRmCommand.InputArgs = args
			podRmCommand.GlobalFlags = MainGlobalOpts
			return podRmCmd(&podRmCommand)
		},
		Args: func(cmd *cobra.Command, args []string) error {
			return checkAllAndLatest(cmd, args, false)
		},
		Example: `podman pod rm mywebserverpod
  podman pod rm -f 860a4b23
  podman pod rm -f -a`,
	}
)

func init() {
	podRmCommand.Command = _podRmCommand
	podRmCommand.SetUsageTemplate(UsageTemplate())
	flags := podRmCommand.Flags()
	flags.BoolVarP(&podRmCommand.All, "all", "a", false, "Remove all running pods")
	flags.BoolVarP(&podRmCommand.Force, "force", "f", false, "Force removal of a running pod by first stopping all containers, then removing all containers in the pod.  The default is false")
	flags.BoolVarP(&podRmCommand.Latest, "latest", "l", false, "Remove the latest pod podman is aware of")
	markFlagHiddenForRemoteClient("latest", flags)
}

// podRmCmd deletes pods
func podRmCmd(c *cliconfig.PodRmValues) error {
	runtime, err := adapter.GetRuntime(&c.PodmanCommand)
	if err != nil {
		return errors.Wrapf(err, "could not get runtime")
	}
	defer runtime.Shutdown(false)
	podRmIds, podRmErrors := runtime.RemovePods(getContext(), c)
	for _, p := range podRmIds {
		fmt.Println(p)
	}
	if len(podRmErrors) == 0 {
		return nil
	}
	// Grab the last error
	lastError := podRmErrors[len(podRmErrors)-1]
	// Remove the last error from the error slice
	podRmErrors = podRmErrors[:len(podRmErrors)-1]

	for _, err := range podRmErrors {
		logrus.Errorf("%q", err)
	}
	return lastError
}
