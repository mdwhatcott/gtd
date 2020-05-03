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
	if this.readErr != nil {
		return 0, this.readErr
	}
	n, err = this.data.Read(p)
	if err != nil {
		return n, err
	}
	return n, nil
}

func (this *FakeReader) Close() error {
	this.closed++
	return this.closeErr
}
