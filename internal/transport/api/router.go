package api

import (
	"github.com/azcov/bookcabin_test/pkg/httpz"
	"github.com/gin-gonic/gin"
)

// NewRouter creates a gin engine and registers routes for the API.
// Pass a previously created handler to wire the endpoints up.
func NewRouter(h *Handler) *gin.Engine {
	r := gin.New()
	// Add our middlewares: request id, recoverer and logger
	r.Use(gin.Recovery()) // still use gin recovery as a baseline
	r.Use(httpz.RequestID())
	r.Use(httpz.Recovery())
	r.Use(httpz.Logger())

	v1 := r.Group("/v1")
	{
		v1.POST("/flights/search", h.SearchFlights)
		v1.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })
	}

	return r
}
