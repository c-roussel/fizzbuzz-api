package main

import (
	_ "github.com/c-roussel/fizzbuzz-api/docs/swagger"
	"github.com/c-roussel/fizzbuzz-api/internal/server"
)

// @title FizzBuzz API
// @version 1.0
// @description This is a custom FizzBuzz HTTP server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath /
// @schemes http
func main() {
	e := server.New()

	e.Logger.Info("Starting fizzbuzz-api server")
	e.Logger.Fatal(e.Start(":3000"))
}
