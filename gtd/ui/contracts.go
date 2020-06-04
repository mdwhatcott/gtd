package ui

type Editor interface {
	EditTempFile(initialContent string) (resultContent string)
}
