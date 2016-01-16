package epazote

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"strconv"
	"strings"
)

const CRLF = "\r\n"

func (self *Epazote) SendEmail(s *Service, body string) {
	// auth Set up authentication information.
	auth := smtp.PlainAuth("",
		self.Config.SMTP.Username,
		self.Config.SMTP.Password,
		self.Config.SMTP.Server,
	)

	// set From
	if _, ok := self.Config.SMTP.Headers["from"]; !ok {
		name, err := os.Hostname()
		if err != nil {
			log.Println(err)
		}
		self.Config.SMTP.Headers["from"] = "epazote@" + name
	}

	// set To
	to := strings.Split(self.Config.SMTP.Headers["to"], " ")

	// add headers
	e.Headers["MIME-Version"] = "1.0"
	e.Headers["Content-Type"] = "text/plain; charset=UTF-8"
	e.Headers["Content-Transfer-Encoding"] = "base64"

	// email Body
	body := ""

	// message template
	msg := ""
	for k, v := range e.Headers {
		if k == "to" {
			msg += fmt.Sprintf("%s: %s %s", strings.Title(k), strings.Join(to, ", "), CRLF)
		} else {
			msg += fmt.Sprintf("%s: %s %s", strings.Title(k), strings.TrimSpace(v), CRLF)
		}
	}
	msg += CRLF + base64.StdEncoding.EncodeToString([]byte(body))

	err := smtp.SendMail(
		e.Server+":"+strconv.Itoa(e.Port),
		auth,
		e.Headers["from"],
		to,
		[]byte(msg),
	)

	if err != nil {
		log.Println("ERROR: attempting to send a mail ", err)
	}

}
