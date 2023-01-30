package view

import (
	"github.com/dyng/ramen/internal/common/conv"
	"github.com/dyng/ramen/internal/service"
	"github.com/dyng/ramen/internal/view/util"
	"github.com/ethereum/go-ethereum/log"
	"github.com/rivo/tview"
)

type Signer struct {
	*tview.Table
	app *App

	signer      *service.Signer
	initialized bool
	address     *util.Section
	balance     *util.Section
}

func NewSigner(app *App) *Signer {
	signer := &Signer{
		Table: tview.NewTable(),
		app:   app,
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
	cell := tview.NewTableCell("[crimson]Not Signed In[-]")
	cell.SetAlign(tview.AlignLeft)
	cell.SetExpansion(1)
	si.Table.SetCell(0, 0, cell)
}

func (si *Signer) layoutSomeSigner() {
	s := si.app.config.Style()

	address := util.NewSectionWithColor("Address:", s.SectionColor2, util.NAValue, s.FgColor)
	address.AddToTable(si.Table, 0, 0)
	si.address = address

	balance := util.NewSectionWithColor("Balance:", s.SectionColor2, util.NAValue, s.FgColor)
	balance.AddToTable(si.Table, 1, 0)
	si.balance = balance
}
