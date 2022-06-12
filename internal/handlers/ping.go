package handlers

import (
	"encoding/json"
	"log"
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

// No need to re-compute json marshalling at every ping.
var pingOut json.RawMessage

func init() {
	var err error
	pingOut, err = json.Marshal(PingOutput{
		Message: "OK",
		GitHash: os.Getenv("GIT_HASH"),
	})
	if err != nil {
		// should never happen
		log.Fatalf("failed to unmarshal ping's output: %v", err)
	}
}

// Ping handles /mon/ping HTTP requests.
//
// It will respond with a 200 HTTP repsonse embedding
// a PingOutput result.
//
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
