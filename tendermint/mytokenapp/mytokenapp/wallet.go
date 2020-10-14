package mytokenapp

import (
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"io/ioutil"
)


type Wallet struct {
	Keys map[string]crypto.PrivKey
}

func (wallet *Wallet) GenNewPrivKey( label string ) crypto.PrivKey {
	privkey := secp256k1.GenPrivKey()
	wallet.Keys[label] = privkey
	return privkey
}

func (wallet *Wallet)GetPrivKeyByLabel(label string)  crypto.PrivKey {
	return wallet.Keys[label]
}

func (wallet *Wallet)GetPubKey(label string) crypto.PubKey {
	return wallet.Keys[label].PubKey()
}

func (wallet *Wallet)GetAddress(label string) crypto.Address {
	return wallet.Keys[label].PubKey().Address()
}

func (wallet *Wallet)Save(filename string)  error {
	bz, err := codec.MarshalJSON( wallet )
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, bz, 0666)
}

func NewWallet() *Wallet {
	return &Wallet{ Keys: map[string]crypto.PrivKey{}}
}


func LoadWalletFromFile(filepath string) *Wallet {
	var wallet Wallet
	bz , err := ioutil.ReadFile(filepath)
	if err != nil {
		panic(err)
	}

	err = codec.UnmarshalJSON(bz, &wallet)
	if err != nil {
		panic(err)
	}

	return &wallet
}





