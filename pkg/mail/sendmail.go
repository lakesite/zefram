package mail

import (
	"fmt"
	valid "github.com/asaskevich/govalidator"
	"net/smtp"
)

// LocalhostSendMail will send an e-mail using localhost, such that your system
// has an smtp daemon configured and running to send mail.
func LocalhostSendMail(from, to, subject, body string) error {
	// from and to validation
	if !valid.IsEmail(from) {
		return fmt.Errorf("From email '%s' is not a valid e-mail.", from)
	}

	// What if we're delivering to a local user? e.g. 'root'?
	if !valid.IsEmail(to) {
		return fmt.Errorf("To email '%s' is not a valid e-mail.", to)
	}

	// subject and body can be blank.
	headers := make(map[string]string)
	headers["From"] = from
	headers["To"] = to
	headers["Subject"] = subject

	// Setup message
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	// Handle plain jane susie mary localhost authenticary
	c, err := smtp.Dial("127.0.0.1:25")
	if err != nil {
		return err
	}
	defer c.Close()

	// from
	if err = c.Mail(from); err != nil {
		return err
	}

	// to
	if err = c.Rcpt(to); err != nil {
		return err
	}

	// Data
	cw, err := c.Data()
	if err != nil {
		return err
	}

	_, err = cw.Write([]byte(message))
	if err != nil {
		return err
	}

	err = cw.Close()
	if err != nil {
		return err
	}
	c.Quit()

	return nil
}
