package handler_test

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"kamari.com/dto"
	"kamari.com/httperror"
	"kamari.com/mocks"
	"kamari.com/server"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_GetSellerDetail(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("should thrown an error when given invalid username param", func(t *testing.T) {
		merchantService := &mocks.MerchantService{}
		router := server.NewRouter(server.RouterConfig{
			MerchantService: merchantService,
		})
		rr := httptest.NewRecorder()
		expectedResponse, _ := json.Marshal(httperror.BadRequestError("failed to get value from the given param key", "INVALID_PARAM"))

		req, _ := http.NewRequest(http.MethodGet, "/api/sellers/ ", nil)
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, string(expectedResponse), rr.Body.String())
	})

	t.Run("should thrown an error when username given not found", func(t *testing.T) {
		merchantService := &mocks.MerchantService{}
		router := server.NewRouter(server.RouterConfig{
			MerchantService: merchantService,
		})
		rr := httptest.NewRecorder()
		expectedResponse, _ := json.Marshal(httperror.BadRequestError("record not found", "NO_ROWS"))
		merchantService.On("GetMerchantDetail", "anonymous").Return(nil, httperror.BadRequestError("record not found", "NO_ROWS"))

		req, _ := http.NewRequest(http.MethodGet, "/api/sellers/anonymous", nil)
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, string(expectedResponse), rr.Body.String())
	})

	t.Run("should return the seller detail when all request is ok", func(t *testing.T) {
		merchantService := &mocks.MerchantService{}
		router := server.NewRouter(server.RouterConfig{
			MerchantService: merchantService,
		})
		rr := httptest.NewRecorder()
		res := &dto.MerchantDetailResponse{
			Name:           "goodseller",
			Description:    "A good seller is nice",
			ProfilePicture: "google.com/fish.jpg",
			FollowerCount:  0,
			ProductCount:   0,
			ShopLocation:   "Ocean St",
		}
		expectedResponse, _ := json.Marshal(res)
		merchantService.On("GetMerchantDetail", "goodseller").Return(res, nil)

		req, _ := http.NewRequest(http.MethodGet, "/api/sellers/goodseller", nil)
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, string(expectedResponse), rr.Body.String())
	})
}
