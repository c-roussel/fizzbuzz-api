package handlers_test

import (
	"fmt"
	"math"
	"net/http"
	"strings"
	"testing"

	"github.com/c-roussel/fizzbuzz-api/internal/handlers"
	"github.com/c-roussel/fizzbuzz-api/internal/server"
	"github.com/maxatome/go-testdeep/helpers/tdhttp"
	"github.com/maxatome/go-testdeep/td"
)

func TestFizzBuzz(t *testing.T) {
	testAPI := tdhttp.NewTestAPI(t, server.New())

	testAPI.Name("valid parameters").
		Get("/fizzbuzz?str1=le&str2=boncoin&limit=10&int1=2&int2=3").
		CmpStatus(http.StatusOK).
		CmpJSONBody(td.JSON(`{"result": ["1", "le", "boncoin", "le", "5", "leboncoin", "7", "le", "boncoin", "le"]}`))
}

func TestFizzBuzzInvalidQuery(t *testing.T) {
	testAPI := tdhttp.NewTestAPI(t, server.New())

	requiredParams := map[string]string{
		"str1":  "str1",
		"str2":  "str2",
		"int1":  "2",
		"int2":  "3",
		"limit": "10",
	}
	for missingKey := range requiredParams {
		queryParams := tdhttp.Q{}
		for key, value := range requiredParams {
			if key != missingKey {
				queryParams[key] = value
			}
		}

		missingKey = strings.Title(missingKey)
		testAPI.Name("missing query parameter", missingKey).
			Get("/fizzbuzz", queryParams).
			CmpStatus(http.StatusBadRequest).
			CmpJSONBody(
				td.JSON(`{"message": $1}`,
					fmt.Sprintf(
						"Key: 'FizzBuzzInput.%s' Error:Field validation for '%s' failed on the 'required' tag",
						missingKey,
						missingKey,
					),
				),
			)
	}

	testCases := []struct {
		name           string
		url            string
		expectedStatus int
		expectedJSON   string
	}{
		{
			name:           "invalid int1 query param",
			url:            "/fizzbuzz?str1=toto&str2=tata&limit=10&int1=-1&int2=3",
			expectedStatus: http.StatusBadRequest,
			expectedJSON:   `{"message": "Key: 'FizzBuzzInput.Int1' Error:Field validation for 'Int1' failed on the 'min' tag"}`,
		},
		{
			name:           "invalid int2 query param",
			url:            "/fizzbuzz?str1=toto&str2=tata&limit=10&int1=1&int2=-3",
			expectedStatus: http.StatusBadRequest,
			expectedJSON:   `{"message": "Key: 'FizzBuzzInput.Int2' Error:Field validation for 'Int2' failed on the 'min' tag"}`,
		},
		{
			name:           "invalid limit query param",
			url:            "/fizzbuzz?str1=toto&str2=tata&limit=-10&int1=1&int2=3",
			expectedStatus: http.StatusBadRequest,
			expectedJSON:   `{"message": "Key: 'FizzBuzzInput.Limit' Error:Field validation for 'Limit' failed on the 'min' tag"}`,
		},
		{
			name:           "invalid limit query param - should be integer",
			url:            "/fizzbuzz?str1=toto&str2=tata&limit=string&int1=1&int2=3",
			expectedStatus: http.StatusBadRequest,
			expectedJSON:   `{"message": "strconv.ParseInt: parsing \"string\": invalid syntax"}`,
		},
		{
			name:           "invalid limit query param - should be lower than threshold",
			url:            "/fizzbuzz?str1=toto&str2=tata&limit=100000000&int1=1&int2=3",
			expectedStatus: http.StatusBadRequest,
			expectedJSON:   `{"message": "limit should be lower than 10000"}`,
		},
	}
	for _, tc := range testCases {
		testAPI.Run(tc.name, func(ta *tdhttp.TestAPI) {
			ta.Get(tc.url).
				CmpStatus(tc.expectedStatus).
				CmpJSONBody(td.JSON(tc.expectedJSON))
		})
	}
}

func BenchmarkFizzBuzz(b *testing.B) {
	defer func(old int) { handlers.FizzBuzzMaxLimit = old }(handlers.FizzBuzzMaxLimit)
	handlers.FizzBuzzMaxLimit = math.MaxInt

	testAPI := tdhttp.NewTestAPI(b, server.New())

	b.ResetTimer()
	testAPI.Name("benchmark", b.N).
		Get(fmt.Sprintf("/fizzbuzz?str1=le&str2=boncoin&limit=%d&int1=7&int2=31", b.N)).
		CmpStatus(http.StatusOK)
	b.StopTimer()
}
