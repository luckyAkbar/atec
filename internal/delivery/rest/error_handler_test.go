package rest_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/luckyAkbar/atec/internal/delivery/rest"
	"github.com/luckyAkbar/atec/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRESTService_UsecaseErrorToRESTResponse(t *testing.T) {
	e := echo.New()

	t.Run("any error should always returning 500", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()
		ectx := e.NewContext(req, rec)
		err := rest.UsecaseErrorToRESTResponse(ectx, assert.AnError)

		require.NoError(t, err)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		assert.Equal(t, "application/json", rec.Header().Get(echo.HeaderContentType))

		assert.JSONEq(t, fmt.Sprintf(`{"status_code":500,"error_code":"Internal Server Error","error_message":"%s"}`, rest.DefaultUknownErrorMessage), rec.Body.String())
	})

	t.Run("usecase returning unknown error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()
		ectx := e.NewContext(req, rec)
		ucErr := usecase.UsecaseError{
			ErrType: assert.AnError,
			Message: "should be passed down to user",
		}
		err := rest.UsecaseErrorToRESTResponse(ectx, ucErr)

		require.NoError(t, err)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		assert.Equal(t, "application/json", rec.Header().Get(echo.HeaderContentType))

		assert.JSONEq(t, fmt.Sprintf(`{"status_code":500,"error_code":"Internal Server Error","error_message":"%s"}`, rest.DefaultInternalErrorMessage), rec.Body.String())
	})

	t.Run("usecase internal error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()
		ectx := e.NewContext(req, rec)
		ucErr := usecase.UsecaseError{
			ErrType: usecase.ErrInternal,
			Message: "should be passed down to user",
		}
		err := rest.UsecaseErrorToRESTResponse(ectx, ucErr)

		require.NoError(t, err)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		assert.Equal(t, "application/json", rec.Header().Get(echo.HeaderContentType))

		assert.JSONEq(t, fmt.Sprintf(`{"status_code":500,"error_code":"Internal Server Error","error_message":"%s"}`, ucErr.Message), rec.Body.String())
	})

	t.Run("usecase bad request error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()
		ectx := e.NewContext(req, rec)
		ucErr := usecase.UsecaseError{
			ErrType: usecase.ErrBadRequest,
			Message: "should be passed down to user",
		}
		err := rest.UsecaseErrorToRESTResponse(ectx, ucErr)

		require.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		assert.Equal(t, "application/json", rec.Header().Get(echo.HeaderContentType))

		assert.JSONEq(t, fmt.Sprintf(`{"status_code":400,"error_code":"%s","error_message":"%s"}`, http.StatusText(http.StatusBadRequest), ucErr.Message), rec.Body.String())
	})

	t.Run("usecase not found error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()
		ectx := e.NewContext(req, rec)
		ucErr := usecase.UsecaseError{
			ErrType: usecase.ErrNotFound,
			Message: "should be passed down to user",
		}
		err := rest.UsecaseErrorToRESTResponse(ectx, ucErr)

		require.NoError(t, err)

		assert.Equal(t, http.StatusNotFound, rec.Code)

		assert.Equal(t, "application/json", rec.Header().Get(echo.HeaderContentType))

		assert.JSONEq(t, fmt.Sprintf(`{"status_code":404,"error_code":"%s","error_message":"%s"}`, http.StatusText(http.StatusNotFound), ucErr.Message), rec.Body.String())
	})

	t.Run("usecase unauthorized error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()
		ectx := e.NewContext(req, rec)
		ucErr := usecase.UsecaseError{
			ErrType: usecase.ErrUnauthorized,
			Message: "should be passed down to user",
		}
		err := rest.UsecaseErrorToRESTResponse(ectx, ucErr)

		require.NoError(t, err)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)

		assert.Equal(t, "application/json", rec.Header().Get(echo.HeaderContentType))

		assert.JSONEq(t, fmt.Sprintf(`{"status_code":401,"error_code":"%s","error_message":"%s"}`, http.StatusText(http.StatusUnauthorized), ucErr.Message), rec.Body.String())
	})

	t.Run("usecase forbidden error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()
		ectx := e.NewContext(req, rec)
		ucErr := usecase.UsecaseError{
			ErrType: usecase.ErrForbidden,
			Message: "should be passed down to user",
		}
		err := rest.UsecaseErrorToRESTResponse(ectx, ucErr)

		require.NoError(t, err)

		assert.Equal(t, http.StatusForbidden, rec.Code)

		assert.Equal(t, "application/json", rec.Header().Get(echo.HeaderContentType))

		assert.JSONEq(t, fmt.Sprintf(`{"status_code":403,"error_code":"%s","error_message":"%s"}`, http.StatusText(http.StatusForbidden), ucErr.Message), rec.Body.String())
	})

	t.Run("usecase too many request error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()
		ectx := e.NewContext(req, rec)
		ucErr := usecase.UsecaseError{
			ErrType: usecase.ErrTooManyRequests,
			Message: "should be passed down to user",
		}
		err := rest.UsecaseErrorToRESTResponse(ectx, ucErr)

		require.NoError(t, err)

		assert.Equal(t, http.StatusTooManyRequests, rec.Code)

		assert.Equal(t, "application/json", rec.Header().Get(echo.HeaderContentType))

		assert.JSONEq(t, fmt.Sprintf(`{"status_code":429,"error_code":"%s","error_message":"%s"}`, http.StatusText(http.StatusTooManyRequests), ucErr.Message), rec.Body.String())
	})
}
