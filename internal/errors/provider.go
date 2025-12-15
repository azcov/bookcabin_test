package errors

import (
	"net/http"

	"github.com/azcov/bookcabin_test/pkg/errorz"
)

var (
	ErrNotFound                         = &errorz.WrappedError{StatusCode: http.StatusNotFound, ErrCode: "not_found", Msg: "Resource not found"}
	ErrLionAirNotFound                  = &errorz.WrappedError{StatusCode: http.StatusNotFound, ErrCode: "not_found", Msg: "Lion Air not found"}
	ErrLionAirInternalError             = &errorz.WrappedError{StatusCode: http.StatusInternalServerError, ErrCode: "internal_error", Msg: "Lion Air internal error"}
	ErrLionAirRateLimitExceeded         = &errorz.WrappedError{StatusCode: http.StatusTooManyRequests, ErrCode: "rate_limit_exceeded", Msg: "Lion Air rate limit exceeded"}
	ErrAirAsiaNotFound                  = &errorz.WrappedError{StatusCode: http.StatusNotFound, ErrCode: "not_found", Msg: "Air Asia not found"}
	ErrAirAsiaInternalError             = &errorz.WrappedError{StatusCode: http.StatusInternalServerError, ErrCode: "internal_error", Msg: "Air Asia internal error"}
	ErrAirAsiaRateLimitExceeded         = &errorz.WrappedError{StatusCode: http.StatusTooManyRequests, ErrCode: "rate_limit_exceeded", Msg: "Air Asia rate limit exceeded"}
	ErrBatikAirNotFound                 = &errorz.WrappedError{StatusCode: http.StatusNotFound, ErrCode: "not_found", Msg: "Batik Air not found"}
	ErrBatikAirRateLimitExceeded        = &errorz.WrappedError{StatusCode: http.StatusTooManyRequests, ErrCode: "rate_limit_exceeded", Msg: "Batik Air rate limit exceeded"}
	ErrGarudaIndonesiaNotFound          = &errorz.WrappedError{StatusCode: http.StatusNotFound, ErrCode: "not_found", Msg: "Garuda Indonesia not found"}
	ErrGarudaIndonesiaRateLimitExceeded = &errorz.WrappedError{StatusCode: http.StatusTooManyRequests, ErrCode: "rate_limit_exceeded", Msg: "Garuda Indonesia rate limit exceeded"}
)
