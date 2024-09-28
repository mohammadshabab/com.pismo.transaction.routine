package apiroutes

import (
	"com.pismo.transaction.routine/controllers"
	"github.com/gin-gonic/gin"
)

func AppRoutes(route *gin.Engine) {

	route.POST("/accounts", controllers.CreateAccount)
	route.GET("/accounts/:accountId", controllers.GetAccount)
	route.POST("/transactions", controllers.CreateTransaction)

}
