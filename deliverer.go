package main

func SendMailFromSavedTemplate(request *Request, templateName string, data interface{}) error {
	template, err := GetTemplateByName(templateName)
	if err != nil {
		return err
	}

	dataTemplate, err := ReadTemplate(template)
	if err != nil {
		return err
	}

	err = request.Send(dataTemplate, data)
	if err != nil {
		return err
	}

	return nil
}

func FastSendFromSavedTemplate(to string, subject string, templateName string, data interface{}) error {
	r := NewRequest([]string{to}, subject)
	return SendMailFromSavedTemplate(r, templateName, data)
}

func FastSendPlainData(to string, subject string, body []byte) error {
	r := NewRequest([]string{to}, subject)
	return r.SendPlain(body)
}
