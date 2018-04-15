package stump

import (
	"github.com/m1ome/stump/package/cli"
	"github.com/m1ome/stump/package/db"
	"github.com/m1ome/stump/package/migrate"
	"github.com/m1ome/stump/package/logger"
)

func (s *Stump) migrate() cli.Command {
	// Create default migrations
	return cli.Command{
		Name:    "migrate",
		Aliases: []string{"m"},
		Usage:   "database migration",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name: "path",
				Value: "./migrations",
				Usage: "Directory with migrations",
			},
		},
		Before: func(c *cli.Context) error {
			s.Logger().Info("Working with migrations")
			return nil
		},
		Subcommands: []cli.Command{
			{
				Name:  "latest",
				Usage: "migrate to latest version",
				Action: func(c *cli.Context) error {
					s.Logger().Info("Migrating to latest version")

					m, err := migrator(s.DB(), s.Logger(), c)
					if err != nil {
						return err
					}

					return m.Up(0)
				},
			},
			{
				Name:  "rollback",
				Usage: "rollback latest migration",
				Action: func(c *cli.Context) error {
					s.Logger().Info("Rolling back migration")

					m, err := migrator(s.DB(), s.Logger(), c)
					if err != nil {
						return err
					}

					return m.Down(1)
				},
			},
			{
				Name:  "version",
				Usage: "see database version",
				Action: func(c *cli.Context) error {
					m, err := migrator(s.DB(), s.Logger(), c)
					if err != nil {
						return err
					}

					v, err :=  m.VersionName()
					if err != nil {
						return err
					}

					s.Logger().Infof("Database version: %s", v)
					return nil
				},
			},
			{
				Name: "create",
				Usage: "create new migration",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name: "name",
						Usage: "Migration name",
					},
				},
				Action: func(c *cli.Context) error {
					path := c.GlobalString("path")
					name := c.Args().Get(0)
					s.Logger().Infow("Creating new migration", "path", path, "name", name)

					return migrate.CreateMigration(path, name)
				},
			},
		},
	}
}

func migrator(db *db.DB, logger *logger.Logger, c *cli.Context) (migrate.Migrator, error) {
	return migrate.New(migrate.Options{
		DB:     db,
		Logger: logger,
		Path:   c.GlobalString("path"),
	})
}
