package main

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

type StandardResponse struct {
	Data  interface{} `json:"data"`
	Error interface{} `json:"error"`
}

type SendEmailRequest struct {
	To      string `json:"to"`
	Subject string `json:"subject"`

	TemplateName    string      `json:"template_name"`
	DataForTemplate interface{} `json:"data_for_template"`
}

func main() {
	engine := gin.Default()
	engine.GET("/login", func(c *gin.Context) {
		credentials := new(AdminCredentials)

		c.BindJSON(credentials)

		credentialIsCorrect := false
		for _, cr := range GlobalConfig.Credentials {
			if cr.Username == credentials.Username {
				if cr.Password == credentials.Password {
					credentialIsCorrect = true
				}
			}
		}

		if !credentialIsCorrect {
			c.JSON(http.StatusUnauthorized, StandardResponse{
				Error: "Invalid credentials",
			})
			return
		}

		c.JSON(http.StatusOK, StandardResponse{
			Data: "ok, you're the boss",
		})

	})

	engine.Use(gin.BasicAuth(gin.Accounts{
		"bregymr": "malpartida1",
	}))

	engine.POST("/upload-new-template", func(c *gin.Context) {
		template, err := c.FormFile("template")
		if err != nil {
			c.JSON(http.StatusInternalServerError, StandardResponse{
				Error: err.Error(),
			})
			return
		}
		templateName := c.PostForm("template_name")

		if templateName == "" {
			c.JSON(http.StatusInternalServerError, StandardResponse{
				Error: "please, put the name of your template",
			})
			return
		}

		fTemplate, err := template.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, StandardResponse{
				Error: err.Error(),
			})
			return
		}

		data, err := ioutil.ReadAll(fTemplate)
		if err != nil {
			c.JSON(http.StatusInternalServerError, StandardResponse{
				Error: err.Error(),
			})
			return
		}

		completedTemplate, err := CreateNewTemplate(templateName, data)
		if err != nil {
			c.JSON(http.StatusInternalServerError, StandardResponse{
				Error: err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, StandardResponse{
			Data:  completedTemplate,
			Error: nil,
		})

	})

	engine.POST("/send-email", func(c *gin.Context) {
		request := new(SendEmailRequest)

		c.BindJSON(request)

		emailRequest := NewRequest([]string{request.To}, request.Subject)

		err := SendMailFromSavedTemplate(emailRequest, request.TemplateName, request.DataForTemplate)

		if err != nil {
			c.JSON(http.StatusInternalServerError, StandardResponse{
				Error: err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, StandardResponse{
			Data:  "Email sent, boss!",
			Error: nil,
		})

	})

	engine.Run(":3300")
}
