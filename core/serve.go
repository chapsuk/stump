package core

import (
	"github.com/m1ome/stump/package/cli"
)

func cliServeHTTP(s *Stump) cli.Command {
	return cli.Command{
		Name:    "serve",
		Aliases: []string{"s"},
		Usage:   "start web server",
		Action: func(c *cli.Context) error {
			address := s.config.GetString("web.address")
			if address == "" {
				s.logger.Info("Binding on default port: 8080")
				address = ":8080"
			} else {
				s.logger.Infof("Start listening on: %v", address)
			}

			return s.web.ListenAndServe(address)
		},
	}
}
