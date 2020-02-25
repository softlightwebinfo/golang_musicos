package routers

import (
	"../controllers"
	"github.com/gin-gonic/gin"
)

func UsersRoute(router *gin.RouterGroup) {
	users := router.Group("/users")
	users.GET("/", controllers.UsersController)
	users.POST("/", controllers.UsersCreateController)
	users.PUT("/:id/", controllers.UsersUpdateController)
	users.GET("/:id/", controllers.UsersGetController)
	users.DELETE("/:id/", controllers.UsersDeleteController)
}
