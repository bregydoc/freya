package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
)

const freyaConfigFileName = "./freya.config.yml"

type FreyaConfig struct {
	Mail MailConfig `yaml:"mail"`
	SMS  SMSConfig  `yaml:"sms"`

	DB DataBaseConfig `yaml:"db_config"`

	Minio MinioConfig `yaml:"minio"`

	Port int64 `yaml:"port"`
}

type MailConfig struct {
	Mime     string `yaml:"mime"`
	Server   string `yaml:"server"`
	Port     int    `yaml:"port"`
	Email    string `yaml:"email"`
	Password string `yaml:"password"`

	MetaData SenderConfig `yaml:"metadata"`
}

type SMSConfig struct {
	Backend  string `yaml:"backend"`
	Endpoint string `yaml:"endpoint"`
	Key      string `yaml:"key"`
	Secret   string `yaml:"secret"`
}

type SenderConfig struct {
	FromName  string `yaml:"from_name"`
	FromEmail string `yaml:"from_email"`
}

type DataBaseConfig struct {
	AbsoluteFolder string `yaml:"absolute_folder"`
	RelativeFolder string `yaml:"relative_folder"`

	PlansDBName     string `yaml:"plans_db_name"`
	LoggerDBName    string `yaml:"logger_db_name"`
	TemplatesDBName string `yaml:"templates_db_name"`
}

type MinioConfig struct {
	Endpoint        string `yaml:"endpoint"`
	AccessKeyID     string `yaml:"access_key_id"`
	SecretAccessKey string `yaml:"secret_access_key"`
	UseSSL          bool   `yaml:"use_ssl"`

	BucketName string `yaml:"bucket_name"`
	Location   string `yaml:"location"`
}

func getEnvVarFromConfig(config string) string {
	envVar := strings.Replace(config, "${", "", 1)
	envVar = strings.Replace(envVar, "}", "", 1)
	return os.Getenv(envVar)
}

func inflateConfigWithEnvs(c *FreyaConfig) {

	if strings.HasPrefix(c.SMS.Key, "$") {
		c.SMS.Key = getEnvVarFromConfig(c.SMS.Key)
	}

	if strings.HasPrefix(c.SMS.Secret, "$") {
		c.SMS.Secret = getEnvVarFromConfig(c.SMS.Secret)
	}

	if strings.HasPrefix(c.Mail.Email, "$") {
		c.Mail.Email = getEnvVarFromConfig(c.Mail.Email)
	}

	if strings.HasPrefix(c.Mail.Password, "$") {
		c.Mail.Password = getEnvVarFromConfig(c.Mail.Password)
	}

}

func ReadConfig(filename string) *FreyaConfig {
	configData, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	c := new(FreyaConfig)
	err = yaml.Unmarshal(configData, c)
	if err != nil {
		panic(err)
	}

	inflateConfigWithEnvs(c)

	return c
}

func GetDefaultConfig() *FreyaConfig {
	c := new(FreyaConfig)

	err := yaml.Unmarshal([]byte(defaultConfigContent), c)
	if err != nil {
		panic(err)
	}

	inflateConfigWithEnvs(c)

	return c
}
