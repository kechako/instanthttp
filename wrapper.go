package main

import "net/http"

type ResponseWriterWrapper struct {
	http.ResponseWriter
	Code int
	Size int64
}

func (w *ResponseWriterWrapper) Header() http.Header {
	return w.ResponseWriter.Header()
}

func (w *ResponseWriterWrapper) Write(b []byte) (n int, err error) {
	n, err = w.ResponseWriter.Write(b)
	if err == nil {
		w.Size += int64(n)
	}
	return
}

func (w *ResponseWriterWrapper) WriteHeader(statusCode int) {
	w.Code = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}
