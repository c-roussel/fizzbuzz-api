package handlers

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

// PintOutput is the result of a /mon/ping call.
//
// The git_hash will only be populated depending on
// a GIT_HASH environment variable.
type PingOutput struct {
	Message string `json:"message"`
	GitHash string `json:"git_hash"`
}

// No need to re-compute at every ping.
var pingOut = PingOutput{
	Message: "OK",
	GitHash: os.Getenv("GIT_HASH"),
}

// Ping handles /mon/ping http queries
//
// It will respond with a 200 HTTP repsonse embedding
// a PingOutput result
// @Summary Show the status of server.
// @Description get the status of server.
// @Tags monitoring
// @Accept */*
// @Produce json
// @Success 200 {object} handlers.PingOutput
// @Router /mon/ping [get]
func Ping(c echo.Context) error {
	return c.JSON(http.StatusOK, pingOut)
}
