package controllers

import (
	"../libs"
	"../models"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
)

// @ID AuthController
// @Summary Categories Api Router
// @Description Description of the api
// @Success 200 {array} models.User "Show Roles"
// @Router /api/v1/auth/ [get]
// @tags Auth
func AuthController(c *gin.Context) {
	u, err := models.GetUsers()

	if err != nil {
		log.Fatal(err)
	}
	libs.JSONResponseOk(c, u)
}

// @ID AuthSignupController
// @Summary Role Create Api Router
// @Description Description of the api
// @Success 200 {object} models.AuthUser "Show Roles"
// @Success 400 {object} libs.ErrorResponse "Error"
// @Success 404 {object} libs.ErrorResponse "Error"
// @Router /api/v1/auth/signup/ [post]
// @Accept  json
// @Produce  json
// @Param auth body models.UserRegister true "Add Role"
// @tags Auth
func AuthSignupController(c *gin.Context) {
	var user models.UserRegister
	//var email = "rafael.gonzalez.1737@gmail.com"
	emailSend := models.EmailModel{}
	var admin = c.Query("admin")

	if err := c.ShouldBindJSON(&user); err != nil {
		libs.JSONResponseError(c, "Hay un error en los campos del formulario", err.Error())
		return
	}
	_, errIsEmailExist := models.IsExistUserEmail(user.Email)
	if errIsEmailExist != nil {
		libs.JSONResponseError(c, "El correo electronico ya existe", errIsEmailExist.Error())
		return
	}
	// Create user
	id, err := models.CreateUser(models.UserCreated{
		Password:        user.Password,
		UserDescription: user.UserDescription,
	})
	if err != nil {
		libs.JSONResponseError(c, "No se ha podido crear el usuario", err.Error())
		return
	}
	if user.FkRoleId == 7 {
		errCreateBusiness := models.CreateBusiness(id, user)
		if errCreateBusiness != nil {
			libs.JSONResponseError(c, "No se ha podido crear la empresa", errCreateBusiness.Error())
			return
		}
	}

	// Get user
	us, errGetUser := models.GetUser(id)
	if errGetUser != nil {
		libs.JSONResponseError(c, "No se ha podido encontrar el usuario registrado", errGetUser.Error())
		return
	}
	auth := models.AuthModel{}
	// Create Token
	token, err := auth.CreateToken(us)
	if err != nil {
		libs.JSONResponseError(c, "Error creacion token", err.Error())
		return
	}
	// Auth User
	response := models.AuthUser{
		User:  us,
		Token: token,
	}
	// Generate Token confirm send email
	tokenEmailConfirm := models.GenerateTokenUserEmail(models.UserCreated{
		Password:        user.Password,
		UserDescription: user.UserDescription,
	})
	// Send email to user
	if admin != "true" {
		emailSend.Register(
			c,
			user.Name,
			user.Email,
			tokenEmailConfirm,
		)
		errSaveToken := models.UserSaveTokenConfirmEmail(int64(us.Id), tokenEmailConfirm, true)
		if errSaveToken != nil {
			libs.JSONResponseError(c, "Error, no se ha podido guardar el token", err)
			return
		}
	}
	libs.JSONResponseOk(c, response)
}

// @ID AuthSigninController
// @Summary Role Create Api Router
// @Description Description of the api
// @Success 200 {object} models.AuthUser "Show Roles"
// @Success 400 {object} libs.ErrorResponse "Error"
// @Success 404 {object} libs.ErrorResponse "Error"
// @Router /api/v1/auth/signin/ [post]
// @Accept  json
// @Produce  json
// @Param user body models.AutCredentials true "Add Role"
// @tags Auth
func AuthSigninController(c *gin.Context) {
	var user models.AutCredentials
	if err := c.ShouldBindJSON(&user); err != nil {
		libs.JSONResponseError(c, "Hay un error en los campos del formulario", err.Error())
		return
	}
	userLogin, err := models.AuthLogin(user)

	if err != nil {
		libs.JSONResponseNotFound(c, "El email y/o contraseña no son validas", err.Error())
		return
	}
	model := models.AuthModel{}
	userToken, err := model.CreateToken(userLogin)
	if err != nil {
		libs.JSONResponseNotFound(c, "Error en la creacion del token de sesión", err.Error())
		return
	}
	response := models.AuthUser{
		User:  userLogin,
		Token: userToken,
	}
	_ = models.UserLog(int64(userLogin.Id))
	libs.JSONResponseOk(c, response)
}

// @ID AuthVerifyController
// @Summary Role Create Api Router
// @Description Description of the api
// @Success 200 {object} models.AuthUser "Show Roles"
// @Success 400 {object} libs.ErrorResponse "Error"
// @Success 404 {object} libs.ErrorResponse "Error"
// @Router /api/v1/auth/verify/ [get]
// @Accept  json
// @Produce  json
// @tags Auth
func AuthVerifyController(c *gin.Context) {
	user, isLogin := models.AuthDecodeUser(c)
	if !isLogin {
		libs.JSONResponseNotFound(c,
			"Error: no tienes permisos para acceder",
			"El token no se ha enviado, caducado...",
		)
		return
	}

	libs.JSONResponseOk(c, user)
}

// @ID AuthUserController
// @Summary Role Create Api Router
// @Description Description of the api
// @Success 200 {object} models.AuthUser "Show Roles"
// @Success 400 {object} libs.ErrorResponse "Error"
// @Success 404 {object} libs.ErrorResponse "Error"
// @Router /api/v1/auth/user/ [get]
// @Accept  json
// @Produce  json
// @tags Auth
func AuthUserController(c *gin.Context) {
	user, isLogin := models.AuthDecodeUser(c)
	if !isLogin {
		libs.JSONResponseNotFound(c,
			"Error: no tienes permisos para acceder",
			"El token no se ha enviado, caducado...",
		)
		return
	}
	us, err := models.GetUser(int64(user.User.Id))
	if err != nil {
		libs.JSONResponseNotFound(c, "El usuario no existe", "El id del usuario no existe en la base de datos")
		return
	}
	auth := models.AuthModel{}
	token, err := auth.CreateToken(us)
	if err != nil {
		libs.JSONResponseError(c, "Error creacion token", err.Error())
		return
	}
	response := models.AuthUser{
		User:  us,
		Token: token,
	}
	_ = models.UserLog(int64(us.Id))
	libs.JSONResponseOk(c, response)
}

// @ID AuthGenerateTokenController
// @Summary Role Create Api Router
// @Description Description of the api
// @Success 200 {object} models.AuthUser "Show Roles"
// @Success 400 {object} libs.ErrorResponse "Error"
// @Success 404 {object} libs.ErrorResponse "Error"
// @Router /api/v1/auth/generate-token/{id}/ [get]
// @Param id path int true "Item id"
// @Accept  json
// @Produce  json
// @tags Auth
func AuthGenerateTokenController(c *gin.Context) {
	//user, _ := models.AuthDecodeUser(c)
	id, _ := strconv.Atoi(c.Param("id"))
	//if user.User.Id != 1 {
	//	libs.JSONResponseNotFound(c,
	//		"Error: no tienes permisos para acceder",
	//		"El token no se ha enviado, caducado...",
	//	)
	//	return
	//}
	us, err := models.GetUser(int64(id))
	if err != nil {
		libs.JSONResponseNotFound(c, "El usuario no existe", "El id del usuario no existe en la base de datos")
		return
	}
	auth := models.AuthModel{}
	token, err := auth.CreateToken(us)
	if err != nil {
		libs.JSONResponseError(c, "Error creacion token", err.Error())
		return
	}
	response := models.AuthUser{
		User:  us,
		Token: token,
	}
	_ = models.UserLog(int64(us.Id))
	libs.JSONResponseOk(c, response)
}

// @ID AuthDecodeController
// @Summary Role Create Api Router
// @Description Description of the api
// @Success 200 {object} models.Role "Show Roles"
// @Success 400 {object} libs.ErrorResponse "Error"
// @Success 404 {object} libs.ErrorResponse "Error"
// @Router /api/v1/auth/decode/ [post]
// @Accept  json
// @Produce  json
// @Param user body models.RoleCreated true "Add Role"
// @tags Auth
func AuthDecodeController(c *gin.Context) {
	libs.JSONResponseOk(c, nil)
}

// @ID AuthEncodeController
// @Summary Role Create Api Router
// @Description Description of the api
// @Success 200 {object} models.AuthUser "Show Roles"
// @Success 400 {object} libs.ErrorResponse "Error"
// @Success 404 {object} libs.ErrorResponse "Error"
// @Router /api/v1/auth/encode/ [post]
// @Accept  json
// @Produce  json
// @Param user body models.UserCreated true "Add Role"
// @tags Auth
func AuthEncodeController(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		libs.JSONResponseError(c, "Hay un error en los campos del formulario", err.Error())
		return
	}

	auth := models.AuthModel{}
	token, err := auth.CreateToken(user)
	if err != nil {
		libs.JSONResponseError(c, "Error creacion token", err.Error())
		return
	}
	response := models.AuthUser{
		User:  user,
		Token: token,
	}

	libs.JSONResponseOk(c, response)
}

// @ID AuthRecoveryController
// @Summary Role Create Api Router
// @Description Description of the api
// @Success 200 {boolean} true "Show Roles"
// @Success 400 {object} libs.ErrorResponse "Error"
// @Success 404 {object} libs.ErrorResponse "Error"
// @Router /api/v1/auth/recovery/ [post]
// @Accept  json
// @Produce  json
// @Param user body models.UserRecovery true "Add Role"
// @tags Auth
func AuthRecoveryController(c *gin.Context) {
	var user models.UserRecovery
	emailSend := models.EmailModel{}
	if err := c.ShouldBindJSON(&user); err != nil {
		libs.JSONResponseError(c, "Hay un error en los campos del formulario", err.Error())
		return
	}
	tokenEmailConfirm := models.GenerateTokenEmail(user.Email)
	emailSend.Recovery(
		c,
		user.Email,
		tokenEmailConfirm,
	)
	err := models.SetUserRecovery(user.Email, tokenEmailConfirm);
	libs.JSONResponseOk(c, err == nil)
}

// @ID AuthRecoveryEmailController
// @Summary Role Create Api Router
// @Description Description of the api
// @Success 200 {boolean} true "Show Roles"
// @Success 400 {object} libs.ErrorResponse "Error"
// @Success 404 {object} libs.ErrorResponse "Error"
// @Router /api/v1/auth/recovery-email/ [post]
// @Accept  json
// @Produce  json
// @Param user body models.UserRecoveryEmail true "Add Role"
// @tags Auth
func AuthRecoveryEmailController(c *gin.Context) {
	var user models.UserRecoveryEmail
	if err := c.ShouldBindJSON(&user); err != nil {
		libs.JSONResponseError(c, "Hay un error en los campos del formulario", err.Error())
		return
	}
	email, _ := models.DeleteUserRecovery(user.Token)
	err := models.ChangePassword(email, user.Password)
	libs.JSONResponseOk(c, err == nil)
}
