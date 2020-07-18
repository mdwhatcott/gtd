package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func GTDREPL(APP *Application, version string) {
	directive := Prompt(version)
	for GTDOnce(APP, directive) {
		directive = Prompt(version)
	}
}

func GTDOnce(APP *Application, directive []string) bool {
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
		log.Println("\nUnrecognized directive:", directive)
	}
	return true
}

func Prompt(version string) []string {
	fmt.Print("\n" + PromptPrefix(version) + " ")
	fields := strings.Fields(ScanLine())
	if len(fields) > 1 && fields[0] == "gtd" {
		fields = fields[1:]
	}
	return fields
}

func PromptPrefix(version string) string {
	return fmt.Sprintf("gtd@%s:", version)
}

func ScanLine() string {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}
