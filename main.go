package main

import (
	"fmt"
	"github.com/bregydoc/freya/freyacon/go"
	"github.com/k0kubun/pp"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

func main() {
	var config *FreyaConfig
	if _, err := os.Stat(freyaConfigFileName); os.IsNotExist(err) {
		config = GetDefaultConfig()
		log.Println("Loading default config")
	} else {
		config = ReadConfig(freyaConfigFileName)
		log.Println("Loading freya.config.yml config file")
	}

	mailJet := NewMailJetMailBackend(config.Mail)
	nexmoSMS, err := NewNexmoSMSBackend(config.SMS)

	if err != nil {
		log.Fatalf("failed to create nexmo client %v", err)
	}

	f, err := NewFreya(config, mailJet, nexmoSMS)

	pp.Println(f.Config)

	if err != nil {
		log.Fatalf("failed to create freya repository %v", err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", config.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	freya.RegisterFreyaServer(grpcServer, &Service{
		repo: f,
	})

	log.Printf("listening on :%d\n", config.Port)

	err = grpcServer.Serve(lis)

	if err != nil {
		log.Fatalf("failed on serve: %v", err)
	}

}
