package handlers_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/c-roussel/fizzbuzz-api/internal/handlers"
	"github.com/c-roussel/fizzbuzz-api/internal/server"
	"github.com/maxatome/go-testdeep/helpers/tdhttp"
	"github.com/maxatome/go-testdeep/td"
)

func TestFizzBuzzStats(t *testing.T) {
	handlers.ExportFizzBuzzGatherer.Reset()

	testAPI := tdhttp.NewTestAPI(t, server.New())

	for idx, params := range []string{
		"str1=le&str2=boncoin&limit=6&int1=2&int2=3",
		"str1=l&str2=bc&limit=6&int1=2&int2=3",
	} {
		for i := 0; i < idx+1; i++ {
			testAPI.Name("/fizzbuzz stat population", params, i).
				Get("/fizzbuzz?" + params).
				CmpStatus(http.StatusOK)
		}
	}

	// gathering is done asynchronously
	time.Sleep(100 * time.Millisecond)

	testAPI.Name("/fizzbuzz stat retrieval").
		Get("/fizzbuzz/stats").
		CmpStatus(http.StatusOK).
		CmpJSONBody(td.JSON(`
[{
  "key": "FizzBuzzInput str1=l str2=bc int1=2 int2=3 limit=6",
  "hit": 2
},{
  "key": "FizzBuzzInput str1=le str2=boncoin int1=2 int2=3 limit=6",
  "hit": 1
}]`))

	for i := 0; i < 100; i++ {
		testAPI.Name("/fizzbuzz stat population up to 102", i).
			Get(fmt.Sprintf("/fizzbuzz?str1=l&str2=bc&limit=6&int1=2&int2=10%d", i)).
			CmpStatus(http.StatusOK)
	}

	// gathering is done asynchronously
	time.Sleep(100 * time.Millisecond)

	testAPI.Name("/fizzbuzz stat retrieval max result number is 100").
		Get("/fizzbuzz/stats").
		CmpStatus(http.StatusOK).
		CmpJSONBody(td.JSON("Len(100)"))
}
