package controllers

import (
	"../libs"
	"github.com/gin-gonic/gin"
)
// ShowAccount godoc
// @ID HomeController
// @Summary Home Api Router
// @Description Description of the api
// @Success 200 {object} libs.HelloWorldResponse "Show messages and provided name"
// @Router /api/v1/ [get]
// @Tags Home
func HomeController(c *gin.Context) {
	h := libs.HelloWorldResponse{
		Message:      "Welcome to restful",
		ProvidedName: "Home",
	}
	libs.JSONResponseOk(c, h)
}
