package main

import (
	"context"

	"fmt"
	"github.com/bregydoc/freya/freyabuf"
	"google.golang.org/grpc"
	"log"
	"net"
)

type FreyaService struct{}

func (s *FreyaService) SendEmail(ctx context.Context, params *freyabuf.SendEmailParams) (*freyabuf.SendEmailResponse, error) {
	to := make([]string, 0) // TODO: Make a priority queque
	for _, i := range params.To {
		to = append(to, i)
	}

	emailRequest := NewRequest(to, params.Subject)

	err := SendMailFromSavedTemplate(emailRequest, params.TemplateName, params.TemplateFill)
	if err != nil {
		return &freyabuf.SendEmailResponse{
			Error: &freyabuf.Error{
				ErrorCode: 1,
				Message:   err.Error(),
			},
		}, err
	}

	return &freyabuf.SendEmailResponse{
		Sended: true,
		Error:  nil,
	}, nil

}

func (s *FreyaService) SaveNewTemplate(ctx context.Context, templateData *freyabuf.TemplateData) (*freyabuf.SaveTemplateResponse, error) {
	_, err := CreateNewTemplate(templateData.TemplateName, templateData.Data)
	if err != nil {
		return &freyabuf.SaveTemplateResponse{
			Saved: false,
		}, err
	}

	return &freyabuf.SaveTemplateResponse{
		Saved: true,
	}, nil
}

func (s *FreyaService) UpdateTemplate(ctx context.Context, templateData *freyabuf.TemplateData) (*freyabuf.UpdateTemplateResponse, error) {
	panic("unimplemented")
}

func main() {
	port := 10000
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()

	freyabuf.RegisterFreyaServer(grpcServer, &FreyaService{})

	log.Printf("Listening on :%d\n", port)
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("failed on serve: %v", err)
	}

}
