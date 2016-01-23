package epazote

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/mail"
	"os"
	"strings"
)

const CRLF = "\r\n"

// extract email address from a list
func GetEmailAddress(s string) (error, []string) {
	var address []string
	for _, v := range strings.Split(s, " ") {
		e, err := mail.ParseAddress(v)
		if err != nil {
			return err, nil
		}
		address = append(address, e.Address)
	}
	return nil, address
}

func (self *Epazote) VerifyEmail() error {
	// set Headers
	if _, ok := self.Config.SMTP.Headers["MIME-Version"]; !ok {
		self.Config.SMTP.Headers["MIME-Version"] = "1.0"
	}
	if _, ok := self.Config.SMTP.Headers["Content-Type"]; !ok {
		self.Config.SMTP.Headers["Content-Type"] = "text/plain; charset=UTF-8"
	}
	if _, ok := self.Config.SMTP.Headers["Content-Transfer-Encoding"]; !ok {
		self.Config.SMTP.Headers["Content-Transfer-Encoding"] = "base64"
	}

	// set From
	if _, ok := self.Config.SMTP.Headers["from"]; !ok {
		name, err := os.Hostname()
		if err != nil {
			return err
		}
		self.Config.SMTP.Headers["from"] = "epazote@" + name
	}

	// check To recipients
	if val, ok := self.Config.SMTP.Headers["to"]; ok {
		err, to := GetEmailAddress(val)
		if err != nil {
			return err
		}
		self.Config.SMTP.Headers["to"] = strings.Join(to, " ")
	}

	// if any serivce needs to notify, check the SMPT settings
	var notify bool

	// check To recipients per service
	for k, v := range self.Services {
		// check for Expect IfNot
		if v.Expect.IfNot.Notify != "" {
			notify = true
			if v.Expect.IfNot.Notify != "yes" {
				err, to := GetEmailAddress(v.Expect.IfNot.Notify)
				if err != nil {
					return fmt.Errorf("Verify notify email addresses for service: %s - %q", k, err)
				}
				v.Expect.IfNot.Notify = strings.Join(to, " ")
			} else if v.Expect.IfNot.Notify == "yes" {
				if _, ok := self.Config.SMTP.Headers["to"]; !ok {
					return fmt.Errorf("Service %q need smtp/headers/to settings to be available to notify.", k)
				}
			}
		}

		// check for Test.IfNot
		if v.IfNot.Notify != "" {
			notify = true
			if v.IfNot.Notify != "yes" {
				err, to := GetEmailAddress(v.IfNot.Notify)
				if err != nil {
					return fmt.Errorf("Verify notify email addresses for service: %s - %q", k, err)
				}
				v.IfNot.Notify = strings.Join(to, " ")
			} else if v.IfNot.Notify == "yes" {
				if _, ok := self.Config.SMTP.Headers["to"]; !ok {
					return fmt.Errorf("Service %q need smtp/headers/to settings to be available to notify.", k)
				}
			}
		}

		// check SMTP.Headers["to"] settings for IfNot
		if len(v.IfStatus) > 0 {
			// key for Service
			for kS, j := range v.IfStatus {
				if j.Notify != "" {
					notify = true
					if j.Notify != "yes" {
						err, to := GetEmailAddress(j.Notify)
						if err != nil {
							return fmt.Errorf("Verify notify email addresses for service [%q if_status: %d]: %q", k, kS, err)
						}
						j.Notify = strings.Join(to, " ")
					} else if j.Notify == "yes" {
						notify = true
						if _, ok := self.Config.SMTP.Headers["to"]; !ok {
							return fmt.Errorf("Service [%q - %d] need smtp/headers/to settings to be available to notify.", k, kS)
						}
					}
				}
			}
		}

		// check SMTP.Headers["to"] settings for IfHeader
		if len(v.IfHeader) > 0 {
			// key for Header
			for kH, j := range v.IfHeader {
				if j.Notify != "" {
					notify = true
					if j.Notify != "yes" {
						err, to := GetEmailAddress(j.Notify)
						if err != nil {
							return fmt.Errorf("Verify notify email addresses for service [%q if_header: %s]: %q", k, kH, err)
						}
						j.Notify = strings.Join(to, " ")
					} else if j.Notify == "yes" {
						notify = true
						if _, ok := self.Config.SMTP.Headers["to"]; !ok {
							return fmt.Errorf("Service [%q - %s] need smtp/headers/to settings to be available to notify.", k, kH)
						}
					}
				}
			}
		}
	}

	if notify {
		if self.Config.SMTP.Server == "" {
			return fmt.Errorf("SMTP server required for been available to send email notifications.")
		}
	}

	return nil
}

func (self *Epazote) SendEmail(m MailMan, to []string, body []byte) {
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

	err := m.Send(to, []byte(msg))
	if err != nil {
		log.Println("ERROR: attempting to send a mail ", err)
	}
}
