package main

import (
	"bufio"
	"fmt"
	"os"
)

func createManyProjects(inputs []string) {
	flags(usageFlagsCreateManyProjects).Parse(inputs)

	fmt.Println("(<CTRL>-C to exit)")
	for {
		var name string
		for name == "" {
			fmt.Print("Enter project name: ")
			name = readLine()
		}
		createProject([]string{"-name", name})
	}
}

func readLine() string {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}
