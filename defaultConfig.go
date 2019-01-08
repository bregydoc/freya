package main

const DefaultMIME = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

const DefaultRelativeDBFolder = "./udb"
const DefaultPlansDBName = "plans"
const DefaultLoggerDBName = "logger"
const DefaultTemplatesDBName = "templates"

const DefaultMinioEndpoint = "127.0.0.1:9000"
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
  },
}`
