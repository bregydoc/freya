package main

import (
	"context"
	"github.com/bregydoc/freya/freyacon/go"
	"github.com/k0kubun/pp"
	"google.golang.org/grpc"
)

func main() {
	endpoint := "127.0.0.1:10000"
	client, err := grpc.Dial(endpoint, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	freyaClient := freya.NewFreyaClient(client)

	in := &freya.SendEmailParams{
		Subject:      "Hello",
		TemplateFill: map[string]string{},
		TemplateName: "welcome_mail",
		To:           map[int32]string{0: "bregy.malpartida@utec.edu.pe"},
	}

	res, err := freyaClient.SendEmail(context.Background(), in)
	if err != nil {
		panic(err)
	}

	pp.Println(res)
}
