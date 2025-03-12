package email_template

import (
	"bytes"
	"html/template"
	"os"
	"path/filepath"
	"runtime"
)

type ConfirmAccountData struct {
	ConfirmationURL string
}
type ForgotPasswordData struct {
	Token string
}

const (
	CONFIRM_ACCOUNT = "register.html"
	FORGOT_PASSWORD = "forgot_password.html"
)

func getCurrentFilePath() string {
	_, filePath, _, _ := runtime.Caller(0)
	return filePath
}

func LoadTemplate(fileName string, data any) (string, error) {
	dir := filepath.Dir(getCurrentFilePath())
	content, err := os.ReadFile(filepath.Join(dir, fileName))
	if err != nil {
		return "", err
	}

	tmpl, err := template.New("email").Parse(string(content))
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
