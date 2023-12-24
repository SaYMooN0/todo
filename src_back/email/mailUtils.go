package dbutils

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/smtp"
)

var (
	smtpHost     string
	smtpPort     string
	smtpUser     string
	smtpPassword string
)

func InitMailUtils(host, port, user, password string) {
	smtpHost = host
	smtpPort = port
	smtpUser = user
	smtpPassword = password
}
func createMessage(to, subject, body string) []byte {
	var email bytes.Buffer
	writer := multipart.NewWriter(&email)
	writer.WriteField("From", "Todo app")
	writer.WriteField("To", to)
	writer.WriteField("Subject", subject)
	writer.WriteField("Content-Type", "text/plain; charset=UTF-8")
	writer.WriteField("Content-Transfer-Encoding", "quoted-printable")
	writer.WriteField("MIME-Version", "1.0")
	writer.WriteField("body", body)
	writer.Close()
	return email.Bytes()
}
func SendEmail(to, subject, body string) error {
	message := createMessage(to, subject, body)
	auth := smtp.PlainAuth("", smtpUser, smtpPassword, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, "Todo app", []string{to}, message)
	return err
}
func SendConfirmationCode(to, confirmationCode string) error {
	htmlBody := fmt.Sprintf("<p>Your confirmation code is: <strong>%s</strong></p>", confirmationCode)
	return SendEmail(to, "Confirmation code", htmlBody)
}
