package commands

import "github.com/urfave/cli/v2"

var BlockchainAuthApp = []*cli.Command{
	vipCommand,
	mintCommand,
	burnCommand,
}

var vipCommand = &cli.Command{
	Name:   "vip",
	Usage:  "Trigger a hello vip message",
	Action: triggerAction("vip"),
}

var mintCommand = &cli.Command{
	Name:   "mint",
	Usage:  "Mint the token pass",
	Action: triggerAction("mint"),
}

var burnCommand = &cli.Command{
	Name:   "burn",
	Usage:  "Burn the token pass",
	Action: triggerAction("burn"),
}
