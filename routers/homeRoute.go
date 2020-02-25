package routers

import (
	"../controllers"
	"github.com/gin-gonic/gin"
)

func HomeRoute(router *gin.RouterGroup) {
	router.GET("/", controllers.HomeController)
}
