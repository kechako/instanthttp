package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/kechako/sigctx"
	flag "github.com/spf13/pflag"
)

var (
	addr string
	cert string
	key  string
	dir  string
)

func init() {
	flag.StringVar(&addr, "http", ":8080", "IP address and port number to bind.")
	flag.StringVar(&cert, "cert", "", "TLS certificate file")
	flag.StringVar(&key, "key", "", "TLS private key file")
	flag.StringVar(&dir, "dir", "", "Directory to serve.")
}

func printError(err error, exit bool) {
	fmt.Fprintf(os.Stderr, "[ERROR] %+v\n", err)
	if exit {
		os.Exit(1)
	}
}

var accessLogger = log.New(os.Stdout, "", log.LstdFlags)

func accessLog(code int, size int64, path string) {
	accessLogger.Printf("[%d]: %s (%d bytes)", code, path, size)
}

func accessLogHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wrapper := ResponseWriterWrapper{ResponseWriter: w}

		next.ServeHTTP(&wrapper, r)

		accessLog(wrapper.Code, wrapper.Size, r.URL.Path)
	})
}

func main() {
	flag.Parse()

	if dir == "" {
		d, err := os.Getwd()
		if err != nil {
			printError(err, true)
		}
		dir = d
	} else {
		var err error
		dir, err = filepath.Abs(dir)
		if err != nil {
			printError(err, true)
		}
	}

	srv := &http.Server{
		Addr:    addr,
		Handler: accessLogHandler(http.FileServer(http.Dir(dir))),
	}

	ctx, cancel := sigctx.WithCancelBySignal(context.Background(), os.Interrupt)
	defer cancel()

	go func() {
		<-ctx.Done()

		sdCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := srv.Shutdown(sdCtx); err != nil {
			printError(err, false)
		}
	}()

	fmt.Printf("Start server [%s], and serve [%s]\n", addr, dir)

	var err error
	if cert != "" && key != "" {
		err = srv.ListenAndServeTLS(cert, key)
	} else {
		err = srv.ListenAndServe()
	}
	if err != nil {
		if err != http.ErrServerClosed {
			printError(err, true)
		}
	}
}
