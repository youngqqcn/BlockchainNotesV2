package mytokenapp

import (
	"bytes"
	"errors"
	"fmt"
	abcitypes "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/kv"
	"math/big"
	"strconv"
	"strings"
)


var _ abcitypes.Application = (*MyTokenApp)(nil)


type MyTokenApp struct {
	//types.BaseApplication   // 组合? 继承?
	//types.Application
	Accounts map[string]int64  // 暂时不做持久化
}


func (app *MyTokenApp) Info(info abcitypes.RequestInfo) abcitypes.ResponseInfo {
	//panic("implement me")
	return abcitypes.ResponseInfo{Version: "v1.0.0yqq"}
}

func (app *MyTokenApp) SetOption(option abcitypes.RequestSetOption) abcitypes.ResponseSetOption {
	//panic("implement me")
	return abcitypes.ResponseSetOption{}
}

func (app *MyTokenApp) Query(query abcitypes.RequestQuery) (resQuery abcitypes.ResponseQuery) {
	//panic("implement me")
	resQuery.Key = query.Data

	if !app.isValidAddress(string( query.Data )) {
		resQuery.Code = 1111
		resQuery.Log = "invalid address"
		return
	}

	balance  := app.Accounts[string(query.Data)]

	resQuery.Code = 0
	resQuery.Value = []byte(strconv.FormatInt(balance, 10))
	resQuery.Log = "query succeed"

	return
}

func (app *MyTokenApp) CheckTx(tx abcitypes.RequestCheckTx) abcitypes.ResponseCheckTx {
	//panic("implement me")
	code := app.isValidTx(tx.Tx)
	return abcitypes.ResponseCheckTx{Code: code}
}

func (app *MyTokenApp) InitChain(chain abcitypes.RequestInitChain) abcitypes.ResponseInitChain {
	//panic("implement me")
	return abcitypes.ResponseInitChain{}
}

func (app *MyTokenApp) BeginBlock(block abcitypes.RequestBeginBlock) abcitypes.ResponseBeginBlock {
	//panic("implement me")


	return abcitypes.ResponseBeginBlock{}
}

func (app *MyTokenApp) DeliverTx(tx abcitypes.RequestDeliverTx) abcitypes.ResponseDeliverTx {

	//panic("implement me")

	//如果是发行
	if strings.Contains( string(tx.Tx) ,  "release") {

        tmpStr  := strings.Replace(string(tx.Tx), "release", "", 1 )

        tmpStr1 := tmpStr
		// parts := bytes.Split( []byte(tmpStr), []byte(","))
		parts := strings.FieldsFunc( tmpStr1, func (c rune) bool {
            if c == ',' {
                return true
            }
            return false
        })

		ownerAddress, toAddress, value := parts[0], parts[1], parts[2]

		if !( app.isValidAddress(string(ownerAddress)) && app.isValidAddress(string(toAddress)) ) {
			fmt.Println("invalid address")
			return abcitypes.ResponseDeliverTx{Code: 1, Log: "invalid release address"}
		}

		//var bgValue big.Int
		//bgValue.SetString(string(value), 10)
		bgValue, _ := strconv.ParseInt(string(value), 10, 64 )

		if _ , err := app.release(ownerAddress, toAddress , bgValue) ; err != nil{
			return abcitypes.ResponseDeliverTx{Code: 2, Log: err.Error()}
		}
		return abcitypes.ResponseDeliverTx{Code: 0, Log: "release succeed"}
	}


	code := app.isValidTx( tx.Tx )
	if code != 0 {
		return abcitypes.ResponseDeliverTx{Code: code }
	}

    // parts := bytes.Split( tx.Tx, []byte(","))
    parts := strings.FieldsFunc( string(tx.Tx) , func (c rune) bool {
            if c == ',' {
                return true
            }
            return false
        })

	fromAddress, toAddress, value :=  parts[0] ,  parts[1] , parts[2]

	//var bigValue big.Int
	//bigValue.SetString(string(value), 10)
	bgValue, _ := strconv.ParseInt(string(value), 10, 64)

	if _, err :=  app.transfer( fromAddress, toAddress, bgValue) ; err != nil {
		return abcitypes.ResponseDeliverTx{Code: 2, Log: "transfer error"}
	}

	events := []abcitypes.Event {
		{
			Type: "transfer",
			Attributes: []kv.Pair {
				{Key: []byte("from"), Value: []byte(fromAddress)},
				{Key: []byte("to"), Value: []byte(toAddress)},
				{Key: []byte("value"), Value: []byte(value)},
			},
		},
	}

	return abcitypes.ResponseDeliverTx{Code: 0, Events: events, Log: "ok"}
}

func (app *MyTokenApp) EndBlock(block abcitypes.RequestEndBlock) abcitypes.ResponseEndBlock {
	//panic("implement me")
	return abcitypes.ResponseEndBlock{}
}

func (app *MyTokenApp) Commit() abcitypes.ResponseCommit {
	//panic("implement me")
	return abcitypes.ResponseCommit{}
}



var SUPER_USER = []byte("hello")

func (app *MyTokenApp)release( owner, receiver string, value int64 ) (bool, error)  {

	if !bytes.Equal([]byte(owner), SUPER_USER )  {
		return false, errors.New("permmition deny")
	}

	// 检查 receiver 是否合法

	//a :=  app.Accounts[receiver.String()]
	//a.SetString("66", 10)
	app.Accounts[receiver] = value

	return true, nil
}


func (app *MyTokenApp)transfer(fromAddress, toAddress string, value int64) (bool, error ) {

	balance := app.Accounts[fromAddress]
	if balance < 0 {
		return false, errors.New("balance is not enough")
	}
	app.Accounts[fromAddress] =  balance - value   //*(balance.Sub(&balance, &value))

	tobalance := app.Accounts[toAddress]
	app.Accounts[toAddress] =  tobalance + value

	return true , nil
}


func (app *MyTokenApp)isValidAddress( address string ) bool {
	return true
}


func (app *MyTokenApp) isValidTx(tx []byte) (code uint32) {

	// check format
	parts := bytes.Split(tx, []byte(",")) // 交易必须包含 >， 类似  yqq:tom:100
	if len(parts) != 3 {
		fmt.Println("invalid tx")
		return 1
	}

	fromAddress :=  string( parts[0] )
	toAddress := string( parts[1] )

	var bv big.Int
	_, ok := bv.SetString( string(parts[2]), 10)
	if !ok {
		fmt.Println("invalid value")
		return 2
	}


	if !( app.isValidAddress(fromAddress) && app.isValidAddress(toAddress) ) {
		fmt.Println("invalid address")
		return 3
	}

	return code
}
