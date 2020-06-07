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
		if fail.Err() == nil {
			return
		}
		log.Println("Instruction failed with err:", fail.Err())
	}
}
