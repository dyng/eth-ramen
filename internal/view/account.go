package view

import (
	"github.com/dyng/ramen/internal/common"
	"github.com/dyng/ramen/internal/common/conv"
	"github.com/dyng/ramen/internal/service"
	serv "github.com/dyng/ramen/internal/service"
	"github.com/dyng/ramen/internal/view/format"
	"github.com/dyng/ramen/internal/view/style"
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
	importABI       *ImportABIDialog
	account         *serv.Account
	contract        *serv.Contract
}

type AccountInfo struct {
	*tview.Flex
	avatar      *util.Avatar
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

	// subscribe to new blocks
	app.eventBus.Subscribe(service.TopicNewBlock, account.onNewBlock)

	return account
}

func (a *Account) SetAccount(account *serv.Account) {
	// change current account
	a.account = account

	// set base account
	base := a.account.GetAddress()
	a.transactionList.SetBaseAccount(&base)

	// populate contract field if account is a contract
	if account.IsContract() {
		contract, err := account.AsContract()
		if err == nil {
			a.contract = contract
			a.methodCall.SetContract(contract)
		} else {
			log.Error("Cannot upgrade account to contract", "account", account.GetAddress(), "error", err)
			a.app.root.NotifyError(format.FineErrorMessage("Cannot upgrade account to contract", err))
		}
	}

	// refresh
	a.refresh()
}

func (a *Account) initLayout() {
	s := a.app.config.Style()

	// AccountInfo
	accountInfo := &AccountInfo{
		Flex:        tview.NewFlex(),
		avatar:      util.NewAvatar(style.AvatarSize),
		address:     util.NewSectionWithStyle("Address", util.NAValue, s),
		accountType: util.NewSectionWithStyle("Type", util.NAValue, s),
		balance:     util.NewSectionWithStyle("Balance", util.NAValue, s),
	}

	info := tview.NewTable()
	accountInfo.address.AddToTable(info, 2, 0)
	accountInfo.accountType.AddToTable(info, 1, 0)
	accountInfo.balance.AddToTable(info, 3, 0)

	accountInfo.SetDirection(tview.FlexRow)
	accountInfo.AddItem(accountInfo.avatar, style.AvatarSize, 0, false)
	accountInfo.AddItem(info, 0, 1, false)
	a.accountInfo = accountInfo

	// MethodCallDialog
	methodCall := NewMethodCallDialog(a.app)
	a.methodCall = methodCall

	// ImportABIDialog
	importABI := NewImportABIDialog(a.app)
	a.importABI = importABI

	// Transactions
	transactions := NewTransactionList(a.app, true)
	transactions.SetTitleColor(s.TitleColor2)
	transactions.SetBorderColor(s.BorderColor2)
	a.transactionList = transactions

	// Root
	flex := tview.NewFlex()
	flex.SetBorder(true)
	flex.SetTitle(style.BoldPadding("Account"))
	flex.SetBorderColor(s.BorderColor)
	flex.SetTitleColor(s.TitleColor)
	flex.SetDirection(tview.FlexRow)
	flex.AddItem(accountInfo, 0, 3, false)
	flex.AddItem(transactions, 0, 7, true)
	a.Flex = flex
}

func (a *Account) initKeymap() {
	InitKeymap(a, a.app)
}

func (a *Account) KeyMaps() util.KeyMaps {
	keymaps := make(util.KeyMaps, 0)

	// KeyC: call a contract
	keymaps = append(keymaps, util.KeyMap{
		Key:         util.KeyC,
		Shortcut:    "C",
		Description: "Call Contract",
		Handler: func(*tcell.EventKey) {
			// TODO: don't show "Call Contract" for wallet account
			if a.account.IsContract() {
				if a.methodCall.contract.HasABI() {
					a.ShowMethodCallDialog()
				} else {
					a.ShowImportABIDialog()
				}
			}
		},
	})

	return keymaps
}

func (a *Account) ShowMethodCallDialog() {
	if !a.account.IsContract() {
		return
	}
	a.methodCall.Clear()
	a.methodCall.Show()
}

func (a *Account) ShowImportABIDialog() {
	a.importABI.Clear()
	a.importABI.Show()
}

func (a *Account) onNewBlock(block *common.Block) {
	if a.account == nil {
		return
	}

	txns, err := a.app.service.GetTransactionsByBlock(block)
	if err != nil {
		log.Error("cannot extract transactions from block", "blockHash", block.Hash(), "error", err)
		return
	}

	// update current account
	a.account.UpdateBalance()

	a.app.QueueUpdateDraw(func() {
		a.refresh()
		a.transactionList.FilterAndPrependTransactions(txns)
	})
}

func (a *Account) refresh() {
	addr := a.account.GetAddress()
	a.accountInfo.address.SetText(addr.Hex())
	a.accountInfo.accountType.SetText(StyledAccountType(a.account.GetType()))

	// avatar
	a.accountInfo.avatar.SetAddress(addr)

	// fetch balance
	bal := a.account.GetBalance()
	a.accountInfo.balance.SetText(conv.ToEther(bal).String())

	// update transaction history asynchronously
	a.transactionList.LoadAsync(a.account.GetTransactions)
}

// Primitive Interface Implementation

// HasFocus implements tview.Primitive
func (a *Account) HasFocus() bool {
	if a.methodCall.HasFocus() {
		return true
	}
	if a.importABI.HasFocus() {
		return true
	}
	return a.Flex.HasFocus()
}

// InputHandler implements tview.Primitive
func (a *Account) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		if a.methodCall.HasFocus() {
			if handler := a.methodCall.InputHandler(); handler != nil {
				handler(event, setFocus)
				return
			}
		}
		if a.importABI.HasFocus() {
			if handler := a.importABI.InputHandler(); handler != nil {
				handler(event, setFocus)
				return
			}
		}
		if a.Flex.HasFocus() {
			if handler := a.Flex.InputHandler(); handler != nil {
				handler(event, setFocus)
				return
			}
		}
	}
}

// SetRect implements tview.SetRect
func (a *Account) SetRect(x int, y int, width int, height int) {
	a.Flex.SetRect(x, y, width, height)
	a.methodCall.SetCentral(a.GetInnerRect())
	a.importABI.SetCentral(a.GetInnerRect())
}

// Draw implements tview.Primitive
func (a *Account) Draw(screen tcell.Screen) {
	a.Flex.Draw(screen)
	a.methodCall.Draw(screen)
	a.importABI.Draw(screen)
}
