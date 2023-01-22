package provider

import (
	"context"
)

type (
	GetAssetTransfersParams struct {
		FromBlock         string   `json:"fromBlock,omitempty"`
		ToBlock           string   `json:"toBlock,omitempty"`
		FromAddress       string   `json:"fromAddress,omitempty"`
		ToAddress         string   `json:"toAddress,omitempty"`
		ContractAddresses []string `json:"contractAddresses,omitempty"`
		Category          []string `json:"category,omitempty"`
		Order             string   `json:"order,omitempty"`
		WithMetadata      bool     `json:"withMetadata,omitempty"`
		ExcludeZeroValue  bool     `json:"excludeZeroValue,omitempty"`
		MaxCount          string   `json:"maxCount,omitempty"`
		PageKey           string   `json:"pageKey,omitempty"`
	}

	GetAssetTransfersResult struct {
		PageKey   string             `json:"pageKey"`
		Transfers []*AlchemyTransfer `json:"transfers"`
	}

	AlchemyRawContract struct {
		Value   string `json:"value"`
		Address string `json:"address"`
		Decimal string `json:"decimal"`
	}

	AlchemyTransfer struct {
		Category    string              `json:"category"`
		BlockNum    string              `json:"blockNum"`
		From        string              `json:"from"`
		To          string              `json:"to"`
		Value       float64             `json:"value"`
		TokenId     string              `json:"tokenId"`
		Asset       string              `json:"asset"`
		UniqueId    string              `json:"uniqueId"`
		Hash        string              `json:"hash"`
		RawContract *AlchemyRawContract `json:"rawContract"`
	}
)

func (p *Provider) GetAssetTransfers(params GetAssetTransfersParams) (*GetAssetTransfersResult, error) {
	var result *GetAssetTransfersResult
	err := p.rpcClient.CallContext(context.Background(), &result, "alchemy_getAssetTransfers", params)
	if err != nil {
		return nil, err
	}
	return result, nil
}
