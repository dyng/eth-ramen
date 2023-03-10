package style

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
)

const (
	// HeaderHeight is the height of header
	HeaderHeight = 5
	// AvatarSize is the width and height of an avatar
	AvatarSize = 4
)

type Style struct {
	// main
	FgColor       tcell.Color
	BgColor       tcell.Color
	SectionColor  tcell.Color
	SectionColor2 tcell.Color

	// help
	HelpKeyColor tcell.Color

	// body
	TitleColor   tcell.Color
	BorderColor  tcell.Color
	TitleColor2  tcell.Color
	BorderColor2 tcell.Color

	// methodCall
	MethResultBorderColor tcell.Color

	// table
	TableHeaderStyle tcell.Style

	// dialog
	DialogBgColor     tcell.Color
	DialogBorderColor tcell.Color
	ButtonBgColor     tcell.Color

	// progress bar
	PrgBarCellColor   tcell.Color
	PrgBarTitleColor  tcell.Color
	PrgBarBorderColor tcell.Color

	// others
	InputFieldLableColor tcell.Color
	InputFieldBgColor    tcell.Color
}

func Bold(s string) string {
	return fmt.Sprintf("[::b]%s[::-]", s)
}

func Padding(s string) string {
	return fmt.Sprintf(" %s ", s)
}

func BoldPadding(s string) string {
	return fmt.Sprintf(" [::b]%s[::-] ", s)
}

func Color(s string, color tcell.Color) string {
	return fmt.Sprintf("[#%06x::]%s[-::]", color.Hex(), s)
}
