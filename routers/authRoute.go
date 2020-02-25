package routers

import (
	"../controllers"
	"github.com/gin-gonic/gin"
)

func AuthRoute(router *gin.RouterGroup) {
	auth := router.Group("/auth")
	auth.GET("/", controllers.AuthController)
	auth.POST("/signup/", controllers.AuthSignupController)
	auth.POST("/signin/", controllers.AuthSigninController)
	auth.GET("/verify/", controllers.AuthVerifyController)
	auth.GET("/user/", controllers.AuthUserController)
	auth.POST("/decode/", controllers.AuthDecodeController)
	auth.POST("/encode/", controllers.AuthEncodeController)
	auth.POST("/recovery/", controllers.AuthRecoveryController)
	auth.POST("/recovery-email/", controllers.AuthRecoveryEmailController)
	auth.GET("/generate-token/:id/", controllers.AuthGenerateTokenController)
}
