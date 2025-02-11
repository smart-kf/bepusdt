package monitor

import (
	"testing"
	"time"
)

func TestTron(t *testing.T) {
	// c := NewTronApiClient(
	// 	&TronConfig{
	// 		ContractAddress: "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t",
	// 		ApiHost:         "https://api.trongrid.io",
	// 		Timeout:         3 * time.Second,
	// 		// Proxy:           "http://localhost:10081",
	// 		APiKey: "",
	// 	}, func() []string {
	// 		return nil
	// 	},
	// )
	// c.doRequest("TQfeVbRzf7tUZkaeYdZrpKefFumVPbc3RM")

	c := NewTronApiClient(
		&TronConfig{
			ContractAddress: "TXYZopYRdj2D9XRtbG411XZZ3kM5VkAeBf",
			ApiHost:         "https://nile.trongrid.io",
			Timeout:         3 * time.Second,
			Proxy:           "http://localhost:10081",
			APiKey:          "",
		},
	)

	c.GetTransactions("TLT2gJpbRx4e2fWuwzUwhkn7eHXfAfGocG", "")

}

// https://api.trongrid.io/v1/accounts/TQfeVbRzf7tUZkaeYdZrpKefFumVPbc3RM/transactions/trc20?limit=1&contract_address=TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t
