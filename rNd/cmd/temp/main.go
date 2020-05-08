package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

func main() {
	file, err := ioutil.TempFile("", "")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(file.Name())

	out, err := exec.Command("subl", file.Name(), "--wait").CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(out))

	err = file.Close()
	if err != nil {
		log.Fatal(err)
	}

	all, err := ioutil.ReadFile(file.Name())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Content:")
	fmt.Println(string(all))

	err = os.Remove(file.Name())
	if err != nil {
		log.Fatal(err)
	}

}
