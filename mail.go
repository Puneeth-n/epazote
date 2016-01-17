package epazote

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/smtp"
	"strconv"
	"strings"
)

const CRLF = "\r\n"

func (self *Epazote) SendEmail(s *Service, to []string, m string) {
	// auth Set up authentication information.
	auth := smtp.PlainAuth("",
		self.Config.SMTP.Username,
		self.Config.SMTP.Password,
		self.Config.SMTP.Server,
	)

	// set To
	if len(to) < 1 {
		to = strings.Split(self.Config.SMTP.Headers["to"], " ")
	}

	if len(to) == 0 {
		log.Println("No recipients set, no email sent")
		return
	}

	// email Body
	body := `
	Service: %q

	`

	// message template
	msg := ""
	for k, v := range self.Config.SMTP.Headers {
		if k == "to" {
			msg += fmt.Sprintf("%s: %s %s", strings.Title(k), strings.Join(to, ", "), CRLF)
		} else {
			msg += fmt.Sprintf("%s: %s %s", strings.Title(k), strings.TrimSpace(v), CRLF)
		}
	}
	msg += CRLF + base64.StdEncoding.EncodeToString([]byte(body))

	err := smtp.SendMail(
		self.Config.SMTP.Server+":"+strconv.Itoa(self.Config.SMTP.Port),
		auth,
		self.Config.SMTP.Headers["from"],
		to,
		[]byte(msg),
	)

	if err != nil {
		log.Println("ERROR: attempting to send a mail ", err)
	}

}
