package controllers

import (
	"../libs"
	"../models"
	"../settings"
	"github.com/gin-gonic/gin"
	"strconv"
)

// @ID AdsController
// @Summary Categories Api Router
// @Description Description of the api
// @Success 200 {object} libs.AdsResponseGetAll "Show Roles"
// @Router /api/v1/ads/ [get]
// @tags Ads
func AdsController(c *gin.Context) {
	user, err := models.AuthDecodeUser(c)
	if !err {
		libs.JSONResponseUnauthorized(c, "User no login", "User no login")
		return
	}
	var userId = int64(user.User.Id)
	//acl := new(libs.ACL)
	//acl.Init(userId)
	settings.InstanceDb = libs.GetConnection()
	ads, _ := models.GetAllADS(userId)
	defer settings.InstanceDb.Close()
	libs.JSONResponseOk(c, ads)
}

// @ID AdsNewFirstController
// @Summary Role Create Api Router
// @Description Description of the api
// @Success 200 {object} models.AclPackageInvoiceSuccess "Show Roles"
// @Success 400 {object} libs.ErrorResponse "Error"
// @Success 404 {object} libs.ErrorResponse "Error"
// @Router  /api/v1/ads/new/first/ [post]
// @Accept  json
// @Produce  json
// @Param auth body models.ADSResponseNewFirst true "Add Role"
// @tags Ads
func AdsNewFirstController(c *gin.Context) {
	first := models.ADSRequestNewFirst{}
	user, err := models.AuthDecodeUser(c)
	settings.InstanceDb = libs.GetConnection()
	defer settings.InstanceDb.Close()
	if !err {
		libs.JSONResponseUnauthorized(c, "User no login", "User no login")
		return
	}
	if err := c.ShouldBindJSON(&first); err != nil {
		libs.JSONResponseError(c, "Hay un error en los campos del formulario", err.Error())
		return
	}
	var userId = int64(user.User.Id)

	id, e := models.PostAdsNewFirst(userId, first)
	if e != nil {
		libs.JSONResponseError(c, "no se ha podido guardar la campaña", e.Error())
		return
	}
	response, er := models.GetADS(userId, id)
	if er != nil {
		libs.JSONResponseError(c, "no se ha podido obtener la campaña", er.Error())
		return
	}
	libs.JSONResponseOk(c, response)
}

// @ID AdsNewThreeController
// @Summary Role Create Api Router
// @Description Description of the api
// @Success 200 {object} models.AclPackageInvoiceSuccess "Show Roles"
// @Success 400 {object} libs.ErrorResponse "Error"
// @Success 404 {object} libs.ErrorResponse "Error"
// @Router  /api/v1/ads/new/three/ [post]
// @Accept  json
// @Produce  json
// @Param auth body models.ADSRequestNewThree true "Add Role"
// @tags Ads
func AdsNewThreeController(c *gin.Context) {
	first := models.ADSRequestNewThree{}
	user, err := models.AuthDecodeUser(c)
	settings.InstanceDb = libs.GetConnection()
	defer settings.InstanceDb.Close()
	if !err {
		libs.JSONResponseUnauthorized(c, "User no login", "User no login")
		return
	}
	if err := c.ShouldBindJSON(&first); err != nil {
		libs.JSONResponseError(c, "Hay un error en los campos del formulario", err.Error())
		return
	}
	var userId = int64(user.User.Id)

	models.AdsChangeStep(first.Id, userId, 4)
	e := models.AdsActiveStep(first.Id, first.Selected)
	if e != nil {
		libs.JSONResponseError(c, "No se puede publicar el anuncio", e.Error())
		return
	}

	response, er := models.GetADS(userId, first.Id)
	if er != nil {
		libs.JSONResponseError(c, "No se encuentra el anuncio", er.Error())
		return
	}
	libs.JSONResponseOk(c, response)
}

// @ID AdsNewSecondController
// @Summary Role Create Api Router
// @Description Description of the api
// @Success 200 {object} models.AclPackageInvoiceSuccess "Show Roles"
// @Success 400 {object} libs.ErrorResponse "Error"
// @Success 404 {object} libs.ErrorResponse "Error"
// @Router  /api/v1/ads/new/second/ [post]
// @Accept  json
// @Produce  json
// @Param auth body models.ADSRequestNewSecond true "Add Role"
// @tags Ads
func AdsNewSecondController(c *gin.Context) {
	first := models.ADSRequestNewSecond{}
	user, err := models.AuthDecodeUser(c)
	settings.InstanceDbTX, _ = libs.GetConnection().Begin()
	if !err {
		libs.JSONResponseUnauthorized(c, "User no login", "User no login")
		return
	}
	if err := c.ShouldBindJSON(&first); err != nil {
		libs.JSONResponseError(c, "Hay un error en los campos del formulario", err.Error())
		return
	}
	var userId = int64(user.User.Id)

	id, e := models.PostAdsNewSecond(userId, first)
	if e != nil {
		libs.JSONResponseError(c, "no se ha podido guardar la campaña", e.Error())
		return
	}
	settings.InstanceDb = libs.GetConnection()
	defer settings.InstanceDb.Close()
	response, er := models.GetADS(userId, id)
	if er != nil {
		libs.JSONResponseError(c, "No se encuentra el anuncio", er.Error())
		return
	}
	libs.JSONResponseOk(c, response)
}

// @ID AdsNewReviewController
// @Summary Role Create Api Router
// @Description Description of the api
// @Success 200 {object} models.AdsResponseGet "Show Roles"
// @Success 400 {object} libs.ErrorResponse "Error"
// @Success 404 {object} libs.ErrorResponse "Error"
// @Router  /api/v1/ads/new/second/ [post]
// @Accept  json
// @Produce  json
// @Param auth body models.ADSRequestNewReview true "Add Role"
// @tags Ads
func AdsNewReviewController(c *gin.Context) {
	first := models.ADSRequestNewReview{}
	user, err := models.AuthDecodeUser(c)
	settings.InstanceDb = libs.GetConnection()
	if !err {
		libs.JSONResponseUnauthorized(c, "User no login", "User no login")
		return
	}
	if err := c.ShouldBindJSON(&first); err != nil {
		libs.JSONResponseError(c, "Hay un error en los campos del formulario", err.Error())
		return
	}
	var userId = int64(user.User.Id)
	e := models.PostAdsNewReview(userId, first)
	if e != nil {
		libs.JSONResponseError(c, "No se ha podido guardar el anuncio", e.Error())
		return
	}
	e2 := models.PostAdsNewReviewActive(first.Id, true)
	if e2 != nil {
		libs.JSONResponseError(c, "No se ha podido activar el anuncio", e2.Error())
		return
	}
	e3 := models.PostAdsNewReviewActiveAds(first.Id, true)
	if e3 != nil {
		libs.JSONResponseError(c, "No se ha podido activar el anuncio", e3.Error())
		return
	}
	response, er := models.GetADS(userId, first.Id)
	if er != nil {
		libs.JSONResponseError(c, "No se ha podido encontrar el anuncio", er.Error())
		return
	}
	libs.JSONResponseOk(c, response)
}

// @ID AdsDeleteController
// @Summary Blogs Api Router s
// @Description Description of the api
// @Success 200 {object} models.Blog "Show Item"
// @Router /api/v1/ads/delete/{id}/ [delete]
// @Param id path int true "Item id"
// @tags ads
func AdsDeleteController(c *gin.Context) {
	user, err := models.AuthDecodeUser(c)
	id, _ := strconv.Atoi(c.Param("id"))
	settings.InstanceDb = libs.GetConnection()
	defer settings.InstanceDb.Close()
	if !err {
		libs.JSONResponseUnauthorized(c, "User no login", "User no login")
		return
	}
	errDelete := models.DeleteAds(int64(id), int64(user.User.Id))
	if errDelete != nil {
		libs.JSONResponseNotFound(c, "El ads no se ha podido eliminar", errDelete.Error())
		return
	}
	libs.JSONResponseOk(c, models.DeleteResponse{
		Id: id,
	})
}
