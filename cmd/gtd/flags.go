package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

var (
	reviewFlags            = flag.NewFlagSet("", flag.ExitOnError)
	tasksFlags             = flag.NewFlagSet("", flag.ExitOnError)
	projectFlags           = flag.NewFlagSet("", flag.ExitOnError)
	listProjectFlags       = flag.NewFlagSet("", flag.ExitOnError)
	updateProjectFlags     = flag.NewFlagSet("", flag.ExitOnError)
	createProjectFlags     = flag.NewFlagSet("", flag.ExitOnError)
	createManyProjectFlags = flag.NewFlagSet("", flag.ExitOnError)
)

const (
	flagUsage = `
Usage of 'gtd':

  Available commands:
    - gtd project
    - gtd tasks
    - gtd review
`

	projectFlagUsage = `
Usage of command 'project':

  Available subcommands:
    - gtd project list
    - gtd project create
    - gtd project create-many
    - gtd project update
`

	listProjectFlagUsage = `
Usage of subcommand 'list':

  When invoked with no other flags, print a listing of all projects and exit.

`

	// TODO: more flag usage messages
)

func init() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, strings.TrimSpace(flagUsage))
		flag.PrintDefaults()
	}
	projectFlags.Usage = func() {
		fmt.Fprintln(os.Stderr, strings.TrimSpace(projectFlagUsage))
		projectFlags.PrintDefaults()
	}
	listProjectFlags.Usage = func() {
		fmt.Fprintln(os.Stderr, strings.TrimSpace(listProjectFlagUsage)+"\n")
		listProjectFlags.PrintDefaults()
	}

	// TODO: more flag usage registration
}
