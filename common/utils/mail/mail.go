package mail

import (
	"crypto/tls"
	"fmt"
	netmail "net/mail"
	"net/smtp"
	"slices"
	"strings"

	"cth.release/common"
)

type Config struct {
	Host               string
	Port               int
	Username           string
	Password           string
	From               string
	AuthIdentity       string
	StartTLS           bool
	InsecureSkipVerify bool
}

type Message struct {
	To       []string
	Cc       []string
	Bcc      []string
	Subject  string
	TextBody string
	HTMLBody string
	ReplyTo  string
	Header   map[string]string
}

type Sender struct {
	config Config
}

func NewSender(cfg Config) *Sender {
	return &Sender{config: cfg}
}

func NewSenderFromAppConfig(cfg *common.Config) *Sender {
	return NewSender(Config{
		Host:               cfg.SMTP.Host,
		Port:               cfg.SMTP.Port,
		Username:           cfg.SMTP.User,
		Password:           cfg.SMTP.Pass,
		From:               cfg.SMTP.From,
		AuthIdentity:       cfg.SMTP.AuthIdentity,
		StartTLS:           cfg.SMTP.StartTLS,
		InsecureSkipVerify: cfg.SMTP.InsecureSkipVerify,
	})
}

func (s *Sender) IsConfigured() bool {
	return s.config.Host != "" && s.config.Port > 0 && s.config.From != ""
}

func (s *Sender) Send(message Message) error {
	if !s.IsConfigured() {
		return fmt.Errorf("smtp sender is not configured")
	}
	if len(message.To) == 0 {
		return fmt.Errorf("at least one recipient is required")
	}

	recipients := collectRecipients(message)
	address := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
	payload, err := buildMessage(s.config.From, message)
	if err != nil {
		return err
	}

	if s.config.StartTLS {
		return s.sendWithTLS(address, recipients, payload)
	}

	auth := s.smtpAuth()
	return smtp.SendMail(address, auth, s.config.From, recipients, payload)
}

func (s *Sender) sendWithTLS(address string, recipients []string, payload []byte) error {
	conn, err := tls.Dial("tcp", address, &tls.Config{
		ServerName:         s.config.Host,
		InsecureSkipVerify: s.config.InsecureSkipVerify,
	})
	if err != nil {
		return err
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, s.config.Host)
	if err != nil {
		return err
	}
	defer client.Close()

	if auth := s.smtpAuth(); auth != nil {
		if err := client.Auth(auth); err != nil {
			return err
		}
	}

	if err := client.Mail(s.config.From); err != nil {
		return err
	}
	for _, recipient := range recipients {
		if err := client.Rcpt(recipient); err != nil {
			return err
		}
	}

	writer, err := client.Data()
	if err != nil {
		return err
	}

	if _, err := writer.Write(payload); err != nil {
		return err
	}
	if err := writer.Close(); err != nil {
		return err
	}

	return client.Quit()
}

func (s *Sender) smtpAuth() smtp.Auth {
	if s.config.Username == "" {
		return nil
	}
	return smtp.PlainAuth(s.config.AuthIdentity, s.config.Username, s.config.Password, s.config.Host)
}

func buildMessage(from string, message Message) ([]byte, error) {
	headers := map[string]string{
		"From":         from,
		"To":           strings.Join(message.To, ", "),
		"Subject":      message.Subject,
		"MIME-Version": "1.0",
	}

	if len(message.Cc) > 0 {
		headers["Cc"] = strings.Join(message.Cc, ", ")
	}
	if message.ReplyTo != "" {
		headers["Reply-To"] = message.ReplyTo
	}
	for key, value := range message.Header {
		headers[key] = value
	}

	var body string
	if message.HTMLBody != "" && message.TextBody != "" {
		boundary := "fiber-base-mail-boundary"
		headers["Content-Type"] = fmt.Sprintf("multipart/alternative; boundary=%q", boundary)
		body = strings.Join([]string{
			"--" + boundary,
			"Content-Type: text/plain; charset=UTF-8",
			"",
			message.TextBody,
			"--" + boundary,
			"Content-Type: text/html; charset=UTF-8",
			"",
			message.HTMLBody,
			"--" + boundary + "--",
			"",
		}, "\r\n")
	} else if message.HTMLBody != "" {
		headers["Content-Type"] = "text/html; charset=UTF-8"
		body = message.HTMLBody
	} else {
		headers["Content-Type"] = "text/plain; charset=UTF-8"
		body = message.TextBody
	}

	var builder strings.Builder
	for key, value := range headers {
		builder.WriteString(fmt.Sprintf("%s: %s\r\n", key, value))
	}
	builder.WriteString("\r\n")
	builder.WriteString(body)

	return []byte(builder.String()), nil
}

func collectRecipients(message Message) []string {
	recipients := make([]string, 0, len(message.To)+len(message.Cc)+len(message.Bcc))
	recipients = append(recipients, message.To...)
	recipients = append(recipients, message.Cc...)
	recipients = append(recipients, message.Bcc...)

	normalized := make([]string, 0, len(recipients))
	for _, recipient := range recipients {
		trimmed := strings.TrimSpace(recipient)
		if trimmed == "" {
			continue
		}
		if _, err := netmail.ParseAddress(trimmed); err == nil && !slices.Contains(normalized, trimmed) {
			normalized = append(normalized, trimmed)
		}
	}
	return normalized
}
