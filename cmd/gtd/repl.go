package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func REPL(APP *Application, version string, directive []string) {
	defer APP.CommitChanges()

	if len(directive) == 0 {
		directive = InitialPrompt(version)
	}

	for {
		switch directive[0] {

		case "projects":
			APP.PresentOutcomesListing(directive[1:])

		case "tasks":
			APP.PresentIncompleteActions(directive[1:])

		case "task":
			APP.PresentIncompleteAction(directive[1:])

		case "contexts":
			APP.PresentContexts()

		case "help":
			Usage()

		case "exit":
			return

		default:
			log.Println("Unrecognized directive:", directive)
		}

		// TODO: Display stalled projects!
		directive = Prompt()
	}
}

func InitialPrompt(version string) []string {
	fmt.Println("Welcome to Getting Things Done!")
	fmt.Println()
	fmt.Println("Version:", version)
	fmt.Println()
	fmt.Println("Enter 'help' for instructions.")

	return Prompt()
}

func Prompt() []string {
	fmt.Print("\n-> ")
	return ScanLine()
}

func ScanLine() []string {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	line := scanner.Text()
	return strings.Fields(line)
}

func Usage() {
	flags := log.Flags()
	log.SetFlags(0)
	defer log.SetFlags(flags)

	log.Println(`Usage of gtd:

This application provides a REPL-style interface for selecting directives.
Here are several examples of commands that can be entered:


# To show this usage documentation:
-> help


# To exit the program:
-> exit


# To present a listing of all active contexts:
-> contexts


# To present a single, random, pending task from an active project:
-> task

## Optionally, draw from tasks matching the provided contexts, home and work:
-> task home work


# To present all pending tasks from active projects:
-> tasks

## Optionally, show only those tasks that match the provided contexts, home and work:
-> tasks home work


# To present all projects (separated by status):
-> projects

## Optionally, only present projects matching the provided statuses, fixed and deferred:
-> projects fixed deferred

## [UNDER CONSTRUCTION] Optionally, only present projects that have no pending tasks:
-> projects fixed deferred ?


The End`)
}
