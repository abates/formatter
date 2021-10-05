package formatter_test

import (
	"log"
	"os"

	"github.com/abates/formatter"
)

type SystemLogger struct {
	*log.Logger
}

func (sl SystemLogger) Log(v ...interface{})                 { sl.Print(v...) }
func (sl SystemLogger) Logf(format string, v ...interface{}) { sl.Printf(format, v...) }

func ExampleColorLogger() {
	sysLogger := log.New(os.Stdout, "", 0)

	logger := formatter.ColorLogger(SystemLogger{sysLogger})
	logger.Log("This is some text with <em>emphasis</em> in the middle")
	logger.Log("This is some text with a <warn>warning</warn> in the middle")
	logger.Log("This is some text with a <fail>failure</fail> in the middle")
	logger.Log("This is something <success>good</success> that happened")
	// Output:
	// This is some text with [36memphasis[0m in the middle
	// This is some text with a [33mwarning[0m in the middle
	// This is some text with a [31mfailure[0m in the middle
	// This is something [32mgood[0m that happened
}
