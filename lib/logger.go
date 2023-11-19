package lib

import (
	"log"
	"net/http"
	"os"
	"time"
)

type Logger struct {
	logger  *log.Logger
	handler http.Handler
}

func NewLogger(handler http.Handler) *Logger {
	logger := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	logger.SetPrefix("INFO: ")
	logger.SetFlags(log.Ldate | log.Lmicroseconds | log.Llongfile)
	return &Logger{logger, handler}
}

func (l *Logger) Info(message string) {
	l.logger.Println(message)
}

func (l *Logger) Error(message string) {
	l.logger.SetPrefix(" ERROR: ")
	l.logger.Println(message)
}

func (l *Logger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	l.handler.ServeHTTP(w, r)
	log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
}
