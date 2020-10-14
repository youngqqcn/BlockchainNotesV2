package mytokenapp

import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestNewWallet(t *testing.T) {
	wallet := NewWallet()
	if wallet == nil {
		t.Error("NewWallet error")
	}

	t.Log("new NewWallet success")
}


func TestWallet_Save(t *testing.T) {

	//t.Error("hellllllllllllllllllll")
	wallet := NewWallet()
	if wallet == nil {
		t.Error("sdfdsfsd")
		return
	}

	rand.Seed(time.Now().UnixNano())
	for n := 0; n < 10; n++ {
		label := fmt.Sprintf("user%d", n)
		privkey := wallet.GenNewPrivKey(label)
		if privkey == nil {
			t.Error("gennew private key error")
			return
		}

		t.Logf("private key: %v", privkey)
		t.Logf("public key: %v", wallet.GetPubKey(label))
		t.Logf("address : %v", wallet.GetAddress(label))
	}

	if err := wallet.Save("walletfile.dat") ; err != nil {
		t.Errorf("save wallet file failed: %v", err)
		return
	}

	t.Logf("save wallet file success")
}



func TestLoadWalletFromFile(t *testing.T) {
	wallet  := LoadWalletFromFile("walletfile.dat")
	if wallet == nil {
		t.Error("load wallet file failed!")
		return
	}

	for label, key  := range wallet.Keys {
		t.Logf("label: %v, private key: %v, address : %v", label ,
				hex.EncodeToString( key.Bytes()), key.PubKey().Address() )
	}

	t.Logf("LoadWalletFromFile successed")
}



