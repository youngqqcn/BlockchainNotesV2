package main

import (
	"bytes"
	abcitypes "github.com/tendermint/tendermint/abci/types"
	"github.com/dgraph-io/badger"  //使用 BadgerDB（类似RocksDB）
	"github.com/tendermint/tendermint/libs/kv"
)

type KVStoreApplication struct {
	db           *badger.DB
	currentBatch *badger.Txn
}

var _ abcitypes.Application = (*KVStoreApplication)(nil)

func NewKVStoreApplication(db *badger.DB) *KVStoreApplication {
	return &KVStoreApplication{
		db: db,
	}
}

func (KVStoreApplication) Info(req abcitypes.RequestInfo) abcitypes.ResponseInfo {
	return abcitypes.ResponseInfo{}
}

func (KVStoreApplication) SetOption(req abcitypes.RequestSetOption) abcitypes.ResponseSetOption {
	return abcitypes.ResponseSetOption{}
}



func (app *KVStoreApplication) isValid(tx []byte) (code uint32) {

	// check format
	parts := bytes.Split(tx, []byte("=")) // 交易必须包含 =， 类似  yqq=100
	if len(parts) != 2 {
		return 1
	}

	key, value := parts[0], parts[1]

	// check if the same key=value already exists
	err := app.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)

		// 其他错误
		if err != nil && err != badger.ErrKeyNotFound {
			return err
		}

		// 如果键值已经存在
		if err == nil {
			return item.Value(func(val []byte) error {
				if bytes.Equal(val, value) {
					code = 2
				}
				return nil
			})
		}

		// 如果 err是 badger.ErrKeyNotFound  ， key不存在， 则正常
		return nil
	})
	if err != nil {
		panic(err)  // 其他错误
	}

	return code
}


func (app *KVStoreApplication) DeliverTx(req abcitypes.RequestDeliverTx) abcitypes.ResponseDeliverTx {
	code := app.isValid(req.Tx)
	if code != 0 {
		return abcitypes.ResponseDeliverTx{Code: code}
	}

	parts := bytes.Split(req.Tx, []byte("="))
	key, value := parts[0], parts[1]

	err := app.currentBatch.Set(key, value)  //添加 kv
	if err != nil {
		panic(err)
	}


	//增加 event
	events := []abcitypes.Event{
		{
			Type: "transfer",
			Attributes: []kv.Pair{
				{Key: []byte("sender"), Value: []byte("Bob")},
				{Key: []byte("recipient"), Value: []byte("Alice")},
				{Key: []byte("balance"), Value: []byte("100")},
				{Key: []byte("note"), Value: []byte("nothing")},
			},
		},
	}

	return abcitypes.ResponseDeliverTx{Code:0, Events: events}
}


func (app *KVStoreApplication) Commit() abcitypes.ResponseCommit {
	app.currentBatch.Commit() //提交事务
	return abcitypes.ResponseCommit{Data: []byte{}}
}

// 调用 Tendermint Core RPC 的 RPC接口， /abci_query , 会调用应用程序实现的 Query 方法
func (app *KVStoreApplication) Query(reqQuery abcitypes.RequestQuery) (resQuery abcitypes.ResponseQuery) {
	resQuery.Key = reqQuery.Data
	err := app.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(reqQuery.Data)
		if err != nil && err != badger.ErrKeyNotFound {
			return err
		}
		if err == badger.ErrKeyNotFound {
			resQuery.Log = "does not exist"
		} else {
			return item.Value(func(val []byte) error {
				resQuery.Log = "exists"
				resQuery.Value = val
				return nil
			})
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return
}

func (KVStoreApplication) InitChain(req abcitypes.RequestInitChain) abcitypes.ResponseInitChain {
	return abcitypes.ResponseInitChain{}
}

// 一个区块的产生分为3个部分 ： BeginBlock , DeliverTx,  EndBlock
// DeliverTx 是异步的
func (app *KVStoreApplication) BeginBlock(req abcitypes.RequestBeginBlock) abcitypes.ResponseBeginBlock {
	app.currentBatch = app.db.NewTransaction(true) // 开启事务
	return abcitypes.ResponseBeginBlock{}
}

func (KVStoreApplication) EndBlock(req abcitypes.RequestEndBlock) abcitypes.ResponseEndBlock {
	return abcitypes.ResponseEndBlock{}
}


//
//
//func (KVStoreApplication) ListSnapshots(abcitypes.RequestListSnapshots) abcitypes.ResponseListSnapshots {
//	return abcitypes.ResponseListSnapshots{}
//}
//
//func (KVStoreApplication) OfferSnapshot(abcitypes.RequestOfferSnapshot) abcitypes.ResponseOfferSnapshot {
//	return abcitypes.ResponseOfferSnapshot{}
//}
//
//func (KVStoreApplication) LoadSnapshotChunk(abcitypes.RequestLoadSnapshotChunk) abcitypes.ResponseLoadSnapshotChunk {
//	return abcitypes.ResponseLoadSnapshotChunk{}
//}
//
//func (KVStoreApplication) ApplySnapshotChunk(abcitypes.RequestApplySnapshotChunk) abcitypes.ResponseApplySnapshotChunk {
//	return abcitypes.ResponseApplySnapshotChunk{}
//}

func (app *KVStoreApplication) CheckTx(req abcitypes.RequestCheckTx) abcitypes.ResponseCheckTx {
	code := app.isValid(req.Tx)
	return abcitypes.ResponseCheckTx{Code: code, GasWanted: 1}
}