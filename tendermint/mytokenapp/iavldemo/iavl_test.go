package iavldemo

import (
	"encoding/hex"
	"fmt"
	"github.com/tendermint/iavl"
	db "github.com/tendermint/tm-db"
	"mytokenapp/mytokenapp"
	"testing"
)

func TestIavl(t *testing.T) {

	ldb, err := db.NewGoLevelDB("account", "./accountdb")
	if err != nil {
		panic(err)
	}

	mt, err := iavl.NewMutableTree(ldb, 128)
	if err != nil {
		panic(err)
	}
	ver, err := mt.Load()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v\n", ver)

	// set balance
	balance := 1000
	bz, err := mytokenapp.MarshalBinaryBare(balance)
	if err != nil {
		panic(err)
	}

	if len(bz) > 0 {
	}

	mt.Set([]byte("name"), []byte("hello"))
	mt.Set([]byte("age"), bz)

	// get balance
	idx, value := mt.Get([]byte("age"))
	fmt.Printf("index : %v\n", idx)

	var b int64
	mytokenapp.UnMarshalBinaryBare(value, &b)

	fmt.Printf("value: %v\n", b)
	hash, ver, err := mt.SaveVersion()
	if err != nil {
		panic(err)
	}
	fmt.Printf("hash: %v\n", hex.EncodeToString(hash))
	fmt.Printf("ver: %v\n", ver)

}
