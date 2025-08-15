package rest

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/luckyAkbar/atec/internal/usecase"
)

// @Summary		Get my profile data
// @Description	Return the profile of the currently authenticated user
// @Tags			Users
// @Accept			json
// @Produce		json
// @Security		ParentLevelAuth
// @Param			Authorization	header		string												true	"JWT Token"
// @Success		200				{object}	StandardSuccessResponse{data=GetMyProfileOutput}	"Successful response"
// @Failure		401				{object}	StandardErrorResponse								"Unauthorized"
// @Failure		404				{object}	StandardErrorResponse								"Not Found"
// @Failure		500				{object}	StandardErrorResponse								"Internal Error"
// @Router			/v1/users/me [get]
func (s *Service) HandleGetMyProfile() echo.HandlerFunc {
	return func(c echo.Context) error {
		output, err := s.usersUsecase.GetMyProfile(c.Request().Context())
		if err != nil {
			return UsecaseErrorToRESTResponse(c, err)
		}

		return c.JSON(http.StatusOK, StandardSuccessResponse{
			StatusCode: http.StatusOK,
			Message:    http.StatusText(http.StatusOK),
			Data: GetMyProfileOutput{
				ID:          output.ID,
				Username:    output.Username,
				IsActive:    output.IsActive,
				Roles:       output.Roles,
				CreatedAt:   output.CreatedAt,
				UpdatedAt:   output.UpdatedAt,
				Email:       output.Email,
				PhoneNumber: output.PhoneNumber,
				Address:     output.Address,
			},
		})
	}
}

// @Summary		Get all therapists
// @Description	Return the list of all users with therapist role
// @Tags			Users
// @Accept			json
// @Produce		json
// @Security		ParentLevelAuth
// @Param			Authorization	header		string												true	"JWT Token"
// @Success		200				{object}	StandardSuccessResponse{data=[]GetTherapistOutput}	"Successful response"
// @Failure		401				{object}	StandardErrorResponse								"Unauthorized"
// @Failure		404				{object}	StandardErrorResponse								"Not Found"
// @Failure		500				{object}	StandardErrorResponse								"Internal Error"
// @Router			/v1/users/therapists [get]
func (s *Service) HandleGetTherapists() echo.HandlerFunc {
	return func(c echo.Context) error {
		output, err := s.usersUsecase.GetTherapistData(c.Request().Context())
		if err != nil {
			return UsecaseErrorToRESTResponse(c, err)
		}

		resp := make([]GetTherapistOutput, 0, len(output))
		for _, therapist := range output {
			resp = append(resp, GetTherapistOutput{
				ID:        therapist.ID,
				Username:  therapist.Username,
				IsActive:  therapist.IsActive,
				Roles:     therapist.Roles,
				CreatedAt: therapist.CreatedAt,
				UpdatedAt: therapist.UpdatedAt,
			})
		}

		return c.JSON(http.StatusOK, StandardSuccessResponse{
			StatusCode: http.StatusOK,
			Message:    http.StatusText(http.StatusOK),
			Data:       resp,
		})
	}
}

// @Summary		Update my profile data
// @Description	Allow authenticated user to update their username, phone number, and address
// @Tags			Users
// @Accept			json
// @Produce		json
// @Security		ParentLevelAuth
// @Param			Authorization			header		string												true	"JWT Token"
// @Param			update_profile_input	body		UpdateMyProfileInput								true	"Update profile input"
// @Success		200						{object}	StandardSuccessResponse{data=UpdateMyProfileOutput}	"Successful response"
// @Failure		400						{object}	StandardErrorResponse								"Bad Request"
// @Failure		401						{object}	StandardErrorResponse								"Unauthorized"
// @Failure		500						{object}	StandardErrorResponse								"Internal Error"
// @Router			/v1/users/me [patch]
func (s *Service) HandleUpdateMyProfile() echo.HandlerFunc {
	return func(c echo.Context) error {
		input := &UpdateMyProfileInput{}
		if err := c.Bind(input); err != nil {
			return c.JSON(http.StatusBadRequest, StandardErrorResponse{
				StatusCode:   http.StatusBadRequest,
				ErrorMessage: "failed to parse input",
				ErrorCode:    http.StatusText(http.StatusBadRequest),
			})
		}

		output, err := s.usersUsecase.UpdateMyProfile(c.Request().Context(), usecase.UpdateMyProfileInput{
			Username:    input.Username,
			PhoneNumber: input.PhoneNumber,
			Address:     input.Address,
		})
		if err != nil {
			return UsecaseErrorToRESTResponse(c, err)
		}

		return c.JSON(http.StatusOK, StandardSuccessResponse{
			StatusCode: http.StatusOK,
			Message:    http.StatusText(http.StatusOK),
			Data: UpdateMyProfileOutput{
				Message: output.Message,
			},
		})
	}
}
