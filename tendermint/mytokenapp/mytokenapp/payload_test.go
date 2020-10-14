package mytokenapp

import (
	"encoding/json"
	"testing"
)

func TestNewReleasePayload(t *testing.T) {

	wallet := NewWallet()
	wallet.GenNewPrivKey("yqq")
	wallet.GenNewPrivKey("tmp")

	fromAddress := wallet.GetAddress("yqq")
	toAddress := wallet.GetAddress("tmp")
	r := NewReleasePayload( fromAddress, toAddress, 1000 )
	if r == nil {
		t.Error("NewReleasePayload is error")
		return
	}

	t.Logf("signer: %v", r.GetSigner())
	t.Logf("signBytes: %v", r.GetSignBytes())
	t.Logf("type: %v", r.GetType())
}


func TestNewTransferPayload(t *testing.T) {

	wallet := NewWallet()
	wallet.GenNewPrivKey("yqq")
	wallet.GenNewPrivKey("tmp")

	fromAddress := wallet.GetAddress("yqq")
	toAddress := wallet.GetAddress("tmp")
	r := NewTransferPayload( fromAddress, toAddress, 1000 )
	if r == nil {
		t.Error("NewReleasePayload is error")
		return
	}

	t.Logf("signer: %v", r.GetSigner())

	signBytes := r.GetSignBytes()
	t.Logf("signBytes: %v", signBytes)

	var tmpTx TransferPayload
	if err := json.Unmarshal(signBytes, &tmpTx ); err != nil {
		t.Errorf("GetSignBytes error : %v", err)
		return
	}

	if r.GetType() != "transfer" {
		t.Errorf("type error: %v", r.GetType())
		return
	}

	t.Logf("type: %v", r.GetType())

}