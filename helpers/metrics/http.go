package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Config interface, example:
//	type Config struct {
//		Metrics *MetricsConfig `yaml:"metrics"`
//	}
//	type MetricsConfig struct {
//		Address string
//	}
//
//	func (c *Config) MetricsAddress() (string, bool) {
//		if c.Metrics == nil || c.Metrics.Address == "" {
//			return "", false
//		}
//		return c.Metrics.Address, true
//	}
type Config interface {
	MetricsAddress() (string, bool)
}

// Logger interface
type Logger interface {
	Info(i ...interface{})
	Errorw(message string, args ...interface{})
}

// ServeOptions struct
type ServeOptions struct {
	Config Config
	Logger Logger
}

// Serve run http.Server for prometheus metrics
func Serve(opts ServeOptions) {
	if addr, ok := opts.Config.MetricsAddress(); ok {
		if err := http.ListenAndServe(addr, promhttp.Handler()); err != nil {
			opts.Logger.Errorw("start metrics server error", "error", err)
		}
	}
	opts.Logger.Info("Prometheus metrics disabled")
}
