package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

func main() {
	name := createTempFile()
	fmt.Println(name)

	all := editFile(name)
	fmt.Println(all)

	deleteFile(name)

}

func createTempFile() string {
	file, err := ioutil.TempFile("", "")
	if err != nil {
		log.Fatal(err)
	}

	err2 := file.Close()
	if err2 != nil {
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
