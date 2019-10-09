package main

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"strings"
	"text/template"
	"text/template/parse"
	"time"
)

type TemplateType string

var MJML TemplateType = "mjml"

type Template struct {
	ID        string            `json:"id" boltholdIndex:"ID"`
	Name      string            `json:"name"`
	Params    map[string]string `json:"params"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
	Filename  string            `json:"filename"`
	Type      TemplateType
	Data      io.Reader `json:"-"`
	Body      []byte
}

const bucketName = "freya"

func listTemplateFields(t *template.Template) []string {
	return listNodeFields(t.Tree.Root, nil)
}

func listNodeFields(node parse.Node, res []string) []string {
	if node.Type() == parse.NodeAction {
		// res = append(res, node.String())
		p := strings.Replace(node.String(), "{{", "", -1)
		p = strings.Replace(p, "}}", "", -1)
		p = strings.Replace(p, ".", "", 1)
		res = append(res, p)
	}

	if ln, ok := node.(*parse.ListNode); ok {
		for _, n := range ln.Nodes {
			res = listNodeFields(n, res)
		}
	}
	return res
}

func getParamsMapFromDataOfTemplate(name string, data []byte) (map[string]string, error) {
	tmpl, err := template.New(name).Parse(string(data))
	if err != nil {
		return nil, err
	}
	params := listTemplateFields(tmpl)
	paramsMap := map[string]string{}

	for _, p := range params {
		paramsMap[p] = "string"
	}

	return paramsMap, nil
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

	data, err := ioutil.ReadAll(t.Data)
	if err != nil {
		return nil, err
	}

	err = f.Storage.SaveObject(bucketName, filename, data, &StorageEngineOptions{
		ContentType: "text/html",
	})

	if err != nil {
		return nil, err
	}

	paramsMap, err := getParamsMapFromDataOfTemplate(t.Name, data)
	if err != nil {
		return nil, err
	}

	t.Filename = filename
	t.ID = id
	t.CreatedAt = time.Now()
	t.Params = paramsMap

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

	oldT, err := f.GetTemplateByName(t.Name)
	if err != nil {
		return nil, err
	}

	t.ID = oldT.ID
	t.Filename = oldT.Filename
	t.Name = oldT.Name
	t.CreatedAt = oldT.CreatedAt

	data, err := ioutil.ReadAll(t.Data)
	if err != nil {
		return nil, err
	}

	params, err := getParamsMapFromDataOfTemplate(t.Name, data)
	if err != nil {
		return nil, err
	}

	t.Params = params
	t.UpdatedAt = time.Now()

	err = f.Db.Write(f.Config.DB.TemplatesDBName, t.ID, t)
	if err != nil {
		return nil, err
	}

	err = f.Storage.SaveObject(bucketName, t.Filename, data, &StorageEngineOptions{
		ContentType: "text/html",
	})

	if err != nil {
		return nil, err
	}

	newTemplate, err := f.GetTemplateByID(t.ID, true)
	if err != nil {
		return nil, err
	}

	return newTemplate, nil
}

func (f *Freya) GetTemplateByName(name string, withData ...bool) (*Template, error) {

	allTemplates, err := f.Db.ReadAll(f.Config.DB.TemplatesDBName)
	if err != nil {
		return nil, err
	}

	for _, template := range allTemplates {
		t := new(Template)
		err := json.Unmarshal([]byte(template), t)
		if err != nil {
			return nil, err
		}

		if t.Name == name {
			if len(withData) > 0 {
				if withData[0] {

					obj, err := f.Storage.GetObject(bucketName, t.Filename)
					if err != nil {
						return nil, err
					}
					t.Data = bytes.NewBuffer(obj)
				}
			}
			return t, err
		}
	}
	return nil, templateNotExistError
}

func (f *Freya) GetTemplateByID(id string, withData ...bool) (*Template, error) {

	t := new(Template)
	err := f.Db.Read(f.Config.DB.TemplatesDBName, id, t)
	if err != nil {
		return nil, err
	}

	if t.ID == "" {
		return nil, templateNotExistError
	}

	if len(withData) > 0 {
		if withData[0] {

			obj, err := f.Storage.GetObject(bucketName, t.Filename)
			if err != nil {
				return nil, err
			}
			t.Data = bytes.NewBuffer(obj)
		}
	}

	return t, nil
}

func (f *Freya) GetAllTemplates(withData ...bool) ([]*Template, error) {

	allTemplates, err := f.Db.ReadAll(f.Config.DB.TemplatesDBName)
	if err != nil {
		return nil, err
	}
	templates := make([]*Template, 0)
	for _, tx := range allTemplates {
		t := new(Template)
		err = json.Unmarshal([]byte(tx), t)
		if err != nil {
			return nil, err
		}
		templates = append(templates, t)
	}

	return templates, nil
}
