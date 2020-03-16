package gtd

import "bytes"

type RecordingCLIResult struct {
	*CLIResult
	Stdout *bytes.Buffer
	Stderr *bytes.Buffer
}

func NewRecordingCLIResult() *RecordingCLIResult {
	out := new(bytes.Buffer)
	err := new(bytes.Buffer)
	return &RecordingCLIResult{
		CLIResult: NewCLIResult(out, err),
		Stdout:    out,
		Stderr:    err,
	}
}
