package cli

import (
	"os"
	"sort"

	"github.com/urfave/cli"
)

type Cli struct {
	app *App
}

type Options struct {
	Name    string
	Usage   string
	Version string
}

func (o *Options) init() {

}

func New(opts *Options) *Cli {
	opts.init()

	// Creating main application
	app := cli.NewApp()

	// Basic stuff
	app.Name = opts.Name
	app.Usage = opts.Usage
	app.Version = opts.Version

	return &Cli{
		app: app,
	}
}

func (c *Cli) Run() error {
	sort.Sort(cli.FlagsByName(c.app.Flags))
	sort.Sort(cli.CommandsByName(c.app.Commands))
	return c.app.Run(os.Args)
}

func (c *Cli) Add(command Command) {
	c.app.Commands = append(c.app.Commands, command)
}
