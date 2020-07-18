package wireupstorage

import (
	"io"
	"os"

	"github.com/google/uuid"
)

func reading(path string) func() io.ReadCloser  { return func() io.ReadCloser { return open(path) } }
func writing(path string) func() io.WriteCloser { return func() io.WriteCloser { return open(path) } }

func open(path string) *os.File {
	FILE, ERR := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if ERR != nil {
		panic(ERR)
	}
	return FILE
}

func GenerateID() string {
	return uuid.New().String()
}
