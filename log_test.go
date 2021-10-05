package formatter

import (
	"strings"
	"testing"

	"github.com/abates/formatter/colors"
)

func TestLog(t *testing.T) {
	oldResetColor := resetColor
	resetColor = colors.Color("RESET")

	oldContextColors := ContextColors
	ContextColors = map[string]colors.Color{
		"em":      colors.Color("EM"),
		"success": colors.Color("SUCCESS"),
	}

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"no tags", "some input", "some input\n"},
		{"one tag", "some <em>input</em>", "some EMinputRESET\n"},
		{"nested tag", "some <em>input<success>with green</success></em>", "some EMinputSUCCESSwith greenRESETEMRESET\n"},
		{"trailing newline", "some text\n", "some text\n"},
		{"error", "some text<foo>\n", "some text<foo> !%ERR expected </foo>\n"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			builder := &strings.Builder{}
			f := ContextFormatter()
			l := ColorLogger(LogWriter(builder), LogFormatter(f))

			l.Logf("%s", test.input)
			got := builder.String()
			if test.want != got {
				t.Errorf("Wanted %q got %q", test.want, got)
			}
		})
	}
	resetColor = oldResetColor
	ContextColors = oldContextColors
}
