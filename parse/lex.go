package parse

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

type itemType int

const eof = -1

const (
	itemError itemType = iota
	itemEndTag
	itemEOF
	itemLiteralLeft
	itemLiteralRight
	itemStartTag
	itemText
)

const (
	endTagMeta       = "</"
	leftLiteralMeta  = "&lt;"
	leftMeta         = "<"
	rightLiteralMeta = "&gt;"
	rightMeta        = ">"
)

type item struct {
	typ itemType
	val string
}

type stateFn func(*lexer) stateFn

type lexer struct {
	state stateFn
	input string
	start int
	pos   int
	width int
	items chan item
}

func lex(input string) *lexer {
	l := &lexer{
		input: input,
		state: lexText,
		items: make(chan item, 2),
	}

	return l
}

func (l *lexer) nextItem() item {
	for {
		select {
		case t := <-l.items:
			return t
		default:
			l.state = l.state(l)
			if l.state == nil {
				break
			}
		}
	}
	// unreachable
}

func (l *lexer) next() (r rune) {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}

	r, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width
	return
}

func (l *lexer) backup() {
	l.pos -= l.width
}

// error returns an error token and terminates the scan
// by passing back a nil pointer that will be the next
// state, terminating l.run.
func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	l.items <- item{
		itemError,
		fmt.Sprintf(format, args...),
	}
	return nil
}

func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

func (l *lexer) ignore() {
	l.start = l.pos
}

func (l *lexer) acceptUntil(invalid string) bool {
	for strings.IndexRune(invalid, l.next()) == -1 {
		if l.peek() == eof {
			return false
		}
	}
	l.backup()
	return true
}

func (l *lexer) emit(itemType itemType) {
	l.items <- item{itemType, l.input[l.start:l.pos]}
	l.start = l.pos
}

func lexText(l *lexer) stateFn {
	emitIf := func() {
		if l.pos > l.start {
			l.emit(itemText)
		}
	}

	for {
		switch {
		case strings.HasPrefix(l.input[l.pos:], leftLiteralMeta):
			emitIf()
			return lexAngle(itemLiteralLeft, leftLiteralMeta)
		case strings.HasPrefix(l.input[l.pos:], rightLiteralMeta):
			emitIf()
			return lexAngle(itemLiteralRight, rightLiteralMeta)
		case strings.HasPrefix(l.input[l.pos:], endTagMeta):
			emitIf()
			return lexEndTag
		case strings.HasPrefix(l.input[l.pos:], leftMeta):
			emitIf()
			return lexStartTag
		}

		if l.next() == eof {
			break
		}
	}

	if l.pos > l.start {
		l.emit(itemText)
	}
	l.emit(itemEOF)
	return nil
}

func lexAngle(typ itemType, meta string) stateFn {
	return func(l *lexer) stateFn {
		l.pos += len(meta)
		l.emit(typ)
		return lexText
	}
}

func lexStartTag(l *lexer) stateFn {
	// ignore the opening angle bracket
	l.next()
	l.ignore()

	if l.acceptUntil(rightMeta) {
		l.emit(itemStartTag)
		// ignore right angle bracket
		l.next()
		l.ignore()
		return lexText
	}

	return l.errorf("expecting >")
}

func lexEndTag(l *lexer) stateFn {
	// ignore opening angle bracket and slash
	l.next()
	l.next()
	l.ignore()

	if l.acceptUntil(rightMeta) {
		l.emit(itemEndTag)
		// ignore right angle bracket
		l.next()
		l.ignore()
		return lexText
	}

	return l.errorf("expecting %q at %q", rightMeta, l.input[l.pos:])
}
