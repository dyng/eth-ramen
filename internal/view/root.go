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

	// header
	chainInfo *ChainInfo
	signer    *Signer
	help      *Help

	// body
	body        *tview.Pages
	home        *Home
	account     *Account
	transaction *TransactionDetail

	// dialogs
	query        *QueryDialog
	notification *Notification
	signin       *SignInDialog
	transfer     *TransferDialog
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

	// signer
	signer := NewSigner(r.app)
	r.signer = signer

	// help
	help := NewHelp(r.app)
	r.help = help

	// header
	header := tview.NewFlex().SetDirection(tview.FlexColumn)
	header.AddItem(chainInfo, 0, 6, false)
	header.AddItem(signer, 0, 6, false)
	header.AddItem(help, 0, 4, false)

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

	// signin dialog
	signin := NewSignInDialog(r.app)
	r.signin = signin

	// transfer dialog
	transfer := NewTransferDialog(r.app)
	r.transfer = transfer

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
		Description: "Search",
		Handler: func(*tcell.EventKey) {
			r.ShowQueryDialog()
		},
	})

	// KeyH: back to home
	keymaps = append(keymaps, util.KeyMap{
		Key:         util.KeyH,
		Shortcut:    "h",
		Description: "Home",
		Handler: func(*tcell.EventKey) {
			r.ShowHomePage()
		},
	})

	// KeyS: signin
	keymaps = append(keymaps, util.KeyMap{
		Key:         util.KeyS,
		Shortcut:    "s",
		Description: "Sign In",
		Handler: func(*tcell.EventKey) {
			r.ShowSignInDialog()
		},
	})

	// KeyM: transfer
	keymaps = append(keymaps, util.KeyMap{
		Key:         util.KeyM,
		Shortcut:    "m",
		Description: "Transfer",
		Handler: func(*tcell.EventKey) {
			r.ShowTransferDialog()
		},
	})

	// KeyCtrlC: quit
	keymaps = append(keymaps, util.KeyMap{
		Key:         tcell.KeyCtrlC,
		Shortcut:    "ctrl-c",
		Description: "Quit",
		Handler: func(*tcell.EventKey) {
			r.app.Stop()
		},
	})

	return keymaps
}

func (r *Root) ShowQueryDialog() {
	r.query.Clear()
	r.query.Show()
}

func (r *Root) NotifyInfo(message string) {
	r.ShowNotification("[lightgreen::b]INFO[-::-]", message)
}

func (r *Root) NotifyError(errmsg string) {
	r.ShowNotification("[crimson::b]ERROR[-::-]", errmsg)
}

func (r *Root) ShowNotification(title string, text string) {
	r.notification.SetContent(title, text)
	r.notification.Show()
}

func (r *Root) ShowSignInDialog() {
	r.signin.Clear()
	r.signin.Show()
}

func (r *Root) ShowTransferDialog() {
	if r.signer.HasSignedIn() {
		r.transfer.ClearAndRefresh()
		r.transfer.Show()
	}
}

func (r *Root) SignIn(signer *service.Signer) {
	log.Debug("Account signed in", "account", signer.GetAddress())
	r.signer.SetSigner(signer)
	r.transfer.SetSender(signer)
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

	log.Debug("Update keys help", "keymaps", keymaps)
	r.help.SetKeyMaps(keymaps)
}

// Primitive Interface Implementation

// HasFocus implements tview.Primitive
func (r *Root) HasFocus() bool {
	if r.query.HasFocus() {
		return true
	}
	if r.signin.HasFocus() {
		return true
	}
	if r.transfer.HasFocus() {
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
		if r.signin.HasFocus() {
			if handler := r.signin.InputHandler(); handler != nil {
				handler(event, setFocus)
				return
			}
		}
		if r.transfer.HasFocus() {
			if handler := r.transfer.InputHandler(); handler != nil {
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
	r.signin.SetCentral(r.GetInnerRect())
	r.transfer.SetCentral(r.GetInnerRect())
	r.notification.SetCentral(r.GetInnerRect())
}

// Draw implements tview.Primitive
func (r *Root) Draw(screen tcell.Screen) {
	r.Flex.Draw(screen)
	r.query.Draw(screen)
	r.signin.Draw(screen)
	r.transfer.Draw(screen)
	r.notification.Draw(screen)
}
