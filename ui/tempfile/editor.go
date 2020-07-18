package tempfile

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

type Editor struct {
	editor string
	args   []string
}

func NewEditor() *Editor {
	EDITOR := os.Getenv("EDITOR")
	if EDITOR == "" {
		log.Panic("$EDITOR environment variable not set.")
	}
	ARGS := strings.Fields(EDITOR)
	return &Editor{
		editor: ARGS[0],
		args:   ARGS[1:],
	}
}

func (this *Editor) EditTempFile(initialContent string) (resultContent_ string) {
	NAME := createTempFile(initialContent)
	defer deleteFile(NAME)
	return this.editFile(NAME)
}

func createTempFile(content string) string {
	FILE, ERR := ioutil.TempFile("", "*.md")
	if ERR != nil {
		log.Fatal(ERR)
	}

	_, ERR2 := io.WriteString(FILE, content)
	if ERR2 != nil {
		log.Fatal(ERR2)
	}

	ERR3 := FILE.Close()
	if ERR3 != nil {
		log.Fatal(ERR2)
	}

	return FILE.Name()
}

func (this *Editor) editFile(name string) string {
	_, ERR1 := exec.Command(this.editor, append(this.args, name)...).CombinedOutput()
	if ERR1 != nil {
		log.Fatal(ERR1)
	}

	ALL, ERR2 := ioutil.ReadFile(name)
	if ERR2 != nil {
		log.Fatal(ERR2)
	}

	return string(ALL)
}

func deleteFile(name string) {
	ERR := os.Remove(name)
	if ERR != nil {
		log.Fatal(ERR)
	}
}
