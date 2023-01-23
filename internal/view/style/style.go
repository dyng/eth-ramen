package style

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
)

const (
	// HeaderHeight is the height of header
	HeaderHeight = 5
)

type Style struct {
	// main
	FgColor      tcell.Color
	BgColor      tcell.Color
	SectionColor tcell.Color

	// help
	HelpKeyColor tcell.Color

	// body
	PrimaryTitleColor    tcell.Color
	PrimaryBorderColor   tcell.Color
	SecondaryTitleColor  tcell.Color
	SecondaryBorderColor tcell.Color

	// prompt
	PromptBgColor tcell.Color
	PromptBorderColor tcell.Color

	// methodCall
	MethNameBorderColor tcell.Color
	MethArgsBorderColor tcell.Color
	MethResultBorderColor tcell.Color

	// table
	TableHeaderStyle tcell.Style

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
