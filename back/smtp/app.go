package smtp

import (
	"app/config"
	"app/log"
	"app/request"
	"bytes"
	"fmt"
	"html/template"
	"time"

	"gopkg.in/gomail.v2"
)

// Типы кодов подтверждения
type ConfirmationCodeType int

const (
	RegistrationCode  ConfirmationCodeType = iota // Код подтверждения при регистрации
	LoginCode                                     // Код для входа
	PasswordResetCode                             // Код для восстановления пароля
)

// SMTPClient структура для отправки писем через SMTP
type SMTPClient struct {
	Host     string
	Port     int
	Username string
	Password string

	RequestHandler *request.RequestHandler
	EmailTemplate  *template.Template // Поле для хранения кэшированного шаблона
}

// NewSMTPClient создает новый экземпляр SMTPClient
func NewSMTPClient(host string, port int, username, password string, repetPause time.Duration) (*SMTPClient, error) {
	requestHandler, err := request.NewRequestHandler()
	if err != nil {
		log.App.Error(err)
	}

	go requestHandler.ProcessRequests(repetPause)

	// Парсим шаблон и сохраняем его в структуре
	emailTemplate, err := template.ParseFiles("smtp/email_template.html")
	if err != nil {
		return nil, fmt.Errorf("ошибка при загрузке шаблона: %w", err)
	}

	return &SMTPClient{
		Host:           host,
		Port:           port,
		Username:       username,
		Password:       password,
		RequestHandler: requestHandler,
		EmailTemplate:  emailTemplate, // Сохраняем шаблон в структуре
	}, nil
}

// SendEmail отправляет письмо с использованием SMTP
func (client *SMTPClient) SendEmail(from, to, subject, body string) {
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(client.Host, client.Port, client.Username, client.Password)

	client.RequestHandler.HandleRequest(func() error {
		return d.DialAndSend(m)
	})
}

// SendConfirmationCodeEmail отправляет письмо с кодом подтверждения
func (client *SMTPClient) SendConfirmationCodeEmail(to string, confirmationCode string, codeType ConfirmationCodeType) error {
	from := client.Username
	subject := "Код подтверждения"

	switch codeType {
	case RegistrationCode:
		subject = "Код подтверждения регистрации"
	case LoginCode:
		subject = "Код для входа"
	case PasswordResetCode:
		subject = "Код для восстановления пароля"
	}

	var body bytes.Buffer
	data := struct {
		ConfirmationCode string
		SiteURL          string
	}{
		ConfirmationCode: confirmationCode,
		SiteURL:          config.File.WebConfig.APPURL,
	}

	// Используем кэшированный шаблон
	if err := client.EmailTemplate.Execute(&body, data); err != nil {
		return fmt.Errorf("ошибка при выполнении шаблона: %w", err)
	}

	client.SendEmail(from, to, subject, body.String())

	return nil
}
