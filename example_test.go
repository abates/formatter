package formatter_test

import (
	"os"

	"github.com/abates/formatter"
)

func ExampleColorLogger() {
	logger := formatter.ColorLogger(formatter.LogWriter(os.Stdout))
	// Only logf parses and evaluates the tagged string and then only
	// on the format string.  This allows angle brackets to appear in
	// the variadic arguments
	logger.Logf("This is some text with <em>emphasis</em> in the middle")
	logger.Logf("This is some text with a <warn>warning</warn> in the middle")
	logger.Logf("This is some text with a <fail>failure</fail> in the middle")
	logger.Logf("This is something <success>good</success> that happened")
	// Output:
	// This is some text with [36memphasis[0m in the middle
	// This is some text with a [33mwarning[0m in the middle
	// This is some text with a [31mfailure[0m in the middle
	// This is something [32mgood[0m that happened
}
