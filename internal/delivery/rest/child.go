package rest

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/luckyAkbar/atec/internal/usecase"
)

// @Summary		Register a new child
// @Description	Register a new child
// @Tags			Childern
// @Accept			json
// @Produce		json
// @Security		ParentLevelAuth
// @Param			Authorization			header		string												true	"JWT Token"
// @Param			register_child_input	body		RegisterChildInput									true	"child details"
// @Success		200						{object}	StandardSuccessResponse{data=RegisterChildOutput}	"Successful response"
// @Failure		400						{object}	StandardErrorResponse								"Bad request"
// @Failure		500						{object}	StandardErrorResponse								"Internal Error"
// @Router			/v1/childern [post]
func (s *Service) HandleRegisterChildern() echo.HandlerFunc {
	return func(c echo.Context) error {
		input := &RegisterChildInput{}
		if err := c.Bind(input); err != nil {
			return c.JSON(http.StatusBadRequest, StandardErrorResponse{
				StatusCode:   http.StatusBadRequest,
				ErrorMessage: "failed to parse input",
				ErrorCode:    http.StatusText(http.StatusBadRequest),
			})
		}

		dateOfBirth, err := time.Parse("2006-01-02", input.DateOfBirth)
		if err != nil {
			return c.JSON(http.StatusBadRequest, StandardErrorResponse{
				StatusCode:   http.StatusBadRequest,
				ErrorMessage: "invalid time format. should be: 2001-11-29 (YYYY-MM-DD)",
				ErrorCode:    http.StatusText(http.StatusBadRequest),
			})
		}

		output, err := s.childUsecase.Register(c.Request().Context(), usecase.RegisterChildInput{
			DateOfBirth: dateOfBirth,
			Gender:      input.Gender,
			Name:        input.Name,
		})

		if err != nil {
			return UsecaseErrorToRESTResponse(c, err)
		}

		return c.JSON(http.StatusOK, StandardSuccessResponse{
			StatusCode: http.StatusOK,
			Message:    http.StatusText(http.StatusOK),
			Data: RegisterChildOutput{
				ID: output.ID,
			},
		})
	}
}

// @Summary		Update child data
// @Description	Update child data
// @Tags			Childern
// @Accept			json
// @Produce		json
// @Security		ParentLevelAuth
// @Param			Authorization		header		string												true	"JWT Token"
//
// @Param			child_id			path		string												true	"Child ID to be updated (UUID v4)"
//
// @Param			update_child_input	body		UpdateChildernInput									true	"new child details"
// @Success		200					{object}	StandardSuccessResponse{data=UpdateChildernOutput}	"Successful response"
// @Failure		400					{object}	StandardErrorResponse								"Bad request"
// @Failure		500					{object}	StandardErrorResponse								"Internal Error"
// @Router			/v1/childern/{child_id} [put]
func (s *Service) HandleUpdateChildern() echo.HandlerFunc {
	return func(c echo.Context) error {
		input := &UpdateChildernInput{}
		if err := c.Bind(input); err != nil {
			return c.JSON(http.StatusBadRequest, StandardErrorResponse{
				StatusCode:   http.StatusBadRequest,
				ErrorMessage: "failed to parse input",
				ErrorCode:    http.StatusText(http.StatusBadRequest),
			})
		}

		ucUpdateChildInput := usecase.UpdateChildInput{
			ChildID:     input.ChildID,
			DateOfBirth: nil,
			Gender:      input.Gender,
			Name:        input.Name,
		}

		if input.DateOfBirth != "" {
			dateOfBirth, err := time.Parse("2006-01-02", input.DateOfBirth)
			if err != nil {
				return c.JSON(http.StatusBadRequest, StandardErrorResponse{
					StatusCode:   http.StatusBadRequest,
					ErrorMessage: "invalid time format. should be: 2001-11-29 (YYYY-MM-DD)",
					ErrorCode:    http.StatusText(http.StatusBadRequest),
				})
			}

			ucUpdateChildInput.DateOfBirth = &dateOfBirth
		}

		output, err := s.childUsecase.Update(c.Request().Context(), ucUpdateChildInput)

		if err != nil {
			return UsecaseErrorToRESTResponse(c, err)
		}

		return c.JSON(http.StatusOK, StandardSuccessResponse{
			StatusCode: http.StatusOK,
			Message:    http.StatusText(http.StatusOK),
			Data: UpdateChildernOutput{
				Message: output.Message,
			},
		})
	}
}

// @Summary		Get all childern registered under this account
// @Description	Get all childern registered under this account
// @Tags			Childern
// @Accept			json
// @Produce		json
// @Security		ParentLevelAuth
// @Param			Authorization	header		string												true	"JWT Token"
// @Param			limit			query		int													true	"limit searching param"
// @Param			offset			query		int													true	"offset searching param"
// @Success		200				{object}	StandardSuccessResponse{data=[]GetMyChildernOutput}	"Successful response"
// @Failure		400				{object}	StandardErrorResponse								"Bad request"
// @Failure		500				{object}	StandardErrorResponse								"Internal Error"
// @Router			/v1/childern [get]
func (s *Service) HandleGetMyChildern() echo.HandlerFunc {
	return func(c echo.Context) error {
		input := &GetMyChildrenInput{}
		if err := c.Bind(input); err != nil {
			return c.JSON(http.StatusBadRequest, StandardErrorResponse{
				StatusCode:   http.StatusBadRequest,
				ErrorMessage: "failed to parse input",
				ErrorCode:    http.StatusText(http.StatusBadRequest),
			})
		}

		children, err := s.childUsecase.GetRegisteredChildren(c.Request().Context(), usecase.GetRegisteredChildrenInput{
			Limit:  input.Limit,
			Offset: input.Offset,
		})

		if err != nil {
			return UsecaseErrorToRESTResponse(c, err)
		}

		output := []GetMyChildernOutput{}

		for _, child := range children {
			output = append(output, GetMyChildernOutput{
				ID:             child.ID,
				ParentUserID:   child.ParentUserID,
				ParentUserName: child.ParentUsername,
				DateOfBirth:    child.DateOfBirth,
				Gender:         child.Gender,
				Name:           child.Name,
				CreatedAt:      child.CreatedAt,
				UpdatedAt:      child.UpdatedAt,
			})
		}

		return c.JSON(http.StatusOK, StandardSuccessResponse{
			StatusCode: http.StatusOK,
			Message:    http.StatusText(http.StatusOK),
			Data:       output,
		})
	}
}

// @Summary		Search childern data
// @Description	Search childern data
// @Tags			Childern
// @Accept			json
// @Produce		json
// @Security		TherapistLevelAuth
// @Param			Authorization		header		string												true	"JWT Token"
// @Param			search_child_input	query		SearchChildrenInput									true	"search parameters"
// @Success		200					{object}	StandardSuccessResponse{data=SearchChildernOutput}	"Successful response"
// @Failure		400					{object}	StandardErrorResponse								"Bad request"
// @Failure		500					{object}	StandardErrorResponse								"Internal Error"
// @Router			/v1/childern/search [get]
func (s *Service) HandleSearchChildern() echo.HandlerFunc {
	return func(c echo.Context) error {
		input := &SearchChildrenInput{}
		if err := c.Bind(input); err != nil {
			return c.JSON(http.StatusBadRequest, StandardErrorResponse{
				StatusCode:   http.StatusBadRequest,
				ErrorMessage: "failed to parse input",
				ErrorCode:    http.StatusText(http.StatusBadRequest),
			})
		}

		children, err := s.childUsecase.Search(c.Request().Context(), usecase.SearchChildInput{
			ParentUserID: input.ParentUserID,
			Name:         input.Name,
			Gender:       input.Gender,
			Limit:        input.Limit,
			Offset:       input.Offset,
		})

		if err != nil {
			return UsecaseErrorToRESTResponse(c, err)
		}

		output := []SearchChildernOutput{}

		for _, child := range children {
			output = append(output, SearchChildernOutput{
				ID:           child.ID,
				ParentUserID: child.ParentUserID,
				DateOfBirth:  child.DateOfBirth,
				Gender:       child.Gender,
				Name:         child.Name,
				CreatedAt:    child.CreatedAt,
				UpdatedAt:    child.UpdatedAt,
			})
		}

		return c.JSON(http.StatusOK, StandardSuccessResponse{
			StatusCode: http.StatusOK,
			Message:    http.StatusText(http.StatusOK),
			Data:       output,
		})
	}
}

// @Summary		Get child ATEC score history
// @Description	The returned data will be JSON but contains sufficient data to be drawn as graph on frontend
// @Tags			Childern
// @Accept			json
// @Produce		json
// @Security		ParentLevelAuth
// @Param			Authorization	header		string														true	"JWT Token"
// @Param			child_id		path		string														true	"Child ID (UUID v4)"
// @Success		200				{object}	StandardSuccessResponse{data=[]usecase.StatisticComponent}	"Successful response"
// @Failure		400				{object}	StandardErrorResponse										"Bad request"
// @Failure		500				{object}	StandardErrorResponse										"Internal Error"
// @Router			/v1/childern/{child_id}/stats [get]
func (s *Service) HandleGetChildStats() echo.HandlerFunc {
	return func(c echo.Context) error {
		input := &GetChildStatInput{}
		if err := c.Bind(input); err != nil {
			return c.JSON(http.StatusBadRequest, StandardErrorResponse{
				StatusCode:   http.StatusBadRequest,
				ErrorMessage: "failed to parse input",
				ErrorCode:    http.StatusText(http.StatusBadRequest),
			})
		}

		stats, err := s.childUsecase.HandleGetStatistic(c.Request().Context(), usecase.GetStatisticInput{
			ChildID: input.ChildID,
		})

		if err != nil {
			return UsecaseErrorToRESTResponse(c, err)
		}

		return c.JSON(http.StatusOK, StandardSuccessResponse{
			StatusCode: http.StatusOK,
			Message:    http.StatusText(http.StatusOK),
			Data:       stats.Statistic, // deliberately return only statistic and avoid any type casting or data conversion
		})
	}
}
