package main

const defaultConfigContent = `{
  "mime": "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n",
  "server": "in-v3.mailjet.com",
  "port": 587,
  "email": "xxx",
  "password": "yyy",
  "meta_data": {
    "from_name": "Freya",
    "from_email": "hello@bombo.pe"
  },
  "db_config": {
    "absolute_folder": "/Users/macBook/Documents/freyabuf/udb",
    "relative_folder": "./udb",

    "plans_db_name": "plans",
    "logger_db_name": "logger",
    "templates_db_name": "templates"
  },
  "minio_storage_config": {
    "endpoint": "127.0.0.1:9000",
    "access_key_id": "AKIAIOSFODNN7EXAMPLE",
    "secret_access_key": "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
    "use_ssl": false,

    "bucket_name": "freyabuf",
    "location": "us-east-1"
  },
  "credentials": [
    {
      "username": "bregymr",
      "password": "malpartida1"
    }
  ]
}`
