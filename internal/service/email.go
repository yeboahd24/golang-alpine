package service

import (
    "fmt"
    "net/smtp"
)

type EmailConfig struct {
    Host     string
    Port     int
    Username string
    Password string
    From     string
}

type EmailService struct {
    config EmailConfig
}

func NewEmailService(config EmailConfig) *EmailService {
    return &EmailService{config: config}
}

func (s *EmailService) SendVerificationEmail(to, token string) error {
    subject := "Verify Your Email"
    verificationLink := fmt.Sprintf("http://your-domain.com/verify?token=%s", token)
    body := fmt.Sprintf("Please click the link below to verify your email:\n%s", verificationLink)

    msg := fmt.Sprintf("From: %s\r\n"+
        "To: %s\r\n"+
        "Subject: %s\r\n"+
        "\r\n"+
        "%s\r\n", s.config.From, to, subject, body)

    auth := smtp.PlainAuth("", s.config.Username, s.config.Password, s.config.Host)
    addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)

    return smtp.SendMail(addr, auth, s.config.From, []string{to}, []byte(msg))
}