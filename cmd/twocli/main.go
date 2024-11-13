package main

import (
	"github.com/bykclk/twocli/internal/cli"
	"github.com/bykclk/twocli/internal/commands"
)

func main() {
	cmds := []cli.Command{
		commands.NewAddCommand(),
		commands.NewListCommand(),
		commands.NewCodeCommand(),
		commands.NewDeleteCommand(),
		commands.NewUpdateCommand(),
	}

	cli.Run(cmds)
}
