package main

import (
	"net/http"

	"github.com/labstack/echo"
)

const helloMessage = "Hello World from Echo."

func main() {
	router := NewRouter()

	router.Start(":3000")
}

func NewRouter() *echo.Echo {
	e := echo.New()

	e.GET("/hello", helloHandler)

	return e
}

func helloHandler(c echo.Context) error {
	return c.String(http.StatusOK, helloMessage)
}
