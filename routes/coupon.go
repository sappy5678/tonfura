package routes

import (
	"github.com/ebubekiryigit/golang-mongodb-rest-api-starter/controllers"
	"github.com/ebubekiryigit/golang-mongodb-rest-api-starter/middlewares/validators"
	"github.com/gin-gonic/gin"
)

func CouponRoute(router *gin.RouterGroup) {
	coupon := router.Group("/coupon")
	{
		coupon.POST(
			"/reserve",
			validators.ReserveValidator(),
			controllers.Reserve,
		)

		// TODO implement get coupon

	}
}
