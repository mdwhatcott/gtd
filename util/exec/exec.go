package exec

import (
	"log"
	"os/exec"
	"strings"
)

func MustDo(where, message string, args ...string) string {
	OUT, ERR := Do(where, message, args...)
	if ERR != nil {
		log.Fatalf(""+
			"# Command: %v\n\n"+
			"# Error:   %v\n\n"+
			"# Output:\n%s",
			args, ERR, OUT,
		)
	}
	return OUT
}

func Do(where, message string, args ...string) (string, error) {
	if len(message) > 0 {
		log.Println(message)
	}
	COMMAND := exec.Command(args[0], args[1:]...)
	COMMAND.Dir = where
	OUT, ERR := COMMAND.CombinedOutput()
	return strings.TrimSpace(string(OUT)), ERR
}
