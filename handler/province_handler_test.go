package handler_test

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"kamari.com/httperror"
	"kamari.com/mocks"
	"kamari.com/model"
	"kamari.com/server"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_GetAllProvinces(t *testing.T) {
	t.Run("Should return all provinces with status code 200 OK by default", func(t *testing.T) {
		provinceService := &mocks.ProvinceService{}
		router := server.NewRouter(server.RouterConfig{
			ProvinceService: provinceService,
		})
		recorder := httptest.NewRecorder()
		data := []model.Province{
			model.Province{
				Id:   1,
				Name: "DKI Jakarta",
			},
			model.Province{
				Id:   2,
				Name: "Banten",
			},
		}
		provinceService.On("GetAllProvinces").Return(data, nil)
		expectedResponse, _ := json.Marshal(data)

		req, _ := http.NewRequest(http.MethodGet, "/api/provinces", nil)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, string(expectedResponse), recorder.Body.String())
	})

	t.Run("Should return 500 Internal Server Error when an error happens", func(t *testing.T) {
		provinceService := &mocks.ProvinceService{}
		router := server.NewRouter(server.RouterConfig{
			ProvinceService: provinceService,
		})
		recorder := httptest.NewRecorder()

		provinceService.On("GetAllProvinces").Return(nil, httperror.InternalServerError("Internal server error"))
		expectedResponse, _ := json.Marshal(httperror.InternalServerError("Internal server error"))

		req, _ := http.NewRequest(http.MethodGet, "/api/provinces", nil)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
		assert.Equal(t, string(expectedResponse), recorder.Body.String())
	})
}
