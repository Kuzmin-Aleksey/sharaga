package util

import "io"

type MultiReadCloser struct {
	readers []io.ReadCloser
	io.Reader
}

func (r *MultiReadCloser) Close() (err error) {
	for _, reader := range r.readers {
		if reader != nil {
			if e := reader.Close(); e != nil {
				err = e
			}
		}
	}
	return
}

func NewMultiReadCloser(readers ...io.ReadCloser) *MultiReadCloser {
	simpleReaders := make([]io.Reader, len(readers))
	for i, reader := range readers {
		simpleReaders[i] = reader
	}

	return &MultiReadCloser{
		readers: readers,
		Reader:  io.MultiReader(simpleReaders...),
	}
}
