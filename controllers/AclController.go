package controllers

import (
	"../libs"
	"../models"
	"../settings"
	"crypto/md5"
	"fmt"
	"github.com/gin-gonic/gin"
)

// @ID AclController
// @Summary Categories Api Router
// @Description Description of the api
// @Success 200 {object} libs.ACLResponse "Show Roles"
// @Router /api/v1/acl/ [get]
// @tags Acl
func AclController(c *gin.Context) {
	user, err := models.AuthDecodeUser(c)
	if !err {
		libs.JSONResponseUnauthorized(c, "User no login", "User no login")
		return
	}
	var userId = int64(user.User.Id)
	acl := new(libs.ACL)
	acl.Init(userId)
	libs.JSONResponseOk(c, acl.GetData())
}

// @ID AclPackagesController
// @Summary Categories Api Router
// @Description Description of the api
// @Success 200 {array} models.Package "Show Roles"
// @Router /api/v1/acl/packages [get]
// @tags Acl
func AclPackagesController(c *gin.Context) {
	settings.InstanceDb = libs.GetConnection()
	pack, err := models.Packages(false)
	defer settings.InstanceDb.Close()
	if err != nil {
		libs.JSONResponseError(c, "No se ha podido recuperar los paquetes", err.Error())
		return
	}
	libs.JSONResponseOk(c, pack)
}

// @ID AclPackagesGroupController
// @Summary Categories Api Router
// @Description Description of the api
// @Success 200 {array} models.PackageGroup "Show Roles"
// @Router /api/v1/acl/packages-group/ [get]
// @tags Acl
func AclPackagesGroupController(c *gin.Context) {
	var business = c.Query("business")
	settings.InstanceDb = libs.GetConnection()
	pack, err := models.PackagesGroup(business != "")
	defer settings.InstanceDb.Close()
	if err != nil {
		libs.JSONResponseError(c, "No se ha podido recuperar los paquetes", err.Error())
		return
	}
	libs.JSONResponseOk(c, pack)
}

// @ID AclPackagesInvoiceController
// @Summary Role Create Api Router
// @Description Description of the api
// @Success 200 {object} models.AclPackageInvoiceSuccess "Show Roles"
// @Success 400 {object} libs.ErrorResponse "Error"
// @Success 404 {object} libs.ErrorResponse "Error"
// @Router /api/v1/acl/packages-invoice/ [post]
// @Accept  json
// @Produce  json
// @Param auth body models.AclPackageInvoice true "Add Role"
// @tags Acl
func AclPackagesInvoiceController(c *gin.Context) {
	aclRequest := models.AclPackageInvoice{}
	settings.InstanceDb = libs.GetConnection()
	defer settings.InstanceDb.Close()
	user, err := models.AuthDecodeUser(c)
	if !err {
		libs.JSONResponseUnauthorized(c, "User no login", "User no login")
		return
	}
	var userId = int64(user.User.Id)
	if err := c.ShouldBindJSON(&aclRequest); err != nil {
		libs.JSONResponseError(c, "Hay un error en los campos del formulario", err.Error())
		return
	}
	off, err1 := models.GetOfferPackageActive(aclRequest.Id)
	if err1 != nil {
		libs.JSONResponseError(c, "Error", err1.Error())
		return
	}
	name := libs.GenerateName(user.User.Name)
	data := []byte(name)
	name = fmt.Sprintf("%x", md5.Sum(data))
	errCreate := models.CreateOrder(userId, off.Id, aclRequest.Value, name)
	if errCreate != nil {
		libs.JSONResponseError(c, "Error", errCreate.Error())
		return
	}
	generate := models.AclPackageInvoiceSuccess{
		Invoice: name,
	}
	libs.JSONResponseOk(c, generate)
}

// @ID AclPackagesActivateController
// @Summary Role Create Api Router
// @Description Description of the api
// @Success 200 {array} AclResponse "Show Roles"
// @Success 400 {object} libs.ErrorResponse "Error"
// @Success 404 {object} libs.ErrorResponse "Error"
// @Router /api/v1/acl/packages-activate/ [post]
// @Accept  json
// @Produce  json
// @Param auth body models.AclPackageActivate true "Add Role"
// @tags Acl
func AclPackagesActivateController(c *gin.Context) {
	aclRequest := models.AclPackageActivate{}
	if err := c.ShouldBindJSON(&aclRequest); err != nil {
		libs.JSONResponseError(c, "Hay un error en los campos del formulario", err.Error())
		return
	}
	settings.InstanceDb = libs.GetConnection()
	defer settings.InstanceDb.Close()
	order, err := models.GetOrder(aclRequest.Invoice)
	if err != nil {
		libs.JSONResponseError(c, "Error no existe el pedido", err.Error())
		return
	}
	errorSubmit := models.ActivateOrder(order, aclRequest)

	if errorSubmit != nil {
		libs.JSONResponseError(c, "Error no se ha podido activar el pedido", errorSubmit.Error())
		return
	}
	errDelete := models.DeleteOrder(order)
	if errDelete != nil {
		libs.JSONResponseError(c, "Error, no se na podido eliminar el pedido", errDelete.Error())
		return
	}
	delUserPackErr := models.DeleteUserPackages(order)
	if delUserPackErr != nil {
		libs.JSONResponseError(c, "Error, no se na podido cambiar de paquete", delUserPackErr.Error())
		return
	}
	errChangeUser := models.ChangeUserPackage(order)
	if errChangeUser != nil {
		libs.JSONResponseError(c, "Error, no se na podido asignar el nuevo paquete", errChangeUser.Error())
		return
	}
	acl := new(libs.ACL)
	acl.Init(order.FkIdUser)
	libs.JSONResponseOk(c, acl.GetData())
}
