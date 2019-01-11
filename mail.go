package main

func (f *Freya) SendMail(templateName string, params interface{}, subject string, to []string) error {
	template, err := f.GetTemplateByName(templateName)
	if err != nil {
		return err
	}

	request := newRequest(to, subject)

	err = request.Send(f.Config.Mail, template, params)
	if err != nil {
		return err
	}

	return nil
}
