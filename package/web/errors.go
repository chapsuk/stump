package web

import (
	"fmt"
	"net/http"
	"runtime"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

func BadRequest(message interface{}) *HTTPError {
	return echo.NewHTTPError(http.StatusBadRequest, message)
}

func BadRequestf(format string, args ...interface{}) *HTTPError {
	return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf(format, args))
}

func getErrorStackTrace(err error) []string {
	ster, ok := err.(interface {
		StackTrace() errors.StackTrace
	})
	if !ok {
		return nil
	}

	trace := make([]string, 0)
	for _, f := range ster.StackTrace() {
		pc := uintptr(f) - 1
		fn := runtime.FuncForPC(pc)
		var file string
		var line int
		if fn != nil {
			file, line = fn.FileLine(pc)
		} else {
			file = "unknown"
		}

		trace = append(trace, fmt.Sprintf("%v:%v", file, line))
	}

	return trace
}
