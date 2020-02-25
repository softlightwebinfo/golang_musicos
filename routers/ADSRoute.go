package routers

import (
	"../controllers"
	"github.com/gin-gonic/gin"
)

func ADSRoute(router *gin.RouterGroup) {
	auth := router.Group("/ads")
	auth.GET("/", controllers.AdsController)
	auth.POST("/new/first/", controllers.AdsNewFirstController)
	auth.PUT("/new/second/", controllers.AdsNewSecondController)
	auth.PUT("/new/review/", controllers.AdsNewReviewController)
	auth.POST("/new/three/", controllers.AdsNewThreeController)
	auth.DELETE("/delete/:id/", controllers.AdsDeleteController)
}
