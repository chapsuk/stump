package workers

import (
	"github.com/m1ome/stump/core"
	"github.com/chapsuk/worker"
	"context"
	"github.com/m1ome/stump/examples/basic/models"
	"github.com/m1ome/stump/package/crud"
	"github.com/m1ome/stump/helpers"
	"time"
)

type Workers struct {
	stump *core.Stump
	g     *worker.Group
}

func New(s *core.Stump) *Workers {
	return &Workers{
		stump: s,
		g:     worker.NewGroup(),
	}
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
	w.stump.Logger().Infow("Starting workers")

	w.g.Add(
		helpers.ScheduleWithLock(w.stump.Redis(), w.UpdateUserRatings, time.Second*30, helpers.LockOptions{
			Key:     "locker",
			TTL:     time.Minute,
			Logger:  w.stump.Logger(),
			Retries: 0,
		}),
	)
	w.g.Run()
}

func (w *Workers) Stop() {
	w.g.Stop()
}
