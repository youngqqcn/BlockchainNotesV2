package mytokenapp

import (
	"encoding/hex"
	"fmt"
	"github.com/stretchr/testify/require"
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

	if err := wallet.Save("walletfile.dat"); err != nil {
		t.Errorf("save wallet file failed: %v", err)
		return
	}

	t.Logf("save wallet file success")
}

func TestLoadWalletFromFile(t *testing.T) {
	wallet := LoadWalletFromFile("walletfile.dat")
	if wallet == nil {
		t.Error("load wallet file failed!")
		return
	}

	for label, key := range wallet.Keys {
		t.Logf("label: %v, private key: %v, address : %v", label,
			hex.EncodeToString(key.Bytes()), key.PubKey().Address())
	}

	t.Logf("LoadWalletFromFile successed")
}

func TestInitWallet(t *testing.T) {

	//wallet := NewWallet()
	//require.NotNil(t, wallet, "new wallet error " )
	//
	//wallet.GenNewPrivKey("superuser")
	//wallet.Save("superuser.wallet")

	nw := LoadWalletFromFile("../bin/wallet.dat")
	require.NotNil(t, nw, "load wallet from file failed")

	//nw.GenNewPrivKey("yqq")
	//require.Nil(t,  nw.Save("superuser.wallet"), "save wallet error" )

	t.Log(nw.GetAddress("yqq"))

}

func TestWallet_GetAddress(t *testing.T) {

	w := LoadWalletFromFile("../bin/wallet.dat")
	require.NotNil(t, w, "load wallet from file failed")
	yqqAddr := w.GetAddress("yqq")
	tomAddr := w.GetAddress("tom")
	superuser := w.GetAddress("superuser")

	t.Log("yqq: ", yqqAddr)
	t.Log("tom: ", tomAddr)
	t.Log("superuser: ", superuser)

	require.NotEqual(t, yqqAddr, tomAddr)
	require.NotEqual(t, tomAddr, superuser)

}
