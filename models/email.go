package models

import (
	"../libs"
	"fmt"
	"github.com/gin-gonic/gin"
)

func EmailRegisterModel(to string) (ok bool, e error) {
	templateData := struct {
		Name string
		URL  string
	}{
		Name: "Dhanush",
		URL:  "http://geektrust.in",
	}
	r := libs.NewRequest(to, "Hello Junk!", "Hello, World!")
	if err := r.ParseTemplate("templates/emails/register.html", templateData); err == nil {
		ok, e = r.SendEmail()
	}
	return
}

type EmailModel struct {
}

/*
Send email Register
*/
func (then EmailModel) Register(c *gin.Context, username, email, token string) {
	model := libs.MailjetModel{}
	model.Init()
	_, err := model.Send(
		username,
		email,
		fmt.Sprintf("Bienvenid@ %s", username),
		map[string]interface{}{
			"username":          username,
			"confirmation_link": fmt.Sprintf("https://www.musicosdelmundo.com/email/confirmation/%s", token),
		},
		980883,
	)

	if err != nil {
		libs.JSONResponseError(c, "Error, no se ha podido enviar el el email de bienvenida", err)
		return
	}
}

/*
RECOVERY email Register
*/
func (then EmailModel) Recovery(c *gin.Context, email, token string) {
	model := libs.MailjetModel{}
	model.Init()
	_, err := model.Send(
		email,
		email,
		fmt.Sprintf("Bienvenid@ %s", email),
		map[string]interface{}{
			"username":          email,
			"confirmation_link": fmt.Sprintf("https://www.musicosdelmundo.com/email/recovery/%s", token),
		},
		1084385,
	)

	if err != nil {
		libs.JSONResponseError(c, "Error, no se ha podido enviar el el email de recuperación", err)
		return
	}
}
