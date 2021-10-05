package parse

import (
	"fmt"
	"strings"
)

const (
	leftLiteral  = "<"
	rightLiteral = ">"
)

type nodeType int

const (
	nodeRoot nodeType = iota
	nodeText
	nodeTag
	nodeLiteral
)

func (t nodeType) typ() nodeType {
	return t
}

type node interface {
	typ() nodeType
	String() string
}

type branchNode interface {
	append(n node)
}

type textNode struct {
	nodeType
	text string
}

func (t *textNode) String() string { return t.text }

type tagNode struct {
	nodeType
	nodes []node
	name  string
}

func (t *tagNode) append(n node) {
	t.nodes = append(t.nodes, n)
}

func (t *tagNode) String() string {
	nodes := []string{}
	for _, n := range t.nodes {
		nodes = append(nodes, n.String())
	}
	if t.typ() == nodeRoot {
		return strings.Join(nodes, "")
	}
	return fmt.Sprintf("<%s>%s</%s>", t.name, strings.Join(nodes, ""), t.name)
}

type tree struct {
	nodes     *tagNode
	lex       *lexer
	peekCount int
	text      string
	token     [3]item
}

func parse(input string) (*tree, error) {
	t := &tree{
		lex:   lex(input),
		nodes: &tagNode{nodeType: nodeRoot},
	}
	return t, t.parse()
}

func (t *tree) parse() error {
	tags := []*tagNode{t.nodes}

	app := func(n node, i item) {
		tags[len(tags)-1].append(n)
	}

	var item item
	for item = t.lex.nextItem(); item.typ != itemEOF && item.typ != itemError; item = t.lex.nextItem() {
		switch item.typ {
		case itemText:
			app(&textNode{nodeText, item.val}, item)
		case itemStartTag:
			tag := &tagNode{nodeType: nodeTag, name: item.val}
			app(tag, item)
			tags = append(tags, tag)
		case itemEndTag:
			if len(tags) == 1 {
				return fmt.Errorf("end tag </%s> without start tag", item.val)
			} else {
				tag := tags[len(tags)-1]
				if tag.name == item.val {
					// pop the tag
					tags = tags[0 : len(tags)-1]
				} else {
					return fmt.Errorf("expecting end tag </%s> got </%s>", tag.name, item.val)
				}
			}
		case itemLiteralLeft:
			app(&textNode{nodeLiteral, item.val}, item)
		case itemLiteralRight:
			app(&textNode{nodeLiteral, item.val}, item)
		}
	}

	if item.typ == itemError {
		return fmt.Errorf("%s", item.val)
	} else if len(tags) > 1 {
		return fmt.Errorf("expected </%s>", tags[len(tags)-1].name)
	}
	return nil
}

func (t *tree) String() string {
	return t.nodes.String()
}
