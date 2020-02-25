package controllers

import (
	"../libs"
	"../models"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
)

// @ID UsersController
// @Summary Users Api Router
// @Description Description of the api
// @Success 200 {array} models.User "Show Users"
// @Router /api/v1/users/ [get]
// @Path "/users/"
// @Name usersAll
// @tags Users
func UsersController(c *gin.Context) {
	u, err := models.GetUsers()

	if err != nil {
		log.Fatal(err)
	}
	libs.JSONResponseOk(c, u)
}

// @ID UsersCreateController
// @Summary Users Create Api Router
// @Description Description of the api
// @Success 200 {object} models.User "Show Users"
// @Success 400 {string} string "error"
// @Router /api/v1/users/ [post]
// @Accept  json
// @Produce  json
// @Param user body models.UserCreated true "Add User account"
// @Path "/users/"
// @Name usersCreate
// @tags Users
func UsersCreateController(c *gin.Context) {
	var user models.UserCreated
	if err := c.ShouldBindJSON(&user); err != nil {
		libs.JSONResponseError(c, "Hay un error en los campos del formulario", err.Error())
		return
	}

	id, err := models.CreateUser(user)
	if err != nil {
		log.Fatal(err)
		libs.JSONResponseError(c, "No se ha podido crear el usuario", err.Error())
		return
	}
	us, errGetUser := models.GetUser(id)
	if errGetUser != nil {
		log.Fatal(errGetUser)
		libs.JSONResponseError(c, "No se ha podido encontrar el usuario registrado", errGetUser.Error())
		return
	}
	fmt.Println("Creado exitosamente: ", id)
	libs.JSONResponseOk(c, us)
}

// @ID UsersUpdateController
// @Summary Users Update Api Router
// @Description Description of the api
// @Success 200 {object} models.User "Show Users"
// @Success 400 {string} string "error"
// @Router /api/v1/users/{id}/ [put]
// @Accept  json
// @Produce  json
// @Param user body models.UserCreated true "Add User account"
// @Param id path int true "User id"
// @Path "/users/"
// @Name usersCreate
// @tags Users
func UsersUpdateController(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		libs.JSONResponseError(c, "Hay un error en los campos del formulario", err.Error())
		return
	}
	user.Id = id
	err := models.UpdateUser(user)
	if err != nil {
		libs.JSONResponseError(c, "No se ha podido crear el usuario", err.Error())
		return
	}
	if user.Password != "" {
		_ = models.ChangePassword(user.Email, user.Password)
	}
	us, errGetUser := models.GetUser(int64(user.Id))
	if errGetUser != nil {
		libs.JSONResponseError(c, "No se ha podido encontrar el usuario registrado", errGetUser.Error())
		return
	}
	fmt.Println("Creado exitosamente: ", id)
	libs.JSONResponseOk(c, us)
}

// @ID UsersGetController
// @Summary Users Api Router s
// @Description Description of the api
// @Success 200 {object} models.User "Show Users"
// @Router /api/v1/users/{id}/ [get]
// @Param id path int true "User id"
// @Path "/users/"
// @Name users
// @tags Users
func UsersGetController(c *gin.Context) {
	user := models.User{}
	id, _ := strconv.Atoi(c.Param("id"))
	user, e := models.GetUser(int64(id))
	if e != nil {
		libs.JSONResponseError(c, "El usuario no se ha encontrado", e.Error())
		return
	}
	libs.JSONResponseOk(c, user)
}

// @ID UsersDeleteController
// @Summary Users Api Router s
// @Description Description of the api
// @Success 200 {object} models.User "Show User"
// @Router /api/v1/users/{id}/ [delete]
// @Param id path int true "User id"
// @Name users
// @tags Users
func UsersDeleteController(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	user, e := models.GetUser(int64(id))
	errDelete := models.DeleteUser(int64(id))
	if errDelete != nil {
		libs.JSONResponseNotFound(c, "El usuario no se ha podido eliminar", e.Error())
		return
	}
	libs.JSONResponseOk(c, user)
}
