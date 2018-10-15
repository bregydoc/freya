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
