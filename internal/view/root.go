package view

import (
	"github.com/dyng/ramen/internal/common"
	"github.com/dyng/ramen/internal/service"
	"github.com/dyng/ramen/internal/view/style"
	"github.com/dyng/ramen/internal/view/util"
	"github.com/ethereum/go-ethereum/log"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type bodyPage interface {
	KeyMaps() util.KeyMaps
}

type Root struct {
	*tview.Flex
	app *App

	chainInfo    *ChainInfo
	help         *Help
	query        *QueryDialog
	notification *Notification
	body         *tview.Pages
	home         *Home
	account      *Account
	transaction  *TransactionDetail
}

func NewRoot(app *App) *Root {
	root := &Root{
		app: app,
	}

	// setup layout
	root.initLayout()

	// setup keymap
	root.initKeymap()

	return root
}

func (r *Root) initLayout() {
	// chainInfo
	chainInfo := NewChainInfo(r.app)
	r.chainInfo = chainInfo

	// help
	help := NewHelp(r.app)
	r.help = help

	// header
	header := tview.NewFlex().SetDirection(tview.FlexColumn)
	header.AddItem(chainInfo, 0, 1, false)
	header.AddItem(help, 0, 1, false)

	// body
	body := tview.NewPages()
	r.body = body

	// home page
	home := NewHome(r.app)
	body.AddPage("home", home, true, true)
	r.home = home

	// account page
	account := NewAccount(r.app)
	body.AddPage("account", account, true, false)
	r.account = account

	// transaction detail page
	transaction := NewTransactionDetail(r.app)
	body.AddPage("transaction", transaction, true, false)
	r.transaction = transaction

	// query dialog
	query := NewQueryDialog(r.app)
	r.query = query

	// notiication bar
	notification := NewNotification(r.app)
	r.notification = notification

	// root
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(header, style.HeaderHeight, 0, false).
		AddItem(body, 0, 1, true)
	r.Flex = flex
}

func (r *Root) initKeymap() {
	InitKeymap(r, r.app)
}

func (r *Root) KeyMaps() util.KeyMaps {
	keymaps := make(util.KeyMaps, 0)

	// KeySlash: show a query dialog
	keymaps = append(keymaps, util.KeyMap{
		Key:         util.KeySlash,
		Shortcut:    "/",
		Description: "Search Account",
		Handler: func(*tcell.EventKey) {
			r.ShowQueryDialog()
		},
	})

	// KeyH: back to home
	keymaps = append(keymaps, util.KeyMap{
		Key:         util.KeyH,
		Shortcut:    "H",
		Description: "Back to Home",
		Handler: func(*tcell.EventKey) {
			r.ShowHomePage()
		},
	})

	// KeyCtrlC: quit
	keymaps = append(keymaps, util.KeyMap{
		Key:         util.KeyQ,
		Shortcut:    "Q",
		Description: "Quit",
		Handler: func(*tcell.EventKey) {
			r.app.Stop()
		},
	})

	return keymaps
}

func (r *Root) ShowQueryDialog() {
	log.Debug("Show query dialog")
	if !r.query.IsDisplay() {
		r.query.Clear()
		r.query.Display(true)
	}
	r.app.SetFocus(r.query)
}

func (r *Root) HideQueryDialog() {
	log.Debug("Hide query dialog")
	if r.query.IsDisplay() {
		r.query.Display(false)
	}
	r.app.SetFocus(r)
}

func (r *Root) NotifyError(errmsg string) {
	r.ShowNotification("[crimson::b]ERROR[-::-]", errmsg)
}

func (r *Root) ShowNotification(title string, text string) {
	log.Debug("Show notification")
	if !r.notification.IsDisplay() {
		r.notification.SetContent(title, text)
		r.notification.Display(true)
	}
	r.app.SetFocus(r.notification)
}

func (r *Root) HideNotification() {
	log.Debug("Hide notification")
	if r.notification.IsDisplay() {
		r.notification.Display(false)
	}
	r.app.SetFocus(r)
}

func (r *Root) ShowHomePage() {
	log.Debug("Switch to home page")
	r.body.SwitchToPage("home")
	r.updateHelp(r.home)
}

func (r *Root) ShowAccountPage(account *service.Account) {
	log.Debug("Switch to account page", "account", account.GetAddress())
	r.account.SetAccount(account)
	r.body.SwitchToPage("account")
	r.updateHelp(r.account)
}

func (r *Root) ShowTransactionPage(transaction common.Transaction) {
	log.Debug("Switch to transaction page", "transaction", transaction.Hash())
	r.transaction.SetTransaction(transaction)
	r.body.SwitchToPage("transaction")
	r.updateHelp(r.transaction)
}

func (r *Root) updateHelp(page bodyPage) {
	keymaps := r.KeyMaps().
		Add(page.KeyMaps())

	log.Debug("Update help with keys from page", "keymaps", keymaps)
	r.help.SetKeyMaps(keymaps)
}

// Primitive Interface Implementation

// HasFocus implements tview.Primitive
func (r *Root) HasFocus() bool {
	if r.query.HasFocus() {
		return true
	}
	if r.notification.HasFocus() {
		return true
	}
	return r.Flex.HasFocus()
}

// InputHandler implements tview.Primitive
func (r *Root) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		if r.Flex.HasFocus() {
			if handler := r.Flex.InputHandler(); handler != nil {
				handler(event, setFocus)
				return
			}
		}
		if r.query.HasFocus() {
			if handler := r.query.InputHandler(); handler != nil {
				handler(event, setFocus)
				return
			}
		}
		if r.notification.HasFocus() {
			if handler := r.notification.InputHandler(); handler != nil {
				handler(event, setFocus)
				return
			}
		}
	}
}

// SetRect implements tview.SetRect
func (r *Root) SetRect(x int, y int, width int, height int) {
	r.Flex.SetRect(x, y, width, height)
	r.query.SetCentral(r.GetInnerRect())
	r.notification.SetCentral(r.GetInnerRect())
}

// Draw implements tview.Primitive
func (r *Root) Draw(screen tcell.Screen) {
	r.Flex.Draw(screen)
	r.query.Draw(screen)
	r.notification.Draw(screen)
}
