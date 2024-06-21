package validators

import (
	"net/http"

	"github.com/ebubekiryigit/golang-mongodb-rest-api-starter/models"
	"github.com/gin-gonic/gin"
)

func ReserveValidator() gin.HandlerFunc {
	return func(c *gin.Context) {

		var request models.ReserveRequest
		_ = c.ShouldBindHeader(&request)

		if err := request.Validate(); err != nil {
			models.SendErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		c.Next()
	}
}

func SnatchValidator() gin.HandlerFunc {
	return func(c *gin.Context) {

		var request models.SnatchRequest
		_ = c.ShouldBindHeader(&request)

		if err := request.Validate(); err != nil {
			models.SendErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		c.Next()
	}
}
