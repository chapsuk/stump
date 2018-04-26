package controllers

import (
	"github.com/m1ome/stump/lib"
	"github.com/m1ome/stump/package/web"
	"github.com/m1ome/stump/examples/basic/models"
	"github.com/m1ome/stump/package/crud"

	"net/http"
)

type Controller struct {
	stump *lib.Stump
}

func New(s *lib.Stump) *Controller {
	return &Controller{stump: s}
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
