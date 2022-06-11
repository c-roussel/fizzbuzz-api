package handlers

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

// FizzBuzzEnvLimit  is the environment variable to override the servers
// maximum limit on GET /fizzbuzz route
const FizzBuzzEnvLimit = "FIZZBUZZ_MAX_LIMIT"

// FizzBuzzMaxLimit is the maximum threshold for GET /fizzbuzz limit parameter
var FizzBuzzMaxLimit = 10000

func init() {
	if envLimitStr := os.Getenv(FizzBuzzEnvLimit); envLimitStr != "" {
		envLimit, err := strconv.Atoi(envLimitStr)
		if err != nil {
			log.Error(err.Error())
			return
		}

		FizzBuzzMaxLimit = envLimit
	}
}

// FizzBuzzInput describes the expected input for the fizzbuzz handler
type FizzBuzzInput struct {
	Str1  string `query:"str1" validate:"required"`
	Str2  string `query:"str2" validate:"required"`
	Int1  int    `query:"int1" validate:"required,min=1"`
	Int2  int    `query:"int2" validate:"required,min=1"`
	Limit int    `query:"limit" validate:"required,min=0"`
}

// Register increments the input parameters in fizzbuzz statistics
// See FizzBuzzStats
func (in FizzBuzzInput) Register() {
	fizzBuzzGatherer.Hit(in.String())
}

// String is the string representing of a FizzBuzzInput
func (in FizzBuzzInput) String() string {
	return fmt.Sprintf("FizzBuzzInput str1=%s str2=%s int1=%d int2=%d limit=%d",
		in.Str1, in.Str2, in.Int1, in.Int2, in.Limit)
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

	if in.Limit > FizzBuzzMaxLimit {
		c.Logger().Warnf("limit %d is higher than threshold %d", in.Limit, FizzBuzzMaxLimit)
		return echo.NewHTTPError(
			http.StatusBadRequest,
			fmt.Sprintf("limit should be lower than %d", FizzBuzzMaxLimit),
		)
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
			slice[i] = strconv.Itoa(v)
		}
	}

	go in.Register()

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
