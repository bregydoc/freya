package main

import (
	"bytes"
	"gopkg.in/njern/gonexmo.v2"
	"io/ioutil"
	"text/template"
)

type NexmoSMSBackend struct {
	client *nexmo.Client
}

func NewNexmoSMSBackend(config *SMSConfig) (*NexmoSMSBackend, error) {
	nexmoClient, err := nexmo.NewClient(config.Key, config.Secret)
	if err != nil {
		return nil, err
	}

	return &NexmoSMSBackend{
		client: nexmoClient,
	}, nil
}

func (n *NexmoSMSBackend) SendSMS(config *SMSConfig, to *PhoneNumber, t *Template, params map[string]string) (SMSSendResponse, error) {
	data, err := ioutil.ReadAll(t.Data)
	if err != nil {
		return "", err
	}

	temp, err := template.New(t.Name).Parse(string(data))

	buffer := new(bytes.Buffer)

	if err = temp.Execute(buffer, params); err != nil {
		return "", err
	}

	message := &nexmo.SMSMessage{
		From:  config.From,
		To:    to.CountryCode + to.Number,
		Text:  buffer.String(),
		Type:  nexmo.Text,
		Class: nexmo.Standard,
	}

	m, err := n.client.SMS.Send(message)
	if err != nil {
		return "", err
	}

	return "", nil
}

func (n *NexmoSMSBackend) GetBalance() (float64, error) {
	return n.client.Account.GetBalance()
}
