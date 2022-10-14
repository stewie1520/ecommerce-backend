package templates

import (
	"bytes"
	"html/template"

	"github.com/stewie1520/ecommerce-backend/internal/tools/path"
)

func ResolveVerifyOTPMailTemplate(otp string, appName string) (string, error) {
	tmpl, err := template.ParseFiles(path.PathCWD("internal/templates/mail/verify-otp.tmpl"))
	if err != nil {
		return "", err
	}

	var wr bytes.Buffer

	tmpl.Execute(&wr, map[string]interface{}{
		"OTP":     otp,
		"AppName": appName,
	})

	return wr.String(), nil
}
