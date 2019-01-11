package main

import (
	"fmt"
	"github.com/minio/minio-go"
	"io"
	"log"
	"time"
)

type Template struct {
	ID        string
	Name      string
	Params    map[string]string
	CreatedAt time.Time
	UpdatedAt time.Time
	Filename  string
	Data      io.Reader
}

func (f *Freya) RegisterTemplate(t *Template) (*Template, error) {

	id, err := f.Gen.GetNewFreyaID()
	if err != nil {
		return nil, err
	}

	hash, err := f.Gen.GetLittleHash()
	if err != nil {
		return nil, err
	}

	filename := t.Name + "_" + hash + ".html"

	length, err := t.Data.Read(nil)
	if err != nil {
		return nil, err
	}

	bucketName := f.Config.Minio.BucketName

	_, err = f.Storage.PutObject(
		bucketName,
		filename,
		t.Data,
		int64(length),
		minio.PutObjectOptions{
			ContentType: "text/html",
		},
	)

	if err != nil {
		return nil, err
	}

	t.Filename = filename
	t.ID = id
	t.CreatedAt = time.Now()

	err = f.Db.Write(f.Config.DB.TemplatesDBName, t.ID, t)
	if err != nil {
		return nil, err
	}

	newTemplate := new(Template)
	err = f.Db.Read(f.Config.DB.TemplatesDBName, t.ID, newTemplate)
	if err != nil {
		return nil, err
	}

	return newTemplate, nil
}

func (f *Freya) UpdateTemplate(t *Template) (*Template, error) {
	t.UpdatedAt = time.Now()

	t, err := f.GetTemplateByName(t.Name)
	if err != nil {
		return nil, err
	}
	fmt.Print("0, ", err)
	err = f.Db.Write(f.Config.DB.TemplatesDBName, t.ID, t)
	if err != nil {
		return nil, err
	}
	fmt.Print("1, ", err)
	length, err := t.Data.Read(nil)
	if err != nil {
		return nil, err
	}

	fmt.Print("2, ", err)
	bucketName := f.Config.Minio.BucketName

	_, err = f.Storage.PutObject(
		bucketName,
		t.Filename,
		t.Data,
		int64(length),
		minio.PutObjectOptions{
			ContentType: "text/html",
		},
	)
	fmt.Print("3, ", err)
	newTemplate := new(Template)
	err = f.Db.Read(f.Config.DB.TemplatesDBName, t.ID, newTemplate)
	if err != nil {
		return nil, err
	}
	fmt.Print("4, ", err)
	return newTemplate, nil
}

func (f *Freya) GetTemplateByName(name string, withData ...bool) (*Template, error) {
	allTemplates, err := f.Db.ReadAll(f.Config.DB.TemplatesDBName)
	if err != nil {
		return nil, err
	}

	log.Println(allTemplates)

	for _, id := range allTemplates {
		t, err := f.GetTemplateByID(id, false)
		if err != nil {
			return nil, err
		}
		if t.Name == name {
			if len(withData) > 0 {
				if withData[0] {
					bucketName := f.Config.Minio.BucketName
					obj, err := f.Storage.GetObject(bucketName, t.Filename, minio.GetObjectOptions{})
					if err != nil {
						return nil, err
					}
					t.Data = obj
				}
			}
			return t, err
		}
	}
	return nil, templateNotExistError
}

func (f *Freya) GetTemplateByID(id string, withData ...bool) (*Template, error) {
	template := new(Template)
	err := f.Db.Read(f.Config.DB.TemplatesDBName, id, template)
	if err != nil {
		return nil, err
	}

	if template.ID == "" {
		return nil, templateNotExistError
	}

	if len(withData) > 0 {
		if withData[0] {
			bucketName := f.Config.Minio.BucketName
			obj, err := f.Storage.GetObject(bucketName, template.Filename, minio.GetObjectOptions{})
			if err != nil {
				return nil, err
			}
			template.Data = obj
		}
	}

	return template, nil
}

func (f *Freya) GetAllTemplates(withData ...bool) ([]*Template, error) {
	allTemplates, err := f.Db.ReadAll(f.Config.DB.TemplatesDBName)
	if err != nil {
		return nil, err
	}
	templates := make([]*Template, 0)
	for _, id := range allTemplates {
		t, err := f.GetTemplateByID(id, withData...)
		if err != nil {
			return nil, err
		}
		templates = append(templates, t)
	}

	return templates, nil
}
