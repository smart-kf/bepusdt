package monitor

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	xlogger "github.com/clearcodecn/log"
)

type TronApiClient struct {
	httpClient *http.Client
	tronConfig *TronConfig
}

type TronConfig struct {
	ContractAddress string
	ApiHost         string
	APiKey          string
	Timeout         time.Duration
	Proxy           string
}

func NewTronApiClient(
	config *TronConfig,
) *TronApiClient {
	var cli = http.Client{
		Timeout: config.Timeout,
	}
	if config.Proxy != "" {
		uri, _ := url.Parse(config.Proxy)
		cli.Transport = &http.Transport{
			Proxy: http.ProxyURL(uri),
		}
	}
	return &TronApiClient{
		httpClient: &cli,
		tronConfig: config,
	}
}

func (c *TronApiClient) GetTransactions(address string, fingerPrint string) ([]Transaction, string, error) {
	uri := fmt.Sprintf("%s/v1/accounts/%s/transactions/trc20", c.tronConfig.ApiHost, address)
	query := make(url.Values)
	query.Add("only_confirmed", "1")
	if fingerPrint != "" {
		query.Add("fingerprint", fingerPrint)
	}
	query.Add("limit", "100")
	query.Add("contract_address", c.tronConfig.ContractAddress)
	uri = uri + "?" + query.Encode()
	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil, "", err
	}
	req.Header.Add("TRON-PRO-API-KEY", c.tronConfig.APiKey)
	req.Header.Add("Content-Type", "application/json")

	xlogger.Info(context.Background(), "请求api", xlogger.Any("url", uri))
	rsp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer rsp.Body.Close()
	if rsp.StatusCode != 200 {
		return nil, "", errors.New("server response: " + rsp.Status)
	}

	var res TransactionResponse
	data, _ := ioutil.ReadAll(rsp.Body)
	err = json.Unmarshal(data, &res)
	if !res.Success {
		return nil, "", errors.New("server response:" + res.Error)
	}

	return res.Transactions, res.Meta.FingerPrint, nil
}

type TransactionResponse struct {
	Transactions []Transaction `json:"data"`
	Success      bool          `json:"success"`
	Meta         struct {
		At          int64  `json:"at"`
		PageSize    int    `json:"page_size"`
		FingerPrint string `json:"fingerprint"`
	} `json:"meta"`
	Timestamp time.Time `json:"timestamp"`
	Status    int       `json:"status"`
	Error     string    `json:"error"`
	Path      string    `json:"path"`
}

type Transaction struct {
	TransactionId string `json:"transaction_id"`
	TokenInfo     struct {
		Symbol   string `json:"symbol"`
		Address  string `json:"address"`
		Decimals int    `json:"decimals"`
		Name     string `json:"name"`
	} `json:"token_info"`
	BlockTimestamp int64  `json:"block_timestamp"`
	From           string `json:"from"`
	To             string `json:"to"`
	Type           string `json:"type"`
	Value          string `json:"value"`
}
