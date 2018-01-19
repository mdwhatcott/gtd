package main

import (
	"flag"
	"os"
)

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

## Prerequisites:

- [Git](https://git-scm.com/) revision control installed.
- [Atlassian Source Tree](https://www.sourcetreeapp.com/) w/ command line tools installed ('stree').
- [Sublime Text](https://www.sublimetext.com/) w/ command line tools installed ('subl').
- Understanding of [GTD®](https://gettingthingsdone.com/), the GTD® Weekly Review process, 
and a willingness to accept the variation on that process adopted herein.

## Basic Commands (some of which have their own subcommands):

- gtd project
- gtd tasks
- gtd review

Invoke any of the above commands with '-help' for additional information. 
What follows are several use cases and the corresponding commands.

### Weekly Review Procedures:

    gtd review

### Synchronize Task Completion Status across context lists and projects

    gtd tasks sync

### Generate fresh task listings sorted by context (log projects that don't have any next action):

    gtd tasks sweep

### List all projects (display <h1> header, or filename if not present):

    gtd project list

### Review each project in turn in a REPL session combined w/ editor sessions:

    gtd project list -review

### Create many projects in a REPL session:

    gtd project create-many

### Create Project and specify each relevant section from CLI (only name is required):

    gtd project create -name "Hi" -outcome "Something" -next-action "Something simple" -info "Even more stuff"

### Create blank project:

    gtd project create -name "Hi" -blank

### Create project from CLI and skip editor session:

    gtd project create -name "Hi" -static

### Renegotiate project status:

    gtd project update -id 42 -status complete
    gtd project update -id 42 -status someday
    gtd project update -id 42 -status maybe
    gtd project update -id 42 -status reject
`

	usageFlagsReview = `
Usage of command 'gtd review':

Procedures:

1. Process Inbox - Create many new projects as necessary
2. Synchronize Task Completion Status across context lists and projects
3. Reject/Defer/Update/Complete Projects and Tasks
4. Generate fresh task listings sorted by context

This action is:

- interactive
- destructive (so much so that the tool encourage committing changes to source control along the way)
`

	usageFlagsTasks = `
Usage of command 'gtd tasks':

  Available subcommands:
    - gtd tasks sync
    - gtd tasks sweep

Invoke any of the above commands with '-help' for additional information.
`

	usageFlagsSyncTasks = `
Usage of subcommand 'gtd tasks sync':

This subcommand scans all tasks found in contextual lists and,
if completed, ensures they are marked as such in the corresponding
project document.

This action is:

- non-interactive
- destructive (consider committing all changes to source control before execution)

Flags:
`

	usageFlagsSweepTasks = `
Usage of subcommand 'gtd tasks sweep':

This subcommand scans all active project files for unfinished 
tasks (lines beginning with '- [ ]') and sorts them by context
into separate files.

This action is:

- non-interactive
- destructive (consider committing all changes to source control before execution)

Flags:

`

	usageFlagsProject = `
Usage of command 'gtd project':

  Available subcommands:
    - gtd project list
    - gtd project create
    - gtd project create-many
    - gtd project update

Invoke any of the above commands with '-help' for additional information. 

`

	usageFlagsListProjects = `
Usage of subcommand 'gtd project list':

  When invoked with no other flags, print a listing of all projects and exit.

This action is:

- interactive when invoked with -review
- non-destructive

Flags:
`

	usageFlagsUpdateProjects = `
Usage of subcommand 'gtd project update':

  Updates the status of a project by changing its location.

This action is:

- non-interactive
- destructive (consider committing all changes to source control before execution)

Flags:
`

	usageFlagsCreateProject = `
Usage of subcommand 'gtd project create':

  Create a new file for a project and (by default) open a text editor to allow deeper clarification.

This action is:

- interactive unless invoked with -static
- destructive (consider committing all changes to source control before execution)

Flags:
`

	usageFlagsCreateManyProjects = `
Usage of subcommand 'gtd project create-many':

  You will be prompted to enter a project name.
  Upon having done so, an editor will open, allowing you to clarify the project.
  The above procedures will repeat until the user submits a blank line.

- interactive unless invoked with -static
- destructive (consider committing all changes to source control before execution)

Flags:
`
)
