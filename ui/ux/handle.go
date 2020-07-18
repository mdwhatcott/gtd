package ux

import (
	"context"
	"log"

	"github.com/mdwhatcott/gtd/v3/core"
	"github.com/mdwhatcott/gtd/v3/core/commands"
)

func handle(handler core.Handler, instructions ...interface{}) {
	handler.Handle(context.Background(), instructions...)

	for _, INSTRUCTION := range instructions {
		FAIL, OK := INSTRUCTION.(commands.Fallible)
		if !OK {
			return
		}
		ERR := FAIL.Err()
		if ERR == nil {
			return
		}
		log.Printf("Instruction failed. Error: [%v] Instruction: %#v", ERR, INSTRUCTION)
	}
}
