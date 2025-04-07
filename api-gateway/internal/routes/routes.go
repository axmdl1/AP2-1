package routes

import (
	"github.com/gin-gonic/gin"
	"net/http/httputil"
	"net/url"
)

func newReverseProxy(targetHost string) *httputil.ReverseProxy {
	targetURL, err := url.Parse(targetHost)
	if err != nil {
		panic(err)
	}
	return httputil.NewSingleHostReverseProxy(targetURL)
}

var (
	inventoryProxy = newReverseProxy("http://localhost:1001")
	orderProxy     = newReverseProxy("http://localhost:1002")
)

func SetupRoutes(router *gin.Engine) {
	router.Any("/products/*any", func(c *gin.Context) {
		inventoryProxy.ServeHTTP(c.Writer, c.Request)
	})

	router.Any("/orders/*any", func(c *gin.Context) {
		orderProxy.ServeHTTP(c.Writer, c.Request)
	})
}
