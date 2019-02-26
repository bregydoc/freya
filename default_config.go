package main

const DefaultMIME = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
const DefaultAbsoluteDBFolder = "/udb/udb"
const DefaultRelativeDBFolder = "./udb"
const DefaultPlansDBName = "plans"
const DefaultLoggerDBName = "logger"
const DefaultTemplatesDBName = "templates"

const DefaultFromSMS = "Freya"
const DefaultStorage = "/storage"

const DefaultPort = 10000

const defaultConfigContent = `---
storage: "/storage" # It's for the docker volume'
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

	if newConfig.DB.AbsoluteFolder == "" {
		newConfig.DB.AbsoluteFolder = DefaultAbsoluteDBFolder
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

	if newConfig.Storage == "" {
		newConfig.Storage = DefaultStorage
	}

	if newConfig.Port == 0 {
		newConfig.Port = DefaultPort
	}

	return &newConfig
}
