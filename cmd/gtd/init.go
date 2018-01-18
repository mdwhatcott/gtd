package main

import (
	"log"
	"os"

	"github.com/mdwhatcott/gtd"
)

func init() {
	log.SetFlags(log.Llongfile | log.Lmicroseconds)
	//ensureRequiredDirectories() // TODO
}
func ensureRequiredDirectories() {
	os.MkdirAll(gtd.FolderRoot, os.ModePerm)
	os.MkdirAll(gtd.FolderActions, os.ModePerm)
	os.MkdirAll(gtd.FolderProjects, os.ModePerm)
	os.MkdirAll(gtd.FolderArchive, os.ModePerm)
	os.MkdirAll(gtd.FolderSomeday, os.ModePerm)
	os.MkdirAll(gtd.FolderTickler, os.ModePerm)
	os.MkdirAll(gtd.FolderMaybe, os.ModePerm)
}
