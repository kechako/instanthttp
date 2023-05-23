package main

import (
	"net/http"
)

type ResponseWriterWrapper struct {
	http.ResponseWriter
	Code        int
	Size        int64
	wroteHeader bool
}

var _ http.ResponseWriter = (*ResponseWriterWrapper)(nil)

func (w *ResponseWriterWrapper) Header() http.Header {
	return w.ResponseWriter.Header()
}

func (w *ResponseWriterWrapper) Write(b []byte) (n int, err error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	n, err = w.ResponseWriter.Write(b)
	if err == nil {
		w.Size += int64(n)
	}
	return
}

func (w *ResponseWriterWrapper) WriteHeader(statusCode int) {
	if !w.wroteHeader {
		w.Code = statusCode
		w.wroteHeader = true
	}
	w.ResponseWriter.WriteHeader(statusCode)
}
