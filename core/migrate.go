package core

import (
	"github.com/m1ome/stump/package/cli"
	"github.com/m1ome/stump/package/db"
	"github.com/m1ome/stump/package/logger"
	"github.com/m1ome/stump/package/migrate"
)

func cliMigrate(s *Stump) cli.Command {
	// Create default migrations
	return cli.Command{
		Name:    "migrate",
		Aliases: []string{"m"},
		Usage:   "database migration",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "path",
				Value: "./migrations",
				Usage: "Directory with migrations",
			},
		},
		Before: func(c *cli.Context) error {
			s.logger.Info("Working with migrations")
			return nil
		},
		Subcommands: []cli.Command{
			{
				Name:  "latest",
				Usage: "migrate to latest version",
				Action: func(c *cli.Context) error {
					s.logger.Info("Migrating to latest version")

					m, err := migrator(s.db, s.logger, c)
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
					s.logger.Info("Rolling back migration")

					m, err := migrator(s.db, s.logger, c)
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
					m, err := migrator(s.db, s.logger, c)
					if err != nil {
						return err
					}

					v, err := m.VersionName()
					if err != nil {
						return err
					}

					s.logger.Infof("Database version: %s", v)
					return nil
				},
			},
			{
				Name:  "create",
				Usage: "create new migration",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "name",
						Usage: "Migration name",
					},
				},
				Action: func(c *cli.Context) error {
					path := c.GlobalString("path")
					name := c.Args().Get(0)
					s.logger.Infow("Creating new migration", "path", path, "name", name)

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
