package rest

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/luckyAkbar/atec/internal/usecase"
	"gopkg.in/guregu/null.v4"
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
		input := &SubmitQuestionnaireInput{}
		if err := c.Bind(input); err != nil {
			return c.JSON(http.StatusBadRequest, StandardErrorResponse{
				StatusCode:   http.StatusBadRequest,
				ErrorMessage: "failed to parse input",
				ErrorCode:    http.StatusText(http.StatusBadRequest),
			})
		}

		output, err := s.questionnaireUsecase.HandleSubmitQuestionnaire(c.Request().Context(), usecase.SubmitQuestionnaireInput{
			ChildID:   input.ChildID,
			Answers:   input.Answers,
			PackageID: input.PackageID,
		})

		if err != nil {
			return usecaseErrorToRESTResponse(c, err)
		}

		return c.JSON(http.StatusOK, StandardSuccessResponse{
			StatusCode: http.StatusOK,
			Message:    http.StatusText(http.StatusOK),
			Data: SubmitQuestionnaireOutput{
				ResultID:  output.ResultID,
				Grade:     output.Result,
				ChildID:   output.ChildID,
				CreatedBy: output.CreatedBy,
			},
		})
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
		input := &DownloadQuestionnaireResultInput{}
		if err := c.Bind(input); err != nil {
			return c.JSON(http.StatusBadRequest, StandardErrorResponse{
				StatusCode:   http.StatusBadRequest,
				ErrorMessage: "failed to parse input",
				ErrorCode:    http.StatusText(http.StatusBadRequest),
			})
		}

		output, err := s.questionnaireUsecase.HandleDownloadQuestionnaireResult(c.Request().Context(), usecase.DownloadQuestionnaireResultInput{
			ResultID: input.ResultID,
		})

		if err != nil {
			return usecaseErrorToRESTResponse(c, err)
		}

		c.Response().Header().Set("Content-Type", output.ContentType)
		c.Response().Header().Set("Content-Length", strconv.Itoa(output.Buffer.Len()))

		return c.Blob(http.StatusOK, output.ContentType, output.Buffer.Bytes())
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
		input := &SearchQUestionnaireResultsInput{}
		if err := c.Bind(input); err != nil {
			return c.JSON(http.StatusBadRequest, StandardErrorResponse{
				StatusCode:   http.StatusBadRequest,
				ErrorMessage: "failed to parse input",
				ErrorCode:    http.StatusText(http.StatusBadRequest),
			})
		}

		output, err := s.questionnaireUsecase.HandleSearchQuestionnaireResult(c.Request().Context(), usecase.SearchQuestionnaireResultInput{
			ID:        input.ResultID,
			PackageID: input.PackageID,
			ChildID:   input.ChildID,
			CreatedBy: input.CreatedByID,
			Limit:     input.Limit,
			Offset:    input.Offset,
		})

		if err != nil {
			return usecaseErrorToRESTResponse(c, err)
		}

		result := []SearchQUestionnaireResultsOutput{}
		for _, val := range output {
			result = append(result, SearchQUestionnaireResultsOutput{
				ID:        val.ID,
				PackageID: val.PackageID,
				ChildID:   val.ChildID,
				CreatedBy: val.CreatedBy,
				Answer:    val.Answer,
				Result:    val.Result,
				CreatedAt: val.CreatedAt,
				UpdatedAt: val.UpdatedAt,
				DeletedAt: null.NewTime(val.DeletedAt.Time, val.DeletedAt.Valid),
			})
		}

		return c.JSON(http.StatusOK, StandardSuccessResponse{
			StatusCode: http.StatusOK,
			Message:    http.StatusText(http.StatusOK),
			Data:       result,
		})
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
		input := &GetMyQUestionnaireResultsInput{}
		if err := c.Bind(input); err != nil {
			return c.JSON(http.StatusBadRequest, StandardErrorResponse{
				StatusCode:   http.StatusBadRequest,
				ErrorMessage: "failed to parse input",
				ErrorCode:    http.StatusText(http.StatusBadRequest),
			})
		}

		output, err := s.questionnaireUsecase.HandleGetUserHistory(c.Request().Context(), usecase.GetUserHistoryInput{
			Limit:  input.Limit,
			Offset: input.Offset,
		})

		if err != nil {
			return usecaseErrorToRESTResponse(c, err)
		}

		result := []SearchQUestionnaireResultsOutput{}
		for _, val := range output {
			result = append(result, SearchQUestionnaireResultsOutput{
				ID:        val.ID,
				PackageID: val.PackageID,
				ChildID:   val.ChildID,
				CreatedBy: val.CreatedBy,
				Answer:    val.Answer,
				Result:    val.Result,
				CreatedAt: val.CreatedAt,
				UpdatedAt: val.UpdatedAt,
				DeletedAt: null.NewTime(val.DeletedAt.Time, val.DeletedAt.Valid),
			})
		}

		return c.JSON(http.StatusOK, StandardSuccessResponse{
			StatusCode: http.StatusOK,
			Message:    http.StatusText(http.StatusOK),
			Data:       result,
		})
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
// @Router			/v1/atec/questionnaires [get]
func (s *service) HandleGetATECQuestionaire() echo.HandlerFunc {
	return func(c echo.Context) error {
		input := &GetATECQuestionnaireInput{}
		if err := c.Bind(input); err != nil {
			return c.JSON(http.StatusBadRequest, StandardErrorResponse{
				StatusCode:   http.StatusBadRequest,
				ErrorMessage: "failed to parse input",
				ErrorCode:    http.StatusText(http.StatusBadRequest),
			})
		}

		output, err := s.questionnaireUsecase.HandleInitializeATECQuestionnaire(c.Request().Context(), usecase.InitializeATECQuestionnaireInput{
			PackageID: input.PackageID,
		})

		if err != nil {
			return usecaseErrorToRESTResponse(c, err)
		}

		return c.JSON(http.StatusOK, StandardSuccessResponse{
			StatusCode: http.StatusOK,
			Message:    http.StatusText(http.StatusOK),
			Data: GetATECQuestionnaireOutput{
				ID:            output.ID,
				Questionnaire: output.Questionnaire,
				Name:          output.Name,
			},
		})
	}
}
