package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func DoREPL(APP *Application, version string) {
	directive := Prompt(version)
	for DoOnce(APP, directive) {
		directive = Prompt(version)
	}
}

func DoOnce(APP *Application, directive []string) bool {
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
		PrintUsage()

	case "exit":
		return false

	default:
		log.Println("Unrecognized directive:", directive)
	}
	return true
}

func PrintBanner(version string) {
	fmt.Println("Welcome to Getting Things Done!")
	fmt.Println()
	fmt.Println("Version:", version)
	fmt.Println()
	fmt.Println("Enter 'help' for instructions.")
}

func Prompt(version string) []string {
	fmt.Printf("\ngtd-%s --> ", version)
	fields := ScanLineFields()
	if len(fields) > 1 && fields[0] == "gtd" {
		fields = fields[1:]
	}
	return fields
}

func ScanLineFields() []string {
	return strings.Fields(ScanLine())
}

func ScanLine() string {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}

func PrintUsage() {
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
