package etherscan

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/dyng/ramen/internal/common"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/log"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

const (
	// DefaultTimeout is the default value for request timeout
	// FIXME: code duplication
	DefaultTimeout = 30 * time.Second
)

type EtherscanClient struct {
	endpoint string
	apiKey   string
}

func NewEtherscanClient(endpoint string, apiKey string) *EtherscanClient {
	return &EtherscanClient{
		endpoint: endpoint,
		apiKey:   apiKey,
	}
}

func (c *EtherscanClient) AccountTxList(address common.Address) (common.Transactions, error) {
	req, err := http.NewRequest(http.MethodGet, c.endpoint, nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// build request
	q := req.URL.Query()
	q.Add("apikey", c.apiKey)
	q.Add("module", "account")
	q.Add("action", "txlist")
	q.Add("address", address.Hex())
	q.Add("startblock", "0")
	q.Add("endblock", "99999999")
	q.Add("sort", "desc")
	q.Add("page", "1")
	q.Add("offset", "100")
	req.URL.RawQuery = q.Encode()

	result, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	esTxns := make([]*esTransaction, 0)
	if err = json.Unmarshal(result, &esTxns); err != nil {
		return nil, errors.WithStack(err)
	}

	txns := make(common.Transactions, len(esTxns))
	for i, et := range esTxns {
		txns[i] = et
	}

	return txns, nil
}

func (c *EtherscanClient) GetSourceCode(address common.Address) (string, *abi.ABI, error) {
	req, err := http.NewRequest(http.MethodGet, c.endpoint, nil)
	if err != nil {
		return "", nil, errors.WithStack(err)
	}

	// build request
	q := req.URL.Query()
	q.Add("apikey", c.apiKey)
	q.Add("module", "contract")
	q.Add("action", "getsourcecode")
	q.Add("address", address.Hex())
	req.URL.RawQuery = q.Encode()

	result, err := c.doRequest(req)
	if err != nil {
		return "", nil, err
	}

	var codes []contractJSON
	if err = json.Unmarshal(result, &codes); err != nil {
		return "", nil, errors.WithStack(err)
	}

	code := codes[0]

	// contract source code not verified
	if code.SourceCode == "" {
		return "", nil, nil
	}

	parsedAbi, err := abi.JSON(strings.NewReader(code.ABI))
	if err != nil {
		return "", nil, errors.WithStack(err)
	}

	return code.SourceCode, &parsedAbi, nil
}

func (c *EtherscanClient) EthPrice() (*decimal.Decimal, error) {
	req, err := http.NewRequest(http.MethodGet, c.endpoint, nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// build request
	q := req.URL.Query()
	q.Add("apikey", c.apiKey)
	q.Add("module", "stats")
	q.Add("action", "ethprice")
	req.URL.RawQuery = q.Encode()

	result, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var ethprice ethpriceJSON
	if err = json.Unmarshal(result, &ethprice); err != nil {
		return nil, errors.WithStack(err)
	}

	price, err := decimal.NewFromString(ethprice.EthUsd)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &price, nil
}

func (c *EtherscanClient) doRequest(request *http.Request) ([]byte, error) {
	ctx, cancel := c.createContext()
	defer cancel()

	// set timeout
	request = request.WithContext(ctx)

	res, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if res.StatusCode != 200 {
		log.Error("HTTP status code is not OK", "status", res.Status, "body", resBody)
		return nil, errors.New("HTTP status code is not OK")
	}

	resMsg := resMessage{}
	if err = json.Unmarshal(resBody, &resMsg); err != nil {
		return nil, errors.WithStack(err)
	}

	if resMsg.Status == "0" {
		msg := string(resMsg.Result)
		log.Warn(fmt.Sprintf("Etherscan API status code is not OK. message is '%s'", msg))
	}

	return resMsg.Result, nil
}

func (c *EtherscanClient) createContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), DefaultTimeout)
}
