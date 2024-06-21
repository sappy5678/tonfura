package controllers

import (
	"net/http"

	"github.com/ebubekiryigit/golang-mongodb-rest-api-starter/models"
	"github.com/ebubekiryigit/golang-mongodb-rest-api-starter/services"
	"github.com/gin-gonic/gin"
)

// Reserve godoc
// @Summary      Reserve
// @Description  Reserve a coupon
// @Tags         coupon
// @Accept       json
// @Produce      json
// @Param		 userID	header	string						true "userID"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Router       /coupon/reserve [post]
func Reserve(c *gin.Context) {
	var request models.ReserveRequest
	_ = c.ShouldBindHeader(&request)

	response := &models.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	_, err := services.Reserve(request.UserID)
	if err != nil {
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	response.StatusCode = http.StatusCreated
	response.Success = true
	response.Data = gin.H{}
	response.SendResponse(c)
}

// Snatch godoc
// @Summary      Snatch
// @Description  Snatch a coupon
// @Tags         coupon
// @Accept       json
// @Produce      json
// @Param		 userID	header	string						true "userID"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Router       /coupon/snatch [post]
func Snatch(c *gin.Context) {
	var request models.SnatchRequest
	_ = c.ShouldBindHeader(&request)

	response := &models.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	_, err := services.Snatch(request.UserID)
	if err != nil {
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	response.StatusCode = http.StatusCreated
	response.Success = true
	response.Data = gin.H{}
	response.SendResponse(c)
}
