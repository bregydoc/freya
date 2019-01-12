package main

type MailBackend interface {
	SendMail(config *MailConfig, template *Template, params interface{}, subject string, to []string) error
}
