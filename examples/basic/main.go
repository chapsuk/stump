package main

import (
	"net/http"

	"github.com/m1ome/stump"
	"github.com/m1ome/stump/package/web"
	"github.com/m1ome/stump/package/crud"

	"github.com/m1ome/stump/examples/basic/models"
	"github.com/m1ome/stump/package/worker"
)

//
// Controllers example
//

type Controller struct {
	stump *stump.Stump
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
	s, err := stump.New(stump.Options{})
	if err != nil {
		s.Logger().Panicf("Error creating Stump instance: %v", err)
	}

	if err := s.Storages(&stump.StorageOptions{
		Postgres: true,
		Redis:    true,
	}); err != nil {
		s.Logger().Panicf("Error connection to storages: %v", err)
	}

	c := &Controller{stump: s}
	s.Web().Engine().POST("/", c.Register)
	s.Web().Engine().GET("/", c.UserList)

	s.Workers().Schedule(worker.Task{
		Name:      "example",
		Exclusive: true,
		Scheduler: "@every 30s",
		Handler: func(ctx worker.Context) error {
			s.Logger().Info("Updating user ratings")

			var users []models.User
			if err := crud.FindAll(s.DB(), &users); err != nil {
				s.Logger().Errorf("Error finding users: %v", err)
				return err
			}

			for _, user := range users {
				s.Logger().Infow("Updating user rating", "user", user)
				user.Rating += 1
				if _, err := s.DB().Model(&user).Where("id=?", user.ID).Update(); err != nil {
					s.Logger().Errorf("Error updating user: %v", err)
					return err
				}
			}

			return nil
		},
	})
	s.Workers().Start()

	if err := s.Start("example", "Example application"); err != nil {
		s.Logger().Errorf("Starting application error: %v", err)
	}
}
