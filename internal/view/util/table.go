package util

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Section struct {
	titleCell *tview.TableCell
	textCell *tview.TableCell
}

func NewSection(title string, text string) *Section {
	return NewSectionWithColor(title, tcell.ColorDefault, text, tcell.ColorDefault)
}

func NewSectionWithColor(title string, titleColor tcell.Color, text string, textColor tcell.Color) *Section {
	// initialize a title cell
	titleCell := tview.NewTableCell(title)
	titleCell.SetAlign(tview.AlignLeft).
		SetTextColor(titleColor)

	// initialize a text cell
	textCell := tview.NewTableCell(text)
	textCell.SetAlign(tview.AlignLeft).
		SetExpansion(1).
		SetTextColor(textColor)

	return &Section{
		titleCell: titleCell,
		textCell: textCell,
	}
}

func (s *Section) GetTitleCell() *tview.TableCell {
	return s.titleCell
}

func (s *Section) GetTextCell() *tview.TableCell {
	return s.textCell
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
