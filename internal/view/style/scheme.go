package style

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Palette
//
//	Color60:         #5F5F87
//	Color69:         #5F87FF,
//	Color73:         #5FAFAF
//	Color147:        #AFAFFF,
//	ColorCoral:      #FF7F50,
//	ColorDimGray:    #696969,
//	ColorSandyBrown: #F4A460,
var Ethereum = &Style{
	FgColor:               tcell.ColorFloralWhite,
	BgColor:               tview.Styles.PrimitiveBackgroundColor,
	SectionColor:          tcell.ColorCoral,
	HelpKeyColor:          tcell.Color69,
	PrimaryTitleColor:     tcell.ColorFloralWhite,
	PrimaryBorderColor:    tcell.Color60,
	SecondaryTitleColor:   tcell.ColorFloralWhite,
	SecondaryBorderColor:  tcell.Color60,
	QueryBgColor:         tview.Styles.PrimitiveBackgroundColor,
	QueryBorderColor:     tcell.Color147,
	MethNameBorderColor:   tcell.Color147,
	MethArgsBorderColor:   tcell.Color147,
	MethResultBorderColor: tcell.Color69,
	TableHeaderStyle:      new(tcell.Style).Foreground(tcell.Color147).Bold(true),
	PrgBarCellColor:       tcell.ColorCoral,
	PrgBarTitleColor:      tcell.ColorFloralWhite,
	PrgBarBorderColor:     tcell.ColorDimGray,
	InputFieldLableColor:  tcell.ColorSandyBrown,
	InputFieldBgColor:     tcell.Color60,
}
