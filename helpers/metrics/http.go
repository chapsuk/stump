package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Config interface {
	MetricsAddress() (string, bool)
}

type Logger interface {
	Info(i ...interface{})
	Errorw(message string, args ...interface{})
}

type ServeOptions struct {
	Config Config
	Logger Logger
}

func Serve(opts ServeOptions) {
	if addr, ok := opts.Config.MetricsAddress(); ok {
		if err := http.ListenAndServe(addr, promhttp.Handler()); err != nil {
			opts.Logger.Errorw("start metrics server error", "error", err)
		}
	}
	opts.Logger.Info("Prometheus metrics disabled")
}
