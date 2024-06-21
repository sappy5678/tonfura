package controllers

import (
	"net/http"
	"strings"

	"github.com/ebubekiryigit/golang-mongodb-rest-api-starter/models"
	db "github.com/ebubekiryigit/golang-mongodb-rest-api-starter/models/db"
	"github.com/ebubekiryigit/golang-mongodb-rest-api-starter/services"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"golang.org/x/crypto/bcrypt"
)

// Register godoc
// @Summary      Register
// @Description  registers a user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        req  body      models.RegisterRequest true "Register Request"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Router       /auth/register [post]
func Register(c *gin.Context) {
	var request models.RegisterRequest
	_ = c.ShouldBindBodyWith(&request, binding.JSON)

	response := &models.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	// is email in use
	err := services.CheckUserMail(request.Email)
	if err != nil {
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	// create user record
	request.Name = strings.TrimSpace(request.Name)
	user, err := services.CreateUser(request.Name, request.Email, request.Password)
	if err != nil {
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	// generate access tokens
	accessToken, refreshToken, err := services.GenerateAccessTokens(user)
	if err != nil {
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	response.StatusCode = http.StatusCreated
	response.Success = true
	response.Data = gin.H{
		"user": user,
		"token": gin.H{
			"access":  accessToken.GetResponseJson(),
			"refresh": refreshToken.GetResponseJson()},
	}
	response.SendResponse(c)
}

// Login godoc
// @Summary      Login
// @Description  login a user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        req  body      models.LoginRequest true "Login Request"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Router       /auth/login [post]
func Login(c *gin.Context) {
	var request models.LoginRequest
	_ = c.ShouldBindBodyWith(&request, binding.JSON)

	response := &models.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	// get user by email
	user, err := services.FindUserByEmail(request.Email)
	if err != nil {
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	// check hashed password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		response.Message = "email and password don't match"
		response.SendResponse(c)
		return
	}

	// generate new access tokens
	accessToken, refreshToken, err := services.GenerateAccessTokens(user)
	if err != nil {
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	response.StatusCode = http.StatusOK
	response.Success = true
	response.Data = gin.H{
		"user": user,
		"token": gin.H{
			"access":  accessToken.GetResponseJson(),
			"refresh": refreshToken.GetResponseJson()},
	}
	response.SendResponse(c)
}

// Refresh godoc
// @Summary      Refresh
// @Description  refreshes a user token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        req  body      models.RefreshRequest true "Refresh Request"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Router       /auth/refresh [post]
func Refresh(c *gin.Context) {
	var request models.RefreshRequest
	_ = c.ShouldBindBodyWith(&request, binding.JSON)

	response := &models.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	// check token validity
	token, err := services.VerifyToken(request.Token, db.TokenTypeRefresh)
	if err != nil {
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	user, err := services.FindUserById(token.User)
	if err != nil {
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	// delete old token
	err = services.DeleteTokenById(token.ID)
	if err != nil {
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	accessToken, refreshToken, err := services.GenerateAccessTokens(user)
	response.StatusCode = http.StatusOK
	response.Success = true
	response.Data = gin.H{
		"user": user,
		"token": gin.H{
			"access":  accessToken.GetResponseJson(),
			"refresh": refreshToken.GetResponseJson()},
	}
	response.SendResponse(c)
}
