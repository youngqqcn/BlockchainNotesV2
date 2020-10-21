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

func NewReleasePayload(from, to crypto.Address, value, seq int64, memo string) *ReleasePayload {
	//return &ReleasePayload{from, to, value}
	return &ReleasePayload{CommanFeild{from, to, value, seq, memo}}
}

type CommanFeild struct {
	FromAddress crypto.Address // 所有者
	ToAddress   crypto.Address // 目的地址
	Value       int64          // 转账金额
	Sequence    int64          // 转账序号
	Memo        string         // 转账备注
}

type ReleasePayload struct {
	CommanFeild
}

func (r ReleasePayload) GetSigner() crypto.Address {
	return r.FromAddress
}

func (r ReleasePayload) GetSignBytes() []byte {
	bz, err := json.Marshal(r)
	if err != nil {
		return []byte{}
	}
	return bz
}

func (r ReleasePayload) GetType() string {
	return "release"
}

func NewTransferPayload(from, to crypto.Address, value, seq int64, memo string) *TransferPayload {
	//return &ReleasePayload{from, to, value}
	return &TransferPayload{CommanFeild{from, to, value, seq, memo}}
}

type TransferPayload struct {
	CommanFeild
}

func (t TransferPayload) GetSigner() crypto.Address {
	return t.FromAddress
}

func (t TransferPayload) GetSignBytes() []byte {
	bz, err := json.Marshal(t)
	if err != nil {
		return []byte{}
	}
	return bz
}

func (t TransferPayload) GetType() string {
	//panic("implement me")
	return "transfer"
}
