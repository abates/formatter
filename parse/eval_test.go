package parse

import (
	"strings"
	"testing"
)

func TestEval(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantText string
		wantTags []string
		wantErr  string
	}{
		{"basic", "before <tag>during</tag> after", "before during after", []string{"tag"}, ""},
		{"moderate", "before <tag1>during <tag2>emphasis</tag2></tag1> after", "before during emphasis after", []string{"tag1", "tag2"}, ""},
		{"extensive", "before <tag1><tag3>during</tag3> between <tag2>emphasis</tag2></tag1> after", "before during between emphasis after", []string{"tag1", "tag3", "tag2"}, ""},
		{"parse error", "before <tag1><tag3>during", "", []string{"tag1", "tag3"}, "expected </tag3>"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotTags := []string{}
			cb := func(tag string, builder *strings.Builder, next func()) {
				gotTags = append(gotTags, tag)
				next()
			}

			gotText, gotErr := Eval(test.input, cb)
			if gotErr == nil {
				if test.wantErr != "" {
					t.Errorf("Wanted error %q got nil", test.wantErr)
					return
				}

				if test.wantText != gotText {
					t.Errorf("Wanted text %q got %q", test.wantText, gotText)
				}

				correct := true
				if len(test.wantTags) == len(gotTags) {
					for i, wantTag := range test.wantTags {
						if wantTag != gotTags[i] {
							correct = false
							break
						}
					}
				} else {
					correct = false
				}

				if !correct {
					t.Errorf("Wanted tags %v got %v", test.wantTags, gotTags)
				}
			} else if gotErr.Error() != test.wantErr {
				t.Errorf("wanted error %q got %q", test.wantErr, gotErr.Error())
			}
		})
	}
}
