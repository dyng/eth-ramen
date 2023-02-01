package util

import (
	"fmt"
	"hash/fnv"

	"github.com/dyng/ramen/internal/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/pkg/errors"
	"github.com/rivo/tview"
	"github.com/rrivera/identicon"
)

const (
	backgroundColor = "dimgrey"
)

var (
	palette = []string{
		"black",
		"maroon",
		"green",
		"olive",
		"navy",
		"purple",
		"teal",
		"silver",
		"red",
		"lime",
		"yellow",
		"blue",
		"aqua",
		"aliceblue",
		"antiquewhite",
		"aquamarine",
		"azure",
		"beige",
		"bisque",
		"blanchedalmond",
		"blueviolet",
		"brown",
		"burlywood",
		"cadetblue",
		"chartreuse",
		"chocolate",
		"coral",
		"cornflowerblue",
		"cornsilk",
		"darkblue",
		"darkcyan",
		"darkgoldenrod",
		"darkgreen",
		"darkkhaki",
		"darkmagenta",
		"darkolivegreen",
		"darkorange",
		"darkorchid",
		"darkred",
		"darksalmon",
		"darkseagreen",
		"darkslateblue",
		"darkturquoise",
		"darkviolet",
		"deeppink",
		"deepskyblue",
		"dodgerblue",
		"firebrick",
		"floralwhite",
		"forestgreen",
		"gainsboro",
		"ghostwhite",
		"gold",
		"goldenrod",
		"greenyellow",
		"honeydew",
		"hotpink",
		"indianred",
		"indigo",
		"ivory",
		"khaki",
		"lavender",
		"lavenderblush",
		"lawngreen",
		"lemonchiffon",
		"lightblue",
		"lightcoral",
		"lightcyan",
		"lightgoldenrodyellow",
		"lightgreen",
		"lightpink",
		"lightsalmon",
		"lightseagreen",
		"lightskyblue",
		"lightsteelblue",
		"lightyellow",
		"limegreen",
		"linen",
		"mediumaquamarine",
		"mediumblue",
		"mediumorchid",
		"mediumpurple",
		"mediumseagreen",
		"mediumslateblue",
		"mediumspringgreen",
		"mediumturquoise",
		"mediumvioletred",
		"midnightblue",
		"mintcream",
		"mistyrose",
		"moccasin",
		"navajowhite",
		"oldlace",
		"olivedrab",
		"orange",
		"orangered",
		"orchid",
		"palegoldenrod",
		"palegreen",
		"paleturquoise",
		"palevioletred",
		"papayawhip",
		"peachpuff",
		"peru",
		"pink",
		"plum",
		"powderblue",
		"rebeccapurple",
		"rosybrown",
		"royalblue",
		"saddlebrown",
		"salmon",
		"sandybrown",
		"seagreen",
		"seashell",
		"sienna",
		"skyblue",
		"slateblue",
		"snow",
		"springgreen",
		"steelblue",
		"tan",
		"thistle",
		"tomato",
		"turquoise",
		"violet",
		"wheat",
		"whitesmoke",
		"yellowgreen",
	}
)

type Avatar struct {
	*tview.Table
	size      int
	generator *identicon.Generator
	address   common.Address
}

func NewAvatar(size int) *Avatar {
	ig, err := identicon.New("ramen", size, 3)
	if err != nil {
		log.Crit("Failed to create identicon generator", "error", errors.WithStack(err))
	}

	return &Avatar{
		Table:     tview.NewTable(),
		size:      size,
		generator: ig,
	}
}

func (a *Avatar) SetAddress(address common.Address) {
	bitmap, color := a.identiconFrom(address)
	log.Debug("Generated avatar for account", "address", address, "bitmap", bitmap, "color", color)
	for i := 0; i < a.size; i++ {
		text := ""
		for j := 0; j < a.size; j++ {
			if bitmap[i][j] > 0 {
				text += fmt.Sprintf("[%s]██[-]", color)
			} else {
				text += fmt.Sprintf("[%s]██[-]", backgroundColor)
			}
		}
		a.SetCell(i, 0, tview.NewTableCell(text))
	}
}

func (a *Avatar) identiconFrom(address common.Address) (bitmap [][]int, color string) {
	icon, err := a.generator.Draw(address.Hex())
	if err != nil {
		log.Error("Failed to generate identicon, fallback to default", "error", errors.WithStack(err))
		return a.defaultIdenticon()
	}

	bitmap = icon.Array()
	color = a.selectColor(address.Hex())
	return
}

func (a *Avatar) defaultIdenticon() (bitmap [][]int, color string) {
	color = backgroundColor
	bitmap = make([][]int, a.size)
	for i := 0; i < a.size; i++ {
		bitmap[i] = make([]int, a.size)
		for j := 0; j < a.size; j++ {
			bitmap[i][j] = 0
		}
	}
	return
}

func (a *Avatar) selectColor(text string) string {
	h := fnv.New32a()
	h.Write([]byte(text))
	i := h.Sum32()
	return palette[i%uint32(len(palette))]
}
