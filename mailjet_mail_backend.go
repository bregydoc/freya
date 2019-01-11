package main

import (
	"io/ioutil"
	"log"
)

type MailJetMailBackend struct {
	req *Request
}

func NewMailJetMailBackend() *MailJetMailBackend {
	return &MailJetMailBackend{
		req: nil,
	}
}

func (mail *MailJetMailBackend) SendMail(config *MailConfig, template *Template, params interface{}, subject string, to []string) error {
	data, err := ioutil.ReadAll(template.Data)
	if err != nil {
		return err
	}
	// TODO: Check this stupid frame of code
	request := newRequest(to, subject)

	mail.req = request

	err = mail.req.parseTemplate(data, params)
	if err != nil {
		log.Fatal(err)
	}
	if err := mail.req.sendMail(config); err != nil {
		log.Printf("Failed to send the email to %s\n", mail.req.to)
		return err
	}
	log.Printf("Email has been sent to %s\n", mail.req.to)
	return nil
}
