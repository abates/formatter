package formatter

import (
	"log"
	"strings"
	"testing"

	"github.com/abates/formatter/colors"
)

type printer struct {
	*log.Logger
}

func (p printer) Log(v ...interface{})                 { p.Print(v...) }
func (p printer) Logf(format string, v ...interface{}) { p.Printf(format, v...) }

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
		{"no tags", "some input", "some input"},
		{"one tag", "some <em>input</em>", "some EMinputRESET"},
		{"nested tag", "some <em>input<success>with green</success></em>", "some EMinputSUCCESSwith greenRESETEMRESET"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			builder := &strings.Builder{}
			logger := log.New(builder, "", 0)
			f := ContextFormatter()
			l := ColorLogger(printer{logger}, LogFormatter(f))

			l.Logf("%s", test.input)
			got := strings.TrimSpace(builder.String())
			if test.want != got {
				t.Errorf("Wanted %q got %q", test.want, got)
			}
		})
	}
	resetColor = oldResetColor
	ContextColors = oldContextColors
}
