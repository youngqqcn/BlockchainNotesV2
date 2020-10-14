package mytokenapp

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewTx(t *testing.T) {

	wallet := NewWallet()
	wallet.GenNewPrivKey("yqq")
	wallet.GenNewPrivKey("tmp")

	addr1 := wallet.GetAddress("yqq")
	addr2 := wallet.GetAddress("tmp")

	txRelease := NewTx(  NewReleasePayload(addr1, addr2, 10, 10, "hello") )
	if txRelease == nil {
		t.Errorf("  NewTx(  NewReleasePayload  failed ")
		return
	}


	txTransfer := NewTx( NewTransferPayload(addr1, addr2, 10, 10, "hello"))
	if txTransfer == nil {
		t.Errorf("  NewTx( NewTransferPayload ")
		return
	}

	t.Log("NewTx test successed")
}

func TestTx_Sign(t *testing.T) {

	wallet := NewWallet()
	wallet.GenNewPrivKey("yqq")
	wallet.GenNewPrivKey("tmp")

	addr1 := wallet.GetAddress("yqq")
	addr2 := wallet.GetAddress("tmp")

	txRelease := NewTx(  NewReleasePayload(addr1, addr2, 10, 10, "hello") )
	if err := txRelease.Sign(  wallet.GetPrivKeyByLabel("yqq") ) ; err != nil {
		t.Errorf("txRelease.Sign failed")
		return
	}

	txTransfer := NewTx( NewTransferPayload(addr1, addr2, 10, 10, "hello"))
	if err := txTransfer.Sign(  wallet.GetPrivKeyByLabel("yqq") ) ; err != nil {
		t.Errorf("txTransfer.Sign failed")
		return
	}

	t.Logf("sign successed")

}

func TestTx_Verify(t *testing.T) {
	wallet := NewWallet()
	wallet.GenNewPrivKey("yqq")
	wallet.GenNewPrivKey("tmp")

	addr1 := wallet.GetAddress("yqq")
	addr2 := wallet.GetAddress("tmp")

	txRelease := NewTx(  NewReleasePayload(addr1, addr2, 10, 10, "hello") )
	if err := txRelease.Sign(  wallet.GetPrivKeyByLabel("yqq") ) ; err != nil {
		t.Errorf("txRelease.Sign failed")
		return
	}

	txTransfer := NewTx( NewTransferPayload(addr1, addr2, 10, 10, "hello"))
	if err := txTransfer.Sign(  wallet.GetPrivKeyByLabel("yqq") ) ; err != nil {
		t.Errorf("txTransfer.Sign failed")
		return
	}

	if err :=  txRelease.Verify(); err != nil {
		t.Errorf("txRelease verified failed")
		return
	}


	err := txTransfer.Verify()
	require.Nil(t, err, "txTransfer verifyed failed : %v", err)


	t.Logf("tx verify test successed")

}

func TestMarshalBinaryBare(t *testing.T) {

	wallet := NewWallet()
	wallet.GenNewPrivKey("yqq")
	wallet.GenNewPrivKey("tmp")

	addr1 := wallet.GetAddress("yqq")
	addr2 := wallet.GetAddress("tmp")

	txRelease := NewTx(  NewReleasePayload(addr1, addr2, 10, 10, "hello") )
	if err := txRelease.Sign(  wallet.GetPrivKeyByLabel("yqq") ) ; err != nil {
		t.Errorf("txRelease.Sign failed")
		return
	}

	txTransfer := NewTx( NewTransferPayload(addr1, addr2, 10, 10, "hello"))






	if err := txTransfer.Sign(  wallet.GetPrivKeyByLabel("yqq") ) ; err != nil {
		t.Errorf("txTransfer.Sign failed")
		return
	}

	if err :=  txRelease.Verify(); err != nil {
		t.Errorf("txRelease verified failed")
		return
	}


	err := txTransfer.Verify()
	require.Nil(t, err, "txTransfer verifyed failed : %v", err)


	t.Logf("tx verify test successed")

}



