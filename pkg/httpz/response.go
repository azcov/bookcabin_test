package httpz

import (
	"net/http"

	"github.com/azcov/bookcabin_test/pkg/errorz"
	"github.com/gin-gonic/gin"
)

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
	Data    any    `json:"data,omitempty"`
	Detail  any    `json:"detail,omitempty"`
}

func NewResponse(status, message string, data any) *Response {
	return &Response{
		Status:  status,
		Message: message,
		Data:    data,
	}
}

type ErrorResponse struct {
	Code    int    `json:"-"`
	Err     string `json:"error"`
	Message string `json:"message,omitempty"`
	Detail  any    `json:"detail,omitempty"`
}

func NewErrorResponse(code int, err, message string, detail any) *ErrorResponse {
	return &ErrorResponse{
		Code:    code,
		Err:     err,
		Message: message,
		Detail:  detail,
	}
}

// Error implements the error interface so ErrorResponse can be passed as error
func (e *ErrorResponse) Error() string {
	return e.Message
}

func JSONResponse(c *gin.Context, res any, err error) {
	if err != nil {
		errorResponse(c, err)
		return
	}
	c.JSON(http.StatusOK, res)
}

// OK sends a usual success payload with json wrapper
func OK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, NewResponse("ok", "", data))
}

// Created sends an object with status created
func Created(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, NewResponse("created", "", data))
}

func errorResponse(c *gin.Context, err error) {
	switch e := err.(type) {
	case *ErrorResponse:
		c.JSON(e.Code, e)
		return
	case *errorz.WrappedError:
		status := e.StatusCode
		if status == 0 {
			status = http.StatusInternalServerError
		}
		c.JSON(status, gin.H{"error": e.ErrCode, "message": e.Msg, "detail": e.Detail})
		return
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_error", "message": err.Error()})
	}
}
