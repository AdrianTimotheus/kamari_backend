package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"kamari.com/dto"
	"kamari.com/httperror"
	"kamari.com/model"
)

func (h Handler) Register(c *gin.Context) {
	var (
		req *dto.RegisterRequest
		err error
	)

	if err = c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(ErrFailedBindingValidation)
		return
	}

	err = h.userService.CreateNewUser(req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "user created successfully",
	})
}

func (h Handler) Login(c *gin.Context) {
	var (
		req *dto.LoginRequest
		err error
	)
	if err = c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(ErrFailedBindingValidation)
		return
	}
	response, err := h.userService.Login(req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, response)
}

func (h Handler) LoginWithGoogle(c *gin.Context) {
	var (
		req *dto.GoogleAuthTokenReq
		err error
	)
	if err = c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(ErrFailedBindingValidation)
		return
	}
	response, err := h.userService.LoginWithGmail(req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, response)
}

func (h Handler) GetNewTokenFromRefresh(c *gin.Context) {
	response, err := h.userService.RefreshToken()
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, response)
}

func (h Handler) Logout(c *gin.Context) {
	err := h.userService.DeleteLogin()
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "logout success",
	})
}

func (h Handler) CheckPasswordAuth(c *gin.Context) {
	userVal, _ := c.Get("user")
	user := userVal.(model.User)
	var payload dto.CheckAuthPassword
	if err := c.ShouldBind(&payload); err != nil {
		_ = c.Error(httperror.BadRequestError("Missing Required Body", "MISSING_REQUIRED_BODY"))
		return
	}
	status, err := h.userService.CheckAuthPassword(payload, user.Id)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dto.StatusMessageResponse{Status: status})
}
func (h Handler) GenerateOTP(c *gin.Context) {
	userVal, _ := c.Get("user")
	user := userVal.(model.User)
	response, err := h.userService.GenerateOTPRedis(user)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h Handler) AuthOTP(c *gin.Context) {
	userVal, _ := c.Get("user")
	user := userVal.(model.User)
	var payload dto.AuthOTPRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		_ = c.Error(httperror.BadRequestError("Missing required body", "MISSING_REQUIRED_BODY"))
		return
	}
	response, err := h.userService.AuthOTP(user, payload)
	if err != nil {
		_ = c.Error(err)
		return
	}
	if response.Status {
		c.JSON(http.StatusOK, response)
		return
	}
	c.JSON(http.StatusOK, response)
}

func (h Handler) GetUserDetail(c *gin.Context) {
	userVal, _ := c.Get("user")
	user := userVal.(model.User)

	response, err := h.userService.GetUserDetail(user.Id)

	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h Handler) EditUserDetail(c *gin.Context) {
	var (
		req dto.EditUserDetailReq
		err error
	)

	userVal, _ := c.Get("user")
	user := userVal.(model.User)

	if err = c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(httperror.BadRequestError("Failed binding validation", "BINDING_FAILED"))
		return
	}

	response, err := h.userService.EditUserDetail(req, user.Id)

	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response)
}
