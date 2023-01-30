package style

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Palette
//
//	Color60:         #5F5F87
//	Color69:         #5F87FF
//	Color73:         #5FAFAF
//	Color147:        #AFAFFF
//	ColorCoral:      #FF7F50
//	ColorDimGray:    #696969
//	ColorSandyBrown: #F4A460
var Ethereum = &Style{
	FgColor:               tcell.ColorFloralWhite,
	BgColor:               tview.Styles.PrimitiveBackgroundColor,
	SectionColor:          tcell.ColorCoral,
	SectionColor2:         tcell.Color73,
	HelpKeyColor:          tcell.Color69,
	TitleColor:            tcell.ColorFloralWhite,
	BorderColor:           tcell.Color60,
	TitleColor2:           tcell.ColorFloralWhite,
	BorderColor2:          tcell.Color60,
	MethResultBorderColor: tcell.Color69,
	TableHeaderStyle:      new(tcell.Style).Foreground(tcell.Color147).Bold(true),
	DialogBgColor:         tview.Styles.PrimitiveBackgroundColor,
	DialogBorderColor:     tcell.Color147,
	ButtonBgColor:         tcell.ColorCoral,
	PrgBarCellColor:       tcell.ColorCoral,
	PrgBarTitleColor:      tcell.ColorFloralWhite,
	PrgBarBorderColor:     tcell.ColorDimGray,
	InputFieldLableColor:  tcell.ColorSandyBrown,
	InputFieldBgColor:     tcell.Color60,
}
