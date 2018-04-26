package metrics

import (
	"strconv"
	"time"

	"github.com/m1ome/stump/package/web"
)

func WebMiddleware(next web.HandlerFunc) web.HandlerFunc {
	return func(ctx web.Context) error {
		start := time.Now()

		err := next(ctx)

		status := ctx.Response().Status
		if err != nil {
			if e, ok := err.(*web.HTTPError); ok {
				status = e.Code
			}
		}

		ControllerSummary.
			WithLabelValues(ctx.Path(), strconv.Itoa(status)).
			Observe(time.Since(start).Seconds())

		return err
	}
}
