package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// FizzBuzzStats responds to GET /fizbuzz/stats http requests
//
// It will respond with a 200 HTTP repsonse embedding
// an array of stats.Count values
//
// The result is computed following the following algorithm:
// - Every succesfful GET /fizzbuzz will increment its parameters's stats
// - Respond with the top 100 stats
//
// @Summary Top 100 /fizzbuzz parameters.
// @Description Get the 100 most used parameters on GET /fizbuzz route.
// @Tags fizzbuzz
// @Accept */*
// @Produce json
// @Success 200 {array} stats.Count
// @Router /fizzbuzz/stats [get]
func FizzBuzzStats(c echo.Context) error {
	res := fizzBuzzGatherer.OrderedValues()
	return c.JSON(http.StatusOK, res[:min(len(res), 100)])
}

func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}
