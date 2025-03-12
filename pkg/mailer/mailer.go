package mailer

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/smtp"

	configs "backend/pkg/config"
)

type Mailer interface {
	SendText(ctx context.Context, to string, subject string, body string) error
	SendHTML(ctx context.Context, to string, subject string, htmlBody string) error
}

type SMTPMailer struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
}

func NewSMTPMailer(appConfig *configs.AppConfig) Mailer {
	config := appConfig.Smtp

	return &SMTPMailer{
		Host:     config.Host,
		Port:     config.Port,
		Username: config.UserName,
		Password: config.Password,
		From:     config.From,
	}
}

func (m *SMTPMailer) sendEmail(to, subject, body, contentType string) error {
	msg := fmt.Sprintf("From: %s\nTo: %s\nSubject: %s\nMIME-version: 1.0;\nContent-Type: %s\n\n%s",
		m.From, to, subject, contentType, body)

	addr := m.Host + ":" + m.Port
	tlsconfig := &tls.Config{
		//InsecureSkipVerify: true,
		ServerName: m.Host,
	}
	conn, err := tls.Dial("tcp", addr, tlsconfig)

	client, err := smtp.NewClient(conn, m.Host)
	if err != nil {
		return fmt.Errorf("lỗi tạo SMTP client: %v", err)
	}
	defer client.Quit()

	auth := smtp.PlainAuth("", m.Username, m.Password, m.Host)

	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("lỗi xác thực SMTP: %v", err)
	}
	if err = client.Mail(m.From); err != nil {
		return fmt.Errorf("lỗi MAIL FROM: %v", err)
	}
	if err = client.Rcpt(to); err != nil {
		return fmt.Errorf("lỗi RCPT TO: %v", err)
	}
	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("lỗi mở writer: %v", err)
	}
	_, err = writer.Write([]byte(msg))
	if err != nil {
		return fmt.Errorf("lỗi ghi nội dung email: %v", err)
	}
	return writer.Close()
}

func (m *SMTPMailer) SendText(ctx context.Context, to, subject, body string) error {
	return m.sendEmail(to, subject, body, "text/plain; charset=UTF-8")
}

func (m *SMTPMailer) SendHTML(ctx context.Context, to, subject, htmlBody string) error {
	return m.sendEmail(to, subject, htmlBody, "text/html; charset=UTF-8")
}
