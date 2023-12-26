package dbutils

import (
	"fmt"
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
	headers := "From: Todo app\r\n"
	headers += "To: " + to + "\r\n"
	headers += "Subject: " + subject + "\r\n"
	headers += "MIME-Version: 1.0\r\n"
	headers += "Content-Type: text/html; charset=UTF-8\r\n"
	headers += "\r\n" + body

	return []byte(headers)
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
