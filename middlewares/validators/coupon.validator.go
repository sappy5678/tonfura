package validators

import (
	"net/http"

	"github.com/ebubekiryigit/golang-mongodb-rest-api-starter/models"
	"github.com/gin-gonic/gin"
)

func ReserveValidator() gin.HandlerFunc {
	return func(c *gin.Context) {

		var registerRequest models.CouponReserveRequest
		_ = c.ShouldBindHeader(&registerRequest)

		if err := registerRequest.Validate(); err != nil {
			models.SendErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		c.Next()
	}
}
