package colors

type Color string

func (c Color) String() string {
	return string(c)
}

const (
	Black       = Color("\u001b[30m")
	BrightWhite = Color("\u001b[97m")
	Red         = Color("\u001b[31m")
	Green       = Color("\u001b[32m")
	Yellow      = Color("\u001b[33m")
	Blue        = Color("\u001b[34m")
	Magenta     = Color("\u001b[35m")
	Cyan        = Color("\u001b[36m")
	White       = Color("\u001b[37m")
	Reset       = Color("\u001b[0m")
)

var colors = map[string]Color{
	"black":   Black,
	"red":     Red,
	"green":   Green,
	"yellow":  Yellow,
	"blue":    Blue,
	"magenta": Magenta,
	"cyan":    Cyan,
	"white":   White,
	"reset":   Reset,
}
