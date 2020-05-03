package eventstore

import "strings"

type FakeReader struct {
	data     *strings.Reader
	readErr  error
	closeErr error
	closed   int
}

func NewFakeReader(_data string) *FakeReader {
	return &FakeReader{data: strings.NewReader(_data)}
}

func (this *FakeReader) Read(_p []byte) (n_ int, err_ error) {
	if this.readErr != nil {
		return 0, this.readErr
	}
	n_, err_ = this.data.Read(_p)
	if err_ != nil {
		return n_, err_
	}
	return n_, nil
}

func (this *FakeReader) Close() error {
	this.closed++
	return this.closeErr
}
