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

	_ "github.com/senayst4745/diamond/docs" // you need to update github.com/rizalgowandy/go-swag-sample with your own project path
	sw "github.com/swaggo/echo-swagger"
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

// @title Diamond And Mine
// @version 1.0
// @description This is diamond play.

// @host localhost:7171
// @BasePath /
// @schemes http
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
	e.GET("/swagger/*", sw.WrapHandler)
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

// healthcheck to check is server alive
// @Summary Show the status of server.
// @Description get the status of server.
// @Accept */*
// @Produce json
// @Success 200
// @Router /health [get]
func healthcheck(e echo.Context) error {
	return e.JSON(http.StatusOK, struct {
		Message string
	}{Message: "OK"})
}

type mineName struct {
	Name string `json:"name"`
}

type mineCount struct {
	Count int `json:"count"`
}

// getDiamonds to get random count of diamonds from mine
// @Summary Get Diamonds from Mine.
// @Description method to get some diamonds from a mine. If mine is empty, this method delete this mine.
// @Accept */*
// @Produce json
// @Success 200 {object} mineCount
// @Param body body mineName true "the name of the mine from which we want to extract diamonds"
// @Router /diamond [post]
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

// getAllMines Show all mines
// @Summary Show all mines.
// @Accept */*
// @Produce json
// @Success 200 {object} []repository.Mine "list of all mines"
// @Router /mine [get]
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

// addDiamondMine add new mine
// @Summary Add new mine.
// @Accept */*
// @Produce json
// @Success 201
// @Param body body repository.Mine true "new mine model"
// @Router /mine [post]
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

// emptyMine gets all the diamonds from the mine and closes it
// @Summary Closes mine.
// @Description gets all the diamonds from the mine and closes it
// @Accept */*
// @Produce json
// @Success 200 {object} mineCount
// @Param body body mineName true "the name of the mine to be closed"
// @Router /mine [delete]
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
