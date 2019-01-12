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
	data := []byte("Hello {{.Name}}, from freya sms engine")
	resp, err := freyaClient.SaveNewTemplate(context.Background(), &freya.TemplateData{
		TemplateName: "sms_hello",
		Data:         data,
	})

	if err != nil {
		panic(err)
	}

	if resp.Saved {
		log.Println("Template sms saved")
	}

	freyaClient.SendSMS()
}
