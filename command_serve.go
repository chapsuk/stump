package stump

import (
	"github.com/m1ome/stump/package/cli"
)

func (s *Stump) serve() cli.Command {
	return cli.Command{
		Name:    "serve",
		Aliases: []string{"s"},
		Usage:   "start web server",
		Action: func(c *cli.Context) error {
			address := s.Config().GetString("web.address")
			if address == "" {
				s.Logger().Info("Binding on default port: 8080")
				address = ":8080"
			} else {
				s.Logger().Infof("Start listening on: %v", address)
			}

			return s.Web().ListenAndServe(address)
		},
	}
}