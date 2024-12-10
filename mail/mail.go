package mail

import (
	"bytes"
	"fmt"
	"net/smtp"
	"os"
	"text/template"
)

type Mailer struct {
	smtpHost string
	smtpPort string
	smtpUser string
	smtpPass string
}

func New(smtpHost string, smtpPort string, smtpUser string, smtpPass string) Mailer {
	return Mailer{
		smtpHost: smtpHost,
		smtpPort: smtpPort,
		smtpUser: smtpUser,
		smtpPass: smtpPass,
	}
}

func (m *Mailer) SendMail(to []string, templateName string, data interface{}) error {
	templateContent := getTemplate(templateName)
	if templateContent == "" {
		return fmt.Errorf("template %s not found", templateName)
	}

	tmpl, err := template.New("email").Parse(templateContent)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	if err := tmpl.Execute(buf, data); err != nil {
		return err
	}

	auth := smtp.PlainAuth(m.smtpUser, m.smtpUser, m.smtpPass, m.smtpHost)
	return smtp.SendMail(m.smtpHost+":"+m.smtpPort, auth, m.smtpUser, to, buf.Bytes())
}

func getTemplate(templateName string) string {
	templatePath := "templates/" + templateName + ".gotmpl"
	content, err := os.ReadFile(templatePath)
	if err != nil {
		fmt.Printf("Error leyendo plantilla: %v\n", err)
		return ""
	}
	return string(content)
}
