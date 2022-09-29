package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"kamari.com/dto"
	"kamari.com/handler"
	"kamari.com/httperror"
	"kamari.com/mocks"
	"kamari.com/server"
)

func TestHandler_Register(t *testing.T) {
	t.Run("should throw an error when data payload is invalid", func(t *testing.T) {
		userService := &mocks.UserService{}
		router := server.NewRouter(server.RouterConfig{
			UserService: userService,
		})
		rr := httptest.NewRecorder()
		data := dto.RegisterRequest{
			Name:     "test",
			Username: "testwaw123",
		}
		expectedResponse, _ := json.Marshal(httperror.InternalServerError(handler.ErrFailedBindingValidation.Error()))
		payload, _ := json.Marshal(data)

		req, _ := http.NewRequest(http.MethodPost, "/api/users", strings.NewReader(string(payload)))
		router.ServeHTTP(rr, req)

		assert.Equal(t, rr.Code, http.StatusInternalServerError)
		assert.Equal(t, rr.Body.String(), string(expectedResponse))
	})

	t.Run("should give JSON response success", func(t *testing.T) {
		userService := &mocks.UserService{}
		router := server.NewRouter(server.RouterConfig{
			UserService: userService,
		})
		rr := httptest.NewRecorder()
		data := &dto.RegisterRequest{
			Name:        "abdul mantab",
			Username:    "abdulwahabgaming",
			Email:       "test@gmail.com",
			Password:    "notinusername",
			DateOfBirth: 1663642290536,
			Phone:       "084216327163",
			GenderId:    1,
		}
		userService.On("CreateNewUser", data).Return(nil)
		expectedResponse, _ := json.Marshal(map[string]any{
			"message": "user created successfully",
		})
		payload, _ := json.Marshal(data)

		req, _ := http.NewRequest(http.MethodPost, "/api/users", strings.NewReader(string(payload)))
		router.ServeHTTP(rr, req)

		assert.Equal(t, rr.Code, http.StatusCreated)
		assert.Equal(t, rr.Body.String(), string(expectedResponse))
	})

	t.Run("should throw an error when password contains username", func(t *testing.T) {
		userService := &mocks.UserService{}
		router := server.NewRouter(server.RouterConfig{
			UserService: userService,
		})
		rr := httptest.NewRecorder()
		data := &dto.RegisterRequest{
			Name:        "abdul",
			Username:    "abdull",
			Email:       "test@gmail.com",
			Password:    "abdull12345",
			DateOfBirth: 1663642290536,
			Phone:       "084216327163",
			GenderId:    1,
		}
		userService.On("CreateNewUser", data).Return(httperror.BadRequestError("password cannot contains username", "INVALID_PASSWORD_PAYLOAD"))
		expectedResponse, _ := json.Marshal(httperror.BadRequestError("password cannot contains username", "INVALID_PASSWORD_PAYLOAD"))
		payload, _ := json.Marshal(&data)

		req, _ := http.NewRequest(http.MethodPost, "/api/users", strings.NewReader(string(payload)))
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, string(expectedResponse), rr.Body.String())
	})
}

func TestHandler_Login(t *testing.T) {
	t.Run("should return an error if request body doesn't match dto", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		userService := &mocks.UserService{}

		router := server.NewRouter(server.RouterConfig{
			UserService: userService,
		})

		data := dto.LoginRequest{
			Email: "adrian@gmail.com",
		}

		expectedResponse, _ := json.Marshal(httperror.InternalServerError(handler.ErrFailedBindingValidation.Error()))
		payload, _ := json.Marshal(data)

		req, _ := http.NewRequest(http.MethodPost, "/api/login", strings.NewReader(string(payload)))
		router.ServeHTTP(recorder, req)

		assert.Equal(t, recorder.Code, http.StatusInternalServerError)
		assert.Equal(t, recorder.Body.String(), string(expectedResponse))
	})

	t.Run("should return 200 OK if email and password match", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		userService := &mocks.UserService{}

		router := server.NewRouter(server.RouterConfig{
			UserService: userService,
		})

		data := dto.LoginRequest{
			Email:    "adrian@gmail.com",
			Password: "password",
		}

		token := dto.TokenResponse{
			IDToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9." +
				"eyJpc3MiOiJrYW1hcmkiLCJleHAiOjE2NjQxMjkyOTIsIm" +
				"lhdCI6MTY2NDEyNTY5MiwidXNlciI6eyJpZCI6NCwibmFt" +
				"ZSI6IkFkcmlhbiIsInVzZXJuYW1lIjoiYWRyaWFuIiwiZW" +
				"1haWwiOiJhZHJpYW5AZ21haWwuY29tIiwicGhvbmUiOiIw" +
				"ODEyMzQ1Njc4OSIsImRhdGVfb2ZfYmlydGgiOiIxOTcwLT" +
				"AxLTAxVDAwOjAwOjAwWiIsInBpY191cmwiOiIiLCJnZW5k" +
				"ZXJfaWQiOjEsInJvbGVfaWQiOjJ9fQ.LZRkf2o4yHxR1-" +
				"9IpZvlW5w8ocavs9XAP5Xw1YyvRRw",
		}

		userService.On("Login", &data).Return(&token, nil)
		expectedResponse, _ := json.Marshal(token)
		payload, _ := json.Marshal(data)

		req, _ := http.NewRequest(http.MethodPost, "/api/login", strings.NewReader(string(payload)))
		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, string(expectedResponse), recorder.Body.String())
	})

	t.Run("should return 401 Unauthorized Error if email and password don't match", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		userService := &mocks.UserService{}

		router := server.NewRouter(server.RouterConfig{
			UserService: userService,
		})

		data := dto.LoginRequest{
			Email:    "adrian@gmail.com",
			Password: "notmatch",
		}

		userService.On("Login", &data).Return(nil, httperror.InvalidLogin("Invalid email or password"))
		expectedResponse, _ := json.Marshal(httperror.InvalidLogin("Invalid email or password"))
		payload, _ := json.Marshal(data)

		req, _ := http.NewRequest(http.MethodPost, "/api/login", strings.NewReader(string(payload)))
		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusUnauthorized, recorder.Code)
		assert.Equal(t, string(expectedResponse), recorder.Body.String())
	})
}
