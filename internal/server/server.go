package server

import (
	"net/http"

	"github.com/c-roussel/fizzbuzz-api/internal/handlers"
	"github.com/go-playground/validator"
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

func New() *echo.Echo {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Default data validation
	e.Validator = &CustomValidator{validator: validator.New()}

	// Routes
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.GET("/mon/ping", handlers.Ping)
	e.GET("/fizzbuzz", handlers.FizzBuzz)

	return e
}
