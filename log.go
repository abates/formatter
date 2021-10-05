package formatter

import (
	"fmt"
	"log"
	"os"
)

type Printer interface {
	Print(v ...interface{})
}

type Logger interface {
	Log(v ...interface{})
	Logf(format string, v ...interface{})
}

type logger struct {
	Logger
	format Formatter
}

type printer struct {
	*log.Logger
}

func (p printer) Log(v ...interface{})                 { p.Print(v...) }
func (p printer) Logf(format string, v ...interface{}) { p.Printf(format, v...) }

type LogOption func(*logger)

func LoggerOption(l *log.Logger) LogOption {
	return func(cl *logger) {
		cl.Logger = printer{l}
	}
}

func LogFormatter(format Formatter) LogOption {
	return func(cl *logger) {
		cl.format = format
	}
}

func ColorLogger(options ...LogOption) Logger {
	f := &logger{
		Logger: printer{log.New(os.Stderr, "", log.LstdFlags)},
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
		l.Logger.Log(s)
	} else {
		l.Logger.Log(fmt.Sprintf("%s !%%ERR %v", input, err))
	}
}
