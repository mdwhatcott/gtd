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
	FolderActions  = join(FolderRoot, "0-next-actions")
	FolderProjects = join(FolderRoot, "1-projects")
	FolderComplete = join(FolderRoot, "1-projects-archive", strconv.Itoa(Now.Year()))
	FolderRejected = join(FolderComplete, "rejected")
	FolderSomeday  = join(FolderRoot, "1-projects-someday")
	FolderTickler  = join(FolderSomeday, "1-tickler")
	FolderMaybe    = join(FolderRoot, "1-projects-tentative")
)

func root(proposed string) string {
	if proposed != "" {
		return proposed
	}
	return filepath.Join(os.Getenv("HOME"), "Documents", "gtd")
}

var join = filepath.Join