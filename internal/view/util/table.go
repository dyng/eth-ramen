package util

import "github.com/rivo/tview"

type Section struct {
	titleCell *tview.TableCell
	textCell *tview.TableCell
}

func NewSection(title string, text string) *Section {
	// initialize a title cell
	titleCell := tview.NewTableCell(title + ":")
	titleCell.SetAlign(tview.AlignLeft)

	// initialize a text cell
	textCell := tview.NewTableCell(text)
	textCell.SetExpansion(2)

	return &Section{
		titleCell: titleCell,
		textCell: textCell,
	}
}

func (s *Section) GetText() string {
	return s.textCell.Text
}

func (s *Section) SetText(text string) {
	s.textCell.SetText(text)
}

func (s *Section) AddToTable(table *tview.Table, row, column int) {
	table.SetCell(row, column, s.titleCell)
	table.SetCell(row, column + 1, s.textCell)
}
