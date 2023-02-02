package view

import (
	"github.com/dyng/ramen/internal/common/conv"
	"github.com/dyng/ramen/internal/service"
	"github.com/dyng/ramen/internal/view/style"
	"github.com/dyng/ramen/internal/view/util"
	"github.com/ethereum/go-ethereum/log"
	"github.com/rivo/tview"
)

type Signer struct {
	tview.Primitive
	app *App

	signer      *service.Signer
	initialized bool
	avatar      *util.Avatar
	table       *tview.Table
	address     *util.Section
	balance     *util.Section
}

func NewSigner(app *App) *Signer {
	signer := &Signer{
		app:    app,
		avatar: util.NewAvatar(style.AvatarSize),
		table:  tview.NewTable(),
	}

	// setup layout
	signer.initLayout()

	return signer
}

func (si *Signer) initLayout() {
	// not signed in by default
	si.layoutNoSigner()
}

func (si *Signer) HasSignedIn() bool {
	return si.signer != nil
}

func (si *Signer) GetSigner() *service.Signer {
	return si.signer
}

func (si *Signer) SetSigner(signer *service.Signer) {
	si.signer = signer
	si.refresh()
}

func (si *Signer) refresh() {
	if !si.initialized {
		si.layoutSomeSigner()
		si.initialized = true
	}

	current := si.signer
	addr := current.GetAddress()

	// update avatar
	si.avatar.SetAddress(addr)

	// update address
	si.address.SetText(addr.Hex())

	// update balance
	bal, err := current.GetBalance()
	if err == nil {
		si.balance.SetText(conv.ToEther(bal).String())
	} else {
		log.Error("Failed to fetch account balance", "account", addr, "error", err)
	}
}

func (si *Signer) layoutNoSigner() {
	cell := tview.NewTableCell("[dimgrey]Not Signed In[-]")
	cell.SetAlign(tview.AlignLeft)
	cell.SetExpansion(1)
	si.table.SetCell(0, 0, cell)
	si.Primitive = si.table
}

func (si *Signer) layoutSomeSigner() {
	s := si.app.config.Style()

	flex := tview.NewFlex()
	flex.SetDirection(tview.FlexColumn)
	flex.AddItem(si.avatar, style.AvatarSize*2+1, 0, false)
	flex.AddItem(si.table, 0, 1, false)

	address := util.NewSectionWithColor("Address:", s.SectionColor2, util.NAValue, s.FgColor)
	address.AddToTable(si.table, 0, 0)
	si.address = address

	balance := util.NewSectionWithColor("Balance:", s.SectionColor2, util.NAValue, s.FgColor)
	balance.AddToTable(si.table, 1, 0)
	si.balance = balance

	si.Primitive = flex
}
