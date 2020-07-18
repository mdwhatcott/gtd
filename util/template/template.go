package template

import (
	"log"
	"strings"
	"text/template"
)

func MustExecute(raw string, data interface{}) string {
	OUT, ERR := Execute(raw, data)
	if ERR != nil {
		log.Panic(ERR)
	}
	return OUT
}

func Execute(raw string, data interface{}) (string, error) {
	BUFFER := new(strings.Builder)
	PARSED := template.Must(template.New("").Parse(raw))
	ERR := PARSED.Execute(BUFFER, data)
	return BUFFER.String(), ERR
}
