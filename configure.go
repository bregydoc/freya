package main

import (
	"encoding/json"
	"github.com/minio/minio-go"
	"github.com/nanobox-io/golang-scribble"
	"io/ioutil"
	"log"
	"os"
)

type FreyaConfig struct {
	Mime     string `json:"mime"`
	Server   string `json:"server"`
	Port     int    `json:"port"`
	Email    string `json:"email"`
	Password string `json:"password"`

	MetaData SenderConfig `json:"meta_data"`

	DBConfig DataBaseConfig `json:"db_config"`

	MinioStorageConfig MinioConfig `json:"minio_storage_config"`

	Credentials []*AdminCredentials `json:"credentials"`
}

type AdminCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
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

func ReadConfig(filename string) *FreyaConfig {
	configData, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	c := new(FreyaConfig)

	err = json.Unmarshal(configData, c)
	if err != nil {
		panic(err)
	}

	return c
}

func GetDefaultConfig() *FreyaConfig {
	c := new(FreyaConfig)

	err := json.Unmarshal([]byte(defaultConfigContent), c)
	if err != nil {
		panic(err)
	}

	return c

}

func init() {

	log.Println("Executing init function...")

	//
	var err error

	// TODO: Bregy, please fill the default config for Freya
	if _, err := os.Stat("./freyabuf.config.json"); os.IsNotExist(err) {
		GlobalConfig = GetDefaultConfig()
	} else {
		GlobalConfig = ReadConfig("./freyabuf.config.json")
	}

	ScribbleDriver, err = scribble.New(GlobalConfig.DBConfig.AbsoluteFolder, nil)
	if err != nil {
		panic(err)
	}

	log.Println("Scribble setup done ✔︎")

	MinioClient, err = minio.New(
		GlobalConfig.MinioStorageConfig.Endpoint,
		GlobalConfig.MinioStorageConfig.AccessKeyID,
		GlobalConfig.MinioStorageConfig.SecretAccessKey,
		GlobalConfig.MinioStorageConfig.UseSSL,
	)

	if err != nil {
		panic(err)
	}
	log.Println("Minio Server setup done ✔︎")

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

	log.Println("Minio bucket created ✔︎")

}
