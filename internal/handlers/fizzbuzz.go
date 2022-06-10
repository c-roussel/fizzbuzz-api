package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// FizzBuzzInput describes the expected input for the fizzbuzz handler
type FizzBuzzInput struct {
	Str1  string `query:"str1" validate:"required"`
	Str2  string `query:"str2" validate:"required"`
	Int1  int    `query:"int1" validate:"required,min=1"`
	Int2  int    `query:"int2" validate:"required,min=1"`
	Limit int    `query:"limit" validate:"required,min=0"`
}

// FizzBuzzOutput describes the response output for the fizzbuzz handler
type FizzBuzzOutput struct {
	Result []string `json:"result"`
}

// FizzBuzz responds to GET /fizbuzz http queries
//
// It will respond with a 200 HTTP repsonse embedding
// a FizzBuzzOutput result
//
// The result is computed following the following algorithm:
// - Return a list from 1 to `limit`
// - Each multiple of int1 is replaced by str1
// - Each multiple of int2 is replaced by str2
// - Each multiple of both is replaced by str1+str2
//
// @Summary Customizable fizzbuzz algorithm.
// @Description Get your own version of the fizzbuzz algortihm.
// @Tags fizzbuzz
// @Accept */*
// @Param int1  query int    true "fizzbuzz's first multiple"     minimum(1)
// @Param int2  query int    true "fizzbuzz's second multiple"    minimum(1)
// @Param str1  query string true "fizzbuzz's first replacement"
// @Param str2  query string true "fizzbuzz's second replacement"
// @Param limit query int    true "fizzbuzz's up-to value"        minimum(0)
// @Produce json
// @Success 200 {object} handlers.FizzBuzzOutput
// @Router /fizzbuzz [get]
func FizzBuzz(c echo.Context) error {
	var in FizzBuzzInput
	err := c.Bind(&in)
	if err != nil {
		c.Logger().Warnf("failed to parse query parameters: %v", err)
		return err
	}

	err = c.Validate(&in)
	if err != nil {
		c.Logger().Warnf("failed to validate query parameters: %v", err)
		return err
	}

	slice := make([]string, in.Limit)
	int3 := lcm(in.Int1, in.Int2)
	str3 := in.Str1 + in.Str2
	for i := range slice {
		v := i + 1 // array shall start with value 1

		switch {
		case v%int3 == 0:
			slice[i] = str3
		case v%in.Int1 == 0:
			slice[i] = in.Str1
		case v%in.Int2 == 0:
			slice[i] = in.Str2
		default:
			// strconv is more efficient than fmt.Sprint
			slice[i] = strconv.FormatInt(int64(v), 10)
		}
	}

	return c.JSON(http.StatusOK, FizzBuzzOutput{Result: slice})
}

// greatest common divisor (GCD) via Euclidean algorithm
func gcd(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

// find Least Common Multiple (LCM) via GCD
func lcm(a, b int) int {
	return a * b / gcd(a, b)
}
