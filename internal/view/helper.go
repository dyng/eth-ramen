package view

import (
	"fmt"

	"github.com/dyng/ramen/internal/common"
	serv "github.com/dyng/ramen/internal/service"
)

func Inc(i *int) int {
	t := *i
	*i++
	return t
}

func StyledAccountType(t serv.AccountType) string {
	switch t {
	case serv.TypeWallet:
		return fmt.Sprintf("[lightgreen::b]%s[-:-:-]", t)
	case serv.TypeContract:
		return fmt.Sprintf("[dodgerblue::b]%s[-:-:-]", t)
	default:
		return t.String()
	}
}

func StyledNetworkName(n serv.Network) string {
	netType := n.NetType()

	if netType == serv.TypeMainnet {
		return "[crimson::b]Mainnet[-:-:-]"
	}

	if netType == serv.TypeTestnet {
		return fmt.Sprintf("[lightgreen::b]%s[-:-:-]", n.Name)
	}

	chainId := n.ChainId.String()

	if chainId == "1337" {
		return "[lightgreen::b]Ganache[-:-:-]"
	}

	if chainId == "31337" {
		return "[lightgreen::b]Hardhat[-:-:-]"
	}

	return n.Name
}

func StyledTxnDirection(base *common.Address, txn common.Transaction) string {
	if base == nil {
		return ""
	}

	if txn.From().String() == base.String() {
		return "[sandybrown]OUT[-]"
	} else {
		return "[lightgreen]IN[-]"
	}
}
