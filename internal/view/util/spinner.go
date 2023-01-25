package util

import (
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Spinner struct {
	*tview.Box
	app *tview.Application

	display   bool
	counter   int
	ticker    *time.Ticker
}

var (
	frames = []rune(`⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏`)
)

func NewSpinner(app *tview.Application) *Spinner {
	spinner := &Spinner{
		Box:       &tview.Box{},
		app:       app,
		display:   false,
		counter:   0,
	}
	spinner.SetBorder(false)
	spinner.SetBackgroundColor(tview.Styles.PrimitiveBackgroundColor)

	return spinner
}

func (s *Spinner) StartAndShow() {
	s.Start()
	s.Display(true)
}

func (s *Spinner) StopAndHide() {
	s.Stop()
	s.Display(false)
}

func (s *Spinner) Start() {
	s.ticker = time.NewTicker(100 * time.Millisecond)

	update := func() {
		for {
			select {
			case <-s.ticker.C:
				s.pulse()
				if s.display {
					s.app.Draw()
				}
			}
		}
	}
	go update()
}

func (s *Spinner) Stop() {
	if s.ticker != nil {
		s.ticker.Stop()
	}
}

func (s *Spinner) pulse() {
	s.counter++
	if s.counter >= len(frames) {
		s.counter = s.counter % len(frames)
	}
}

func (s *Spinner) Display(display bool) {
	s.display = display
}

func (s *Spinner) IsDisplay() bool {
	return s.display
}

func (s *Spinner) SetCentral(x, y, width, height int) {
	spinnerWidth := 1
	spinnerHeight := 1
	spinnerX := x + ((width - spinnerWidth) / 2)
	spinnerY := y + ((height - spinnerHeight) / 2)
	s.Box.SetRect(spinnerX, spinnerY, spinnerWidth, spinnerHeight)
}

// SetRect implements tview.SetRect
func (s *Spinner) SetRect(x, y, width, height int) {
	s.Box.SetRect(x, y, 1, 1)
}

// Draw implements tview.Draw
func (s *Spinner) Draw(screen tcell.Screen) {
	if !s.display {
		return
	}

	s.Box.DrawForSubclass(screen, s)
	x, y, _, _ := s.Box.GetInnerRect()
	frame := string(frames[s.counter])
	tview.Print(screen, frame, x, y, 1, tview.AlignLeft, tcell.ColorDefault)
}
