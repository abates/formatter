package formatter

import (
	"testing"

	"github.com/abates/formatter/colors"
)

func TestColorFormatterPeekPop(t *testing.T) {
	tests := []struct {
		name  string
		input []colors.Color
	}{
		{"empty stack", nil},
		{"one item in stack", []colors.Color{colors.Color("EMRED")}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			f := &colorFormatter{colorStack: test.input}
			want := colors.Color("")
			if len(test.input) > 0 {
				want = test.input[len(test.input)-1]
			}
			got := f.peek()
			if want != got {
				t.Errorf("Wanted %q got %q", want, got)
			} else {
				wantLen := len(test.input)
				want = got
				got = f.pop()
				gotLen := len(test.input)
				if want != got {
					t.Errorf("Wanted pop to return last element %q got %q", want, got)
				}

				if wantLen != gotLen {
					t.Errorf("Wanted color stack to be reduced by one")
				}
			}
		})
	}
}

func TestContextFormatter(t *testing.T) {
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
			f := ContextFormatter()
			got, err := f(test.input)
			if err == nil {
				if test.want != got {
					t.Errorf("Wanted %q got %q", test.want, got)
				}
			} else {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
	resetColor = oldResetColor
	ContextColors = oldContextColors
}
