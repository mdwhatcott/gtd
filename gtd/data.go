package gtd

import (
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var (
	Now            = time.Now()
	NotesRoot      = os.Getenv("NOTESPATH")
	FolderRoot     = root(os.Getenv("GTDPATH"))
	FolderActions  = filepath.Join(FolderRoot, "0-next-actions")
	FolderProjects = filepath.Join(FolderRoot, "1-projects")
	FolderComplete = filepath.Join(FolderRoot, "1-projects-archive", strconv.Itoa(Now.Year()))
	FolderRejected = filepath.Join(FolderComplete, "rejected")
	FolderSomeday  = filepath.Join(FolderRoot, "1-projects-someday")
	FolderTickler  = filepath.Join(FolderSomeday, "1-tickler")
	FolderMaybe    = filepath.Join(FolderRoot, "1-projects-tentative")
)

func root(proposed string) string {
	if proposed != "" {
		return proposed
	}
	return filepath.Join(os.Getenv("HOME"), "Documents", "gtd")
}
