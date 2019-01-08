package main

const DefaultMIME = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

const DefaultRelativeDBFolder = "./udb"
const DefaultPlansDBName = "plans"
const DefaultLoggerDBName = "logger"
const DefaultTemplatesDBName = "templates"

const DefaultMinioEndpoint = "minio:9000"
const DefaultMinioAccessKey = "AKIAIOSFODNN7EXAMPLE"
const DefaultMinioSecretKey = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
const DefaultMinioUseSSL = false
const DefaultMinioBucketName = "freya"
const DefaultMinioLocation = "us-east-1"

const DefaultAdminUsername = "admin"
const DefaultAdminPassword = "admin"

const defaultConfigContent = `{
  "server": "in-v3.mailjet.com",
  "port": 587,
  "email": "xxx",
  "password": "yyy",
  "meta_data": {
    "from_name": "Freya",
    "from_email": "hello@bombo.pe"
  }
}`

func FillConfigWithDefaults(config *FreyaConfig) *FreyaConfig {
	newConfig := *config
	if newConfig.Mime == "" {
		newConfig.Mime = DefaultMIME
	}

	voidDatabaseConfig := DataBaseConfig{}
	if newConfig.DBConfig == voidDatabaseConfig {
		newConfig.DBConfig = DataBaseConfig{
			RelativeFolder:  DefaultRelativeDBFolder,
			PlansDBName:     DefaultPlansDBName,
			LoggerDBName:    DefaultLoggerDBName,
			TemplatesDBName: DefaultTemplatesDBName,
		}
	}

	if newConfig.DBConfig.RelativeFolder == "" {
		newConfig.DBConfig.RelativeFolder = DefaultRelativeDBFolder
	}
	if newConfig.DBConfig.PlansDBName == "" {
		newConfig.DBConfig.PlansDBName = DefaultPlansDBName
	}
	if newConfig.DBConfig.LoggerDBName == "" {
		newConfig.DBConfig.LoggerDBName = DefaultLoggerDBName
	}
	if newConfig.DBConfig.TemplatesDBName == "" {
		newConfig.DBConfig.TemplatesDBName = DefaultTemplatesDBName
	}

	voidMinioConfig := MinioConfig{}
	if newConfig.MinioStorageConfig == voidMinioConfig {
		newConfig.MinioStorageConfig = MinioConfig{
			Endpoint:        DefaultMinioEndpoint,
			AccessKeyID:     DefaultMinioAccessKey,
			SecretAccessKey: DefaultMinioSecretKey,

			BucketName: DefaultMinioBucketName,
			Location:   DefaultMinioLocation,
		}
	}

	if newConfig.MinioStorageConfig.Endpoint == "" {
		newConfig.MinioStorageConfig.Endpoint = DefaultMinioEndpoint
	}
	if newConfig.MinioStorageConfig.AccessKeyID == "" {
		newConfig.MinioStorageConfig.AccessKeyID = DefaultMinioAccessKey
	}
	if newConfig.MinioStorageConfig.SecretAccessKey == "" {
		newConfig.MinioStorageConfig.SecretAccessKey = DefaultMinioSecretKey
	}
	if newConfig.MinioStorageConfig.BucketName == "" {
		newConfig.MinioStorageConfig.BucketName = DefaultMinioBucketName
	}
	if newConfig.MinioStorageConfig.Location == "" {
		newConfig.MinioStorageConfig.Location = DefaultMinioLocation
	}

	if newConfig.Credentials == nil {
		newConfig.Credentials = []*AdminCredentials{{
			Username: DefaultAdminUsername,
			Password: DefaultAdminPassword,
		},
		}
	}

	return &newConfig

}
