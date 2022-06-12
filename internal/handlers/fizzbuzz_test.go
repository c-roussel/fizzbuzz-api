package handlers_test

import (
	"fmt"
	"math"
	"net/http"
	"net/url"
	"testing"

	"github.com/c-roussel/fizzbuzz-api/internal/handlers"
	"github.com/c-roussel/fizzbuzz-api/internal/server"
	"github.com/maxatome/go-testdeep/helpers/tdhttp"
	"github.com/maxatome/go-testdeep/td"
)

func TestFizzBuzz(t *testing.T) {
	testAPI := tdhttp.NewTestAPI(t, server.New())

	testCases := []struct {
		name           string
		url            string
		expectedStatus int
		expectedJSON   string
	}{
		{
			name:           "valid call with parameters",
			url:            "/fizzbuzz?str1=le&str2=boncoin&limit=10&int1=2&int2=3",
			expectedStatus: http.StatusOK,
			expectedJSON:   `{"result": ["1", "le", "boncoin", "le", "5", "leboncoin", "7", "le", "boncoin", "le"]}`,
		}, {
			name:           "valid call without parameters",
			url:            "/fizzbuzz",
			expectedStatus: http.StatusOK,
			expectedJSON: `{"result": [
			"1","2","fizz","4","buzz","fizz","7","8","fizz","buzz","11","fizz","13","14","fizzbuzz","16","17","fizz","19",
			"buzz","fizz","22","23","fizz","buzz","26","fizz","28","29","fizzbuzz","31","32","fizz","34","buzz","fizz","37",
			"38","fizz","buzz","41","fizz","43","44","fizzbuzz","46","47","fizz","49","buzz","fizz","52","53","fizz","buzz",
			"56","fizz","58","59","fizzbuzz","61","62","fizz","64","buzz","fizz","67","68","fizz","buzz","71","fizz","73",
			"74","fizzbuzz","76","77","fizz","79","buzz","fizz","82","83","fizz","buzz","86","fizz","88","89","fizzbuzz",
			"91","92","fizz","94","buzz","fizz","97","98","fizz","buzz"
			]}`,
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

func TestFizzBuzzInvalidQuery(t *testing.T) {
	testAPI := tdhttp.NewTestAPI(t, server.New())

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

func FuzzFizzBuzz(f *testing.F) {
	f.Add(1, 2, 3, "str1", "str2")
	f.Fuzz(func(t *testing.T, int1, int2, limit int, str1, str2 string) {
		expectedStatus := http.StatusOK
		if int1 < 1 || int2 < 1 || limit < 0 {
			expectedStatus = http.StatusBadRequest
		}

		testAPI := tdhttp.NewTestAPI(t, server.New())
		testAPI.Get(
			"/fizzbuzz",
			tdhttp.Q{
				"str1":  url.QueryEscape(str1),
				"str2":  url.QueryEscape(str2),
				"int1":  int1,
				"int2":  int2,
				"limit": limit,
			}).
			CmpStatus(expectedStatus)
	})
}
