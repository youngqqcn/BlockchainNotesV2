package mytokenapp

import (
	"encoding/hex"
	"errors"
	"fmt"
	abcitypes "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/libs/kv"
	"strconv"
	"strings"
)

var _ abcitypes.Application = (*MyTokenApp)(nil)

type MyTokenApp struct {
	//types.BaseApplication   // 组合? 继承?
	//types.Application

	abcitypes.BaseApplication
	//Accounts map[string]int64 // 暂时不做持久化

	store *Store // 基于 iavl
}

func NewMyTokenApp(accDbDirPath string) *MyTokenApp {
	return &MyTokenApp{store: NewStore(accDbDirPath)}
}

func (app *MyTokenApp) Info(info abcitypes.RequestInfo) abcitypes.ResponseInfo {
	return abcitypes.ResponseInfo{Version: "v1.0.0yqq"}
}

func (app *MyTokenApp) SetOption(option abcitypes.RequestSetOption) abcitypes.ResponseSetOption {
	return abcitypes.ResponseSetOption{}
}

func (app *MyTokenApp) Query(query abcitypes.RequestQuery) (resQuery abcitypes.ResponseQuery) {
	resQuery.Key = query.Data

	if !app.isValidAddress(string(query.Data)) {
		resQuery.Code = 1111
		resQuery.Log = "invalid address"
		return
	}

	upStr := strings.ToUpper(string(query.Data))
	//balance := app.Accounts[upStr]
	addr, err := hex.DecodeString(upStr)
	if err != nil {
		resQuery.Code = 1111
		resQuery.Log = "invalid address"
		return
	}

	balance, err := app.store.GetBalance(addr)
	if err != nil {
		resQuery.Code = 2
		resQuery.Log = fmt.Sprintf("getbalance error: %v\n", err)
		return
	}

	resQuery.Code = 0
	resQuery.Value = []byte(strconv.FormatInt(balance, 10))
	resQuery.Log = fmt.Sprintf("%s balance is %d\n", string(query.Data), balance)

	// 客户端如何使用proof进行验证?
	//resQuery.Proof = &merkle.Proof{Ops: []merkle.ProofOp{app.getProofOp(string(query.Data)).ProofOp()}}

	return
}

func (app *MyTokenApp) CheckTx(tx abcitypes.RequestCheckTx) abcitypes.ResponseCheckTx {
	var trx Tx
	if err := codec.UnmarshalBinaryBare(tx.Tx, &trx); err != nil {
		return abcitypes.ResponseCheckTx{Code: 1, Log: "invalid transaction data"}
	}

	if trx.Payload.GetType() == "transfer" {
		txp := trx.Payload.(*TransferPayload)

		// 判断交易的发送方是否是交易签名方
		if txp.FromAddress.String() != trx.PubKey.Address().String() {
			return abcitypes.ResponseCheckTx{Code: 4, Log: fmt.Sprintf("signature is not matched")}
		}
	} else if trx.Payload.GetType() == "release" {
		txp := trx.Payload.(*ReleasePayload)
		// 判断交易的发送方是否是交易签名方
		if txp.FromAddress.String() != trx.PubKey.Address().String() {
			return abcitypes.ResponseCheckTx{Code: 4, Log: fmt.Sprintf("signature is not matched")}
		}
	}

	return abcitypes.ResponseCheckTx{Code: 0}
}

func (app *MyTokenApp) InitChain(chain abcitypes.RequestInitChain) abcitypes.ResponseInitChain {
	return abcitypes.ResponseInitChain{}
}

func (app *MyTokenApp) BeginBlock(block abcitypes.RequestBeginBlock) abcitypes.ResponseBeginBlock {
	return abcitypes.ResponseBeginBlock{}
}

func (app *MyTokenApp) DeliverTx(tx abcitypes.RequestDeliverTx) abcitypes.ResponseDeliverTx {

	var transaction Tx
	if err := codec.UnmarshalBinaryBare(tx.Tx, &transaction); err != nil {
		return abcitypes.ResponseDeliverTx{Code: 1, Log: "UnmarshalBinaryBare failed, invalid tx data"}
	}

	// 验证交易的签名
	if err := transaction.Verify(); err != nil {
		return abcitypes.ResponseDeliverTx{Code: 2, Log: "transaction verified failed"}
	}

	var events []abcitypes.Event

	if transaction.Payload.GetType() == "transfer" {
		txp := transaction.Payload.(*TransferPayload)

		// 判断交易的发送方是否是交易签名方
		if txp.FromAddress.String() != transaction.PubKey.Address().String() {
			return abcitypes.ResponseDeliverTx{Code: 4, Log: fmt.Sprintf("signature is not matched")}
		}

		if ok, err := app.transfer(txp.FromAddress, txp.ToAddress, txp.Value); !ok {
			return abcitypes.ResponseDeliverTx{Code: 3, Log: fmt.Sprintf("error:%v", err)}
		}
		events = []abcitypes.Event{
			{
				Type: "transfer",
				Attributes: []kv.Pair{
					{Key: []byte("from"), Value: []byte(txp.FromAddress)},
					{Key: []byte("to"), Value: []byte(txp.ToAddress)},
					{Key: []byte("value"), Value: []byte(strconv.FormatInt(txp.Value, 10))},
					{Key: []byte("memo"), Value: []byte(txp.Memo)},
				},
			},
		}
	} else if transaction.Payload.GetType() == "release" {
		txp := transaction.Payload.(*ReleasePayload)

		// 判断交易的发送方是否是交易签名方
		if txp.FromAddress.String() != transaction.PubKey.Address().String() {
			return abcitypes.ResponseDeliverTx{Code: 4, Log: fmt.Sprintf("signature is not matched")}
		}

		if ok, err := app.release(txp.FromAddress, txp.ToAddress, txp.Value); !ok {
			return abcitypes.ResponseDeliverTx{Code: 3, Log: fmt.Sprintf("error:%v", err)}
		}
		events = []abcitypes.Event{
			{
				Type: "release",
				Attributes: []kv.Pair{
					{Key: []byte("from"), Value: txp.FromAddress},
					{Key: []byte("to"), Value: txp.ToAddress},
					{Key: []byte("value"), Value: []byte(strconv.FormatInt(txp.Value, 10))},
					{Key: []byte("memo"), Value: []byte(txp.Memo)},
				},
			},
		}
	}

	return abcitypes.ResponseDeliverTx{Code: 0, Log: "operation ok", Events: events}
}

func (app *MyTokenApp) EndBlock(block abcitypes.RequestEndBlock) abcitypes.ResponseEndBlock {
	return abcitypes.ResponseEndBlock{}
}

func (app *MyTokenApp) Commit() abcitypes.ResponseCommit {
	merkleRoot := app.getRootHash() // merkle tree root hash
	app.store.Commit()              // iavl
	return abcitypes.ResponseCommit{Data: merkleRoot}
}

//var SUPER_USER string //= "365EA5222D2F08A8A1EBF992B0628B1459527400"

func (app *MyTokenApp) release(owner, receiver crypto.Address, value int64) (bool, error) {

	wallet := LoadWalletFromFile("wallet.dat")
	if wallet == nil {
		panic("load wallet error")
	}

	if owner.String() != wallet.GetAddress("superuser").String() {
		return false, errors.New("sender is not super user")
	}

	//app.Accounts[receiver] += value

	balance, _ := app.store.GetBalance(receiver)
	err := app.store.SetBalance(receiver, balance+value)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (app *MyTokenApp) transfer(fromAddress, toAddress crypto.Address, value int64) (bool, error) {

	//balance := app.Accounts[fromAddress]
	balance, _ := app.store.GetBalance(fromAddress)
	if balance < value {
		return false, errors.New("balance is not enough")
	}

	//app.Accounts[fromAddress] -= value //*(balance.Sub(&balance, &value))
	err := app.store.SetBalance(fromAddress, balance-value)
	if err != nil {
		return false, err
	}

	toBalance, _ := app.store.GetBalance(toAddress)
	err = app.store.SetBalance(toAddress, toBalance+value)
	if err != nil {
		return false, err
	}

	//app.Accounts[toAddress] += value

	return true, nil
}

func (app *MyTokenApp) isValidAddress(address string) bool {
	return true
}
