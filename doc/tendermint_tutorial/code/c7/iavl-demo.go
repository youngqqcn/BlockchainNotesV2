package main

import (
	"fmt"
	"math/rand"
	"time"
	"github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/iavl"
	"github.com/tendermint/go-amino"
	"./lib"
)

func main(){
	gdb,err := db.NewGoLevelDB("account",".")
	if err != nil { panic(err) }
	tree := iavl.NewMutableTree(gdb,128)
	ver,err := tree.Load()
	if err != nil { panic(err) }
	fmt.Printf("tree version => %v\n",ver)

	codec := amino.NewCodec()

	rand.Seed(time.Now().Unix())

	getBalance("michael",tree,codec)
	getBalance("britney",tree,codec)


	setBalance("michael",rand.Intn(1000),tree,codec)
	setBalance("britney",rand.Intn(1000),tree,codec)

	hash,ver,err := tree.SaveVersion()
	if err != nil { panic(err) }
	fmt.Printf("tree hash => %v\n",hash)
	fmt.Printf("tree version => %v\n",ver)

	getBalanceVersioned("michael",ver-1,tree,codec)
	getBalanceVersioned("britney",ver-1,tree,codec)
}

func setBalance(name string,balance int,tree *iavl.MutableTree,codec *amino.Codec) {
	fmt.Printf("set %v's balance => %v\n",name,balance)
	bz,err := lib.MarshalBinary(balance)
	if err != nil { panic(err) }
	tree.Set([]byte(name), bz)
}

func getBalanceVersioned(name string,version int64,tree *iavl.MutableTree,codec *amino.Codec) int{
	var val int
	_,bz := tree.GetVersioned([]byte(name),version)
	if bz == nil {
		val = 0
	} else {
		err := lib.UnmarshalBinary(bz,&val)
		if err != nil { panic(err) }
	}
	fmt.Printf("%v's balance@%v => %v\n",name,version,val)
	return val
}

func getBalance(name string,tree *iavl.MutableTree,codec *amino.Codec) int{
	var val int
	_,bz := tree.Get([]byte(name))
	if bz == nil {
		val = 0
	} else {
		err := lib.UnmarshalBinary(bz,&val)
		if err != nil { panic(err) }
	}
	fmt.Printf("%v's balance@workspace => %v\n",name,val)
	return val
}
