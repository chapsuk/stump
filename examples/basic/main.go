package main

import (
	"github.com/m1ome/stump"
	"github.com/m1ome/stump/helpers"
	"github.com/m1ome/stump/examples/basic/controllers"
	"github.com/m1ome/stump/examples/basic/workers"
)

//
// Controllers example
//


//
// Workers handler
//

//
// Main function
//
func main() {
	s := stump.MustSetup()

	serve := helpers.CliCommandServe(s, func() error {
		// Initializing controller & workers
		c := controllers.New(s)
		w := workers.New(s)

		if err := s.Init(true, true); err != nil {
			return err
		}

		// Binding web
		s.Web().Engine().POST("/", c.Register)
		s.Web().Engine().GET("/", c.UserList)

		// Starting workers
		w.Start()

		return nil
	})

	// Adding serve command
	s.Cli().Add(serve)

	// Running instance
	if err := s.Run(); err != nil {
		s.Logger().Errorf("Starting application error: %v", err)
	}
}
