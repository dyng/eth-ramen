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

	chainInfo *ChainInfo
	help      *Help
	prompt    *PromptDialog
	body      *tview.Pages
	home      *Home
	account   *Account
	transaction *TransactionDetail
}

func NewRoot(app *App) *Root {
	root := &Root{
		app: app,
	}

	// setup layout
	root.initLayout()

	// setup keys
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
	header := tview.NewGrid().
		SetSize(1, 2, -1, -1)
	header.AddItem(chainInfo, 0, 0, 1, 1, 0, 0, false)
	header.AddItem(help, 0, 1, 1, 1, 0, 0, false)

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

	// prompt dialog
	prompt := NewPromptDialog(r.app)
	r.prompt = prompt

	// root
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(header, style.HeaderHeight, 0, false).
		AddItem(body, 0, 1, true)
	r.Flex = flex
}

func (r *Root) initKeymap() {
	keymaps := r.KeyMaps()
	r.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		handler, ok := keymaps.FindHandler(util.AsKey(event))
		if ok {
			handler(event)
			return nil
		} else {
			return event
		}
	})
}

func (r *Root) KeyMaps() util.KeyMaps {
	keymaps := make(util.KeyMaps, 0)

	// KeySpace: show a prompt dialog
	keymaps = append(keymaps, util.KeyMap{
		Key:         util.KeySpace,
		Shortcut:    "Space",
		Description: "Start Query",
		Handler: func(*tcell.EventKey) {
			r.ShowPrompt()
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

	return keymaps
}

func (r *Root) ShowPrompt() {
	log.Debug("Show prompt window")
	if !r.prompt.IsDisplay() {
		r.prompt.Clear()
		r.prompt.SetRect(r.GetInnerRect())
		r.prompt.Display(true)
	}
	r.app.SetFocus(r.prompt)
}

func (r *Root) HidePrompt() {
	log.Debug("Hide prompt window")
	if r.prompt.IsDisplay() {
		r.prompt.Display(false)
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
	if r.prompt.HasFocus() {
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
		if r.prompt.HasFocus() {
			if handler := r.prompt.InputHandler(); handler != nil {
				handler(event, setFocus)
				return
			}
		}
	}
}

// Draw implements tview.Primitive
func (r *Root) Draw(screen tcell.Screen) {
	r.Flex.Draw(screen)
	r.prompt.Draw(screen)
}
