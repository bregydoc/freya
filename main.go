package main

import (
	"fmt"
	"google.golang.org/grpc/credentials"
	"log"
	"net"
	"os"

	"gopkg.in/yaml.v2"

	freya "github.com/bregydoc/freya/proto"
	"google.golang.org/grpc"
)

func main() {
	config := new(FreyaConfig)
	composeConfig := os.Getenv("FREYA_CONFIG")
	if composeConfig != "" {
		err := yaml.Unmarshal([]byte(composeConfig), config)
		if err != nil {
			log.Fatalf("failed to read FREYA_CONFIG from docker-compose.yml\n%v", err)
		}
	} else {
		if _, err := os.Stat(freyaConfigFileName); os.IsNotExist(err) {
			config = GetDefaultConfig()
			log.Println("Loading default config")
		} else {
			config = ReadConfig(freyaConfigFileName)
			log.Println("Loading freya.config.yml config file")
		}
	}

	mailJet := NewMailJetMailBackend(config.Mail)
	nexmoSMS, err := NewNexmoSMSBackend(config.SMS)

	if err != nil {
		log.Fatalf("failed to create nexmo client %v", err)
	}

	f, err := NewFreya(config, mailJet, nexmoSMS)

	cfg, _ := yaml.Marshal(config)
	fmt.Printf("===FREYA CONFIG===\n%s\n", string(cfg))

	if err != nil {
		log.Fatalf("failed to create freya repository %v", err)
	}

	// Check if we have TLS a secure connection
	withTLS := true
	certificate := "/run/secrets/grpc_cert"
	key := "/run/secrets/grpc_key"

	_, err = os.Open(certificate)
	if err != nil || os.IsNotExist(err) {
		withTLS = false
	}
	if withTLS {
		_, err = os.Open(key)
		if err != nil || os.IsNotExist(err) {
			withTLS = false
		}
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", config.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var grpcServer *grpc.Server
	if withTLS {
		log.Println("setting with TLS from \"secrets/grpc_cert\" and \"secrets/grpc_key\"")
		c, err := credentials.NewServerTLSFromFile(certificate, key)
		if err != nil {
			log.Fatalf("Failed to setup tls: %v", err)
		}
		grpcServer = grpc.NewServer(
			grpc.Creds(c),
		)
	} else {
		log.Println("setting without any security")
		grpcServer = grpc.NewServer()
	}

	freya.RegisterFreyaServer(grpcServer, &Service{
		repo: f,
	})

	log.Printf("listening on :%d\n", config.Port)

	err = grpcServer.Serve(lis)

	if err != nil {
		log.Fatalf("failed on serve: %v", err)
	}

}
