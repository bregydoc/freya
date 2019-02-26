package main

import (
	"log"

	"github.com/nanobox-io/golang-scribble"
)

type Repository interface {
	RegisterTemplate(t *Template) (*Template, error)
	UpdateTemplate(t *Template) (*Template, error)
	GetTemplateByName(name string, withData ...bool) (*Template, error)
	GetTemplateByID(id string, withData ...bool) (*Template, error)
	GetAllTemplates(withData ...bool) ([]*Template, error)

	SendMail(templateName string, params interface{}, subject string, to []string) error
	SendSMS(templateName string, params map[string]string, to *PhoneNumber) error
}

type Freya struct {
	Config      *FreyaConfig
	Db          *scribble.Driver
	Storage     *StorageEngine
	Gen         *FreyaIDGenerator
	SMSBackend  SMSBackend
	MailBackend MailBackend
}

func NewFreya(config *FreyaConfig, mailBackend MailBackend, smsBackend SMSBackend) (*Freya, error) {
	config = FillConfigWithDefaults(config)
	f := new(Freya)
	f.Config = config

	var err error

	f.Gen = &FreyaIDGenerator{
		alphabet: "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ",
	}

	f.Db, err = scribble.New(config.DB.RelativeFolder, nil)
	if err != nil {
		return nil, err
	}

	log.Println("Scribble setup done ✔︎")

	f.Storage, err = NewStorageEngine(config.Storage)

	log.Println("Storage setup done ✔︎")

	err = f.Storage.Init()
	if err != nil {
		return nil, err
	}

	log.Println("Storage init correctly ✔︎")

	f.SMSBackend = smsBackend
	f.MailBackend = mailBackend
	return f, nil

}
