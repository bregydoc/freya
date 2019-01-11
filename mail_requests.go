package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/smtp"
	"text/template"
)

type MailBackend interface {
	SendMail(config *MailConfig, template *Template, params interface{}, subject string, to []string) error
}

type Request struct {
	from    string
	to      []string
	subject string
	body    string
}

func newRequest(to []string, subject string) *Request {
	return &Request{
		to:      to,
		subject: subject,
	}
}

func (r *Request) parseTemplate(templateData []byte, data interface{}) error {

	t, err := template.New("freya").Parse(string(templateData))
	if err != nil {
		return err
	}

	buffer := new(bytes.Buffer)

	if err = t.Execute(buffer, data); err != nil {
		return err
	}

	r.body = buffer.String()

	return nil
}

func (r *Request) sendMail(config *MailConfig) error {
	body := "From: " +
		config.MetaData.FromName + "" +
		" <" + config.MetaData.FromEmail +
		">\r\nTo: " + r.to[0] + "\r\nSubject: " +
		r.subject + "\r\n" + config.Mime +
		"\r\n" + r.body
	// body := "To: " + r.to[0] + "\r\nSubject: " + r.subject + "\r\n" + MIME + "\r\n" + r.body
	SMTPClient := fmt.Sprintf("%s:%d", config.Server, config.Port)
	auth := smtp.PlainAuth("", config.Email, config.Password, config.Server)

	if err := smtp.SendMail(SMTPClient, auth, config.Email, r.to, []byte(body)); err != nil {

		return err
	}
	return nil
}

func (r *Request) Send(config *MailConfig, template *Template, params interface{}) error {
	data, err := ioutil.ReadAll(template.Data)
	if err != nil {
		return err
	}
	err = r.parseTemplate(data, params)
	if err != nil {
		log.Fatal(err)
	}
	if err := r.sendMail(config); err != nil {
		log.Printf("Failed to send the email to %s\n", r.to)
		return err
	}
	log.Printf("Email has been sent to %s\n", r.to)
	return nil
}
