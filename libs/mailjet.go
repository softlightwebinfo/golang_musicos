package libs

import (
	"../settings"
	"github.com/mailjet/mailjet-apiv3-go"
)

type MailjetModel struct {
	config settings.MailjetConfig
	client *mailjet.Client
}

func (then *MailjetModel) Init() {
	config := settings.GetMailjetConfig()
	then.client = mailjet.NewMailjetClient(config.Public, config.Private)
}
func (then MailjetModel) Send(username string, email, subject string, variables map[string]interface{}, template int) (res *mailjet.ResultsV31, err error) {
	messagesInfo := []mailjet.InfoMessagesV31{
		mailjet.InfoMessagesV31{
			From: &mailjet.RecipientV31{
				Email: then.config.Email,
				Name:  then.config.Name,
			},
			To: &mailjet.RecipientsV31{
				mailjet.RecipientV31{
					Email: email,
					Name:  username,
				},
			},
			TemplateID:       template,
			TemplateLanguage: true,
			Subject:          subject,
			Variables:        variables,
		},
	}
	messages := mailjet.MessagesV31{Info: messagesInfo}
	res, err = then.client.SendMailV31(&messages)
	return
}
