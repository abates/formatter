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
	l.Log(fmt.Sprintf(format, v...))
}

func (l *logger) Log(v ...interface{}) {
	input := fmt.Sprint(v...)
	if s, err := l.format(input); err == nil {
		if strings.HasSuffix(s, "\n") {
			fmt.Fprint(l.writer, s)
		} else {
			fmt.Fprintln(l.writer, s)
		}
	} else {
		fmt.Fprintf(l.writer, "%s !%%ERR %v\n", strings.TrimSpace(input), err)
	}
}
