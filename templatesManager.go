package main

import (
	"bytes"
	"errors"
	"github.com/minio/minio-go"
	"io/ioutil"
	"time"
)

type Template struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Filename  string    `json:"filename"`
}

func CreateNewTemplate(templateName string, data []byte) (*Template, error) {
	bucketName := GlobalConfig.MinioStorageConfig.BucketName
	id, err := GetNewFreyaID()
	if err != nil {
		return nil, err
	}

	protoTemplate := Template{
		ID:        id,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	suffix, err := GetLittleHash()
	if err != nil {
		return nil, err
	}

	filename := templateName + suffix

	buffer := bytes.NewBuffer(data)

	_, err = MinioClient.PutObject(
		bucketName,
		filename,
		buffer,
		int64(buffer.Len()),
		minio.PutObjectOptions{
			ContentType: "text/html",
		},
	)

	if err != nil {
		return nil, err
	}

	protoTemplate.Filename = filename

	err = ScribbleDriver.Write(GlobalConfig.DBConfig.TemplatesDBName, protoTemplate.ID, protoTemplate)
	if err != nil {
		return nil, err
	}

	newTemplate := new(Template)
	err = ScribbleDriver.Read(GlobalConfig.DBConfig.TemplatesDBName, protoTemplate.ID, newTemplate)
	if err != nil {
		return nil, err
	}

	return newTemplate, nil
}

func GetTemplateByID(ID string) ([]byte, error) {
	bucketName := GlobalConfig.MinioStorageConfig.BucketName
	template := new(Template)
	err := ScribbleDriver.Read(GlobalConfig.DBConfig.TemplatesDBName, ID, template)
	if err != nil {
		return nil, err
	}

	if template.ID == "" {
		return nil, errors.New("template not exist")
	}

	obj, err := MinioClient.GetObject(bucketName, template.Filename, minio.GetObjectOptions{})

	dataTemplate, err := ioutil.ReadAll(obj)
	if err != nil {
		return nil, err
	}

	return dataTemplate, nil
}

func UpdateTemplateByID(ID string, templateData []byte) (*Template, error) {
	bucketName := GlobalConfig.MinioStorageConfig.BucketName

	template := new(Template)
	err := ScribbleDriver.Read(GlobalConfig.DBConfig.TemplatesDBName, ID, template)
	if err != nil {
		return nil, err
	}
	if template.ID == "" {
		return nil, errors.New("template not exist")
	}

	buffer := bytes.NewBuffer(templateData)

	_, err = MinioClient.PutObject(
		bucketName,
		template.Filename,
		buffer,
		int64(buffer.Len()),
		minio.PutObjectOptions{
			ContentType: "text/html",
		},
	)

	if err != nil {
		return nil, err
	}

	template.UpdatedAt = time.Now()
	err = ScribbleDriver.Write(GlobalConfig.DBConfig.TemplatesDBName, template.ID, template)
	if err != nil {
		return nil, err
	}

	return template, nil
}
