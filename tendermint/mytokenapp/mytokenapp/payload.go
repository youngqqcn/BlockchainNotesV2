package mytokenapp

import (
	"encoding/json"
	"github.com/tendermint/tendermint/crypto"
)

type Payload interface {
	GetSigner() crypto.Address
	GetSignBytes() []byte
	GetType() string
}

func NewReleasePayload(from, to crypto.Address, value int) *ReleasePayload {
	//return &ReleasePayload{from, to, value}
	return &ReleasePayload{CommanFeild{from, to, value}}
}


type CommanFeild struct {
	FromAddress 	crypto.Address	// 所有者
	ToAddress 		crypto.Address	// 目的地址
	Value 			int   			// 转账金额
	//Sequence 		int				//
}

type ReleasePayload struct {
	CommanFeild
}

func (r *CommanFeild) GetSigner() crypto.Address {
	return r.FromAddress
}

func (r *CommanFeild) GetSignBytes() []byte {
	bz, err := json.Marshal(r)
	if err != nil {
		return []byte{}
	}
	return bz
}

func (r *CommanFeild) GetType() string {
	return "release"
}

func NewTransferPayload(from, to crypto.Address, value int) *TransferPayload {
	return &TransferPayload{ CommanFeild{from, to, value} }
}

type TransferPayload struct {
	CommanFeild
}

func (t *TransferPayload) GetSigner() crypto.Address {
	return t.FromAddress
}

func (t *TransferPayload) GetSignBytes() []byte {
	bz, err := json.Marshal(t)
	if err != nil {
		return  []byte{}
	}
	return bz
}

func (t *TransferPayload) GetType() string {
	return "transfer"
}





