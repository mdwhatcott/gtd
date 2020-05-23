package main

import (
	"log"
	"path/filepath"

	"github.com/mdwhatcott/gtd/legacy/external"
	"github.com/mdwhatcott/gtd/legacy/gtd"
)

func init() {
	log.SetFlags(log.Lshortfile)
	ensureRequiredDirectories()
}

func ensureRequiredDirectories() {
	external.MakeDirectory(gtd.FolderRoot)
	external.MakeDirectory(gtd.FolderActions)
	external.MakeDirectory(gtd.FolderProjects)
	external.MakeDirectory(gtd.FolderComplete)
	external.MakeDirectory(gtd.FolderSomeday)
	external.MakeDirectory(gtd.FolderTickler)
	external.MakeDirectory(gtd.FolderMaybe)
	external.MakeDirectory(gtd.FolderRejected)
}

var join = filepath.Join
