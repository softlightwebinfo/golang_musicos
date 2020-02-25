package routers

import (
	"../controllers"
	"github.com/gin-gonic/gin"
)

func UserExtendsRoute(router *gin.RouterGroup) {
	users := router.Group("/user-extends")
	users.GET("/items/", controllers.UserExtendsItemsController)
	users.GET("/blogs/", controllers.UserExtendsBlogsController)
	confirmation := users.Group("/confirmation")
	confirmation.GET("/email/token/:token", controllers.UserExtendsConfirmationTokenController)
}
