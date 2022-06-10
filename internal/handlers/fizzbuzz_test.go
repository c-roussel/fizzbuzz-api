package handlers_test

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

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
		var first bool
		var path strings.Builder
		path.WriteString("/fizzbuzz?")
		for key, value := range requiredParams {
			if key != missingKey {
				if !first {
					path.WriteRune('&')
				}
				fmt.Fprintf(&path, "%s=%s", key, value)
			}
		}
		testAPI.Name("missing query parameter", missingKey).
			Get(path.String()).
			CmpStatus(http.StatusBadRequest).
			CmpJSONBody(
				td.JSON(`{"message": $1}`,
					fmt.Sprintf(
						"Key: 'FizzBuzzInput.%s' Error:Field validation for '%s' failed on the 'required' tag",
						strings.Title(missingKey),
						strings.Title(missingKey),
					),
				),
			)
	}
	testAPI.Name("invalid int1 query param").
		Get("/fizzbuzz?str1=toto&str2=tata&limit=10&int1=-1&int2=3").
		CmpStatus(http.StatusBadRequest).
		CmpJSONBody(td.JSON(`{"message": "Key: 'FizzBuzzInput.Int1' Error:Field validation for 'Int1' failed on the 'min' tag"}`))

	testAPI.Name("invalid int2 query param").
		Get("/fizzbuzz?str1=toto&str2=tata&limit=10&int1=1&int2=-3").
		CmpStatus(http.StatusBadRequest).
		CmpJSONBody(td.JSON(`{"message": "Key: 'FizzBuzzInput.Int2' Error:Field validation for 'Int2' failed on the 'min' tag"}`))

	testAPI.Name("invalid limit query param").
		Get("/fizzbuzz?str1=toto&str2=tata&limit=-10&int1=1&int2=3").
		CmpStatus(http.StatusBadRequest).
		CmpJSONBody(td.JSON(`{"message": "Key: 'FizzBuzzInput.Limit' Error:Field validation for 'Limit' failed on the 'min' tag"}`))

	testAPI.Name("invalid limit query param").
		Get("/fizzbuzz?str1=toto&str2=tata&limit=string&int1=1&int2=3").
		CmpStatus(http.StatusBadRequest).
		CmpJSONBody(td.JSON(`{"message": "strconv.ParseInt: parsing \"string\": invalid syntax"}`))
}

func BenchmarkFizzBuzz(b *testing.B) {
	testAPI := tdhttp.NewTestAPI(b, server.New())

	testAPI.Name("benchmark", b.N).
		Get(fmt.Sprintf("/fizzbuzz?str1=le&str2=boncoin&limit=%d&int1=7&int2=31", b.N)).
		CmpStatus(http.StatusOK)
}
