package rest_test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
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

func TestQuestionnaireService_HandleSubmitQuestionnaire(t *testing.T) {
	e := echo.New()
	group := e.Group("")

	mockQuestionnaireUsecase := usecase_mock.NewQuestionnaireUsecaseIface(t)
	service := rest.NewService(group, nil, nil, nil, mockQuestionnaireUsecase)

	packageID := uuid.New()
	childID := uuid.New()

	testCases := []struct {
		name   string
		reqCtx func() (*httptest.ResponseRecorder, echo.Context)
		expect func(rec *httptest.ResponseRecorder)
		mockFn func(ectx echo.Context)
	}{
		{
			name: "invalid input body",
			reqCtx: func() (*httptest.ResponseRecorder, echo.Context) {
				req := httptest.NewRequest(http.MethodPost, "/v1/package", strings.NewReader(`{,}`))
				req.Header.Set("Content-Type", "application/json")

				rec := httptest.NewRecorder()
				ectx := e.NewContext(req, rec)

				return rec, ectx
			},
			expect: func(rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, rec.Code)
				assert.Equal(t, "application/json", rec.Header().Get(echo.HeaderContentType))
			},
		},
		{
			name: "usecase return error",
			reqCtx: func() (*httptest.ResponseRecorder, echo.Context) {
				req := httptest.NewRequest(http.MethodPost, "/v1/package", strings.NewReader(fmt.Sprintf(`{
					"package_id": "%s",
					"child_id": "%s",
					"answers": {}
				}`, packageID.String(), childID.String())))
				req.Header.Set("Content-Type", "application/json")

				rec := httptest.NewRecorder()
				ectx := e.NewContext(req, rec)

				return rec, ectx
			},
			expect: func(rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, rec.Code)
				assert.Equal(t, "application/json", rec.Header().Get(echo.HeaderContentType))
			},
			mockFn: func(ectx echo.Context) {
				mockQuestionnaireUsecase.EXPECT().HandleSubmitQuestionnaire(ectx.Request().Context(), usecase.SubmitQuestionnaireInput{
					PackageID: packageID,
					ChildID:   childID,
					Answers:   model.AnswerDetail{},
				}).Return(nil, usecase.UsecaseError{
					ErrType: usecase.ErrBadRequest,
				}).Once()
			},
		},
		{
			name: "ok",
			reqCtx: func() (*httptest.ResponseRecorder, echo.Context) {
				req := httptest.NewRequest(http.MethodPost, "/v1/package", strings.NewReader(fmt.Sprintf(`{
					"package_id": "%s",
					"child_id": "%s",
					"answers": {}
				}`, packageID.String(), childID.String())))
				req.Header.Set("Content-Type", "application/json")

				rec := httptest.NewRecorder()
				ectx := e.NewContext(req, rec)

				return rec, ectx
			},
			expect: func(rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, rec.Code)
				assert.Equal(t, "application/json", rec.Header().Get(echo.HeaderContentType))
			},
			mockFn: func(ectx echo.Context) {
				mockQuestionnaireUsecase.EXPECT().HandleSubmitQuestionnaire(ectx.Request().Context(), usecase.SubmitQuestionnaireInput{
					PackageID: packageID,
					ChildID:   childID,
					Answers:   model.AnswerDetail{},
				}).Return(&usecase.SubmitQuestionnaireOutput{}, nil).Once()
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec, ectx := tc.reqCtx()

			if tc.mockFn != nil {
				tc.mockFn(ectx)
			}

			err := service.HandleSubmitQuestionnaire()(ectx)

			require.NoError(t, err)
			tc.expect(rec)
		})
	}
}

func TestQuestionnaireService_HandleDownloadQuestionnaireResult(t *testing.T) {
	e := echo.New()
	group := e.Group("")

	mockQuestionnaireUsecase := usecase_mock.NewQuestionnaireUsecaseIface(t)
	service := rest.NewService(group, nil, nil, nil, mockQuestionnaireUsecase)

	resultID := uuid.New()

	testCases := []struct {
		name   string
		reqCtx func() (*httptest.ResponseRecorder, echo.Context)
		expect func(rec *httptest.ResponseRecorder)
		mockFn func(ectx echo.Context)
	}{
		{
			name: "invalid input body",
			reqCtx: func() (*httptest.ResponseRecorder, echo.Context) {
				req := httptest.NewRequest(http.MethodGet, "/v1/package", nil)
				req.Header.Set("Content-Type", "application/json")

				rec := httptest.NewRecorder()
				ectx := e.NewContext(req, rec)

				ectx.SetParamNames("result_id")
				ectx.SetParamValues("!@#$*()")

				return rec, ectx
			},
			expect: func(rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, rec.Code)
				assert.Equal(t, "application/json", rec.Header().Get(echo.HeaderContentType))
			},
		},
		{
			name: "usecase return error",
			reqCtx: func() (*httptest.ResponseRecorder, echo.Context) {
				req := httptest.NewRequest(http.MethodGet, "/v1/package", strings.NewReader(`{
					"package_name": "string",
					"image_result_attribute_key": {},
					"indication_categories": [],
					"questionnaire": {}
				}`))
				req.Header.Set("Content-Type", "application/json")

				rec := httptest.NewRecorder()
				ectx := e.NewContext(req, rec)

				ectx.SetParamNames("result_id")
				ectx.SetParamValues(resultID.String())

				return rec, ectx
			},
			expect: func(rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, rec.Code)
				assert.Equal(t, "application/json", rec.Header().Get(echo.HeaderContentType))
			},
			mockFn: func(ectx echo.Context) {
				mockQuestionnaireUsecase.EXPECT().HandleDownloadQuestionnaireResult(ectx.Request().Context(), usecase.DownloadQuestionnaireResultInput{
					ResultID: resultID,
				}).Return(nil, usecase.UsecaseError{
					ErrType: usecase.ErrBadRequest,
				}).Once()
			},
		},
		{
			name: "ok",
			reqCtx: func() (*httptest.ResponseRecorder, echo.Context) {
				req := httptest.NewRequest(http.MethodGet, "/v1/package", strings.NewReader(`{
					"package_name": "string",
					"image_result_attribute_key": {},
					"indication_categories": [],
					"questionnaire": {}
				}`))
				req.Header.Set("Content-Type", "application/json")

				rec := httptest.NewRecorder()
				ectx := e.NewContext(req, rec)

				ectx.SetParamNames("result_id")
				ectx.SetParamValues(resultID.String())

				return rec, ectx
			},
			expect: func(rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, rec.Code)
				assert.Equal(t, "image/png", rec.Header().Get(echo.HeaderContentType))
			},
			mockFn: func(ectx echo.Context) {
				mockQuestionnaireUsecase.EXPECT().HandleDownloadQuestionnaireResult(ectx.Request().Context(), usecase.DownloadQuestionnaireResultInput{
					ResultID: resultID,
				}).Return(&usecase.DownloadQuestionnaireResultOutput{
					ContentType: "image/png",
					Buffer:      *bytes.NewBufferString("abc"),
				}, nil).Once()
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec, ectx := tc.reqCtx()

			if tc.mockFn != nil {
				tc.mockFn(ectx)
			}

			err := service.HandleDownloadQuestionnaireResult()(ectx)

			require.NoError(t, err)
			tc.expect(rec)
		})
	}
}

func TestQuestionnaireService_HandleSearchQUestionnaireResults(t *testing.T) {
	e := echo.New()
	group := e.Group("")

	mockQuestionnaireUsecase := usecase_mock.NewQuestionnaireUsecaseIface(t)
	service := rest.NewService(group, nil, nil, nil, mockQuestionnaireUsecase)

	testCases := []struct {
		name   string
		reqCtx func() (*httptest.ResponseRecorder, echo.Context)
		expect func(rec *httptest.ResponseRecorder)
		mockFn func(ectx echo.Context)
	}{
		{
			name: "invalid input body",
			reqCtx: func() (*httptest.ResponseRecorder, echo.Context) {
				req := httptest.NewRequest(http.MethodGet, "/v1/questionnaire/search", strings.NewReader(`{,}`))
				req.Header.Set("Content-Type", "application/json")

				q := req.URL.Query()
				q.Add("limit", "abc")
				q.Add("offset", "0")
				req.URL.RawQuery = q.Encode()

				rec := httptest.NewRecorder()
				ectx := e.NewContext(req, rec)

				return rec, ectx
			},
			expect: func(rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, rec.Code)
				assert.Equal(t, "application/json", rec.Header().Get(echo.HeaderContentType))
			},
		},
		{
			name: "usecase return error",
			reqCtx: func() (*httptest.ResponseRecorder, echo.Context) {
				req := httptest.NewRequest(http.MethodGet, "/v1/questionnaire/search", nil)
				req.Header.Set("Content-Type", "application/json")

				q := req.URL.Query()
				q.Add("limit", "100")
				q.Add("offset", "0")
				req.URL.RawQuery = q.Encode()

				rec := httptest.NewRecorder()
				ectx := e.NewContext(req, rec)

				return rec, ectx
			},
			expect: func(rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, rec.Code)
				assert.Equal(t, "application/json", rec.Header().Get(echo.HeaderContentType))
			},
			mockFn: func(ectx echo.Context) {
				mockQuestionnaireUsecase.EXPECT().HandleSearchQuestionnaireResult(ectx.Request().Context(), usecase.SearchQuestionnaireResultInput{
					Limit:  100,
					Offset: 0,
				}).Return(nil, usecase.UsecaseError{
					ErrType: usecase.ErrBadRequest,
				}).Once()
			},
		},
		{
			name: "ok",
			reqCtx: func() (*httptest.ResponseRecorder, echo.Context) {
				req := httptest.NewRequest(http.MethodGet, "/v1/questionnaire/search", nil)
				req.Header.Set("Content-Type", "application/json")

				q := req.URL.Query()
				q.Add("limit", "100")
				q.Add("offset", "0")
				req.URL.RawQuery = q.Encode()

				rec := httptest.NewRecorder()
				ectx := e.NewContext(req, rec)

				return rec, ectx
			},
			expect: func(rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, rec.Code)
				assert.Equal(t, "application/json", rec.Header().Get(echo.HeaderContentType))
			},
			mockFn: func(ectx echo.Context) {
				mockQuestionnaireUsecase.EXPECT().HandleSearchQuestionnaireResult(ectx.Request().Context(), usecase.SearchQuestionnaireResultInput{
					Limit:  100,
					Offset: 0,
				}).Return([]usecase.SearchQuestionnaireResultOutput{
					{
						ID: uuid.New(),
					},
				}, nil).Once()
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec, ectx := tc.reqCtx()

			if tc.mockFn != nil {
				tc.mockFn(ectx)
			}

			err := service.HandleSearchQUestionnaireResults()(ectx)

			require.NoError(t, err)
			tc.expect(rec)
		})
	}
}

func TestQuestionnaireService_HandleGetMyQUestionnaireResults(t *testing.T) {
	e := echo.New()
	group := e.Group("")

	mockQuestionnaireUsecase := usecase_mock.NewQuestionnaireUsecaseIface(t)
	service := rest.NewService(group, nil, nil, nil, mockQuestionnaireUsecase)

	testCases := []struct {
		name   string
		reqCtx func() (*httptest.ResponseRecorder, echo.Context)
		expect func(rec *httptest.ResponseRecorder)
		mockFn func(ectx echo.Context)
	}{
		{
			name: "invalid input body",
			reqCtx: func() (*httptest.ResponseRecorder, echo.Context) {
				req := httptest.NewRequest(http.MethodGet, "/v1/questionnaire/", strings.NewReader(`{,}`))
				req.Header.Set("Content-Type", "application/json")

				q := req.URL.Query()
				q.Add("limit", "abc")
				q.Add("offset", "0")
				req.URL.RawQuery = q.Encode()

				rec := httptest.NewRecorder()
				ectx := e.NewContext(req, rec)

				return rec, ectx
			},
			expect: func(rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, rec.Code)
				assert.Equal(t, "application/json", rec.Header().Get(echo.HeaderContentType))
			},
		},
		{
			name: "usecase return error",
			reqCtx: func() (*httptest.ResponseRecorder, echo.Context) {
				req := httptest.NewRequest(http.MethodGet, "/v1/questionnaire/", nil)
				req.Header.Set("Content-Type", "application/json")

				q := req.URL.Query()
				q.Add("limit", "100")
				q.Add("offset", "0")
				req.URL.RawQuery = q.Encode()

				rec := httptest.NewRecorder()
				ectx := e.NewContext(req, rec)

				return rec, ectx
			},
			expect: func(rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, rec.Code)
				assert.Equal(t, "application/json", rec.Header().Get(echo.HeaderContentType))
			},
			mockFn: func(ectx echo.Context) {
				mockQuestionnaireUsecase.EXPECT().HandleGetUserHistory(ectx.Request().Context(), usecase.GetUserHistoryInput{
					Limit:  100,
					Offset: 0,
				}).Return(nil, usecase.UsecaseError{
					ErrType: usecase.ErrBadRequest,
				}).Once()
			},
		},
		{
			name: "ok",
			reqCtx: func() (*httptest.ResponseRecorder, echo.Context) {
				req := httptest.NewRequest(http.MethodGet, "/v1/questionnaire/", nil)
				req.Header.Set("Content-Type", "application/json")

				q := req.URL.Query()
				q.Add("limit", "100")
				q.Add("offset", "0")
				req.URL.RawQuery = q.Encode()

				rec := httptest.NewRecorder()
				ectx := e.NewContext(req, rec)

				return rec, ectx
			},
			expect: func(rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, rec.Code)
				assert.Equal(t, "application/json", rec.Header().Get(echo.HeaderContentType))
			},
			mockFn: func(ectx echo.Context) {
				mockQuestionnaireUsecase.EXPECT().HandleGetUserHistory(ectx.Request().Context(), usecase.GetUserHistoryInput{
					Limit:  100,
					Offset: 0,
				}).Return([]usecase.GetUserHistoryOutput{
					{
						ID: uuid.New(),
					},
				}, nil).Once()
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec, ectx := tc.reqCtx()

			if tc.mockFn != nil {
				tc.mockFn(ectx)
			}

			err := service.HandleGetMyQUestionnaireResults()(ectx)

			require.NoError(t, err)
			tc.expect(rec)
		})
	}
}

func TestQuestionnaireService_HandleGetATECQuestionaire(t *testing.T) {
	e := echo.New()
	group := e.Group("")

	mockQuestionnaireUsecase := usecase_mock.NewQuestionnaireUsecaseIface(t)
	service := rest.NewService(group, nil, nil, nil, mockQuestionnaireUsecase)
	packageID := uuid.New()

	testCases := []struct {
		name   string
		reqCtx func() (*httptest.ResponseRecorder, echo.Context)
		expect func(rec *httptest.ResponseRecorder)
		mockFn func(ectx echo.Context)
	}{
		{
			name: "invalid input body",
			reqCtx: func() (*httptest.ResponseRecorder, echo.Context) {
				req := httptest.NewRequest(http.MethodGet, "/v1/questionnaire/", strings.NewReader(`{,}`))
				req.Header.Set("Content-Type", "application/json")

				q := req.URL.Query()
				q.Add("package_id", "!@#$")
				req.URL.RawQuery = q.Encode()

				rec := httptest.NewRecorder()
				ectx := e.NewContext(req, rec)

				return rec, ectx
			},
			expect: func(rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, rec.Code)
				assert.Equal(t, "application/json", rec.Header().Get(echo.HeaderContentType))
			},
		},
		{
			name: "usecase return error",
			reqCtx: func() (*httptest.ResponseRecorder, echo.Context) {
				req := httptest.NewRequest(http.MethodGet, "/v1/questionnaire/", nil)
				req.Header.Set("Content-Type", "application/json")

				q := req.URL.Query()
				q.Add("package_id", packageID.String())
				req.URL.RawQuery = q.Encode()

				rec := httptest.NewRecorder()
				ectx := e.NewContext(req, rec)

				return rec, ectx
			},
			expect: func(rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, rec.Code)
				assert.Equal(t, "application/json", rec.Header().Get(echo.HeaderContentType))
			},
			mockFn: func(ectx echo.Context) {
				mockQuestionnaireUsecase.EXPECT().HandleInitializeATECQuestionnaire(ectx.Request().Context(), usecase.InitializeATECQuestionnaireInput{
					PackageID: packageID,
				}).Return(nil, usecase.UsecaseError{
					ErrType: usecase.ErrBadRequest,
				}).Once()
			},
		},
		{
			name: "ok",
			reqCtx: func() (*httptest.ResponseRecorder, echo.Context) {
				req := httptest.NewRequest(http.MethodGet, "/v1/questionnaire/", nil)
				req.Header.Set("Content-Type", "application/json")

				q := req.URL.Query()
				q.Add("package_id", packageID.String())
				req.URL.RawQuery = q.Encode()

				rec := httptest.NewRecorder()
				ectx := e.NewContext(req, rec)

				return rec, ectx
			},
			expect: func(rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, rec.Code)
				assert.Equal(t, "application/json", rec.Header().Get(echo.HeaderContentType))
			},
			mockFn: func(ectx echo.Context) {
				mockQuestionnaireUsecase.EXPECT().HandleInitializeATECQuestionnaire(ectx.Request().Context(), usecase.InitializeATECQuestionnaireInput{
					PackageID: packageID,
				}).Return(&usecase.InitializeATECQuestionnaireOutput{
					ID: packageID,
				}, nil).Once()
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec, ectx := tc.reqCtx()

			if tc.mockFn != nil {
				tc.mockFn(ectx)
			}

			err := service.HandleGetATECQuestionaire()(ectx)

			require.NoError(t, err)
			tc.expect(rec)
		})
	}
}
