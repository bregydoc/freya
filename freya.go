package main

import (
	"github.com/minio/minio-go"
	"github.com/nanobox-io/golang-scribble"
	"log"
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
	Storage     *minio.Client
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

	f.Storage, err = minio.New(
		f.Config.Minio.Endpoint,
		f.Config.Minio.AccessKeyID,
		f.Config.Minio.SecretAccessKey,
		f.Config.Minio.UseSSL,
	)

	bucketName := f.Config.Minio.BucketName

	log.Println("Minio Client setup done ✔︎")

	err = f.Storage.MakeBucket(bucketName, f.Config.Minio.Location)

	if err != nil {
		exists, err := f.Storage.BucketExists(bucketName)
		if err == nil && exists {
			log.Printf("%s bucket exist\n", bucketName)
		} else {
			return nil, err
		}
	}

	log.Println("Minio bucket created ✔︎")

	f.SMSBackend = smsBackend
	f.MailBackend = mailBackend
	return f, nil

}
