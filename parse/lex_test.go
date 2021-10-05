package parse

import (
	"reflect"
	"testing"
)

func TestLexing(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantTyps []itemType
		wantVals []string
		wantErr  bool
	}{
		{
			name:     "single tag",
			input:    "<tag></tag>",
			wantTyps: []itemType{itemStartTag, itemEndTag},
			wantVals: []string{"tag", "tag"},
			wantErr:  false,
		},
		{
			name:     "eof error",
			input:    "<tag></tag",
			wantTyps: []itemType{itemStartTag},
			wantVals: []string{"tag"},
			wantErr:  true,
		},
		{
			name:     "text wrapped tag",
			input:    "<tag>text text text</tag>",
			wantTyps: []itemType{itemStartTag, itemText, itemEndTag},
			wantVals: []string{"tag", "text text text", "tag"},
			wantErr:  false,
		},
		{
			name:     "before and after",
			input:    "before<tag>during</tag>after",
			wantTyps: []itemType{itemText, itemStartTag, itemText, itemEndTag, itemText},
			wantVals: []string{"before", "tag", "during", "tag", "after"},
			wantErr:  false,
		},
		{
			name:     "dangling start",
			input:    "<tag during",
			wantTyps: []itemType{},
			wantVals: []string{},
			wantErr:  true,
		},
		{
			name:     "nested",
			input:    "<tag1><tag2></tag2></tag1>",
			wantTyps: []itemType{itemStartTag, itemStartTag, itemEndTag, itemEndTag},
			wantVals: []string{"tag1", "tag2", "tag2", "tag1"},
			wantErr:  false,
		},
		{
			name:     "gt/lt",
			input:    "<tag1>&gt;<tag2></tag2>&lt;</tag1>",
			wantTyps: []itemType{itemStartTag, itemLiteralRight, itemStartTag, itemEndTag, itemLiteralLeft, itemEndTag},
			wantVals: []string{"tag1", "&gt;", "tag2", "tag2", "&lt;", "tag1"},
			wantErr:  false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			l := lex(test.input)
			gotTyps := []itemType{}
			gotVals := []string{}
			var item item
			for item = l.nextItem(); item.typ != itemEOF && item.typ != itemError; item = l.nextItem() {
				gotTyps = append(gotTyps, item.typ)
				gotVals = append(gotVals, item.val)
			}

			if item.typ == itemError && !test.wantErr {
				t.Errorf("Unexpected error: %v", item.val)
			} else if item.typ != itemError && test.wantErr {
				t.Errorf("Expected error but got none")
			} else {
				if !reflect.DeepEqual(test.wantTyps, gotTyps) {
					t.Errorf("Wanted items %v got %v", test.wantTyps, gotTyps)
				}

				if !reflect.DeepEqual(test.wantVals, gotVals) {
					t.Errorf("Wanted values %v got %v", test.wantVals, gotVals)
				}
			}
		})
	}
}
