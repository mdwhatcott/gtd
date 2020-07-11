package ux

import (
	"log"

	"github.com/mdwhatcott/gtd/core"
	"github.com/mdwhatcott/gtd/core/commands"
)

func handle(handler core.Handler, instructions ...interface{}) {
	handler.Handle(instructions...)

	for _, instruction := range instructions {
		fail, ok := instruction.(commands.Fallible)
		if !ok {
			return
		}
		err := fail.Err()
		if err == nil {
			return
		}
		log.Printf("Instruction failed. Error: [%v] Instruction: %#v", err, instruction)
	}
}
