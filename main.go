package main

import (
	"net/http"
	"os"
	"snowflake-service/snowflake"

	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/labstack/echo/middleware"
)

func main() {
	e := echo.New()

	os.Setenv("MACHINE_ID", "2343")

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Route => handler
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!\n")
	})

	snowflake.RegisterService(e)

	// Start server
	e.Run(standard.New(":1323"))
}
