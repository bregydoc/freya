package main

import (
	"context"

	"fmt"
	"github.com/bregydoc/freya/freyacon/go"
	"google.golang.org/grpc"
	"log"
	"net"
)

type FreyaService struct{}

func (s *FreyaService) SendEmail(ctx context.Context, params *freya.SendEmailParams) (*freya.SendEmailResponse, error) {
	to := make([]string, 0) // TODO: Make a priority queque
	for _, i := range params.To {
		to = append(to, i)
	}

	emailRequest := NewRequest(to, params.Subject)

	err := SendMailFromSavedTemplate(emailRequest, params.TemplateName, params.TemplateFill)
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
	t, err := CreateNewTemplate(templateData.TemplateName, templateData.Data)
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
	panic("unimplemented")
}

func (s *FreyaService) GetAllTemplates(ctx context.Context, void *freya.Void) (*freya.TemplatesList, error) {
	templates, err := GetAllTemplates()
	if err != nil {
		return &freya.TemplatesList{
			Templates: nil,
		}, err
	}

	templatesList := map[string]*freya.TemplateFields{}
	for _, t := range templates {
		templatesList[t.Name] = &freya.TemplateFields{Fields: map[string]string{}}
	}
	return &freya.TemplatesList{
		Templates: templatesList,
	}, nil
}

func main() {
	port := 10000
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()

	freya.RegisterFreyaServer(grpcServer, &FreyaService{})

	log.Printf("Listening on :%d\n", port)
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("failed on serve: %v", err)
	}

}
