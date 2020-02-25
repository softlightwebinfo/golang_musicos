package routers

import (
	"../controllers"
	"github.com/gin-gonic/gin"
)

func AclRoute(router *gin.RouterGroup) {
	auth := router.Group("/acl")
	auth.GET("/", controllers.AclController)
	auth.GET("/packages/", controllers.AclPackagesController)
	auth.GET("/packages-group/", controllers.AclPackagesGroupController)
	auth.POST("/packages-invoice/", controllers.AclPackagesInvoiceController)
	auth.POST("/packages-activate/", controllers.AclPackagesActivateController)
}
