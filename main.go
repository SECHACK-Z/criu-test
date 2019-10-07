package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	// "net/http"
	"log"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"time"
)

func main() {
	proxy := echo.New()
	proxy.Server.SetKeepAlivesEnabled(false)

	// e.GET("/ping", func(c echo.Context) error {
	// 	return c.String(http.StatusOK, "pong")
	// })

	// proxy := e.Group("/")
	url, _ := url.Parse("http://localhost:3000")

	cmd := exec.Command("docker", "rm", "-f", "target")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Println(err)
	}
	cmd = exec.Command("docker", "run", "-d", "-p", "3000:3000", "--name", "target", "target")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Println(err)
	}

	timer := time.NewTimer(10 * time.Minute)
	alive := true
	count := 0
	go func() {
		for {
			select {
			case <-timer.C:
				count++
				cmd = exec.Command("docker", "checkpoint", "create", "--leave-running=true", "target", "ch"+strconv.Itoa(count))
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr

				if err := cmd.Run(); err != nil {
					log.Println(err)
				}
				cmd = exec.Command("docker", "stop", "target")
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr

				if err := cmd.Run(); err != nil {
					log.Println(err)
				}
				alive = false
			}
		}

	}()

	proxy.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if alive {
				timer.Reset(10 * time.Second)
			} else {
				cmd = exec.Command("docker", "start", "--checkpoint", "ch"+strconv.Itoa(count), "target")
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr

				if err := cmd.Run(); err != nil {
					log.Println(err)
				}
				alive = true
				timer.Reset(10 * time.Second)
			}
			return next(c)
		}
	})
	proxy.Use(middleware.Proxy(middleware.NewRoundRobinBalancer(
		[]*middleware.ProxyTarget{
			{
				URL: url,
			},
		})))

	proxy.Start(":8888")
}
