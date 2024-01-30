package sendgrid

import (
    "github.com/sendgrid/sendgrid-go"
    "github.com/sendgrid/sendgrid-go/helpers/mail"
    "os"
)

// SendEmail sends an email using SendGrid.
func SendEmail(toEmail, subject, message string) error {
	from := mail.NewEmail("MacBobby with Paritie Innovation Hub", "theghostmac@gmail.com")
    to := mail.NewEmail("Recipient", toEmail)
    content := mail.NewContent("text/plain", message)
    m := mail.NewV3MailInit(from, subject, to, content)

    apiKey := os.Getenv("SENDGRID_API_KEY")
    client := sendgrid.NewSendClient(apiKey)
    _, err := client.Send(m)
    return err
}