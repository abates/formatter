package formatter_test

import (
	"os"

	"github.com/abates/formatter"
)

func ExampleColorLogger() {
	logger := formatter.ColorLogger(formatter.LogWriter(os.Stdout))
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
