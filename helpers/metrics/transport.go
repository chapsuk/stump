package metrics

import (
	"net/http"
	"time"
)

type transport struct {
	name    string
	tripper http.RoundTripper
}

func applyTransport(name string, client *http.Client) *http.Client {
	client.Transport = &transport{
		name:    name,
		tripper: client.Transport,
	}

	return client
}

// RoundTrip defer metric.write and apply base tripper
func (t *transport) RoundTrip(req *http.Request) (res *http.Response, err error) {
	defer func(dt time.Time) {
		WriteExternalCall(ExternalOptions{
			Name:     t.name,
			Start:    dt,
			Request:  req,
			Response: res,
		})
	}(time.Now())

	if t.tripper == nil {
		t.tripper = http.DefaultTransport
	}

	res, err = t.tripper.RoundTrip(req)

	return
}

func WrapClient(name string, client *http.Client) *http.Client {
	return applyTransport(name, client)
}
