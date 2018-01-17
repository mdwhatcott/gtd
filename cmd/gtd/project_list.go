package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func listProjects(inputs []string) {
	review := listProjectFlags.Bool("review", false, "When set, review each project via a REPL and text editor sessions.")
	listProjectFlags.Parse(inputs)

	if *review {
		fmt.Println("Not implemented") // TODO
		return
	}
	dir, err := ioutil.ReadDir(filepath.Join(os.Getenv("HOME"), "Documents/gtd/1-projects"))
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range dir {
		// TODO: parse and display the h1
		fmt.Println(file.Name())
	}
}
