package main

import (
	"github.com/minio/minio-go"
	"github.com/nanobox-io/golang-scribble"
	"log"
)

type FreyaConfig struct {
	Mime     string `json:"mime"`
	Server   string `json:"server"`
	Port     int    `json:"port"`
	Email    string `json:"email"`
	Password string `json:"password"`

	MetaData SenderConfig `json:"meta_data"`

	DBConfig DataBaseConfig `json:"db_config"`

	MinioStorageConfig MinioConfig `json:"minio_config"`
}

type SenderConfig struct {
	FromName  string `json:"from_name"`
	FromEmail string `json:"from_email"`
}

type DataBaseConfig struct {
	AbsoluteFolder string `json:"absolute_folder"`
	RelativeFolder string `json:"relative_folder"`

	PlansDBName     string `json:"plans_db_name"`
	LoggerDBName    string `json:"logger_db_name"`
	TemplatesDBName string `json:"templates_db_name"`
}

type MinioConfig struct {
	Endpoint        string `json:"endpoint"`
	AccessKeyID     string `json:"access_key_id"`
	SecretAccessKey string `json:"secret_access_key"`
	UseSSL          bool   `json:"use_ssl"`

	BucketName string `json:"bucket_name"`
	Location   string `json:"location"`
}

var ScribbleDriver *scribble.Driver

var GlobalConfig *FreyaConfig

var MinioClient *minio.Client

func init() {

	var err error

	// TODO: Bregy, please fill the default config for Freya
	GlobalConfig = &FreyaConfig{}

	ScribbleDriver, err = scribble.New(GlobalConfig.DBConfig.AbsoluteFolder, nil)
	if err != nil {
		panic(err)
	}

	MinioClient, err = minio.New(
		GlobalConfig.MinioStorageConfig.Endpoint,
		GlobalConfig.MinioStorageConfig.AccessKeyID,
		GlobalConfig.MinioStorageConfig.SecretAccessKey,
		GlobalConfig.MinioStorageConfig.UseSSL,
	)

	if err != nil {
		panic(err)
	}

	bucketName := GlobalConfig.MinioStorageConfig.BucketName

	err = MinioClient.MakeBucket(bucketName, GlobalConfig.MinioStorageConfig.Location)
	if err != nil {

		exists, err := MinioClient.BucketExists(bucketName)
		if err == nil && exists {
			log.Printf("We already own %s\n", bucketName)
		} else {
			log.Fatalln(err)
		}
	}

}
