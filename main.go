package main

import (
	"fmt"
	"github.com/bregydoc/freya/freyacon/go"
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

	f, err := NewFreya(config)
	if err != nil {
		log.Fatalf("failed to create freya repository %v", err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", config.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	freya.RegisterFreyaServer(grpcServer, &FreyaService{
		repo: f,
	})

	log.Printf("listening on :%d\n", config.Port)

	err = grpcServer.Serve(lis)

	if err != nil {
		log.Fatalf("failed on serve: %v", err)
	}

}
