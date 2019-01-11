package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/k0kubun/pp"
	"io/ioutil"

	"github.com/bregydoc/freya/freyacon/go"
)

type FreyaService struct {
	repo FreyaRepository
}

func (s *FreyaService) SendEmail(ctx context.Context, params *freya.SendEmailParams) (*freya.SendEmailResponse, error) {
	to := make([]string, 0) // TODO: Make a priority queque
	for _, i := range params.To {
		to = append(to, i)
	}

	err := s.repo.SendMail(params.TemplateName, params.TemplateFill, params.Subject, to)

	if err != nil {
		return &freya.SendEmailResponse{
			Error: &freya.Error{
				ErrorCode: 1,
				Message:   err.Error(),
			},
		}, err
	}

	return &freya.SendEmailResponse{
		Send:  true,
		Error: nil,
	}, nil

}

func (s *FreyaService) SaveNewTemplate(ctx context.Context, templateData *freya.TemplateData) (*freya.SaveTemplateResponse, error) {
	template := &Template{
		Name: templateData.TemplateName,
		Data: bytes.NewBuffer(templateData.Data),
	}

	t, err := s.repo.RegisterTemplate(template)

	if err != nil {
		return &freya.SaveTemplateResponse{
			Saved: false,
		}, err
	}

	return &freya.SaveTemplateResponse{
		Saved:        true,
		TemplateName: t.Name,
	}, nil
}

func (s *FreyaService) UpdateTemplate(ctx context.Context, templateData *freya.TemplateData) (*freya.UpdateTemplateResponse, error) {
	template := &Template{
		Name: templateData.TemplateName,
		Data: bytes.NewBuffer(templateData.Data),
	}
	pp.Println(template)
	t, err := s.repo.UpdateTemplate(template)
	fmt.Println(err)
	if err != nil {
		return &freya.UpdateTemplateResponse{
			Template: nil,
			Updated:  false,
			Error: &freya.Error{
				Message:   err.Error(),
				ErrorCode: 1,
			},
		}, err
	}

	data, err := ioutil.ReadAll(t.Data)
	return &freya.UpdateTemplateResponse{
		Template: &freya.TemplateData{
			Data:         data,
			TemplateName: t.Name,
		},
		Updated: true,
		Error:   nil,
	}, nil
}

func (s *FreyaService) GetAllTemplates(ctx context.Context, void *freya.Void) (*freya.TemplatesList, error) {

	templates, err := s.repo.GetAllTemplates(true)

	if err != nil {
		return &freya.TemplatesList{
			Templates: nil,
		}, err
	}

	templatesList := map[string]*freya.TemplateFields{}
	for _, t := range templates {
		templatesList[t.Name] = &freya.TemplateFields{Fields: t.Params}
	}
	return &freya.TemplatesList{
		Templates: templatesList,
	}, nil
}
