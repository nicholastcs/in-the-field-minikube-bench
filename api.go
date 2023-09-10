package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func getGinServer() *http.Server {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	r.GET("/healthz", func(c *gin.Context) {
		if startup && health {
			c.String(200, "OK")
			logger.Trace("returned 200", "path", "/healthz")
			return
		}

		c.String(500, "FAIL")
		logger.Trace("application is unhealthy, returned 500", "path", "/healthz")
	})

	r.GET("/startupz", func(c *gin.Context) {
		if startup {
			c.String(200, "OK")
			logger.Trace("returned 200", "path", "/startup")
			return
		}

		c.String(500, "FAIL")
		logger.Trace("pending application startup, returned 500", "path", "/startup")
	})

	r.GET("/hello", func(c *gin.Context) {
		c.String(200, "hello!")
		logger.Debug("return back to client!")
	})

	r.POST("/hello", func(c *gin.Context) {
		err := createJob(clientset)
		if err != nil {
			logger.Error("unable to create job...")
			c.String(500, err.Error())
			return
		}

		c.String(200, "hello!")
		logger.Debug("return back to client!")
	})

	return &http.Server{
		Addr:        ":8080",
		Handler:     r,
		IdleTimeout: 10 * time.Second,
	}
}
