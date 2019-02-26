package main

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/k0kubun/pp"
	"io/ioutil"
	"log"

	freya "github.com/bregydoc/freya/freyacon/go"
	"google.golang.org/grpc"
)

func main() {
	endpoint := "127.0.0.1:10000"
	client, err := grpc.Dial(endpoint, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer client.Close()

	freyaClient := freya.NewFreyaClient(client)
	welcomeData, err := ioutil.ReadFile("/Users/macbook/go/src/github.com/bregydoc/micro-culqi/r.html")
	if err != nil {
		panic(err)
	}

	resp, err := freyaClient.SaveNewTemplate(context.Background(), &freya.TemplateData{
		TemplateName: "bombo",
		Data:         welcomeData,
	})

	if err != nil {
		panic(err)
	}

	if resp.Saved {
		log.Println("Template saved")
	}

	templates, err := freyaClient.GetAllTemplates(context.Background(), &empty.Empty{})
	if err != nil {
		panic(err)
	}

	pp.Println(templates)

}
