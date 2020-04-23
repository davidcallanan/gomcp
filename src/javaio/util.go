package javaio

import "bufio"

///////////////////////////////////////
// Bufio reader slice
///////////////////////////////////////

type readerSlice struct {
	underlyingReader *bufio.Reader
	readLimit int
	readCount int
}

func (r *readerSlice) Read(buf []byte) (n int, err error) {
	canRead := r.readLimit - r.readCount
	if canRead > len(buf) {
		canRead = len(buf)
	}

	n, err = r.underlyingReader.Read(buf[:canRead])
	r.readCount += n
	return
}

func newReaderSlice(underlyingReader *bufio.Reader, readLimit int) *bufio.Reader {
	return bufio.NewReader(&readerSlice {
		underlyingReader: underlyingReader,
		readLimit: readLimit,
	})
}
