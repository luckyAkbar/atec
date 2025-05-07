package rest

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/luckyAkbar/atec/internal/model"
	"github.com/luckyAkbar/atec/internal/usecase"
)

// AuthMiddleware will perform authentication check based on the user access token from Authorization header.
// If the token is valid, it will set the user information in the context and call the next handler.
// If the token is invalid or missing, it will return an error response when allowUnauthorized is false.
// Otherwise, it will call the next handler without authentication check.
func (s *service) AuthMiddleware(allowUnauthorized bool) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := getAccessToken(c.Request())
			if token == "" {
				if allowUnauthorized {
					return next(c)
				}

				return c.JSON(http.StatusUnauthorized, StandardErrorResponse{
					StatusCode:   http.StatusUnauthorized,
					ErrorMessage: "missing required auth token",
					ErrorCode:    http.StatusText(http.StatusUnauthorized),
				})
			}

			output, err := s.authUsecase.AuthenticateAccessToken(c.Request().Context(), usecase.AuthenticateAccessTokenInput{
				Token: token,
			})

			if err != nil {
				return UsecaseErrorToRESTResponse(c, err)
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
