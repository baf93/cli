// +build rabinary

package container

import (
	"fmt"
	"strings"
	"os/exec"

	"github.com/docker/cli/cli"
	"github.com/docker/cli/cli/command"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"github.com/docker/docker/api/types"
)

type raInfoOptions struct {
	time        int
	timeChanged bool

	containers []string
}

// NewRaInfoCommand creates a new cobra.Command for `docker stop`
func NewRaInfoCommand(dockerCli *command.DockerCli) *cobra.Command {
	var opts raInfoOptions

	cmd := &cobra.Command{
		Use:   "raInfo [OPTIONS] CONTAINER [CONTAINER...]",
		Short: "Get info for Ra for one or more running containers",
		Args:  cli.RequiresMinArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.containers = args
			opts.timeChanged = cmd.Flags().Changed("time")
			return runRaInfo(dockerCli, &opts)
		},
	}

	flags := cmd.Flags()
	flags.IntVarP(&opts.time, "time", "t", 10, "Seconds to wait for stop before killing it")
	return cmd
}


func runRaInfo(dockerCli *command.DockerCli, opts *raInfoOptions) error {

	ctx := context.Background()

	options := &types.ContainerListOptions{
		Quiet: true,
	}

	var (
		cmdOut []byte
		err    error
	)
	cmdName := "lsblk"
	cmdArgs := []string{}

	if cmdOut, err = exec.Command(cmdName, cmdArgs...).Output(); err != nil {
		fmt.Fprintln(dockerCli.Out(), "There was an error running lslk command: ", err)
		return nil
	}

	out := string(cmdOut)
	//fmt.Fprintln(dockerCli.Out(), out)
	
	
	var errs []string

	containers, err := dockerCli.Client().ContainerList(ctx, *options)
	
	if err != nil{
		errs = append(errs, err.Error())
	}

/*	for _, container := range containers{
		c, err := dockerCli.Client().ContainerInspect(ctx, container.ID)
		if err != nil{
			errs = append(errs, err.Error())
			continue
		}
		
		data := c.GraphDriver.Data["DeviceName"]

		if len(data)<0 {
			errs = append(errs, "DeviceMapper not in use")
			continue
		}

		
	
		
		fmt.Fprintln(dockerCli.Out(), container.ID[0:12] + " " + data + " " + strings.Split(out, data)[1][0:7])	
	}*/

	for _, container := range containers{
		c, err := dockerCli.Client().ContainerInspect(ctx, container.ID)
		if err != nil{
			errs = append(errs, err.Error())
			continue
		}
		
		data := c.GraphDriver.Data["DeviceName"]

		if len(data)<0 {
			errs = append(errs, "DeviceMapper not in use")
			continue
		}

		fmt.Fprintln(dockerCli.Out(), container.ID[0:12] + " " + strings.Split(strings.Split(out, data)[1][1:], " ")[0])	
	}

	return nil
}
