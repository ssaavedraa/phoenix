package email_factory

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"hex/phoenix/config"
	"html/template"
	"log"
	"net/smtp"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var auth smtp.Auth

type EmailFactoryImpl struct{}

func NewEmail() EmailFactory {
	return &EmailFactoryImpl{}
}

func (ef *EmailFactoryImpl) Send(
	email Email) error {
	ef.authenticate()

	body, err := ef.getEmailBody(email.TemplateName, email.Locale, email.TemplateData)

	if err != nil {
		fmt.Printf("[ERROR]: Failed to send email: %+v\n", err)
		return err
	}

	msg := []byte(
		"From: " + email.SenderAddress + "\r\n" +
			"To: " + email.ReceiverAddress + "\r\n" +
			"Subject: " + email.Subject + "\r\n" +
			"MIME-version: 1.0;\n" +
			"Content-Type: text/html; charset=\"UTF-8\";\n\n" +
			body)

	smtpAddress := fmt.Sprintf("%s:%d", config.SmtpServer, config.SmtpPort)
	err = smtp.SendMail(
		smtpAddress,
		auth,
		email.SenderAddress,
		[]string{email.ReceiverAddress},
		msg,
	)

	if err != nil {
		fmt.Printf("[ERROR]: Failed to send email: %+v\n", err)
		return err
	}

	fmt.Printf("Email sent successfully to %s\n", email.ReceiverAddress)
	return nil
}

func (ef *EmailFactoryImpl) authenticate() {
	if auth == nil {
		auth = smtp.PlainAuth("", config.SmtpUsername, config.SmtpPassword, config.SmtpServer)
	}
}

func (ef *EmailFactoryImpl) getEmailBody(
	templateName, locale string,
	templateData map[string]string,
) (string, error) {
	localizationKeys, err := ef.getLocalizationKeys(templateName)

	if err != nil {
		fmt.Printf("[ERROR]: Failed to get localization keys: %+v", err)
		return "", err
	}

	localizer, err := ef.getI18nLocalizer(locale, templateName)

	if err != nil {
		fmt.Printf("[ERROR]: Failed to get i18n localizer: %v\n", err)
		return "", err
	}

	templateValues := make(map[string]string)

	for field, key := range templateData {
		templateValues[field] = key
	}

	for field, key := range localizationKeys {
		localizedString := localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: key})

		templateValues[field] = localizedString
	}

	templatePath := fmt.Sprintf("templates/%s/index.html", templateName)
	t, err := template.ParseFiles(templatePath)

	if err != nil {
		fmt.Printf("[ERROR] Failed to parse HTML template: %+v", err)
		return "", err
	}

	buf := new(bytes.Buffer)

	if err := t.Execute(buf, templateValues); err != nil {
		fmt.Printf("[ERROR] Failed to populate HTML template: %+v", err)
		return "", err
	}

	return buf.String(), nil
}

func (ef *EmailFactoryImpl) getLocalizationKeys(templateName string) (map[string]string, error) {
	switch templateName {
	case "user_invite_mvp":
		templateFields := map[string]string{
			"Subject":        "WelcomeSubject",
			"WelcomeMessage": "WelcomeMessage",
			"Greeting":       "Greeting",
			"Intro":          "Intro",
			"Participation":  "Participation",
			"Expectations":   "Expectations",
			"Benefit1":       "Benefit1",
			"Benefit2":       "Benefit2",
			"Benefit3":       "Benefit3",
			"Benefit4":       "Benefit4",
			"ReachOut":       "ReachOut",
			"WelcomeAboard":  "WelcomeAboard",
			"BestRegards":    "BestRegards",
			"Founders":       "Founders",
			"FoundersTitle":  "FoundersTitle",
			"JoinUs":         "JoinUs",
		}

		return templateFields, nil

	default:
		return map[string]string{}, errors.New("invalid email template")
	}
}

func (ef *EmailFactoryImpl) getI18nLocalizer(locale, templateName string) (*i18n.Localizer, error) {
	var i18nLocale language.Tag

	if locale == "es" {
		i18nLocale = language.Spanish
	} else {
		i18nLocale = language.English
	}

	log.Printf("locale: %v", i18nLocale)

	bundle := i18n.NewBundle(i18nLocale)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	i18nFilePath := fmt.Sprintf("i18n/%s/%s.json", templateName, i18nLocale)
	_, err := bundle.LoadMessageFile(i18nFilePath)

	if err != nil {
		fmt.Printf("[ERROR]: Failed to load message file - i18n: %+v\n", err)
		return &i18n.Localizer{}, err
	}

	localizer := i18n.NewLocalizer(bundle, locale)

	return localizer, nil
}
