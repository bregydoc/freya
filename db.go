package main

import (
	"time"

	bh "github.com/timshannon/bolthold"
)

type PersistanceDatabase struct {
	initialized bool
	path        string
	store       *bh.Store
}

func (db *PersistanceDatabase) Init(filename string) error {
	store, err := bh.Open(filename, 0666, nil)
	if err != nil {
		return err
	}

	db.store = store
	db.path = filename
	db.initialized = true

	return nil
}

func (db *PersistanceDatabase) Close() error {
	db.initialized = false
	return db.store.Close()
}

func (db *PersistanceDatabase) RegisterNewTemplate(payload NewTemplate) (*Template, error) {
	id, err := IDgen.GetNewFreyaID()
	if err != nil {
		return nil, err
	}

	t := MJML
	if payload.Type != nil {
		t = *payload.Type
	}

	params := map[string]string{}
	if payload.Params != nil {
		params = *payload.Params
	}

	newTemplate := &Template{
		ID:        id,
		Name:      payload.Name,
		Type:      t,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Params:    params,
		Body:      payload.Body,
	}

	err = db.store.Insert(newTemplate.ID, newTemplate)
	if err != nil {
		return nil, err
	}

	return newTemplate, nil
}

func (db *PersistanceDatabase) GetTemplateByID(id string) (*Template, error) {
	template := new(Template)
	err := db.store.Get(id, template)
	if err != nil {
		return nil, err
	}

	return template, nil
}

func (db *PersistanceDatabase) GetAllTemplates() ([]*Template, error) {
	result := new([]Template)
	err := db.store.Find(result, bh.Where("ID").Not().IsNil())
	if err != nil {
		return nil, err
	}

	templates := make([]*Template, 0)
	for i := 0; i < len(*result); i++ {
		templates = append(templates, &(*result)[i])
	}

	return templates, nil
}

func (db *PersistanceDatabase) EditTemplate(id string, update UpdateTemplate) (*Template, error) {
	template, err := db.GetTemplateByID(id)
	if err != nil {
		return nil, err
	}

	if update.Name != nil {
		template.Name = *update.Name
	}

	if update.Body != nil {
		template.Body = *update.Body
	}

	if update.Params != nil {
		template.Params = *update.Params
	}

	if update.Type != nil {
		template.Type = *update.Type
	}

	template.UpdatedAt = time.Now()

	err = db.store.Update(id, template)
	if err != nil {
		return nil, err
	}

	return template, nil
}
