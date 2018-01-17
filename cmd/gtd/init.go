package main

import (
	"os"

	"github.com/mdwhatcott/gtd"
)

func init() {
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
	os.MkdirAll(gtd.FolderFocus, os.ModePerm)
	os.MkdirAll(gtd.FolderGoals, os.ModePerm)
	os.MkdirAll(gtd.FolderVision, os.ModePerm)
	os.MkdirAll(gtd.FolderPurpose, os.ModePerm)
}
