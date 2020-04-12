package external

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

type YesNoResponse int

const (
	NO YesNoResponse = iota
	YES
)

func PromptYesNo(assumed YesNoResponse, question string) YesNoResponse {
	if assumed == YES {
		question += " [Y/n] "
	} else if assumed == NO {
		question += " [y/N] "
	}
	response := Prompt(question)
	response = strings.TrimSpace(response)
	response = strings.ToUpper(response)

	if response == "Y" {
		return YES
	} else if response == "N" {
		return NO
	} else {
		return assumed
	}
}

func Prompt(message string) string {
	fmt.Print(message)
	fmt.Print(" ")
	return ReadLine()
}

func PromptDuration(prompt string) time.Duration {
	answer := "invalid"
	for len(answer) > 0 {
		duration, err := time.ParseDuration(answer)
		if err == nil {
			return duration
		}
		answer = Prompt(prompt)
	}
	return 0
}

func ReadLine() string {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}

func Flags(usage string) *flag.FlagSet {
	var flagLog = log.New(os.Stderr, "", 0)
	set := flag.NewFlagSet("", flag.ExitOnError)
	set.Usage = func() {
		flagLog.Println(usage)
		set.PrintDefaults()
	}
	return set
}

func ExecuteTemplate(template *template.Template, data interface{}) string {
	var content bytes.Buffer
	if err := template.Execute(&content, data); err != nil {
		log.Fatalln("Could not execute template:", err)
	}
	return content.String()
}

func MakeDirectory(folder string) {
	if err := os.MkdirAll(folder, os.ModePerm); err != nil {
		log.Fatalln("Could not create directory:", folder, err)
	}
}

func CreateFile(path, content string) {
	if err := ioutil.WriteFile(path, []byte(content), 0644); err != nil {
		log.Fatalln("Could not create file:", err)
	}
}

func MoveFile(oldPath, newPath string) {
	if err := os.Rename(oldPath, newPath); err != nil {
		log.Fatalln("Could not move file:", oldPath, newPath, err)
	}
}

func OpenTextEditorAndWait(path string) {
	if err := exec.Command("subl", "-wait", path).Run(); err != nil {
		log.Fatalln("Could not open Sublime Text:", err)
	}
}

func OpenTextEditor(path string) {
	if err := exec.Command("subl", path).Run(); err != nil {
		log.Fatalln("Could not edit file:", err)
	}
}

func ListDirectory(folder string) []os.FileInfo {
	dir, err := ioutil.ReadDir(folder)
	if err != nil {
		log.Fatalln("Could not read directory:", folder, err)
	}
	return dir
}

func ScanFile(path string) *bufio.Scanner {
	all, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalln("Could not read file:", path, err)
	}
	return bufio.NewScanner(bytes.NewReader(all))
}

func Commit(path string) {
	if err := exec.Command("smerge", path).Run(); err != nil {
		log.Fatalln("Could not open sublime merge:", err)
	}
	Prompt("<Enter> to continue...")
}

func DeleteContents(folder string) {
	for _, file := range ListDirectory(folder) {
		path := filepath.Join(folder, file.Name())
		if err := os.Remove(path); err != nil {
			log.Fatalln("Could not remove specified path:", path, err)
		}
	}
}

func Navigate(address string) {
	if err := exec.Command("open", address).Run(); err != nil {
		log.Fatalln("Could not open browser:", err)
	}
}