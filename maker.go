package main

import (
	"bytes"
	"fmt"
	"github.com/k0kubun/pp"
	"html/template"
	"io/ioutil"
	"log"
	"net/smtp"
	"os"
)

type Request struct {
	from    string
	to      []string
	subject string
	body    string
}

func NewRequest(to []string, subject string) *Request {
	return &Request{
		to:      to,
		subject: subject,
	}
}

func (r *Request) parseTemplate(templateData []byte, data interface{}) error {
	f, err := ioutil.TempFile(os.TempDir(), "freya")
	if err != nil {
		return err
	}
	_, err = f.Write(templateData)
	if err != nil {
		return err
	}
	t, err := template.ParseFiles(f.Name())
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

func (r *Request) sendMail() error {
	body := "From: " +
		GlobalConfig.MetaData.FromName + "" +
		" <" + GlobalConfig.MetaData.FromEmail +
		">\r\nTo: " + r.to[0] + "\r\nSubject: " +
		r.subject + "\r\n" + GlobalConfig.Mime +
		"\r\n" + r.body
	// body := "To: " + r.to[0] + "\r\nSubject: " + r.subject + "\r\n" + MIME + "\r\n" + r.body
	SMTPClient := fmt.Sprintf("%s:%d", GlobalConfig.Server, GlobalConfig.Port)
	auth := smtp.PlainAuth("", GlobalConfig.Email, GlobalConfig.Password, GlobalConfig.Server)

	if err := smtp.SendMail(SMTPClient, auth, GlobalConfig.Email, r.to, []byte(body)); err != nil {
		pp.Println("ERROR", err.Error())
		return err
	}
	return nil
}

func (r *Request) Send(templateData []byte, items interface{}) error {
	err := r.parseTemplate(templateData, items)
	if err != nil {
		log.Fatal(err)
	}
	if err := r.sendMail(); err != nil {
		log.Printf("Failed to send the email to %s\n", r.to)
		return err
	}
	log.Printf("Email has been sent to %s\n", r.to)
	return nil
}
