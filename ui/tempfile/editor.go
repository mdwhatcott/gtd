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
	args := strings.Fields(EDITOR)
	return &Editor{
		editor: args[0],
		args:   args[1:],
	}
}

func (this *Editor) EditTempFile(initialContent string) (resultContent string) {
	name := createTempFile(initialContent)
	defer deleteFile(name)
	return this.editFile(name)
}

func createTempFile(content string) string {
	file, err := ioutil.TempFile("", "*.md")
	if err != nil {
		log.Fatal(err)
	}

	_, err2 := io.WriteString(file, content)
	if err2 != nil {
		log.Fatal(err2)
	}

	err3 := file.Close()
	if err3 != nil {
		log.Fatal(err2)
	}

	return file.Name()
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
