package main

import (
	"log"

	"github.com/mdwhatcott/gtd/v3/util/date"
	"github.com/mdwhatcott/gtd/v3/util/template"
)

func PrintUsage() {
	FLAGS := log.Flags()
	log.SetFlags(0)
	defer log.SetFlags(FLAGS)

	log.Println(template.MustExecute(usageTemplate, UsageData{
		Prompt: PromptPrefix(Version),
		Today:  date.Today(),
	}))
}

type UsageData struct {
	Prompt string
	Today  string
}

var usageTemplate = `Usage of gtd:

This application provides both a command-line and a REPL-style
interface for managing projects and tasks.

## Required Environment Variables

- GTDPATH (the folder containing the "Event Database")


## Event Database

Because this application is purely event-sourced, the only
state required is a database of serialized events. If missing
at startup a new database will be created. It is assumed that
the parent folder of the database is under revision control 
in a 'git' repository with a valid remote that has already
been configured. As a final step at shutdown any new events
added to the database are committed and pushed using the
following automated commands:

    $ git add <root>/events.csv
    $ git commit -m "{{.Today}}"
    $ git push


## Directives

Here are several examples of directives that can be entered
at the REPL, or as non-flag command-line arguments (which,
if included, will bypass the REPL):


To show this usage documentation:

    {{.Prompt}} help


To exit the program:

    {{.Prompt}} exit


To present a listing of all active contexts:


    {{.Prompt}} contexts


To present a single, random, pending task from an active 
project:

    {{.Prompt}} task


// Optionally, draw from tasks matching the provided contexts,
// home and work:

    {{.Prompt}} task home work


To present all pending tasks from active projects:

    {{.Prompt}} tasks


Optionally, show only those tasks that match the provided
contexts, home and work:

    {{.Prompt}} tasks home work


To present all projects (separated by status):

    {{.Prompt}} projects


Optionally, only present projects matching the provided
statuses, fixed and deferred:

    {{.Prompt}} projects fixed deferred


[UNDER CONSTRUCTION] Optionally, only present projects 
that have no pending tasks:

    {{.Prompt}} projects fixed deferred ?


The End`
