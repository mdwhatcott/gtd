package tempfile

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
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
	file, err := ioutil.TempFile("", "")
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
	_, err2 := exec.Command("subl", name, "--wait").CombinedOutput()
	if err2 != nil {
		log.Fatal(err2)
	}

	all, err3 := ioutil.ReadFile(name)
	if err3 != nil {
		log.Fatal(err3)
	}

	return string(all)
}

func deleteFile(name string) {
	err4 := os.Remove(name)
	if err4 != nil {
		log.Fatal(err4)
	}
}
