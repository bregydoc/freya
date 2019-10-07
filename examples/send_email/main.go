package main

import (
	"context"
	freya "github.com/bregydoc/freya/proto"
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

	in := &freya.SendEmailParams{
		Subject: "Hello",
		Params: map[string]string{
			"ChargedAmount": "13.23 PEN",
			"Discount":      "-24.21 PEN",
			"Description":   "Bombo test",
			"Name":          "Bombo",
			"ChargedCard":   "VISA ****2312",
			"Product":       "BomboCard",
			"SubTotal":      "123.213 PEN",
			"TotalCost":     "123.213 PEN",
			"Code":          "AS12398OJW1",
			"Date":          "Tue Feb 26, 2019",
		},
		TemplateName: "bombo",
		To:           map[int32]string{0: "bregy.malpartida@utec.edu.pe", 1: "mateo@bombo.pe"},
	}

	res, err := freyaClient.SendEmail(context.Background(), in)

	if err != nil {
		panic(err)
	}

	log.Printf("%v\n", res)

}
