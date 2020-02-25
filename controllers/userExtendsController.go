package controllers

import (
	"../libs"
	"../models"
	"github.com/gin-gonic/gin"
	"log"
)

// @ID CategoriesController
// @Summary Categories Api Router
// @Description Description of the api
// @Success 200 {object} models.ItemExtendArray "Show Categories"
// @Router /api/v1/user-extends/items/ [get]
// @tags User Extends
func UserExtendsItemsController(c *gin.Context) {
	user, isLogin := models.AuthDecodeUser(c)
	if !isLogin {
		libs.JSONResponseUnauthorized(c, "El usuario no esta logeado", "El usuario ha caducado, no hay token...")
		return
	}
	filter := models.FilterPage{
		All:    true,
		IdUser: user.User.Id,
	}
	u, err := models.GetItemsExtend(filter)

	if err != nil {
		libs.JSONResponseNotFound(c, "El usuario no tiene anuncios", err.Error())
		return
	}
	libs.JSONResponseOk(c, u)
}

// @ID UserExtendsBlogsController
// @Summary Blogs Api Router
// @Description Description of the api
// @Success 200 {array} models.BlogExtend "Show Blogs"
// @Router /api/v1/user-extends/blogs/ [get]
// @tags User Extends
func UserExtendsBlogsController(c *gin.Context) {
	User, isLogin := models.AuthDecodeUser(c)
	if !isLogin {
		libs.JSONResponseUnauthorized(c, "El usuario esta caducado o no existe", "El usuario esta caducado o no existe")
		return
	}
	u, err := models.GetBlogsExtendUser(int64(User.User.Id))
	if err != nil {
		log.Fatal(err)
	}
	libs.JSONResponseOk(c, u)
}

// @ID UserExtendsConfirmationTokenController
// @Summary Blogs Api Router
// @Description Description of the api
// @Success 200 {object} models.UserConfirmEmailToken "Show Blogs"
// @Success 404 {object} models.UserConfirmEmailToken "Show Blogs"
// @Router /api/v1/user-extends/confirmation/email/token/{token}/ [get]
// @Param token path string true "Token"
// @tags User Extends
func UserExtendsConfirmationTokenController(c *gin.Context) {
	token := c.Param("token")
	data := models.UserConfirmEmailToken{
		Success: false,
		Token:   token,
	}
	err := models.UserSaveTokenConfirmEmailConfirmed(token)
	if err != nil {
		error := err.Error()
		data.Err = &error
		c.JSON(404, data)
		return
	}
	data.Success = true
	c.JSON(200, data)
}
