package handlers_test

import (
	"net/http"
	"testing"

	"github.com/c-roussel/fizzbuzz-api/internal/server"
	"github.com/maxatome/go-testdeep/helpers/tdhttp"
	"github.com/maxatome/go-testdeep/td"
)

func TestPing(t *testing.T) {
	testAPI := tdhttp.NewTestAPI(t, server.New())

	testAPI.Name("ping").
		Get("/mon/ping").
		CmpStatus(http.StatusOK).
		CmpJSONBody(td.JSON(`{"message": "OK", "git_hash": ""}`))
}
