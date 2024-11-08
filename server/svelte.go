package server

import (
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func svelte(e *echo.Echo) {
	g := e.Group("/app")
	svelteDevUrl, err := url.Parse(svelteDevAddress)
	if err != nil {
		panic(err)
	}

	_, err = http.Get(svelteDevAddress)
	if err == nil {
		g.Use(middleware.Proxy(middleware.NewRoundRobinBalancer([]*middleware.ProxyTarget{{
			URL: svelteDevUrl,
		}})))
		e.Static("/public", "web/public/")
	} else {
		g.Use(middleware.StaticWithConfig(middleware.StaticConfig{
			HTML5: true,
			Root:  "web/build/",
		}))
		e.Static("/public", "web/build/")
	}
}
