package betterlog

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

type timeWriter struct {
	writer io.Writer
}

func (w timeWriter) Write(p []byte) (int, error) {
	date := time.Now().Format("[2006-01-02 15:04:05] ")
	p = append([]byte(date), p...)
	return w.writer.Write(p)
}

func makeDateWriter(w io.Writer) io.Writer {
	return &timeWriter{w}
}

func makeLogger(path string) (string, *os.File, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", nil, err
	}
	err = os.MkdirAll(filepath.Dir(abs), os.ModePerm)
	if err != nil {
		return "", nil, err
	}
	f, err := os.OpenFile(abs, os.O_WRONLY+os.O_APPEND+os.O_CREATE, os.ModePerm)
	return abs, f, err
}

// MakeDateLogger configure the default logger to write both
// to stderr and a log file.
// Note: don't forget to call 'defer f.Close()' if the returned error is nil
func MakeDateLogger(logFile string) (*os.File, error) {
	log.SetFlags(0)
	abs, f, err := makeLogger(filepath.Join(filepath.Dir(os.Args[0]), logFile))
	if err == nil {
		log.SetOutput(makeDateWriter(io.MultiWriter(f, os.Stderr)))
	}
	log.Println("[info] logging to:", abs)
	return f, err
}
