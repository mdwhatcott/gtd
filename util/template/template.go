package template

import (
	"log"
	"strings"
	"text/template"
)

func MustExecute(_raw string, _data interface{}) string {
	OUT, ERR := Execute(_raw, _data)
	if ERR != nil {
		log.Panic(ERR)
	}
	return OUT
}

func Execute(_raw string, _data interface{}) (string, error) {
	BUFFER := new(strings.Builder)
	PARSED := template.Must(template.New("").Parse(_raw))
	ERR := PARSED.Execute(BUFFER, _data)
	return BUFFER.String(), ERR
}
