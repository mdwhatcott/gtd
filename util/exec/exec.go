package exec

import (
	"log"
	"os/exec"
	"strings"
)

func MustDo(_where, _message string, _args ...string) string {
	OUT, ERR := Do(_where, _message, _args...)
	if ERR != nil {
		log.Fatalf(""+
			"# Command: %v\n\n"+
			"# Error:   %v\n\n"+
			"# Output:\n%s",
			_args, ERR, OUT,
		)
	}
	return OUT
}

func Do(_where, _message string, _args ...string) (string, error) {
	if len(_message) > 0 {
		log.Println(_message)
	}
	COMMAND := exec.Command(_args[0], _args[1:]...)
	COMMAND.Dir = _where
	OUT, ERR := COMMAND.CombinedOutput()
	return strings.TrimSpace(string(OUT)), ERR
}
