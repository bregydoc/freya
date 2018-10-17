package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/k0kubun/pp"
	"github.com/minio/minio-go"
	"io/ioutil"
	"time"
)

type Template struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
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
		Name:      templateName,
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

func ReadTemplate(template *Template) ([]byte, error) {
	bucketName := GlobalConfig.MinioStorageConfig.BucketName
	obj, err := MinioClient.GetObject(bucketName, template.Filename, minio.GetObjectOptions{})

	dataTemplate, err := ioutil.ReadAll(obj)
	if err != nil {
		return nil, err
	}

	return dataTemplate, nil
}
func GetTemplateByID(ID string) (*Template, error) {

	template := new(Template)
	err := ScribbleDriver.Read(GlobalConfig.DBConfig.TemplatesDBName, ID, template)
	if err != nil {
		return nil, err
	}

	if template.ID == "" {
		return nil, errors.New("template not exist")
	}

	return template, nil

}

func GetTemplateByName(name string) (*Template, error) {
	templates, err := ScribbleDriver.ReadAll(GlobalConfig.DBConfig.TemplatesDBName)
	pp.Println(templates)
	if err != nil {
		return nil, err
	}
	for _, temp := range templates {
		t := new(Template)
		err = json.Unmarshal([]byte(temp), t)
		//err = ScribbleDriver.Read(GlobalConfig.DBConfig.TemplatesDBName, id, t)
		if err != nil {
			return nil, err
		}
		if t.Name == name {
			return t, nil
		}
	}

	return nil, errors.New("template not found")
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
