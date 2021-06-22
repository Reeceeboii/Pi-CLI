package cli

import (
	"github.com/Reeceeboii/Pi-CLI/pkg/api"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
)

/*
	Enable the Pi-Hole if it is not already enabled,
*/
func RunEnablePiHoleCommand(*cli.Context) error {
	InitialisePICLI()
	api.LiveSummary.Update(nil)

	if api.LiveSummary.Status == "enabled" {
		color.Yellow("Pi-Hole is already enabled!")
	} else {
		api.EnablePiHole()
		color.Green("Pi-Hole enabled")
	}

	return nil
}

/*
	Disable the Pi-Hole. This command also takes an optional timeout parameter in seconds.
	If given and within constraints, the Pi-Hole will automatically re-enable after this
	time period has elapsed
*/
func RunDisablePiHoleCommand(c *cli.Context) error {
	InitialisePICLI()
	api.LiveSummary.Update(nil)

	if api.LiveSummary.Status == "disabled" {
		color.Yellow("Pi-Hole is already disabled!")
	} else {
		timeout := c.Int64("timeout")
		if timeout == 0 {
			api.DisablePiHole(false, 0)
			color.Green("Pi-Hole disabled until explicitly re-enabled")
		} else {
			api.DisablePiHole(true, timeout)
			color.Green("Pi-Hole disabled. Will re-enable in %d seconds\n", timeout)
		}
	}

	return nil
}
