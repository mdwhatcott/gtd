package tempfile

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

type Editor struct{}

func NewEditor() *Editor {
	return &Editor{}
}

func (*Editor) EditTempFile(initialContent string) (resultContent string) {
	name := createTempFile(initialContent)
	defer deleteFile(name)
	return editFile(name)
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

func editFile(name string) string {
	EDITOR := os.Getenv("EDITOR")
	if EDITOR == "" {
		log.Panic("$EDITOR environment variable not set.")
	}
	ARGS := strings.Fields(EDITOR) // Splitting on space assumes a very simple value in the $EDITOR variable...
	ARGS = append(ARGS, name)
	_, ERR1 := exec.Command(ARGS[0], ARGS[1:]...).CombinedOutput()
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
