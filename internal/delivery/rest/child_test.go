package rest_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/luckyAkbar/atec/internal/delivery/rest"
	"github.com/luckyAkbar/atec/internal/usecase"
	usecase_mock "github.com/luckyAkbar/atec/mocks/internal_/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChildService_HandleRegisterChildern(t *testing.T) {
	e := echo.New()
	group := e.Group("")

	mockChildUsecase := usecase_mock.NewChildUsecaseIface(t)
	service := rest.NewService(group, nil, nil, mockChildUsecase, nil)

	dateOfBirth, err := time.Parse("2006-01-02", "2021-12-29")
	require.NoError(t, err)

	testCases := []struct {
		name   string
		reqCtx func() (*httptest.ResponseRecorder, echo.Context)
		expect func(rec *httptest.ResponseRecorder)
		mockFn func(ectx echo.Context)
	}{
		{
			name: "invalid input body",
			reqCtx: func() (*httptest.ResponseRecorder, echo.Context) {
				req := httptest.NewRequest(http.MethodPost, "/v1/children", strings.NewReader(`{,}`))
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
				req := httptest.NewRequest(http.MethodPost, "/v1/children", strings.NewReader(`{
					"name": "username",
					"date_of_birth": "2021-12-29",
					"gender": false
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
				mockChildUsecase.EXPECT().Register(ectx.Request().Context(), usecase.RegisterChildInput{
					DateOfBirth: dateOfBirth,
					Gender:      false,
					Name:        "username",
				}).Return(nil, usecase.UsecaseError{
					ErrType: usecase.ErrBadRequest,
				}).Once()
			},
		},
		{
			name: "ok",
			reqCtx: func() (*httptest.ResponseRecorder, echo.Context) {
				req := httptest.NewRequest(http.MethodPost, "/v1/children", strings.NewReader(`{
					"name": "username",
					"date_of_birth": "2021-12-29",
					"gender": false
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
				mockChildUsecase.EXPECT().Register(ectx.Request().Context(), usecase.RegisterChildInput{
					DateOfBirth: dateOfBirth,
					Gender:      false,
					Name:        "username",
				}).Return(&usecase.RegisterChildOutput{}, nil).Once()
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec, ectx := tc.reqCtx()

			if tc.mockFn != nil {
				tc.mockFn(ectx)
			}

			err := service.HandleRegisterChildern()(ectx)

			require.NoError(t, err)
			tc.expect(rec)
		})
	}
}

func TestChildService_HandleUpdateChildern(t *testing.T) {
	e := echo.New()
	group := e.Group("")

	mockChildUsecase := usecase_mock.NewChildUsecaseIface(t)
	service := rest.NewService(group, nil, nil, mockChildUsecase, nil)

	id := "315be606-54b9-436d-a84b-90d118c745e7"
	childID, err := uuid.Parse(id)
	require.NoError(t, err)

	dateOfBirth, err := time.Parse("2006-01-02", "2021-12-29")
	require.NoError(t, err)

	gender := false
	username := "username"

	testCases := []struct {
		name   string
		reqCtx func() (*httptest.ResponseRecorder, echo.Context)
		expect func(rec *httptest.ResponseRecorder)
		mockFn func(ectx echo.Context)
	}{
		{
			name: "invalid input body",
			reqCtx: func() (*httptest.ResponseRecorder, echo.Context) {
				req := httptest.NewRequest(http.MethodPost, "/v1/children", strings.NewReader(`{,}`))
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
				req := httptest.NewRequest(http.MethodPost, "/v1/children", strings.NewReader(`{
					"name": "username",
					"date_of_birth": "2021-12-29",
					"gender": false
				}`))
				req.Header.Set("Content-Type", "application/json")

				rec := httptest.NewRecorder()
				ectx := e.NewContext(req, rec)

				ectx.SetParamNames("child_id")
				ectx.SetParamValues(id)

				return rec, ectx
			},
			expect: func(rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, rec.Code)
				assert.Equal(t, "application/json", rec.Header().Get(echo.HeaderContentType))
			},
			mockFn: func(ectx echo.Context) {
				mockChildUsecase.EXPECT().Update(ectx.Request().Context(), usecase.UpdateChildInput{
					DateOfBirth: &dateOfBirth,
					Gender:      &gender,
					Name:        &username,
					ChildID:     childID,
				}).Return(nil, usecase.UsecaseError{
					ErrType: usecase.ErrBadRequest,
				}).Once()
			},
		},
		{
			name: "ok",
			reqCtx: func() (*httptest.ResponseRecorder, echo.Context) {
				req := httptest.NewRequest(http.MethodPost, "/v1/children", strings.NewReader(`{
					"name": "username",
					"date_of_birth": "2021-12-29",
					"gender": false
				}`))
				req.Header.Set("Content-Type", "application/json")

				rec := httptest.NewRecorder()
				ectx := e.NewContext(req, rec)

				ectx.SetParamNames("child_id")
				ectx.SetParamValues(id)

				return rec, ectx
			},
			expect: func(rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, rec.Code)
				assert.Equal(t, "application/json", rec.Header().Get(echo.HeaderContentType))
			},
			mockFn: func(ectx echo.Context) {
				mockChildUsecase.EXPECT().Update(ectx.Request().Context(), usecase.UpdateChildInput{
					DateOfBirth: &dateOfBirth,
					Gender:      &gender,
					Name:        &username,
					ChildID:     childID,
				}).Return(&usecase.UpdateChildOutput{}, nil).Once()
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec, ectx := tc.reqCtx()

			if tc.mockFn != nil {
				tc.mockFn(ectx)
			}

			err := service.HandleUpdateChildern()(ectx)

			require.NoError(t, err)
			tc.expect(rec)
		})
	}
}

func TestChildService_HandleGetMyChildren(t *testing.T) {
	e := echo.New()
	group := e.Group("")

	mockChildUsecase := usecase_mock.NewChildUsecaseIface(t)
	service := rest.NewService(group, nil, nil, mockChildUsecase, nil)

	testCases := []struct {
		name   string
		reqCtx func() (*httptest.ResponseRecorder, echo.Context)
		expect func(rec *httptest.ResponseRecorder)
		mockFn func(ectx echo.Context)
	}{
		{
			name: "invalid input body",
			reqCtx: func() (*httptest.ResponseRecorder, echo.Context) {
				req := httptest.NewRequest(http.MethodGet, "/v1/children", strings.NewReader(`{,}`))
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
				req := httptest.NewRequest(http.MethodGet, "/v1/children", nil)
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
				mockChildUsecase.EXPECT().GetRegisteredChildren(ectx.Request().Context(), usecase.GetRegisteredChildrenInput{
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
				req := httptest.NewRequest(http.MethodGet, "/v1/children", nil)
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
				mockChildUsecase.EXPECT().GetRegisteredChildren(ectx.Request().Context(), usecase.GetRegisteredChildrenInput{
					Limit:  100,
					Offset: 0,
				}).Return([]usecase.GetRegisteredChildrenOutput{
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

			err := service.HandleGetMyChildern()(ectx)

			require.NoError(t, err)
			tc.expect(rec)
		})
	}
}

func TestChildService_HandleSearchChildern(t *testing.T) {
	e := echo.New()
	group := e.Group("")

	mockChildUsecase := usecase_mock.NewChildUsecaseIface(t)
	service := rest.NewService(group, nil, nil, mockChildUsecase, nil)

	testCases := []struct {
		name   string
		reqCtx func() (*httptest.ResponseRecorder, echo.Context)
		expect func(rec *httptest.ResponseRecorder)
		mockFn func(ectx echo.Context)
	}{
		{
			name: "invalid input body",
			reqCtx: func() (*httptest.ResponseRecorder, echo.Context) {
				req := httptest.NewRequest(http.MethodGet, "/v1/children/search", strings.NewReader(`{,}`))
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
				req := httptest.NewRequest(http.MethodGet, "/v1/children/search", nil)
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
				mockChildUsecase.EXPECT().Search(ectx.Request().Context(), usecase.SearchChildInput{
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
				req := httptest.NewRequest(http.MethodGet, "/v1/children/search", nil)
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
				mockChildUsecase.EXPECT().Search(ectx.Request().Context(), usecase.SearchChildInput{
					Limit:  100,
					Offset: 0,
				}).Return([]usecase.SearchChildOutput{
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

			err := service.HandleSearchChildern()(ectx)

			require.NoError(t, err)
			tc.expect(rec)
		})
	}
}

func TestChildService_HandleGetChildStats(t *testing.T) {
	e := echo.New()
	group := e.Group("")

	mockChildUsecase := usecase_mock.NewChildUsecaseIface(t)
	service := rest.NewService(group, nil, nil, mockChildUsecase, nil)

	id := "315be606-54b9-436d-a84b-90d118c745e7"
	childID, err := uuid.Parse(id)
	require.NoError(t, err)

	testCases := []struct {
		name   string
		reqCtx func() (*httptest.ResponseRecorder, echo.Context)
		expect func(rec *httptest.ResponseRecorder)
		mockFn func(ectx echo.Context)
	}{
		{
			name: "invalid input body",
			reqCtx: func() (*httptest.ResponseRecorder, echo.Context) {
				req := httptest.NewRequest(http.MethodPost, "/v1/children/stats", nil)
				req.Header.Set("Content-Type", "application/json")

				rec := httptest.NewRecorder()
				ectx := e.NewContext(req, rec)

				ectx.SetParamNames("child_id")
				ectx.SetParamValues("!@#\"")

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
				req := httptest.NewRequest(http.MethodPost, "/v1/children/stats", nil)
				req.Header.Set("Content-Type", "application/json")

				rec := httptest.NewRecorder()
				ectx := e.NewContext(req, rec)

				ectx.SetParamNames("child_id")
				ectx.SetParamValues(id)

				return rec, ectx
			},
			expect: func(rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, rec.Code)
				assert.Equal(t, "application/json", rec.Header().Get(echo.HeaderContentType))
			},
			mockFn: func(ectx echo.Context) {
				mockChildUsecase.EXPECT().HandleGetStatistic(ectx.Request().Context(), usecase.GetStatisticInput{
					ChildID: childID,
				}).Return(nil, usecase.UsecaseError{
					ErrType: usecase.ErrBadRequest,
				}).Once()
			},
		},
		{
			name: "ok",
			reqCtx: func() (*httptest.ResponseRecorder, echo.Context) {
				req := httptest.NewRequest(http.MethodPost, "/v1/children/stats", nil)
				req.Header.Set("Content-Type", "application/json")

				rec := httptest.NewRecorder()
				ectx := e.NewContext(req, rec)

				ectx.SetParamNames("child_id")
				ectx.SetParamValues(id)

				return rec, ectx
			},
			expect: func(rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, rec.Code)
				assert.Equal(t, "application/json", rec.Header().Get(echo.HeaderContentType))
			},
			mockFn: func(ectx echo.Context) {
				mockChildUsecase.EXPECT().HandleGetStatistic(ectx.Request().Context(), usecase.GetStatisticInput{
					ChildID: childID,
				}).Return(&usecase.GetStatisticOutput{}, nil).Once()
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec, ectx := tc.reqCtx()

			if tc.mockFn != nil {
				tc.mockFn(ectx)
			}

			err := service.HandleGetChildStats()(ectx)

			require.NoError(t, err)
			tc.expect(rec)
		})
	}
}
