package rest_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/luckyAkbar/atec/internal/delivery/rest"
	"github.com/luckyAkbar/atec/internal/model"
	"github.com/luckyAkbar/atec/internal/usecase"
	usecase_mock "github.com/luckyAkbar/atec/mocks/internal_/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUsersService_HandleGetMyProfile(t *testing.T) {
	e := echo.New()
	group := e.Group("")

	mockUsersUC := usecase_mock.NewUsersUsecaseIface(t)

	svc := rest.NewService(group, nil, nil, nil, nil, mockUsersUC)

	t.Run("unauthorized mapped from usecase", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/v1/me", nil)
		ctx := e.NewContext(req, rec)

		mockUsersUC.EXPECT().GetMyProfile(ctx.Request().Context()).Return(nil, usecase.UsecaseError{ErrType: usecase.ErrUnauthorized}).Once()

		err := svc.HandleGetMyProfile()(ctx)
		require.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("success", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/v1/me", nil)
		ctx := e.NewContext(req, rec)

		output := &usecase.GetMyProfileOutput{
			ID:       uuid.New(),
			Username: "user",
			IsActive: true,
			Roles:    model.RolesParent,
		}
		mockUsersUC.EXPECT().GetMyProfile(ctx.Request().Context()).Return(output, nil).Once()

		err := svc.HandleGetMyProfile()(ctx)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "application/json", rec.Header().Get(echo.HeaderContentType))
	})
}
