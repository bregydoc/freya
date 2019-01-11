package main

const DefaultMIME = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

const DefaultRelativeDBFolder = "./udb"
const DefaultPlansDBName = "plans"
const DefaultLoggerDBName = "logger"
const DefaultTemplatesDBName = "templates"

const DefaultMinioEndpoint = "minio:9000"
const DefaultMinioAccessKey = "AKIAIOSFODNN7EXAMPLE"
const DefaultMinioSecretKey = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"

//const DefaultMinioUseSSL = false
const DefaultMinioBucketName = "freya"
const DefaultMinioLocation = "us-east-1"

const DefaultFromSMS = "Freya"

const defaultConfigContent = `---
sms:
  backend: "nexmo" # Now Only support Nexmo Backend
  endpoint: "https://rest.nexmo.com/sms/json"
  key: "${NEXMO_KEY}" # USe $ to extract from env variable
  secret: "${NEXMO_SECRET}"
mail:
  server: "in-v3.mailjet.com"
  port: 587
  email: "${MAILJET_EMAIL}"
  password: "${MAILJET_PASSWORD}"
  metadata:
    from_name: "Freya"
    from_email: "example@example.com"
minio:
  endpoint: "minio:9000"
`

func FillConfigWithDefaults(config *FreyaConfig) *FreyaConfig {
	newConfig := *config
	if newConfig.Mail.Mime == "" {
		newConfig.Mail.Mime = DefaultMIME
	}

	if newConfig.SMS.From == "" {
		newConfig.SMS.From = DefaultFromSMS
	}

	if newConfig.DB == nil {
		newConfig.DB = &DataBaseConfig{
			RelativeFolder:  DefaultRelativeDBFolder,
			PlansDBName:     DefaultPlansDBName,
			LoggerDBName:    DefaultLoggerDBName,
			TemplatesDBName: DefaultTemplatesDBName,
		}
	}

	if newConfig.DB.RelativeFolder == "" {
		newConfig.DB.RelativeFolder = DefaultRelativeDBFolder
	}
	if newConfig.DB.PlansDBName == "" {
		newConfig.DB.PlansDBName = DefaultPlansDBName
	}
	if newConfig.DB.LoggerDBName == "" {
		newConfig.DB.LoggerDBName = DefaultLoggerDBName
	}
	if newConfig.DB.TemplatesDBName == "" {
		newConfig.DB.TemplatesDBName = DefaultTemplatesDBName
	}

	if newConfig.Minio == nil {
		newConfig.Minio = &MinioConfig{
			Endpoint:        DefaultMinioEndpoint,
			AccessKeyID:     DefaultMinioAccessKey,
			SecretAccessKey: DefaultMinioSecretKey,

			BucketName: DefaultMinioBucketName,
			Location:   DefaultMinioLocation,
		}
	}

	if newConfig.Minio.Endpoint == "" {
		newConfig.Minio.Endpoint = DefaultMinioEndpoint
	}
	if newConfig.Minio.AccessKeyID == "" {
		newConfig.Minio.AccessKeyID = DefaultMinioAccessKey
	}
	if newConfig.Minio.SecretAccessKey == "" {
		newConfig.Minio.SecretAccessKey = DefaultMinioSecretKey
	}
	if newConfig.Minio.BucketName == "" {
		newConfig.Minio.BucketName = DefaultMinioBucketName
	}
	if newConfig.Minio.Location == "" {
		newConfig.Minio.Location = DefaultMinioLocation
	}

	return &newConfig

}
