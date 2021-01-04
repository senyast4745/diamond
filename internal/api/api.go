package api

import (
	"context"
	"errors"
	"github.com/senayst4745/diamond/internal/repository"
	"github.com/senayst4745/diamond/internal/service"
	"net/http"
	"reflect"
	"time"

	"github.com/labstack/echo/v4"
	zl "github.com/rs/zerolog/log"
)

type Context struct {
	echo.Context
	Ctx context.Context
}

type API struct {
	e    *echo.Echo
	s    service.Service
	addr string
}

type Config struct {
	Addr string
}

func New(ctx context.Context, cfg *Config, s service.Service) (*API, error) {
	e := echo.New()
	a := &API{
		s:    s,
		e:    e,
		addr: cfg.Addr,
	}
	// inject global cancellable context to every request
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &Context{
				Context: c,
				Ctx:     ctx,
			}
			return next(cc)
		}
	})
	e.Use(logger())
	e.GET("/health", healthcheck)
	e.POST("/diamonds", a.getDiamonds)
	e.GET("/mine", a.getAllMines)
	e.PUT("/mine", a.addDiamondMine)
	e.DELETE("/mine", a.emptyMine)
	return a, nil
}

func logger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			res := c.Response()
			start := time.Now()

			err := next(c)
			stop := time.Now()

			zl.Debug().
				Str("remote", req.RemoteAddr).
				Str("user_agent", req.UserAgent()).
				Str("method", req.Method).
				Str("path", c.Path()).
				Int("status", res.Status).
				Dur("duration", stop.Sub(start)).
				Str("duration_human", stop.Sub(start).String()).
				Msgf("called url %s", req.URL)

			return err
		}
	}
}

func healthcheck(e echo.Context) error {
	return e.String(http.StatusOK, "ok")
}

type mineName struct {
	Name string `json:"name"`
}

type mineCount struct {
	Count int `json:"count"`
}

func (a *API) getDiamonds(e echo.Context) error {
	cc, err := getParentContext(e)
	if err != nil {
		return err
	}

	b := &mineName{}
	if err := e.Bind(b); err != nil {
		return err
	}
	zl.Debug().Interface("body", b).Msg("bind body")
	count, err := a.s.GetDiamonds(cc.Ctx, b.Name)
	if err != nil {
		zl.Error().Err(err).Str("mine name", b.Name).Msg("can not find diamonds")
		return err
	}
	res := &mineCount{Count: count}
	return e.JSON(http.StatusOK, res)
}

func (a *API) getAllMines(e echo.Context) error {
	cc, err := getParentContext(e)
	if err != nil {
		return err
	}

	m, err := a.s.GetAllMines(cc.Ctx)
	if err != nil {
		zl.Error().Err(err).Msg("can not find mines")
		return err
	}
	return e.JSON(http.StatusOK, m)
}

func (a *API) addDiamondMine(e echo.Context) error {
	cc, err := getParentContext(e)
	if err != nil {
		return err
	}
	b := &repository.Mine{}
	if err := e.Bind(b); err != nil {
		return err
	}
	if err := a.s.AddDiamondMine(cc.Ctx, b); err != nil {
		return err
	}
	return e.NoContent(http.StatusCreated)
}

func (a *API) emptyMine(e echo.Context) error {
	cc, err := getParentContext(e)
	if err != nil {
		return err
	}

	b := &mineName{}
	if err := e.Bind(b); err != nil {
		return err
	}
	count, err := a.s.EmptyMine(cc.Ctx, b.Name)
	if err != nil {
		zl.Error().Err(err).Str("mine name", b.Name).Msg("can not find diamonds")
		return err
	}
	res := &mineCount{Count: count}
	return e.JSON(http.StatusOK, res)
}

func getParentContext(e echo.Context) (*Context, error) {
	cc, ok := e.(*Context)
	if !ok {
		zl.Error().Interface("type", reflect.TypeOf(e)).Msg("can not cast to custom context")
		return nil, echo.ErrInternalServerError
	}
	return cc, nil
}

func (a *API) Start() error {
	zl.Debug().
		Msgf("listening on %v", a.addr)
	err := a.e.Start(a.addr)
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}

func (a *API) Close() error {
	return a.e.Close()
}
