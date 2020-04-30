package eventstore

import "strings"

type FakeReader struct {
	data     *strings.Reader
	readErr  error
	closeErr error
	closed   int
}

func NewFakeReader(data string) *FakeReader {
	return &FakeReader{data: strings.NewReader(data)}
}

func (this *FakeReader) Read(p []byte) (n int, err error) {
	n, err = this.data.Read(p)
	if err != nil {
		return n, err
	}
	return n, this.readErr
}

func (this *FakeReader) Close() error {
	this.closed++
	return this.closeErr
}
