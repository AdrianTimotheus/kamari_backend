package handler_test

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"kamari.com/dto"
	"kamari.com/httperror"
	"kamari.com/mocks"
	"kamari.com/model"
	"kamari.com/server"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHandler_GetCartHoverItems(t *testing.T) {
	t.Run("Should return 5 recent cart items by default", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		cartService := &mocks.CartService{}

		router := server.NewRouter(server.RouterConfig{
			CartService: cartService,
		})

		data := []dto.CartHoverRes{
			dto.CartHoverRes{
				Name:   "T-shirt",
				Price:  100000,
				PicUrl: "https://picsum.photos/id/237/200",
			},
			dto.CartHoverRes{
				Name:   "T-shirt",
				Price:  100000,
				PicUrl: "https://picsum.photos/id/237/200",
			},
			dto.CartHoverRes{
				Name:   "T-shirt",
				Price:  100000,
				PicUrl: "https://picsum.photos/id/237/200",
			},
			dto.CartHoverRes{
				Name:   "T-shirt",
				Price:  100000,
				PicUrl: "https://picsum.photos/id/237/200",
			},
			dto.CartHoverRes{
				Name:   "T-shirt",
				Price:  100000,
				PicUrl: "https://picsum.photos/id/237/200",
			},
		}

		user := model.User{
			Id:          4,
			Name:        "Adrian",
			Username:    "adrian",
			Email:       "adrian@gmail.com",
			Phone:       "08123456789",
			DateOfBirth: time.Unix(0, 0).UTC(),
			GenderId:    1,
			RoleId:      2,
		}

		cartService.On("GetCartHoverItems", user).Return(data, nil)
		expectedResponse, _ := json.Marshal(data)
		bearer := "Bearer " + "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJrYW1hcmkiLCJleHAiOjE2NjQxMzM1MjQsImlhdCI6MTY2NDEyOTkyNCwidXNlciI6eyJpZCI6NCwibmFtZSI6IkFkcmlhbiIsInVzZXJuYW1lIjoiYWRyaWFuIiwiZW1haWwiOiJhZHJpYW5AZ21haWwuY29tIiwicGhvbmUiOiIwODEyMzQ1Njc4OSIsImRhdGVfb2ZfYmlydGgiOiIxOTcwLTAxLTAxVDAwOjAwOjAwWiIsInBpY191cmwiOiIiLCJnZW5kZXJfaWQiOjEsInJvbGVfaWQiOjJ9fQ.utJvBaupEAcuzvCOrpmd96zNIq7sL6fQ_DHYM6tI4Lc"

		req, _ := http.NewRequest(http.MethodGet, "/api/carts/hover", nil)
		req.Header.Add("Authorization", bearer)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, string(expectedResponse), recorder.Body.String())
	})

	t.Run("Should return unauthorized error when user not logged in", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		cartService := &mocks.CartService{}

		router := server.NewRouter(server.RouterConfig{
			CartService: cartService,
		})

		user := model.User{
			Id:          4,
			Name:        "Adrian",
			Username:    "adrian",
			Email:       "adrian@gmail.com",
			Phone:       "08123456789",
			DateOfBirth: time.Unix(0, 0).UTC(),
			GenderId:    1,
			RoleId:      2,
		}

		cartService.On("GetCartHoverItems", user).Return(nil, httperror.UnauthorizedError())
		expectedResponse, _ := json.Marshal(httperror.UnauthorizedError())

		req, _ := http.NewRequest(http.MethodGet, "/api/carts/hover", nil)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusUnauthorized, recorder.Code)
		assert.Equal(t, string(expectedResponse), recorder.Body.String())
	})

	t.Run("Should return internal server error when an error happens", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		cartService := &mocks.CartService{}

		router := server.NewRouter(server.RouterConfig{
			CartService: cartService,
		})

		user := model.User{
			Id:          4,
			Name:        "Adrian",
			Username:    "adrian",
			Email:       "adrian@gmail.com",
			Phone:       "08123456789",
			DateOfBirth: time.Unix(0, 0).UTC(),
			GenderId:    1,
			RoleId:      2,
		}

		cartService.On("GetCartHoverItems", user).Return(nil, httperror.InternalServerError("Internal Server Error"))
		expectedResponse, _ := json.Marshal(httperror.InternalServerError("Internal Server Error"))
		bearer := "Bearer " + "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJrYW1hcmkiLCJleHAiOjE2NjQxMzM1MjQsImlhdCI6MTY2NDEyOTkyNCwidXNlciI6eyJpZCI6NCwibmFtZSI6IkFkcmlhbiIsInVzZXJuYW1lIjoiYWRyaWFuIiwiZW1haWwiOiJhZHJpYW5AZ21haWwuY29tIiwicGhvbmUiOiIwODEyMzQ1Njc4OSIsImRhdGVfb2ZfYmlydGgiOiIxOTcwLTAxLTAxVDAwOjAwOjAwWiIsInBpY191cmwiOiIiLCJnZW5kZXJfaWQiOjEsInJvbGVfaWQiOjJ9fQ.utJvBaupEAcuzvCOrpmd96zNIq7sL6fQ_DHYM6tI4Lc"

		req, _ := http.NewRequest(http.MethodGet, "/api/carts/hover", nil)
		req.Header.Add("Authorization", bearer)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
		assert.Equal(t, string(expectedResponse), recorder.Body.String())
	})
}
