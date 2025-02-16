package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	xlogger "github.com/clearcodecn/log"

	"usdtpay/config"
	event2 "usdtpay/domain/event"
	monitor "usdtpay/domain/tron"
	address2 "usdtpay/infr/mysql/address"
	"usdtpay/infr/mysql/dao"
	"usdtpay/infr/mysql/transaction"
)

type TransactionService struct {
	tronApiClient *monitor.TronApiClient
}

func NewTransactionService() *TransactionService {
	var (
		conf = config.Setting.Tron
	)
	cli := monitor.NewTronApiClient(
		&monitor.TronConfig{
			ContractAddress: conf.UsdtContractAddress,
			ApiHost:         conf.ApiHost,
			APiKey:          conf.ApiKey,
			Timeout:         time.Duration(conf.Timeout) * time.Second,
			Proxy:           conf.Proxy,
		},
	)

	return &TransactionService{
		tronApiClient: cli,
	}
}

func (s *TransactionService) RunOnce() {
	fmt.Println("service start-------")
	ctx := context.Background()
	addressList, err := address2.GetAddressList(config.Setting.MysqlClient)
	if err != nil {
		xlogger.Error(ctx, "GetAddressList failed", xlogger.Err(err))
		return
	}
	for _, address := range addressList {
		var fingerPrint string
		lastTx := transaction.GetLastAddressTransactionId(config.Setting.MysqlClient, address.Address)
		if lastTx != nil {
			fingerPrint = lastTx.FingerPrint
		}
		txs, fingerPrint, err := s.tronApiClient.GetTransactions(address.Address, fingerPrint, 100)
		if err != nil {
			xlogger.Error(ctx, "RunOnce-GetTransactions failed", xlogger.Err(err))
			continue
		}
		if len(txs) == 0 {
			continue
		}
		var (
			ats            []*dao.AddressTransaction
			transactionIds []string
			transactionMap = make(map[string]*dao.AddressTransaction)
		)
		for _, tx := range txs {
			at := convertTonTransactionToModel(tx, fingerPrint)
			transactionMap[tx.TransactionId] = at
			transactionIds = append(transactionIds, at.TransactionId)
		}
		var exists []*dao.AddressTransaction
		err = config.Setting.MysqlClient.Where("transaction_id in ?", transactionIds).Find(&exists).Error
		if err != nil {
			continue
		}
		if len(exists) > 0 {
			for _, e := range exists {
				delete(transactionMap, e.TransactionId)
			}
		}
		if len(transactionMap) > 0 {
			for _, tx := range transactionMap {
				ats = append(ats, tx)
			}
		}

		if len(ats) > 0 {
			config.Setting.MysqlClient.CreateInBatches(ats, len(ats))
			xlogger.Info(ctx, "监控到交易数据: "+address.Address, xlogger.Any("data", len(ats)))
			s.publish(ats)
		}
	}
}

func (s *TransactionService) publish(data []*dao.AddressTransaction) {
	for _, d := range data {
		event := &event2.BlockChainEvent{
			TransactionId: d.TransactionId,
		}
		eventBytes, _ := json.Marshal(event)
		topic := config.Setting.NSQ.BlockChainTopic
		if err := config.Setting.NsqProducer.Publish(topic, eventBytes); err != nil {
			xlogger.Error(context.Background(), "PublishMessage Failed: BlockChainTopic", xlogger.Err(err))
		}
	}
}

func convertTonTransactionToModel(tx monitor.Transaction, fingerPrint string) *dao.AddressTransaction {
	v, _ := strconv.ParseInt(tx.Value, 10, 64)
	return &dao.AddressTransaction{
		Address:       tx.To,
		TransactionId: tx.TransactionId,
		BlockTime:     tx.BlockTimestamp,
		FromAddress:   tx.From,
		Type:          tx.Type,
		Token:         tx.TokenInfo.Symbol,
		Value:         v,
		Decimal:       tx.TokenInfo.Decimals,
		FingerPrint:   fingerPrint,
		CreateTime:    time.Now(),
	}
}
