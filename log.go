package formatter

import (
	"fmt"
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

type LogOption func(*logger)

func LogFormatter(format Formatter) LogOption {
	return func(cl *logger) {
		cl.format = format
	}
}

func ColorLogger(l Logger, options ...LogOption) Logger {
	f := &logger{
		Logger: l,
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
