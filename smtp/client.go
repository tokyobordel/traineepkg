// Package smtp предоставляет клиент для отправки email-уведомлений.
package smtp

import (
	"crypto/tls"
	"log"
	"net/smtp"
	"strings"
)

// SmtpClient хранит параметры подключения к SMTP-серверу.
type SmtpClient struct {
	// Email — адрес почтового ящика отправителя.
	Email string
	// Password — пароль почтового ящика отправителя.
	Password string
	// Host — хост SMTP-сервера.
	Host string
	// Port — порт SMTP-сервера.
	Port string
}

// NewSmtpClient создаёт клиент SMTP с заданными учётными данными и адресом сервера.
func NewSmtpClient(e string, p string, h string, port string) *SmtpClient {
	return &SmtpClient{
		Email:    e,
		Password: p,
		Host:     h,
		Port:     port,
	}
}

// SendMessage отправляет email-уведомление списку получателей receiverEmails.
// Параметр notify_type определяет тему письма.
func (s SmtpClient) SendMessage(receiverEmails []string, message string, notify_type string) {
	tlsConfig := &tls.Config{
		ServerName: s.Host,
	}

	conn, err := tls.Dial("tcp", s.Host+":"+s.Port, tlsConfig)
	if err != nil {
		log.Print("подключение TLS:", err.Error())
		return
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, s.Host)
	if err != nil {
		log.Printf("создание клиента: %v", err)
		return
	}
	defer client.Quit()

	auth := smtp.PlainAuth("", s.Email, s.Password, s.Host)
	if err = client.Auth(auth); err != nil {
		log.Printf("аутентификация: %v", err)
		return
	}

	if err = client.Mail(s.Email); err != nil {
		log.Printf("отправитель: %v", err)
		return
	}

	for _, rcpt := range receiverEmails {
		if err = client.Rcpt(rcpt); err != nil {
			log.Printf("получатель %s: %v", rcpt, err)
			return
		}
	}

	w, err := client.Data()
	if err != nil {
		log.Printf("открытие Data: %v", err)
		return
	}
	defer w.Close()

	var topic string
	switch notify_type {
	case "user_register":
		topic = "Регистрация аккаунта"
	case "user_login":
		topic = "Вход в аккаунт"
	case "admin_newImg":
		topic = "Новое изображение для модерации"
	case "user_imgVerdict":
		topic = "Модерация вашего поста"
	case "user_confirm":
		topic = "Подтверждение почты"
	default:
		topic = "Служебное сообщение"
	}

	var temp_message = "Subject: " + topic + "\r\n" + "\r\n"
	temp_message += message
	temp_message = strings.ReplaceAll(temp_message, "<b>", "")
	temp_message = strings.ReplaceAll(temp_message, "</b>", "")
	byte_temp_message := []byte(temp_message)

	_, err = w.Write(byte_temp_message)
	if err != nil {
		log.Printf("запись письма: %v", err)
		return
	}
}
