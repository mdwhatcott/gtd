package main

import (
	"flag"
	"log"
	"os"
	"strings"
)

func flags(usage string) *flag.FlagSet {
	var flagLog = log.New(os.Stderr, "", 0)
	set := flag.NewFlagSet("", flag.ExitOnError)
	set.Usage = func() {
		flagLog.Println(strings.TrimSpace(usage))
		set.PrintDefaults()
	}
	return set
}

func firstAndRemaining(args []string) (string, []string) {
	return firstOrDefault(args), skipFirst(args)
}
func firstOrDefault(args []string) string {
	if len(args) == 0 {
		return ""
	}
	return args[0]
}
func skipFirst(args []string) []string {
	if len(args) < 1 {
		return nil
	}
	return args[1:]
}

func exit(flags *flag.FlagSet) {
	flags.Usage()
	os.Exit(1)
}

const (
	usageFlag = `
# Usage of 'gtd':

Commands:

- project
- tasks
- review

## Usage of 'gtd review':

> Weekly Review Procedures:

> 1. Process Inbox - Create many new projects as necessary
> 2. Synchronize Task Completion Status across context lists and projects
> 3. Reject/Defer/Update/Complete Projects and Tasks
> 4. Generate fresh task listings sorted by context

## Synchronize Task Completion Status across context lists and projects

    gtd tasks sync

## Generate fresh task listings sorted by context (log projects that don't have any next action):

    gtd tasks sweep

## List all projects (display <h1> header, or filename if not present):

    gtd project list

## Review each project in turn in a REPL session combined w/ editor sessions:

    gtd project list -review

## Create many projects in a REPL session:

    gtd project create-many

## Create Project and specify each relevant section from CLI (only name is required):

    gtd project create -name "Hi" -outcome "Something" -next-action "Something simple" -info "Even more stuff"

## Create blank project:

    gtd project create -name "Hi" -blank

## Create project from CLI and skip editor session:

    gtd project create -name "Hi" -static

## Renegotiate project status:

    gtd project update -id 42 -someday
    gtd project update -id 42 -maybe
    gtd project update -id 42 -reject
    gtd project update -id 42 -complete
`

	usageFlagsReview = `
Usage of command 'review':

Coming soon...
`

	usageFlagsTasks = `
Usage of command 'tasks':

  Available subcommands:
    - gtd tasks sync
    - gtd tasks sweep
`

	usageFlagsSyncTasks = `
Usage of subcommand 'sync':

Coming soon...
`

	usageFlagsSweepTasks = `
Usage of subcommand 'sweep':

This subcommand scans all active project files for unfinished 
tasks (lines beginning with '- [ ]') and emits them to stdout.
Do with said output what you will...
`

	usageFlagsProject = `
Usage of command 'project':

  Available subcommands:
    - gtd project list
    - gtd project create
    - gtd project create-many
    - gtd project update
`

	usageFlagsListProjects = `
Usage of subcommand 'list':

  When invoked with no other flags, print a listing of all projects and exit.

`

	usageFlagsUpdateProjects = `
Usage of subcommand 'update':

Coming soon...
`

	usageFlagsCreateProject = `
Usage of subcommand 'create':
`

	usageFlagsCreateManyProjects = `
Usage of subcommand 'create-many':

  You will be prompted to enter a project name.
  Upon having done so, an editor will open, allowing you to clarify the project.
  The above procedures will repeat until the user submits a keyboard interrupt (<CTRL>-c).
`
)
