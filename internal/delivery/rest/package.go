package rest

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// @Summary		Create new ATEC questionaire package
// @Description	Create new ATEC questionaire package
// @Tags			ATEC Package
// @Accept			json
// @Produce		json
// @Security		AdminLevelAuth
// @Param			Authorization			header		string												true	"JWT Token"
// @Param			create_package_input	body		CreatePackageInput									true	"ATEC questionnarie package details"
// @Success		200						{object}	StandardSuccessResponse{data=CreatePackageOutput}	"Successful response"
// @Failure		400						{object}	StandardErrorResponse								"Bad request"
// @Failure		500						{object}	StandardErrorResponse								"Internal Error"
// @Router			/v1/atec/packages [post]
func (s *service) HandleCreatePackage() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.NoContent(http.StatusNotImplemented)
	}
}

// @Summary		Update existing ATEC questionnarie package
// @Description	Update existing ATEC questionnarie package
// @Tags			ATEC Package
// @Accept			json
// @Produce		json
// @Security		AdminLevelAuth
// @Param			Authorization			header		string												true	"JWT Token"
// @Param			update_package_input	body		UpdatePackageInput									true	"ATEC questionnarie package details"
// @Success		200						{object}	StandardSuccessResponse{data=UpdatePackageOutput}	"Successful response"
// @Failure		400						{object}	StandardErrorResponse								"Bad request"
// @Failure		500						{object}	StandardErrorResponse								"Internal Error"
// @Router			/v1/atec/packages [put]
func (s *service) HandleUpdatePackage() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.NoContent(http.StatusNotImplemented)
	}
}

// @Summary		Update package activation status
// @Description	Update existing ATEC questionnarie package activation status
// @Tags			ATEC Package
// @Accept			json
// @Produce		json
// @Security		AdminLevelAuth
// @Param			Authorization				header		string													true	"JWT Token"
//
// @Param			package_id					path		string													true	"package ID to be activated/deactivated (UUID v4)"
// @Param			activation_package_input	body		ActivationPackageInput									true	"activation status"
// @Success		200							{object}	StandardSuccessResponse{data=ActivationPackageOutput}	"Successful response"
// @Failure		400							{object}	StandardErrorResponse									"Bad request"
// @Failure		500							{object}	StandardErrorResponse									"Internal Error"
// @Router			/v1/atec/packages/{package_id} [patch]
func (s *service) HandleActivationPackage() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.NoContent(http.StatusNotImplemented)
	}
}

// @Summary		Delete ATEC questionnaire package
// @Description	Delete ATEC questionnaire package
// @Tags			ATEC Package
// @Accept			json
// @Produce		json
// @Security		AdminLevelAuth
// @Param			Authorization	header		string						true	"JWT Token"
//
// @Param			package_id		path		string						true	"package ID to be deleted (UUID v4)"
// @Success		200				{object}	StandardSuccessResponse{}	"Successful response"
// @Failure		400				{object}	StandardErrorResponse		"Bad request"
// @Failure		500				{object}	StandardErrorResponse		"Internal Error"
// @Router			/v1/atec/packages/{package_id} [delete]
func (s *service) HandleDeletePackage() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.NoContent(http.StatusNotImplemented)
	}
}

// @Summary		Get all active packages
// @Description	Get all active packages
// @Tags			ATEC Package
// @Accept			json
// @Produce		json
// @Success		200	{object}	StandardSuccessResponse{data=SearchActivePackageOutput}	"Successful response"
// @Failure		400	{object}	StandardErrorResponse									"Bad request"
// @Failure		500	{object}	StandardErrorResponse									"Internal Error"
// @Router			/v1/atec/packages/active [get]
func (s *service) HandleSearchActivePackage() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.NoContent(http.StatusNotImplemented)
	}
}
