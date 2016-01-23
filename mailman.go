package epazote

import (
	"net/smtp"
	"strconv"
)

// MailMan to simplify tests
type MailMan interface {
	Send(to []string, body []byte) error
}

// emailRecorder for testing
type emailRecorder struct {
	addr string
	auth smtp.Auth
	from string
	to   []string
	msg  []byte
}

func mockSend(errToReturn error) (func(string, smtp.Auth, string, []string, []byte) error, *emailRecorder) {
	r := new(emailRecorder)
	return func(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
		*r = emailRecorder{addr, a, from, to, msg}
		return errToReturn
	}, r
}

func NewMailMan(conf *Email) MailMan {
	return &mailMan{
		conf,
		smtp.SendMail,
	}
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
	//		log.Println("ERROR: attempting to send a mail ", err)
}
