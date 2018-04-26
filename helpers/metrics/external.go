package metrics

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type ExternalSummary = prometheus.SummaryVec

type ExternalOptions struct {
	Name     string
	Metric   ExternalSummary
	Start    time.Time
	Request  *http.Request
	Response *http.Response
}

// WriteExternalCall writes external call summary,
// if response is nil set status label to `undefined`
// if request is nil set handler label to `undefined``
func WriteExternalCall(opts ExternalOptions) {
	status := "undefined"
	if opts.Response != nil {
		status = strconv.Itoa(opts.Response.StatusCode)
	}
	handler := "undefined"
	if opts.Request != nil && opts.Request.URL != nil {
		handler = opts.Request.URL.Path
	}
	opts.Metric.
		WithLabelValues(opts.Name, handler, status).
		Observe(time.Since(opts.Start).Seconds())
}
