package metrics

import (
	"net/http"
	"strconv"
	"time"
)

type ExternalOptions struct {
	Name     string
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
	method := "undefined"
	if opts.Request != nil && opts.Request.URL != nil {
		method = opts.Request.Method
	}
	ExternalSummary.
		WithLabelValues(opts.Name, method, status).
		Observe(time.Since(opts.Start).Seconds())
}
