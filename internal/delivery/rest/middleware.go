package rest

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/luckyAkbar/atec/internal/model"
	"github.com/luckyAkbar/atec/internal/usecase"
)

func (s *service) AuthMiddleware(allowAllAuthorized, allowAdminOnly bool) echo.MiddlewareFunc {
	if !allowAdminOnly && !allowAllAuthorized {
		panic("invalid configuration for auth middleware. one of the params must be true")
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := getAccessToken(c.Request())
			if token == "" {
				return c.JSON(http.StatusUnauthorized, StandardErrorResponse{
					StatusCode:   http.StatusUnauthorized,
					ErrorMessage: "missing required auth token",
					ErrorCode:    http.StatusText(http.StatusUnauthorized),
				})
			}

			output, err := s.authUsecase.AllowAccess(c.Request().Context(), usecase.AllowAccessInput{
				Token:              token,
				AllowAllAuthorized: allowAllAuthorized,
				AllowAdminOnly:     allowAdminOnly,
			})

			if err != nil {
				return usecaseErrorToRESTResponse(c, err)
			}

			authUser := model.AuthUser{
				ID:   output.UserID,
				Role: output.UserRole,
			}

			ctx := c.Request().Context()
			newCtx := model.SetUserToCtx(ctx, authUser)

			c.SetRequest(c.Request().WithContext(newCtx))

			return next(c)
		}
	}
}

func getAccessToken(req *http.Request) string {
	authHeaders := strings.Split(req.Header.Get("Authorization"), " ")

	if len(authHeaders) != 1 {
		return ""
	}

	return strings.TrimSpace(authHeaders[0])
}
