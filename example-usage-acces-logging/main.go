package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/rogierlommers/logrus-redis-hook"
	"github.com/sirupsen/logrus"
)

const debug = true

func exampleHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "A message was received...")
	})
}

func main() {
	accessLog := logrus.New()

	if !debug {
		accessLog.Out = ioutil.Discard
	}

	hookConfig := logredis.HookConfig{
		Host:     "localhost",
		Key:      "my_redis_key",
		Format:   "access",
		App:      "escher-serve",
		Hostname: "my_app_hostname",
		Port:     6379,
		DB:       0,
		TTL:      3600,
	}

	hook, err := logredis.NewHook(hookConfig)
	if err == nil {
		accessLog.AddHook(hook)
	} else {
		logrus.Errorf("logredis error: %q", err)
	}

	http.Handle("/", LoggedHandler(accessLog, exampleHandler()))
	logrus.Info("listening on :8080")
	logrus.Fatal(http.ListenAndServe(":8080", nil))
}

func LoggedHandler(logger *logrus.Logger, wrappedHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		startTime := time.Now()

		// do request
		lrw := NewLoggingResponseWriter(w)
		wrappedHandler.ServeHTTP(lrw, req)

		// send metrics
		logger.WithFields(logrus.Fields{

			// dynamic stuff
			"agent":       req.UserAgent(),
			"client":      strings.Split(req.RemoteAddr, ":")[0],
			"httpversion": req.Proto,
			"response":    lrw.statusCode,
			"bytes":       lrw.contentLength,
			"time_in_sec": time.Since(startTime).Nanoseconds() / 1e6,
			"referrer":    req.Referer(),
			"verb":        req.Method,
			"site":        req.Host,
			"message":     req.RequestURI,
		}).Info()

	})
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode    int
	contentLength int
}

func NewLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK, 0}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *loggingResponseWriter) Write(b []byte) (int, error) {
	n, err := lrw.ResponseWriter.Write(b)
	lrw.contentLength = lrw.contentLength + n
	return n, err
}
