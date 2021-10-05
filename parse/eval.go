package parse

import "strings"

type TagCallback func(tagName string, builder *strings.Builder, next func())

type formatter struct {
	cb TagCallback
}

func Eval(text string, cb TagCallback) (string, error) {
	f := &formatter{cb}
	return f.evalText(text)
}

func (f formatter) evalText(text string) (string, error) {
	t, err := parse(text)
	if err == nil {
		return f.evalTree(t), nil
	}
	return "", err
}

func (f formatter) evalTree(t *tree) string {
	builder := &strings.Builder{}
	f.evalTag(builder, t.nodes)
	return builder.String()
}

func (f formatter) evalTag(builder *strings.Builder, node *tagNode) {
	for _, node := range node.nodes {
		switch n := node.(type) {
		case *textNode:
			builder.WriteString(n.text)
		case *tagNode:
			next := func() { f.evalTag(builder, n) }
			f.cb(n.name, builder, next)
		}
	}
}
