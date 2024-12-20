package rest

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// @Summary		Submit questionnaire result
// @Description	When submiting questionnaire for a child, ensure using the parent's account or Admin level account. otherwise will be blocked.
// @Tags			Questionnaire
// @Accept			json
// @Produce		json
//
// @Param			Authorization				header		string													false	"Optional jwt auth token to be used if want to fill questionnaire based on a registered child"
//
// @Param			submit_questionnaire_input	body		SubmitQuestionnaireInput								true	"full questionnaire answer"
// @Success		200							{object}	StandardSuccessResponse{data=SubmitQuestionnaireOutput}	"Successful response"
// @Failure		400							{object}	StandardErrorResponse									"Bad request"
// @Failure		500							{object}	StandardErrorResponse									"Internal Error"
// @Router			/v1/atec/questionnaires [post]
func (s *service) HandleSubmitQuestionnaire() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.NoContent(http.StatusNotImplemented)
	}
}

// @Summary		Download quesionnaire result as image
// @Description	Download quesionnaire result as image
// @Tags			Questionnaire
// @Accept			json
// @Produce		jpeg
// @Param			Authorization	header		string					false	"Optional jwt auth token to be used if want to download result owned by a specific owner"
// @Param			result_id		path		string					true	"ID of the submitted questionnaire (UUID v4)"
// @Failure		400				{object}	StandardErrorResponse	"Bad request"
// @Failure		500				{object}	StandardErrorResponse	"Internal Error"
// @Router			/v1/atec/questionnaires/results/{result_id} [get]
func (s *service) HandleDownloadQuestionnaireResult() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.NoContent(http.StatusNotImplemented)
	}
}

// @Summary		Search questionnaire result
// @Description	Search through all the submitted ATEC questionnaires to the systems
// @Tags			Questionnaire
// @Accept			json
// @Produce		json
// @Security		AdminLevelAuth
// @Param			Authorization					header		string															true	"JWT token to prove that you're admin"
// @Param			search_questionnaire_results	query		SearchQUestionnaireResultsInput									true	"param to search"
// @Success		200								{object}	StandardSuccessResponse{data=SearchQUestionnaireResultsOutput}	"success response"
// @Failure		400								{object}	StandardErrorResponse											"Bad request"
// @Failure		500								{object}	StandardErrorResponse											"Internal Error"
// @Router			/v1/atec/questionnaires/results [get]
func (s *service) HandleSearchQUestionnaireResults() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.NoContent(http.StatusNotImplemented)
	}
}

// @Summary		Get my quesionnaires results
// @Description	Get all questionnaires result submitted by the requester account
// @Tags			Questionnaire
// @Accept			json
// @Produce		json
// @Security		UserLevelAuth
// @Param			Authorization					header		string															true	"JWT token from auth process"
// @Param			get_my_questionnaire_results	query		GetMyQUestionnaireResultsInput									true	"param to search"
// @Success		200								{object}	StandardSuccessResponse{data=SearchQUestionnaireResultsOutput}	"success response"
// @Failure		400								{object}	StandardErrorResponse											"Bad request"
// @Failure		500								{object}	StandardErrorResponse											"Internal Error"
// @Router			/v1/atec/questionnaires/results/my [get]
func (s *service) HandleGetMyQUestionnaireResults() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.NoContent(http.StatusNotImplemented)
	}
}

// @Summary		Initialize ATEC questionnaire or get one
// @Description	Used when a user wants to get an active ATEC questionnaire to be filled later
// @Tags			Questionnaire
// @Accept			json
// @Produce		json
// @Param			get_atec_questionnaire	query		GetATECQuestionnaireInput									false	"optional to get questionnaire from its package id. if empty, a default one will be returned"
// @Success		200						{object}	StandardSuccessResponse{data=GetATECQuestionnaireOutput}	"success response"
// @Failure		400						{object}	StandardErrorResponse										"Bad request"
// @Failure		500						{object}	StandardErrorResponse										"Internal Error"
// @Router			/v1/atec/questionnaires/ [get]
func (s *service) HandleGetATECQuestionaire() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.NoContent(http.StatusNotImplemented)
	}
}
