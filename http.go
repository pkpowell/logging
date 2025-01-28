package logging

import (
	"net/http"
	"strconv"
	"time"

	"log"

	"github.com/pkpowell/humanize/units"
)

type wrappedWriter struct {
	http.ResponseWriter
	statusCode int
}

var (
	start    time.Time
	duration time.Duration
	length   int
	ww       *wrappedWriter
	cl       string
	d        string
)

func HTTPHandler(h http.Handler, verbose *bool, webLogs *bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		start = time.Now()
		length = 0

		ww = &wrappedWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		h.ServeHTTP(ww, r)

		cl = w.Header().Get("Content-Length")
		if cl != "" {
			length, err = strconv.Atoi(cl)
			if err != nil {
				Error("strconv.Atoi error", err.Error())
				// length = -1
			}
		}

		duration = time.Since(start)

		d = r.Proto + " from " + r.RemoteAddr + " " + strconv.Itoa(ww.statusCode) + " " + r.Method + " " + r.RequestURI + " " + duration.String() + " " + units.Int(length).String()

		if *webLogs && !*verbose {
			Info(d)
		} else {
			Debug(d)
		}
	})
}

func HTTPHandlerFunc(h http.HandlerFunc, verbose *bool, webLogs *bool) http.HandlerFunc {
	log.Printf("HTTPHandlerFunc %v\n", h)
	hf := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("HTTPHandlerFunc %v\n", h)
		var err error
		start = time.Now()
		length = 0

		ww = &wrappedWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}
		h.ServeHTTP(ww, r) // serve the original request
		duration = time.Since(start)

		if r.Method == http.MethodPost {
			cl = w.Header().Get("Content-Length")
			if cl != "" {
				length, err = strconv.Atoi(cl)
				if err != nil {
					Error("strconv.Atoi error", err.Error())
				}
			}
		}

		// avoid fmt.Sprintf for performance
		d = r.Proto + " from " + r.RemoteAddr + " " + r.Method + " " + r.RequestURI + " " + duration.String() + " " + units.Int(length).String()

		if *webLogs && !*verbose {
			Info(d)
		} else {
			Debug(d)
		}
	})

	return http.HandlerFunc(hf)
}
