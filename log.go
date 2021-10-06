package formatter

import (
	"fmt"
	"io"
	"os"
	"strings"
)

type Logger interface {
	Log(v ...interface{})
	Logf(format string, v ...interface{})
}

type logger struct {
	writer io.Writer
	format Formatter
}

type LogOption func(*logger)

func LogWriter(writer io.Writer) LogOption {
	return func(cl *logger) {
		cl.writer = writer
	}
}

func LogFormatter(format Formatter) LogOption {
	return func(cl *logger) {
		cl.format = format
	}
}

func ColorLogger(options ...LogOption) Logger {
	f := &logger{
		writer: os.Stdout,
		format: ContextFormatter(),
	}

	for _, option := range options {
		option(f)
	}
	return f
}

func (l *logger) Logf(format string, v ...interface{}) {
	f, err := l.format(format)
	if err == nil {
		s := fmt.Sprintf(f, v...)
		if strings.HasSuffix(s, "\n") {
			fmt.Fprint(l.writer, s)
		} else {
			fmt.Fprintln(l.writer, s)
		}
	} else {
		// trim space to prevent a trailing newline from
		// moving the error message to the next line
		fmt.Fprintf(l.writer, strings.TrimSpace(format), v...)
		fmt.Fprintf(l.writer, " !%%ERR %v\n", err)
	}
}

func (l *logger) Log(v ...interface{}) {
	fmt.Fprint(l.writer, v...)
}
