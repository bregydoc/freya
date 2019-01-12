package main

import (
	"context"
	"github.com/bregydoc/freya/freyacon/go"
	"google.golang.org/grpc"
	"log"
)

func main() {
	endpoint := "127.0.0.1:10000"
	client, err := grpc.Dial(endpoint, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	freyaClient := freya.NewFreyaClient(client)

	data := []byte("Hello {{.Name}}, I'm {{.engine}}")
	resp, err := freyaClient.SaveNewTemplate(context.Background(), &freya.TemplateData{
		TemplateName: "hello",
		Data:         data,
	})

	if err != nil {
		panic(err)
	}

	if resp.Saved {
		log.Println("Template sms saved")
	}

	r, err := freyaClient.SendSMS(context.Background(), &freya.SendSMSParams{
		TemplateName: "hello",
		Phone: &freya.PhoneNumber{
			CountryCode: "51",
			Number:      "957821858",
		},
		Params: map[string]string{
			"Name":   "Bregy",
			"engine": "freya",
		},
	})

	if err != nil {
		panic(err)
	}

	if r.Sent {
		log.Println("SMS has been sent")
	}
}
