package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"os/exec"
	"text/template"
)

func executeTemplate(template *template.Template, data interface{}) string {
	var content bytes.Buffer
	if err := template.Execute(&content, data); err != nil {
		log.Fatalln("Could not execute template:", err)
	}
	return content.String()
}

func create(path, content string) {
	if err := ioutil.WriteFile(path, []byte(content), 0644); err != nil {
		log.Fatalln("Could not create file:", err)
	}
}

func edit(path string) {
	if err := exec.Command("subl", "-wait", path).Run(); err != nil {
		log.Fatalln("Could not edit file:", err)
	}
}

func commit(path string) {
	if err := exec.Command("stree", path).Run(); err != nil {
		log.Fatalln("Could not open source tree:", err)
	}
}
