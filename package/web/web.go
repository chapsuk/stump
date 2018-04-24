package web

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/pkg/errors"
)

type Options struct {
	Debug               bool
	RemoveTrailingSlash bool
}

type (
	Engine             = echo.Echo
	Context            = echo.Context
	HandlerFunc        = echo.HandlerFunc
	MiddlewareFunc     = echo.MiddlewareFunc
	Group              = echo.Group
	BasicAuthValidator = middleware.BasicAuthValidator
	HTTPError          = echo.HTTPError
)

type Web struct {
	opts   *Options
	engine *Engine
}

func New(opts *Options) *Web {
	w := new(Web)

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Debug = opts.Debug
	e.HTTPErrorHandler = w.DefaultHTTPErrorHandler

	if !opts.RemoveTrailingSlash {
		e.Pre(middleware.AddTrailingSlash())
	}

	w.engine = e
	w.opts = opts
	return w
}

func (w *Web) Engine() *Engine {
	return w.engine
}

func (w *Web) ListenAndServe(address string) error {
	startFailure := make(chan error)
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	go func() {
		if err := w.engine.Start(address); err != nil {
			startFailure <- err
		}
	}()

	select {
	case <-quit:
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := w.engine.Shutdown(ctx); err != nil {
			return errors.Wrap(err, "error shutdowning server")
		}
	case err := <-startFailure:
		return errors.Wrap(err, "error starting web server")
	}

	return nil
}

func (w *Web) DefaultHTTPErrorHandler(err error, c Context) {
	var (
		code = http.StatusInternalServerError
		msg  interface{}
	)

	if he, ok := err.(*HTTPError); ok {
		code = he.Code
		msg = he.Message
		if he.Internal != nil {
			msg = fmt.Sprintf("%v, %v", err, he.Internal)
		}
	} else if w.opts.Debug {
		msg = err.Error()
	} else {
		msg = http.StatusText(code)
	}

	res := Error{Error: "Internal Server Error", Code: code}
	if err, ok := msg.(string); ok {
		res.Error = err
	}

	if w.opts.Debug {
		res.StackTrace = getErrorStackTrace(err)
	}

	// Send response
	if !c.Response().Committed {
		if c.Request().Method == echo.HEAD {
			err = c.NoContent(code)
		} else {
			err = c.JSON(code, Response{Error: res})
		}
	}
}
