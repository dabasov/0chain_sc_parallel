package smartcontractinterface

import (
	"encoding/json"

	c_state "0chain.net/chaincore/chain/state"
	"0chain.net/chaincore/smartcontractstate"
	"0chain.net/chaincore/transaction"
)

const Seperator = ":"

type SmartContract struct {
	DB smartcontractstate.SCStateI
	ID string
}

func NewSC(db smartcontractstate.SCStateI, id string) *SmartContract {
	return &SmartContract{DB: db, ID: id}
}

type SmartContractTransactionData struct {
	FunctionName string          `json:"name"`
	InputData    json.RawMessage `json:"input"`
}

type SmartContractInterface interface {
	Execute(t *transaction.Transaction, funcName string, input []byte, balances c_state.StateContextI) (string, error)
	SetSC(sc *SmartContract)
}