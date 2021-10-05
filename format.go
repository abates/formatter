package formatter

import (
	"strings"

	"github.com/abates/formatter/colors"
	"github.com/abates/formatter/parse"
)

var resetColor = colors.Reset

type Formatter func(input string) (output string, err error)

var ContextColors = map[string]colors.Color{
	"hl":      colors.BrightWhite,
	"em":      colors.Cyan,
	"warn":    colors.Yellow,
	"fail":    colors.Red,
	"success": colors.Green,
}

// ContextColorFormatter will format text based on
// some simple context tags.  The colors are
// looked up from ContextColors
//  <hl>highlight</hl> - colors.BrightWhite
//   <em>emphasis</em> - colors.Cyan
//       <warn></warn> - colors.Yellow
//       <fail></fail> - colors.Red
// <success></success> - colors.Green
func ContextFormatter() Formatter {
	return ColorFormatter(ContextColors)
}

type colorFormatter struct {
	tagMap     map[string]colors.Color
	colorStack []colors.Color
}

func (cf *colorFormatter) peek() colors.Color {
	c := colors.Color("")
	l := len(cf.colorStack) - 1
	if l >= 0 {
		c = cf.colorStack[0]
	}
	return c
}

func (cf *colorFormatter) pop() colors.Color {
	c := colors.Color("")
	l := len(cf.colorStack) - 1
	if l >= 0 {
		c = cf.colorStack[l]
		cf.colorStack = cf.colorStack[0:l]
	}
	return c
}

func (cf *colorFormatter) cb(tag string, builder *strings.Builder, next func()) {
	color, found := cf.tagMap[tag]
	if found {
		builder.WriteString(color.String())
		cf.colorStack = append(cf.colorStack, color)
	}
	next()
	if found {
		builder.WriteString(resetColor.String())
		cf.pop()
		builder.WriteString(cf.peek().String())
	}
}

// ColorFormatter returns a formatter that will map
// ANSI color codes into the output based on matching
// tags in the tag map.  For instance:
//   colors := map[string]colors.Color{
//     "hl":      colors.BrightWhite,
//     "em":      colors.Cyan,
//     "warn":    colors.Yellow,
//     "fail":    colors.Red,
//     "success": colors.Green,
//   }
//   formatter := ColorFormatter(colors)
func ColorFormatter(tagMap map[string]colors.Color) Formatter {
	cf := &colorFormatter{tagMap, []colors.Color{}}

	return func(input string) (string, error) {
		return parse.Eval(input, cf.cb)
	}
}
