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
// Main function
//

func main() {
	s := stump.MustSetup()
	c := &Controller{stump: s}

	server := s.ServeCommand(func() error {
		if err := s.Init(true, true); err != nil {
			return err
		}

		// Binding web
		s.Web().Engine().POST("/", c.Register)
		s.Web().Engine().GET("/", c.UserList)

		job := func(ctx context.Context) {
			s.Logger().Info("Updating user ratings")

			var users []models.User
			if err := crud.FindAll(s.DB(), &users); err != nil {
				s.Logger().Errorf("Error finding users: %v", err)
				return
			}

			for _, user := range users {
				s.Logger().Infow("Updating user rating", "user", user)
				user.Rating += 1
				if _, err := s.DB().Model(&user).Where("id=?", user.ID).Update(); err != nil {
					s.Logger().Errorf("Error updating user: %v", err)
					return
				}
			}

			return
		}

		wg := worker.NewGroup()
		wg.Add(
			worker.New(job).WithBsmRedisLock(worker_helpers.RedisLockOptions(s.Redis(), worker_helpers.Options{
				Key:     "example",
				TTL:     time.Second,
				Retries: 0,
			})).ByTicker(time.Second * 30),
		)

		wg.Run()
		return nil
	})

	s.Cli().Add(server)
	if err := s.Run(); err != nil {
		s.Logger().Errorf("Starting application error: %v", err)
	}
}
