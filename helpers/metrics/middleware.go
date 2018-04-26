package metrics

import (
	"strconv"
	"time"

	"github.com/cryptopay-dev/yaga/web"
	"github.com/labstack/echo"
	"github.com/prometheus/client_golang/prometheus"
)

type ControllerSummary = prometheus.ObserverVec

func Middleware(metric ControllerSummary) web.MiddlewareFunc {
	return func(next web.HandlerFunc) web.HandlerFunc {
		return func(ctx web.Context) error {
			start := time.Now()

			err := next(ctx)

			status := ctx.Response().Status
			if err != nil {
				if e, ok := err.(*echo.HTTPError); ok {
					status = e.Code
				}
			}

			metric.
				WithLabelValues(ctx.Path(), strconv.Itoa(status)).
				Observe(time.Since(start).Seconds())

			return err
		}
	}
}
