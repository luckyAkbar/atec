package rest

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/luckyAkbar/atec/internal/usecase"
)

func usecaseErrorToRESTResponse(c echo.Context, err error) error {
	switch e := err.(type) {
	default:
		return c.JSON(http.StatusInternalServerError, StandardErrorResponse{
			StatusCode:   http.StatusInternalServerError,
			ErrorCode:    http.StatusText(http.StatusInternalServerError),
			ErrorMessage: "server received unexpected error",
		})
	case usecase.UsecaseError:
		switch e.ErrType {
		default:
			return c.JSON(http.StatusInternalServerError, StandardErrorResponse{
				StatusCode:   http.StatusInternalServerError,
				ErrorCode:    http.StatusText(http.StatusInternalServerError),
				ErrorMessage: "server unable to produce appropriate response",
			})
		case usecase.ErrInternal:
			return c.JSON(http.StatusInternalServerError, StandardErrorResponse{
				StatusCode:   http.StatusInternalServerError,
				ErrorCode:    http.StatusText(http.StatusInternalServerError),
				ErrorMessage: e.Message,
			})
		case usecase.ErrBadRequest:
			return c.JSON(http.StatusBadRequest, StandardErrorResponse{
				StatusCode:   http.StatusBadRequest,
				ErrorMessage: e.Message,
				ErrorCode:    http.StatusText(http.StatusBadRequest),
			})
		case usecase.ErrNotFound:
			return c.JSON(http.StatusNotFound, StandardErrorResponse{
				StatusCode:   http.StatusNotFound,
				ErrorMessage: e.Message,
				ErrorCode:    http.StatusText(http.StatusNotFound),
			})
		case usecase.ErrUnauthorized:
			return c.JSON(http.StatusUnauthorized, StandardErrorResponse{
				StatusCode:   http.StatusUnauthorized,
				ErrorCode:    http.StatusText(http.StatusUnauthorized),
				ErrorMessage: e.Message,
			})
		case usecase.ErrForbidden:
			return c.JSON(http.StatusForbidden, StandardErrorResponse{
				StatusCode:   http.StatusForbidden,
				ErrorCode:    http.StatusText(http.StatusForbidden),
				ErrorMessage: e.Message,
			})
		case usecase.ErrTooManyRequests:
			return c.JSON(http.StatusTooManyRequests, StandardErrorResponse{
				StatusCode:   http.StatusTooManyRequests,
				ErrorCode:    http.StatusText(http.StatusTooManyRequests),
				ErrorMessage: e.Message,
			})
		}
	}
}
