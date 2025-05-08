package rest_test

import (
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

func TestPackageService_HandleCreatePackage(t *testing.T) {
	e := echo.New()
	group := e.Group("")

	mockPackageUsecase := usecase_mock.NewPackageUsecaseIface(t)
	service := rest.NewService(group, nil, mockPackageUsecase, nil, nil)

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
				req := httptest.NewRequest(http.MethodPost, "/v1/package", strings.NewReader(`{
					"package_name": "string",
					"image_result_attribute_key": {},
					"indication_categories": [],
					"questionnaire": {}
				}`))
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
				mockPackageUsecase.EXPECT().Create(ectx.Request().Context(), usecase.CreatePackageInput{
					PackageName:             "string",
					Questionnaire:           model.Questionnaire{},
					IndicationCategories:    model.IndicationCategories{},
					ImageResultAttributeKey: model.ImageResultAttributeKey{},
				}).
					Return(nil, usecase.UsecaseError{
						ErrType: usecase.ErrBadRequest,
					}).Once()
			},
		},
		{
			name: "ok",
			reqCtx: func() (*httptest.ResponseRecorder, echo.Context) {
				req := httptest.NewRequest(http.MethodPost, "/v1/package", strings.NewReader(`{
					"package_name": "string",
					"image_result_attribute_key": {},
					"indication_categories": [],
					"questionnaire": {}
				}`))
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
				mockPackageUsecase.EXPECT().Create(ectx.Request().Context(), usecase.CreatePackageInput{
					PackageName:             "string",
					Questionnaire:           model.Questionnaire{},
					IndicationCategories:    model.IndicationCategories{},
					ImageResultAttributeKey: model.ImageResultAttributeKey{},
				}).
					Return(&usecase.CreatePackageOutput{}, nil).Once()
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec, ectx := tc.reqCtx()

			if tc.mockFn != nil {
				tc.mockFn(ectx)
			}

			err := service.HandleCreatePackage()(ectx)

			require.NoError(t, err)
			tc.expect(rec)
		})
	}
}

func TestPackageService_HandleUpdatePackage(t *testing.T) {
	e := echo.New()
	group := e.Group("")

	mockPackageUsecase := usecase_mock.NewPackageUsecaseIface(t)
	service := rest.NewService(group, nil, mockPackageUsecase, nil, nil)

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
				req := httptest.NewRequest(http.MethodPut, "/v1/package", strings.NewReader(`{,}`))
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
				req := httptest.NewRequest(http.MethodPut, "/v1/package", strings.NewReader(`{
					"package_name": "string",
					"image_result_attribute_key": {},
					"indication_categories": [],
					"questionnaire": {}
				}`))
				req.Header.Set("Content-Type", "application/json")

				rec := httptest.NewRecorder()
				ectx := e.NewContext(req, rec)

				ectx.SetParamNames("package_id")
				ectx.SetParamValues(packageID.String())

				return rec, ectx
			},
			expect: func(rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, rec.Code)
				assert.Equal(t, "application/json", rec.Header().Get(echo.HeaderContentType))
			},
			mockFn: func(ectx echo.Context) {
				mockPackageUsecase.EXPECT().Update(ectx.Request().Context(), usecase.UpdatePackageInput{
					PackageName:   "string",
					Questionnaire: model.Questionnaire{},
					PackageID:     packageID,
				}).
					Return(nil, usecase.UsecaseError{
						ErrType: usecase.ErrBadRequest,
					}).Once()
			},
		},
		{
			name: "ok",
			reqCtx: func() (*httptest.ResponseRecorder, echo.Context) {
				req := httptest.NewRequest(http.MethodPut, "/v1/package", strings.NewReader(`{
					"package_name": "string",
					"image_result_attribute_key": {},
					"indication_categories": [],
					"questionnaire": {}
				}`))
				req.Header.Set("Content-Type", "application/json")

				rec := httptest.NewRecorder()
				ectx := e.NewContext(req, rec)

				ectx.SetParamNames("package_id")
				ectx.SetParamValues(packageID.String())

				return rec, ectx
			},
			expect: func(rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, rec.Code)
				assert.Equal(t, "application/json", rec.Header().Get(echo.HeaderContentType))
			},
			mockFn: func(ectx echo.Context) {
				mockPackageUsecase.EXPECT().Update(ectx.Request().Context(), usecase.UpdatePackageInput{
					PackageName:   "string",
					Questionnaire: model.Questionnaire{},
					PackageID:     packageID,
				}).
					Return(&usecase.UpdatePackageOutput{}, nil).Once()
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec, ectx := tc.reqCtx()

			if tc.mockFn != nil {
				tc.mockFn(ectx)
			}

			err := service.HandleUpdatePackage()(ectx)

			require.NoError(t, err)
			tc.expect(rec)
		})
	}
}

func TestPackageService_HandleActivationPackage(t *testing.T) {
	e := echo.New()
	group := e.Group("")

	mockPackageUsecase := usecase_mock.NewPackageUsecaseIface(t)
	service := rest.NewService(group, nil, mockPackageUsecase, nil, nil)

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
				req := httptest.NewRequest(http.MethodPatch, "/v1/package", strings.NewReader(`{,}`))
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
				req := httptest.NewRequest(http.MethodPatch, "/v1/package", strings.NewReader(`{
					"status": false
				}`))
				req.Header.Set("Content-Type", "application/json")

				rec := httptest.NewRecorder()
				ectx := e.NewContext(req, rec)

				ectx.SetParamNames("package_id")
				ectx.SetParamValues(packageID.String())

				return rec, ectx
			},
			expect: func(rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, rec.Code)
				assert.Equal(t, "application/json", rec.Header().Get(echo.HeaderContentType))
			},
			mockFn: func(ectx echo.Context) {
				mockPackageUsecase.EXPECT().ChangeActiveStatus(ectx.Request().Context(), usecase.ChangeActiveStatusInput{
					ActiveStatus: false,
					PackageID:    packageID,
				}).
					Return(nil, usecase.UsecaseError{
						ErrType: usecase.ErrBadRequest,
					}).Once()
			},
		},
		{
			name: "ok",
			reqCtx: func() (*httptest.ResponseRecorder, echo.Context) {
				req := httptest.NewRequest(http.MethodPatch, "/v1/package", strings.NewReader(`{
					"status": false
				}`))
				req.Header.Set("Content-Type", "application/json")

				rec := httptest.NewRecorder()
				ectx := e.NewContext(req, rec)

				ectx.SetParamNames("package_id")
				ectx.SetParamValues(packageID.String())

				return rec, ectx
			},
			expect: func(rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, rec.Code)
				assert.Equal(t, "application/json", rec.Header().Get(echo.HeaderContentType))
			},
			mockFn: func(ectx echo.Context) {
				mockPackageUsecase.EXPECT().ChangeActiveStatus(ectx.Request().Context(), usecase.ChangeActiveStatusInput{
					ActiveStatus: false,
					PackageID:    packageID,
				}).
					Return(&usecase.ChangeActiveStatusOutput{}, nil).Once()
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec, ectx := tc.reqCtx()

			if tc.mockFn != nil {
				tc.mockFn(ectx)
			}

			err := service.HandleActivationPackage()(ectx)

			require.NoError(t, err)
			tc.expect(rec)
		})
	}
}

func TestPackageService_HandleDeletePackage(t *testing.T) {
	e := echo.New()
	group := e.Group("")

	mockPackageUsecase := usecase_mock.NewPackageUsecaseIface(t)
	service := rest.NewService(group, nil, mockPackageUsecase, nil, nil)

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
				req := httptest.NewRequest(http.MethodDelete, "/v1/package", nil)
				req.Header.Set("Content-Type", "application/json")

				rec := httptest.NewRecorder()
				ectx := e.NewContext(req, rec)

				ectx.SetParamNames("package_id")
				ectx.SetParamValues("!@#$%^&*()")

				return rec, ectx
			},
			expect: func(rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, rec.Code)
			},
		},
		{
			name: "usecase return error",
			reqCtx: func() (*httptest.ResponseRecorder, echo.Context) {
				req := httptest.NewRequest(http.MethodDelete, "/v1/package", strings.NewReader(`{
					"status": false
				}`))
				req.Header.Set("Content-Type", "application/json")

				rec := httptest.NewRecorder()
				ectx := e.NewContext(req, rec)

				ectx.SetParamNames("package_id")
				ectx.SetParamValues(packageID.String())

				return rec, ectx
			},
			expect: func(rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, rec.Code)
				assert.Equal(t, "application/json", rec.Header().Get(echo.HeaderContentType))
			},
			mockFn: func(ectx echo.Context) {
				mockPackageUsecase.EXPECT().Delete(ectx.Request().Context(), packageID).
					Return(usecase.UsecaseError{
						ErrType: usecase.ErrBadRequest,
					}).Once()
			},
		},
		{
			name: "ok",
			reqCtx: func() (*httptest.ResponseRecorder, echo.Context) {
				req := httptest.NewRequest(http.MethodDelete, "/v1/package", strings.NewReader(`{
					"status": false
				}`))
				req.Header.Set("Content-Type", "application/json")

				rec := httptest.NewRecorder()
				ectx := e.NewContext(req, rec)

				ectx.SetParamNames("package_id")
				ectx.SetParamValues(packageID.String())

				return rec, ectx
			},
			expect: func(rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, rec.Code)
			},
			mockFn: func(ectx echo.Context) {
				mockPackageUsecase.EXPECT().Delete(ectx.Request().Context(), packageID).
					Return(nil).Once()
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec, ectx := tc.reqCtx()

			if tc.mockFn != nil {
				tc.mockFn(ectx)
			}

			err := service.HandleDeletePackage()(ectx)

			require.NoError(t, err)
			tc.expect(rec)
		})
	}
}

func TestPackageService_HandleSearchActivePackage(t *testing.T) {
	e := echo.New()
	group := e.Group("")

	mockPackageUsecase := usecase_mock.NewPackageUsecaseIface(t)
	service := rest.NewService(group, nil, mockPackageUsecase, nil, nil)

	testCases := []struct {
		name   string
		reqCtx func() (*httptest.ResponseRecorder, echo.Context)
		expect func(rec *httptest.ResponseRecorder)
		mockFn func(ectx echo.Context)
	}{
		{
			name: "usecase return error",
			reqCtx: func() (*httptest.ResponseRecorder, echo.Context) {
				req := httptest.NewRequest(http.MethodGet, "/v1/package/search", nil)
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
				mockPackageUsecase.EXPECT().FindActiveQuestionnaires(ectx.Request().Context()).Return(nil, usecase.UsecaseError{
					ErrType: usecase.ErrBadRequest,
				}).Once()
			},
		},
		{
			name: "ok",
			reqCtx: func() (*httptest.ResponseRecorder, echo.Context) {
				req := httptest.NewRequest(http.MethodGet, "/v1/children/search", nil)
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
				mockPackageUsecase.EXPECT().FindActiveQuestionnaires(ectx.Request().Context()).Return([]usecase.FindActiveQuestionnaireOutput{
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

			err := service.HandleSearchActivePackage()(ectx)

			require.NoError(t, err)
			tc.expect(rec)
		})
	}
}
