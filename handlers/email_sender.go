package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hex/phoenix/config"
	"html/template"
	"net/smtp"
	"reflect"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

type Email struct {
	To      string
	From    string
	Subject string
	Body    string
}

type EmailData struct {
	Lang           string
	RecipientName  string
	Subject        string
	WelcomeMessage string
	Greeting       string
	Intro          string
	Participation  string
	Expectations   string
	Benefit1       string
	Benefit2       string
	Benefit3       string
	Benefit4       string
	ReachOut       string
	WelcomeAboard  string
	BestRegards    string
	Founders       string
	FoundersTitle  string
	JoinUs         string
}

func SendEmail(email Email) error {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	_, err := bundle.LoadMessageFile("i18n/invite/en.json")

	if err != nil {
		fmt.Printf("[ERROR]: Failed to load message file - i18n: %+v\n", err)
		return err
	}

	localizer := i18n.NewLocalizer(bundle, "en")

	var localizationKeys = map[string]string{
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

	emailData := EmailData{
		Lang:          "en",
		RecipientName: email.To,
	}

	emailDataValue := reflect.ValueOf(&emailData).Elem()
	for field, key := range localizationKeys {
		localozedString := localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: key})
		emailDataValue.FieldByName(field).SetString(localozedString)
	}

	// email templating
	t, err := template.ParseFiles("templates/invite/invite.html")

	if err != nil {
		fmt.Printf("[ERROR] Failed to parse HTML template: %+v", err)
		return err
	}

	buf := new(bytes.Buffer)

	if err := t.Execute(buf, emailDataValue); err != nil {
		fmt.Printf("[ERROR] Failed to populate HTML template: %+v", err)
		return err
	}

	auth := smtp.PlainAuth("", config.SmtpUsername, config.SmtpPassword, config.SmtpServer)

	msg := []byte(
		"From: " + email.From + "\r\n" +
			"To: " + email.To + "\r\n" +
			"Subject: " + email.Subject + "\r\n" +
			"MIME-version: 1.0;\n" +
			"Content-Type: text/html; charset=\"UTF-8\";\n\n" +
			buf.String())

	smtpAddress := fmt.Sprintf("%s:%d", config.SmtpServer, config.SmtpPort)

	err = smtp.SendMail(smtpAddress, auth, email.From, []string{email.To}, msg)

	if err != nil {
		fmt.Printf("Failed to send email: %v\n", err)
		return err
	}

	fmt.Printf("Email sent successfully to %s\n", email.To)
	return nil
}
