package main

import (
	"context"
	"github.com/bregydoc/freya/freyacon/go"
	"github.com/k0kubun/pp"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
)

func main() {
	endpoint := "127.0.0.1:10000"
	client, err := grpc.Dial(endpoint, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	freyaClient := freya.NewFreyaClient(client)
	welcomeData, err := ioutil.ReadFile("/Users/macBook/Documents/welcome_template.html.txt")
	if err != nil {
		panic(err)
	}

	resp, err := freyaClient.UpdateTemplate(context.Background(), &freya.TemplateData{
		TemplateName: "welcome",
		Data:         welcomeData,
	})

	if err != nil {
		pp.Println(resp)
		panic(err)
	}

	if resp.Updated {
		log.Println("Template saved")
	}

	pp.Println(resp.Template)

}
