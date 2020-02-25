package libs

import (
	"../settings"
	"bytes"
	"gopkg.in/gomail.v2"
	"html/template"
)

/**
templateData := struct {
	Name string
	URL  string
}{
	Name: "Dhanush",
	URL:  "http://geektrust.in",
}
r := NewRequest([]string{"junk@junk.com"}, "Hello Junk!", "Hello, World!")
err := r.ParseTemplate("template.html", templateData)
if err := r.ParseTemplate("template.html", templateData); err == nil {
	ok, _ := r.SendEmail()
	fmt.Println(ok)
}
*/

type EmailRequest struct {
	from    string
	to      string
	subject string
	body    string
	config  settings.SmtpEmail
}

func NewRequest(to string, subject, body string) *EmailRequest {
	config := settings.GetSmtEmailConfig()
	return &EmailRequest{
		to:      to,
		subject: subject,
		body:    body,
		config:  config,
		from:    config.From,
	}
}

func (r *EmailRequest) SendEmail() (bool, error) {
	m := gomail.NewMessage()
	m.SetHeader("From", r.from)
	m.SetHeader("To", r.to)
	m.SetHeader("Subject", r.subject)
	m.SetBody("text/html", r.body)

	// Send the email to Bob
	d := gomail.NewDialer(r.config.Host, r.config.Port, r.config.Email, r.config.Password)
	if err := d.DialAndSend(m); err != nil {
		panic(err)
		return false, err
	}
	return true, nil
}

func (r *EmailRequest) ParseTemplate(templateFileName string, data interface{}) error {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return err
	}
	r.body = buf.String()
	return nil
}
