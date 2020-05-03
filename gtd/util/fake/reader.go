package fake

import "strings"

type Reader struct {
	data     *strings.Reader
	ReadErr  error
	CloseErr error
	Closed   int
}

func NewReader(_data string) *Reader {
	return &Reader{data: strings.NewReader(_data)}
}

func (this *Reader) Read(_p []byte) (n_ int, err_ error) {
	if this.ReadErr != nil {
		return 0, this.ReadErr
	}
	n_, err_ = this.data.Read(_p)
	if err_ != nil {
		return n_, err_
	}
	return n_, nil
}

func (this *Reader) Close() error {
	this.Closed++
	return this.CloseErr
}
