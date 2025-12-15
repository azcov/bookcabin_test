package api

import (
	"net/http"

	"github.com/azcov/bookcabin_test/pkg/httpz"
	"github.com/gin-gonic/gin"

	"github.com/azcov/bookcabin_test/internal/domain"
	"github.com/azcov/bookcabin_test/internal/service"
)

// Handler groups dependencies for API handlers
type Handler struct {
	FlightSvc service.FlightInterface
}

// NewHandler returns a new API handler instance
func NewHandler(fsvc service.FlightInterface) *Handler {
	return &Handler{FlightSvc: fsvc}
}

// SearchFlights handles POST /v1/flights/search
func (h *Handler) SearchFlights(c *gin.Context) {
	var req domain.SearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		eresp := httpz.NewErrorResponse(http.StatusBadRequest, "invalid_request", err.Error(), nil)
		httpz.JSONResponse(c, nil, eresp)
		return
	}

	// call service
	resp, err := h.FlightSvc.SerchFlight(c.Request.Context(), &req)
	if err != nil {
		httpz.JSONResponse(c, nil, err)
		return
	}

	httpz.JSONResponse(c, resp, nil)
}
