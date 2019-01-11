package main

import "errors"

func (f *Freya) SendSMS(templateName string, params interface{}, to *PhoneNumber) error {
	t, err := f.GetTemplateByName(templateName)
	if err != nil {
		return err
	}

	r, err := f.SMSBackend.SendSMS(f.Config.SMS, to, t, params)
	if err != nil {
		return err
	}

	if r != "" {
		return errors.New(string(r))
	}

	return nil

}
