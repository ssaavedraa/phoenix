package email_factory

import "github.com/nicksnyder/go-i18n/v2/i18n"

type Email struct {
	ReceiverAddress string
	SenderAddress   string
	Subject         string
	TemplateName    string
	Locale          string
	TemplateData    map[string]string
}

type EmailFactory interface {
	Send(email Email) error
	authenticate()
	getEmailBody(templateName, locale string, templateData map[string]string) (string, error)
	getLocalizationKeys(templateName string) (map[string]string, error)
	getI18nLocalizer(locale, templateName string) (*i18n.Localizer, error)
}
