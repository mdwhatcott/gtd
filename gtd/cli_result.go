package gtd

import (
	"encoding/json"
	"fmt"
	"io"
)

type CLIResult struct {
	Stdout io.Writer
	Stderr io.Writer
	Exit   int
}

func NewCLIResult(out, err io.Writer) *CLIResult {
	return &CLIResult{Stdout: out, Stderr: err}
}

func (this *CLIResult) String(s fmt.Stringer) {
	_, _ = io.WriteString(this, s.String())
}

func (this *CLIResult) JSON(v interface{}) {
	encoder := json.NewEncoder(this)
	encoder.SetIndent("", "  ")
	_ = encoder.Encode(v)
}

func (this *CLIResult) Write(v []byte) (int, error) {
	n, err := this.Stdout.Write(v)
	if err != nil {
		this.Exit++
		_, _ = fmt.Fprintf(this.Stderr, "Write err: %v", err)
	}
	return n, err
}

