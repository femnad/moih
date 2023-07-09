package cmd

import (
	"bytes"
	"log"
	"os"
	"text/template"
)

func getHostname() (string, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return "", err
	}

	return hostname, nil
}

func getUsername() (string, error) {
	return os.Getenv("USER"), nil
}

func expandTemplate(text string) (string, error) {
	text = os.ExpandEnv(text)

	tmpl := template.New("moih")
	tmpl.Funcs(map[string]interface{}{
		"hostname": getHostname,
		"username": getUsername,
	})

	parsed, err := tmpl.Parse(text)
	if err != nil {
		log.Fatalf("error parsing template %s: %v", text, err)
	}

	out := bytes.Buffer{}
	err = parsed.Execute(&out, struct{}{})
	if err != nil {
		log.Fatalf("error executing template %s: %v", text, err)
	}

	return out.String(), nil
}
