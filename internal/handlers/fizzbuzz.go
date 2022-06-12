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
// maximum limit on GET /fizzbuzz route.
const FizzBuzzEnvLimit = "FIZZBUZZ_MAX_LIMIT"

// FizzBuzzMaxLimit is the maximum threshold for GET /fizzbuzz limit parameter.
var FizzBuzzMaxLimit = 10000
var defaultFizzBuzzInput FizzBuzzInput

func init() {
	// setup default FizzBuzzInput values
	defaultStr1 := "fizz"
	defaultStr2 := "buzz"
	defaultInt1 := 3
	defaultInt2 := 5
	defaultLimit := 100

	defaultFizzBuzzInput.Str1 = &defaultStr1
	defaultFizzBuzzInput.Str2 = &defaultStr2
	defaultFizzBuzzInput.Int1 = &defaultInt1
	defaultFizzBuzzInput.Int2 = &defaultInt2
	defaultFizzBuzzInput.Limit = &defaultLimit

	// setup max FizzBuzzInput.Limit value
	if envLimitStr := os.Getenv(FizzBuzzEnvLimit); envLimitStr != "" {
		envLimit, err := strconv.Atoi(envLimitStr)
		if err != nil {
			log.Error(
				"failed to load custom fizzbuzz limit from env",
				err.Error(),
			)
			return
		}

		FizzBuzzMaxLimit = envLimit
	}
}

// FizzBuzzInput describes the expected input for the fizzbuzz handler.
type FizzBuzzInput struct {
	Str1  *string `query:"str1" validate:"required"`
	Str2  *string `query:"str2" validate:"required"`
	Int1  *int    `query:"int1" validate:"required,min=1"`
	Int2  *int    `query:"int2" validate:"required,min=1"`
	Limit *int    `query:"limit" validate:"required,min=0"`
}

// SetDefault converts non-provided inputs to fizzbuzz's algorithm
// default values.
func (in *FizzBuzzInput) SetDefault() {
	if in.Str1 == nil {
		in.Str1 = defaultFizzBuzzInput.Str1
	}
	if in.Str2 == nil {
		in.Str2 = defaultFizzBuzzInput.Str2
	}
	if in.Int1 == nil {
		in.Int1 = defaultFizzBuzzInput.Int1
	}
	if in.Int2 == nil {
		in.Int2 = defaultFizzBuzzInput.Int2
	}
	if in.Limit == nil {
		in.Limit = defaultFizzBuzzInput.Limit
	}
}

// Register increments the input parameters in fizzbuzz statistics.
//
// It assumes that SetDefault method was called on the FizzBuzzInput instance
// so that all values are non-nil.
func (in FizzBuzzInput) Register() {
	fizzBuzzGatherer.Hit(
		fmt.Sprintf("FizzBuzzInput str1=%s str2=%s int1=%d int2=%d limit=%d",
			*in.Str1, *in.Str2, *in.Int1, *in.Int2, *in.Limit),
	)
}

// FizzBuzzOutput describes the response output for the fizzbuzz handler.
type FizzBuzzOutput struct {
	Result []string `json:"result"`
}

// FizzBuzz responds to GET /fizbuzz HTTP requests.
//
// It will respond with a 200 HTTP repsonse embedding
// a FizzBuzzOutput result.
//
// The result is computed using the following algorithm:
//  - Return a list from 1 to `limit`
//  - Each multiple of int1 is replaced by str1
//  - Each multiple of int2 is replaced by str2
//  - Each multiple of both is replaced by str1+str2
//
// @Summary Customizable fizzbuzz algorithm.
// @Description Get your own version of the fizzbuzz algortihm.
// @Tags fizzbuzz
// @Accept */*
// @Param int1  query int    false "fizzbuzz's first multiple"     minimum(1) default(3)
// @Param int2  query int    false "fizzbuzz's second multiple"    minimum(1) default(5)
// @Param str1  query string false "fizzbuzz's first replacement"             default(fizz)
// @Param str2  query string false "fizzbuzz's second replacement"            default(buzz)
// @Param limit query int    false "fizzbuzz's up-to value"        minimum(0) default(100)
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

	in.SetDefault()

	err = c.Validate(&in)
	if err != nil {
		c.Logger().Warnf("failed to validate query parameters: %v", err)
		return err
	}

	if *in.Limit > FizzBuzzMaxLimit {
		c.Logger().Warnf("limit %d is higher than threshold %d", in.Limit, FizzBuzzMaxLimit)
		return echo.NewHTTPError(
			http.StatusBadRequest,
			fmt.Sprintf("limit should be lower than %d", FizzBuzzMaxLimit),
		)
	}

	slice := make([]string, *in.Limit)
	int1, int2, str1, str2 := *in.Int1, *in.Int2, *in.Str1, *in.Str2
	int3 := lcm(int1, int2)
	str3 := str1 + str2

	var v int // avoid v's allocation at every loop
	for i := range slice {
		v = i + 1 // array shall start with value 1

		switch {
		case v%int3 == 0:
			slice[i] = str3
		case v%int1 == 0:
			slice[i] = str1
		case v%int2 == 0:
			slice[i] = str2
		default:
			// strconv is more efficient than fmt.Sprint
			slice[i] = strconv.Itoa(v)
		}
	}

	// inputs are valid, add this request to fizzbuzz's stats
	go in.Register()

	return c.JSON(http.StatusOK, FizzBuzzOutput{Result: slice})
}

// gcd computes the Greatest Common Divisor (GCD) via Euclidean algorithm.
func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

// lcm computes the Least Common Multiple (LCM) via GCD.
func lcm(a, b int) int {
	return a * b / gcd(a, b)
}
