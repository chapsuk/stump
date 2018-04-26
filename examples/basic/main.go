package main

import (
	"context"
	"net/http"
	"time"

	"github.com/chapsuk/worker"

	"github.com/m1ome/stump"
	"github.com/m1ome/stump/examples/basic/models"
	"github.com/m1ome/stump/lib"
	"github.com/m1ome/stump/package/crud"
	"github.com/m1ome/stump/package/web"
	"github.com/m1ome/stump/package/worker_helpers"
	"github.com/m1ome/stump/helpers"
)

//
// Controllers example
//

type Controller struct {
	stump *lib.Stump
}

func (c *Controller) Register(ctx web.Context) error {
	req := models.User{}
	if err := ctx.Bind(&req); err != nil {
		return web.BadRequestf("error binding data: %v", err)
	}

	user, err := models.NewUser(req.Name, req.Email)
	if err != nil {
		return web.BadRequestf("error validating user: %v", err)
	}

	if err := crud.Create(c.stump.DB(), user); err != nil {
		return err
	}

	c.stump.Logger().Infow("Created new user", "user", user)
	return ctx.JSON(http.StatusOK, user)
}

func (c *Controller) UserList(ctx web.Context) error {
	var users []models.User
	if err := crud.FindAll(c.stump.DB(), &users); err != nil {
		return err
	}

	c.stump.Logger().Infow("User list", "users", users)
	return ctx.JSON(http.StatusOK, users)
}

//
// Workers handler
//

type Workers struct {
	stump *lib.Stump
	g     *worker.Group
}

func (w *Workers) UpdateUserRatings(ctx context.Context) {
	w.stump.Logger().Info("Updating user ratings")

	var users []models.User
	if err := crud.FindAll(w.stump.DB(), &users); err != nil {
		w.stump.Logger().Errorf("Error finding users: %v", err)
		return
	}

	for _, user := range users {
		w.stump.Logger().Infow("Updating user rating", "user", user)
		user.Rating += 1
		if _, err := w.stump.DB().Model(&user).Where("id=?", user.ID).Update(); err != nil {
			w.stump.Logger().Errorf("Error updating user: %v", err)
			return
		}
	}

	return
}

func (w *Workers) Start() {
	wg := worker.NewGroup()
	wg.Add(
		helpers.ScheduleWithLock(w.stump.Redis(), w.UpdateUserRatings, time.Second*30, helpers.LockOptions{
			Key:     "locker",
			TTL:     time.Minute,
			Logger:  w.stump.Logger(),
			Retries: 0,
		}),
	)

	w.g = wg
	w.g.Run()
}

func (w *Workers) Stop() {
	w.g.Stop()
}

//
// Main function
//

func main() {
	s := stump.MustSetup()
	c := &Controller{stump: s}
	w := &Workers{stump: s}

	serve := helpers.CliCommandServe(s, func() error {
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
