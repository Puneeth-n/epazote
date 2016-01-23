package epazote

import (
	"net/smtp"
	"strconv"
)

// MailMan to simplify tests
type MailMan interface {
	Send(to []string, body []byte) error
}

type mailMan struct {
	conf *Email
	send func(string, smtp.Auth, string, []string, []byte) error
}

func (self *mailMan) Send(to []string, body []byte) error {
	// x.x.x.x:25
	addr := self.conf.Server + ":" + strconv.Itoa(self.conf.Port)
	// auth Set up authentication information.
	auth := smtp.PlainAuth("",
		self.conf.Username,
		self.conf.Password,
		self.conf.Server,
	)
	return self.send(addr, auth, self.conf.Headers["from"], to, body)
}

func NewMailMan(conf *Email) MailMan {
	return &mailMan{
		conf,
		smtp.SendMail,
	}
}
