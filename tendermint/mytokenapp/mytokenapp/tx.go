package mytokenapp

import (
	"fmt"
	"github.com/tendermint/tendermint/crypto"
)

type Tx struct {
	Payload   Payload       // 交易内容
	Signature []byte        // 签名
	PubKey    crypto.PubKey // 公钥
}

func NewTx(payload Payload) *Tx {
	return &Tx{Payload: payload}
}

func (tx *Tx) Sign(privKey crypto.PrivKey) error {

	data := tx.Payload.GetSignBytes()

	sig, err := privKey.Sign(data)
	if err != nil {
		return err
	}

	tx.Signature = sig
	tx.PubKey = privKey.PubKey()

	return nil
}

func (tx *Tx) Verify() error {
	msg := tx.Payload.GetSignBytes()
	sig := tx.Signature
	if !tx.PubKey.VerifyBytes(msg, sig) {
		return fmt.Errorf("verify bytes failed . msg : %v , sig: %v", msg, sig)
	}

	return nil
}
