package rest

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// @Summary		Register a new child
// @Description	Register a new child
// @Tags			Childern
// @Accept			json
// @Produce		json
// @Security		UserLevelAuth
// @Param			Authorization			header		string												true	"JWT Token"
// @Param			register_child_input	body		RegisterChildInput									true	"child details"
// @Success		200						{object}	StandardSuccessResponse{data=RegisterChildOutput}	"Successful response"
// @Failure		400						{object}	StandardErrorResponse								"Bad request"
// @Failure		500						{object}	StandardErrorResponse								"Internal Error"
// @Router			/v1/childern [post]
func (s *service) HandleRegisterChildern() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.NoContent(http.StatusNotImplemented)
	}
}

// @Summary		Update child data
// @Description	Update child data
// @Tags			Childern
// @Accept			json
// @Produce		json
// @Security		UserLevelAuth
// @Param			Authorization		header		string												true	"JWT Token"
//
// @Param			child_id			path		string												true	"Child ID to be updated (UUID v4)"
//
// @Param			update_child_input	body		UpdateChildernInput									true	"new child details"
// @Success		200					{object}	StandardSuccessResponse{data=UpdateChildernOutput}	"Successful response"
// @Failure		400					{object}	StandardErrorResponse								"Bad request"
// @Failure		500					{object}	StandardErrorResponse								"Internal Error"
// @Router			/v1/childern/{child_id} [put]
func (s *service) HandleUpdateChildern() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.NoContent(http.StatusNotImplemented)
	}
}

// @Summary		Get all childern registered under this account
// @Description	Get all childern registered under this account
// @Tags			Childern
// @Accept			json
// @Produce		json
// @Security		UserLevelAuth
// @Param			Authorization	header		string												true	"JWT Token"
// @Success		200				{object}	StandardSuccessResponse{data=GetMyChildernOutput}	"Successful response"
// @Failure		400				{object}	StandardErrorResponse								"Bad request"
// @Failure		500				{object}	StandardErrorResponse								"Internal Error"
// @Router			/v1/childern [get]
func (s *service) HandleGetMyChildern() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.NoContent(http.StatusNotImplemented)
	}
}

// @Summary		Search childern data
// @Description	Search childern data
// @Tags			Childern
// @Accept			json
// @Produce		json
// @Security		AdminLevelAuth
// @Param			Authorization	header		string												true	"JWT Token"
// @Param			name			query		string												false	"search by child name using ILIKE query"
// @Param			user_id			query		string												false	"search childern registered by this account (user_id is UUID v4)"
// @Success		200				{object}	StandardSuccessResponse{data=SearchChildernOutput}	"Successful response"
// @Failure		400				{object}	StandardErrorResponse								"Bad request"
// @Failure		500				{object}	StandardErrorResponse								"Internal Error"
// @Router			/v1/childern/search [get]
func (s *service) HandleSearchChildern() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.NoContent(http.StatusNotImplemented)
	}
}

// @Summary		Get child ATEC score history
// @Description	The returned data will be JSON but contains sufficient data to be drawn as graph on frontend
// @Tags			Childern
// @Accept			json
// @Produce		json
// @Security		UserLevelAuth
// @Param			Authorization	header		string												true	"JWT Token"
// @Param			child_id		path		string												true	"Child ID (UUID v4)"
// @Success		200				{object}	StandardSuccessResponse{data=GetChildStatOutput}	"Successful response"
// @Failure		400				{object}	StandardErrorResponse								"Bad request"
// @Failure		500				{object}	StandardErrorResponse								"Internal Error"
// @Router			/v1/childern/{child_id}/stats [get]
func (s *service) HandleGetChildStats() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.NoContent(http.StatusNotImplemented)
	}
}