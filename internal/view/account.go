package view

import (
	"github.com/dyng/ramen/internal/common/conv"
	serv "github.com/dyng/ramen/internal/service"
	"github.com/dyng/ramen/internal/view/util"
	"github.com/ethereum/go-ethereum/log"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Account struct {
	*tview.Flex
	app *App

	accountInfo     *AccountInfo
	transactionList *TransactionList
	methodCall      *MethodCallDialog
	account         *serv.Account
}

type AccountInfo struct {
	*tview.Table
	address     *util.Section
	accountType *util.Section
	balance     *util.Section
}

func NewAccount(app *App) *Account {
	account := &Account{
		app: app,
	}

	// setup layout
	account.initLayout()

	// setup keymap
	account.initKeymap()

	return account
}

func (a *Account) SetAccount(account *serv.Account) {
	// change current account
	a.account = account

	// populate contract field if account is a contract
	if account.IsContract() {
		contract, err := account.AsContract()
		if err == nil {
			a.methodCall.SetContract(contract)
		} else {
			// TODO: notify error
			log.Error("Failed to upgrade account to contract", "account", account.GetAddress(), "error", err)
		}
	}

	// refresh
	a.refresh()
}

func (a *Account) initLayout() {
	// AccountInfo
	accountInfo := &AccountInfo{
		Table:       tview.NewTable(),
		address:     util.NewSection("Address", util.NAValue),
		accountType: util.NewSection("Type", util.NAValue),
		balance:     util.NewSection("Balance", util.NAValue),
	}
	accountInfo.address.AddToTable(accountInfo.Table, 0, 0)
	accountInfo.accountType.AddToTable(accountInfo.Table, 0, 2)
	accountInfo.balance.AddToTable(accountInfo.Table, 1, 0)
	a.accountInfo = accountInfo

	// MethodCallDialog
	methodCall := NewMethodCallDialog(a.app)
	a.methodCall = methodCall

	// Transactions
	transactions := NewTransactionList(a.app)
	a.transactionList = transactions

	// Root
	flex := tview.NewFlex()
	flex.SetBorder(true)
	flex.SetTitle("[::b] Account [::-]")
	flex.SetDirection(tview.FlexRow)
	flex.AddItem(accountInfo, 0, 1, false)
	flex.AddItem(transactions, 0, 1, true)
	a.Flex = flex
}

func (a *Account) initKeymap() {
	keymaps := a.KeyMaps()
	a.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		handler, ok := keymaps.FindHandler(util.AsKey(event))
		if ok {
			handler(event)
			return nil
		} else {
			return event
		}
	})
}

func (a *Account) KeyMaps() util.KeyMaps {
	keymaps := make(util.KeyMaps, 0)

	// KeyC: call a contract
	keymaps = append(keymaps, util.KeyMap{
		Key:         util.KeyC,
		Shortcut:    "C",
		Description: "Call Contract",
		Handler: func(*tcell.EventKey) {
			a.ShowMethodCallDialog()
		},
	})

	return keymaps
}

func (a *Account) ShowMethodCallDialog() {
	log.Debug("Show method call dialog")

	if !a.account.IsContract() {
		return
	}

	if !a.methodCall.IsDisplay() {
		a.methodCall.Clear()
		a.methodCall.SetRect(a.GetInnerRect())
		a.methodCall.Display(true)
	}

	a.app.SetFocus(a.methodCall)
}

func (a *Account) HideMethodCallDialog() {
	log.Debug("Hide method call dialog")

	if a.methodCall.IsDisplay() {
		a.methodCall.Display(false)
	}

	a.app.SetFocus(a)
}

func (a *Account) refresh() {
	addr := a.account.GetAddress()
	a.accountInfo.address.SetText(addr.Hex())
	a.accountInfo.accountType.SetText(a.account.GetType().String())

	bal, err := a.account.GetBalance()
	if err == nil {
		a.accountInfo.balance.SetText(conv.ToEther(bal).String())
	} else {
		log.Error("Failed to fetch account balance", "account", addr, "error", err)
	}

	txns, err := a.account.GetTransactions()
	if err == nil {
		a.transactionList.SetTransactions(txns)
	} else {
		log.Error("Failed to fetch account transactions", "account", addr, "error", err)
	}
}

// Primitive Interface Implementation

// HasFocus implements tview.Primitive
func (a *Account) HasFocus() bool {
	if a.methodCall.HasFocus() {
		return true
	}
	return a.Flex.HasFocus()
}

// InputHandler implements tview.Primitive
func (a *Account) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		if a.Flex.HasFocus() {
			if handler := a.Flex.InputHandler(); handler != nil {
				handler(event, setFocus)
				return
			}
		}
		if a.methodCall.HasFocus() {
			if handler := a.methodCall.InputHandler(); handler != nil {
				handler(event, setFocus)
				return
			}
		}
	}
}

// Draw implements tview.Primitive
func (a *Account) Draw(screen tcell.Screen) {
	a.Flex.Draw(screen)
	a.methodCall.Draw(screen)
}
