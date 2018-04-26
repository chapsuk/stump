package cli

import (
	"github.com/m1ome/stump/package/cli"
	"github.com/m1ome/stump/core"
)

func CliCommandServe(s *core.Stump, fn func() error) cli.Command {
	return cli.Command{
		Name:    "serve",
		Aliases: []string{"s"},
		Usage:   "start web server",
		Action: func(c *cli.Context) error {
			// Callback to handle before we start to serve actually
			if err := fn(); err != nil {
				return err
			}

			return s.ServeHTTP()
		},
	}
}
