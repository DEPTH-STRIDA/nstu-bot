package smtp

import (
	"app/internal/logger"
	"app/internal/model"
	"crypto/tls"
	"fmt"
	"net/smtp"
)

// SendAuthorizationCode отправляет код авторизации на указанный email
func SendAuthorizationCode(email, code string) error {
	conf := model.ConfigFile.SMTPConfig

	logger.Log.Info("Отправляем код на ", email)

	// Получаем настройки из окружения
	smtpMail := conf.Mail
	smtpPassword := conf.Password
	smtpHost := conf.Host
	smtpPort := conf.Port

	// Настройка TLS
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         smtpHost,
	}

	// Подключение
	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", smtpHost, smtpPort), tlsConfig)
	if err != nil {
		return fmt.Errorf("connection error: %v", err)
	}
	defer conn.Close()

	// Создание клиента
	client, err := smtp.NewClient(conn, smtpHost)
	if err != nil {
		return fmt.Errorf("client creation error: %v", err)
	}
	defer client.Close()

	// Аутентификация
	auth := smtp.PlainAuth("", smtpMail, smtpPassword, smtpHost)
	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("authentication error: %v", err)
	}

	// Формирование сообщения
	subject := "Авторизация через телеграмм"
	body := fmt.Sprintf("Ваш код подтверждения: %s", code)
	message := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s\r\n",
		smtpMail, email, subject, body)

	// Отправка
	if err := client.Mail(smtpMail); err != nil {
		return fmt.Errorf("sender error: %v", err)
	}

	if err := client.Rcpt(email); err != nil {
		return fmt.Errorf("recipient error: %v", err)
	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("data preparation error: %v", err)
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		return fmt.Errorf("message writing error: %v", err)
	}

	err = w.Close()
	if err != nil {
		return fmt.Errorf("writer closing error: %v", err)
	}

	return nil
}
