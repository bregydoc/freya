package main

type NewTemplate struct {
	Name   string
	Body   []byte
	Params *map[string]string
	Type   *TemplateType
}

type UpdateTemplate struct {
	Name   *string
	Body   *[]byte
	Params *map[string]string
	Type   *TemplateType
}

type Bucket interface {
	Init(string) error
	Close() error
	RegisterNewTemplate(payload NewTemplate) (*Template, error)
	GetTemplateByID(id string) (*Template, error)
	GetAllTemplates() ([]*Template, error)
	EditTemplate(id string, update UpdateTemplate) (*Template, error)
}
