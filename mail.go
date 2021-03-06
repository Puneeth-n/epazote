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
	//initialize Headers in case they don't exists
	if self.Config.SMTP.Headers == nil {
		self.Config.SMTP.Headers = make(map[string]string)
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
					return fmt.Errorf(Red("Verify notify email addresses for service: %s - %q"), k, err)
				}
				v.Expect.IfNot.Notify = strings.Join(to, " ")
			} else if v.Expect.IfNot.Notify == "yes" {
				if _, ok := self.Config.SMTP.Headers["to"]; !ok {
					return fmt.Errorf(Red("Service %q need smtp/headers/to settings to be available to notify."), k)
				}
			}
		}

		// check for Test.IfNot
		if v.IfNot.Notify != "" {
			notify = true
			if v.IfNot.Notify != "yes" {
				err, to := GetEmailAddress(v.IfNot.Notify)
				if err != nil {
					return fmt.Errorf(Red("Verify notify email addresses for service: %s - %q"), k, err)
				}
				v.IfNot.Notify = strings.Join(to, " ")
			} else if v.IfNot.Notify == "yes" {
				if _, ok := self.Config.SMTP.Headers["to"]; !ok {
					return fmt.Errorf(Red("Service %q need smtp/headers/to settings to be available to notify."), k)
				}
			}
		}

		// check SMTP.Headers["to"] settings for IfNot
		if v.IfStatus != nil {
			// key for Service
			for kS, j := range v.IfStatus {
				if j.Notify != "" {
					notify = true
					if j.Notify != "yes" {
						err, to := GetEmailAddress(j.Notify)
						if err != nil {
							return fmt.Errorf(Red("Verify notify email addresses for service [%q if_status: %d]: %q"), k, kS, err)
						}
						j.Notify = strings.Join(to, " ")
					} else if j.Notify == "yes" {
						notify = true
						if _, ok := self.Config.SMTP.Headers["to"]; !ok {
							return fmt.Errorf(Red("Service [%q - %d] need smtp/headers/to settings to be available to notify."), k, kS)
						}
					}
				}
			}
		}

		// check SMTP.Headers["to"] settings for IfHeader
		if v.IfHeader != nil {
			// key for Header
			for kH, j := range v.IfHeader {
				if j.Notify != "" {
					notify = true
					if j.Notify != "yes" {
						err, to := GetEmailAddress(j.Notify)
						if err != nil {
							return fmt.Errorf(Red("Verify notify email addresses for service [%q if_header: %s]: %q"), k, kH, err)
						}
						j.Notify = strings.Join(to, " ")
					} else if j.Notify == "yes" {
						notify = true
						if _, ok := self.Config.SMTP.Headers["to"]; !ok {
							return fmt.Errorf(Red("Service [%q - %s] need smtp/headers/to settings to be available to notify."), k, kH)
						}
					}
				}
			}
		}
	}

	if notify || self.Config.SMTP.Server != "" {
		if self.Config.SMTP.Server == "" {
			return fmt.Errorf(Red("SMTP server required for been available to send email notifications."))
		}

		// default to port 25
		if self.Config.SMTP.Port == 0 {
			self.Config.SMTP.Port = 25
		}

		// set Headers
		if _, ok := self.Config.SMTP.Headers["MIME-Version"]; !ok {
			self.Config.SMTP.Headers["MIME-Version"] = "1.0"
		}
		if _, ok := self.Config.SMTP.Headers["Content-Type"]; !ok {
			self.Config.SMTP.Headers["Content-Type"] = "text/plain; charset=\"utf-8\""
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

		// set subject
		if _, ok := self.Config.SMTP.Headers["subject"]; !ok {
			self.Config.SMTP.Headers["subject"] = "[name, because]"
		}

		// check To recipients
		if val, ok := self.Config.SMTP.Headers["to"]; ok {
			err, to := GetEmailAddress(val)
			if err != nil {
				return fmt.Errorf(Red("Verify recipient's email address: %s"), err)
			}
			self.Config.SMTP.Headers["to"] = strings.Join(to, " ")
		}

		// enable SMTP
		// This is to avoid an error if new services added via scan need to send email
		// but no smtp is defined
		self.Config.SMTP.enabled = true
	}

	return nil
}

func (self *Epazote) SendEmail(m MailMan, to []string, subject string, body []byte) {
	// message template
	msg := ""
	for k, v := range self.Config.SMTP.Headers {
		if k == "to" {
			msg += fmt.Sprintf("%s: %s%s", strings.Title(k), strings.Join(to, ", "), CRLF)
		} else if k == "subject" {
			msg += fmt.Sprintf("%s: %s%s", strings.Title(k), strings.TrimSpace(subject), CRLF)
		} else {
			msg += fmt.Sprintf("%s: %s%s", strings.Title(k), strings.TrimSpace(v), CRLF)
		}
	}

	msg += CRLF + base64.StdEncoding.EncodeToString(body)

	err := m.Send(to, []byte(msg))
	if err != nil {
		log.Printf("ERROR attempting to send a mail: %q", err)
	}
}
