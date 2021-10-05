package parse

import "testing"

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr string
	}{
		{"basic", "one two three", "one two three", ""},
		{"tagged basic", "one <tag>two</tag> three", "one <tag>two</tag> three", ""},
		{"literals", "one &lt;<tag>two</tag>&gt; three", "one &lt;<tag>two</tag>&gt; three", ""},
		{"missing closing tag", "<tag>two three", "", "expected </tag>"},
		{"end tag without start tag", "two three</tag>", "", "end tag </tag> without start tag"},
		{"mismatched tags", "<tag1></tag2>", "", "expecting end tag </tag1> got </tag2>"},
		{"invalid input", "<tag", "", "expecting >"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tree, err := parse(test.input)
			if err == nil {
				if test.wantErr == "" {
					got := tree.String()
					if test.want != got {
						t.Errorf("Wanted %q got %q", test.want, got)
					}
				} else {
					t.Errorf("Wanted error, got none")
				}
			} else if test.wantErr != err.Error() {
				t.Errorf("Wanted err %q got %q", test.wantErr, err)
			}
		})
	}
}
