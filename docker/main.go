package main

import (
	"github.com/labstack/echo/v4"

	"fmt"
	"net/http"
	"strconv"
)

var (
	maxPrime = 2
)

func main() {
	e := echo.New()
	e.Server.SetKeepAlivesEnabled(false)

	go func() {
		for i := 3; ; i++ {
			for j := i - 1; j > 1; j-- {
				if i%j == 0 {
					break
				}
				if j == 2 {
					maxPrime = i
					fmt.Println(maxPrime)
				}
			}
		}
	}()

	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	count := 0
	e.GET("/count", func(c echo.Context) error {
		count++
		return c.String(http.StatusOK, strconv.Itoa(count)+"\n")
	})

	e.GET("/maxPrime", func(c echo.Context) error {
		return c.String(http.StatusOK, strconv.Itoa(maxPrime)+"\n")
	})

	e.Start(":3000")
}
