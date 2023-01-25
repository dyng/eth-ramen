package util

import (
	"time"

	"github.com/dyng/ramen/internal/view/style"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	ProgressBarWidth = 3

	ProgressBarHeight = 1

	ProgressBarCell = "▉"
)

type Loader struct {
	*tview.Box
	app *tview.Application

	display   bool
	counter   int
	cellColor tcell.Color
	ticker    *time.Ticker
}

func NewLoader(app *tview.Application) *Loader {
	loader := &Loader{
		Box:       &tview.Box{},
		app:       app,
		display:   false,
		counter:   0,
		cellColor: tcell.ColorDarkOrange,
	}
	loader.SetBorder(true)
	loader.SetTitle(style.Padding("LOADING"))
	loader.SetTitleAlign(tview.AlignCenter)
	loader.SetBackgroundColor(tview.Styles.PrimitiveBackgroundColor)

	return loader
}

func (l *Loader) SetCellColor(color tcell.Color) {
	l.cellColor = color
}

func (l *Loader) Start() {
	l.ticker = time.NewTicker(500 * time.Millisecond)

	update := func() {
		for {
			select {
			case <-l.ticker.C:
				l.pulse()
				if l.display {
					l.app.Draw()
				}
			}
		}
	}
	go update()
}

func (l *Loader) Stop() {
	if l.ticker != nil {
		l.ticker.Stop()
	}
}

func (l *Loader) pulse() {
	l.counter++
}

func (l *Loader) Display(display bool) {
	l.display = display
}

func (l *Loader) IsDisplay() bool {
	return l.display
}

// SetRect implements tview.SetRect
func (l *Loader) SetCentral(x, y, width, height int) {
	loaderWidth := 20
	loaderHeight := ProgressBarHeight + 2
	loaderX := x + ((width - loaderWidth) / 2)
	loaderY := y + ((height - loaderHeight) / 2)
	l.Box.SetRect(loaderX, loaderY, loaderWidth, loaderHeight)
}

// Draw implements tview.Draw
func (l *Loader) Draw(screen tcell.Screen) {
	if !l.display {
		return
	}

	l.Box.DrawForSubclass(screen, l)
	x, y, width, height := l.Box.GetInnerRect()

	if l.counter >= width-ProgressBarWidth {
		l.counter = l.counter % (width - ProgressBarWidth)
	}
	for i := 0; i < height; i++ {
		for j := 0; j < ProgressBarWidth; j++ {
			tview.Print(screen, ProgressBarCell, x+l.counter+j, y+i, width, tview.AlignLeft, l.cellColor)
		}
	}
}
