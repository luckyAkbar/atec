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
				ErrorMessage: "internal server error",
			})
		case usecase.ErrBadRequest:
			return c.JSON(http.StatusBadRequest, StandardErrorResponse{
				StatusCode:   http.StatusBadRequest,
				ErrorMessage: err.Error(),
				ErrorCode:    http.StatusText(http.StatusBadRequest),
			})
		}
	}
}
