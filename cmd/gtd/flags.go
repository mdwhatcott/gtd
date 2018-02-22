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
- gtd tickler

Invoke any of the above commands with '-help' for additional information. 
What follows are several use cases and the corresponding commands.

### Weekly Review Procedures:

    gtd review

### Synchronize and Generate Tasks organized by context.

    gtd tasks

### Create many projects in a REPL session:

    gtd projects

### List all projects (display <h1> header, or filename if not present):

    gtd project list

### Review each project in turn in a REPL session combined w/ editor sessions:

    gtd project review

### Renegotiate project status:

    gtd project update -id 42 -status complete
    gtd project update -id 42 -status someday
    gtd project update -id 42 -status maybe
    gtd project update -id 42 -status reject

### Scan Tickler folder for projects that are now due:

    gtd tickler
`

	usageFlagsReview = `
Usage of command 'gtd review':

Procedures:

1. Synchronize Task Completion Status across context lists and projects
2. Process Inbox - Create many new projects as necessary
3. Reject/Defer/Update/Complete Projects and Tasks
4. Generate fresh task listings sorted by context

This action is:

- interactive
- destructive (so much so that the tool encourage committing changes to source control along the way)
`

	usageFlagsTasks = `
Usage of command 'gtd tasks':

This subcommand scans all tasks found in contextual lists and,
if completed, ensures they are marked as such in the corresponding
project document. Then it scans all active project files for unfinished 
tasks (lines beginning with '- [ ]') and sorts them by context
into separate files.

This action is:

- non-interactive
- destructive (consider committing all changes to source control before execution)

Flags:

`

	usageFlagsProject = `
Usage of command 'gtd project':

  You will be prompted to enter a project name.
  Upon having done so, an editor will open, allowing you to clarify the project.
  The above procedures will repeat until the user submits a blank line.

- interactive unless invoked with -static
- destructive (consider committing all changes to source control before execution)

  Available subcommands:
    - gtd project list
    - gtd project review
    - gtd project update

Invoke any of the above commands with '-help' for additional information. 

`

	usageFlagsUpdateProjects = `
Usage of subcommand 'gtd project update':

  Updates the status of a project by changing its location.

This action is:

- non-interactive
- destructive (consider committing all changes to source control before execution)

Flags:
`
)
