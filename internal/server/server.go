package server

import (
	"net/http"
	"strings"

	"github.com/c-roussel/fizzbuzz-api/internal/handlers"
	"github.com/go-playground/validator"
	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

// urlSkipper middleware ignores metrics on some route
func urlSkipper(c echo.Context) bool {
	if strings.HasPrefix(c.Path(), "/mon") {
		return true
	}
	return false
}

func New() *echo.Echo {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Enable metrics middleware
	p := prometheus.NewPrometheus("echo", nil)
	p.MetricsPath = "/mon/metrics"
	p.Use(e)

	// Default data validation
	e.Validator = &CustomValidator{validator: validator.New()}

	// Routes
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.GET("/mon/ping", handlers.Ping)
	e.GET("/fizzbuzz", handlers.FizzBuzz)
	e.GET("/fizzbuzz/stats", handlers.FizzBuzzStats)

	return e
}
