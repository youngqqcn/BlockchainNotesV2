package main

import (
	"fmt"
	"github.com/tendermint/tendermint/crypto/merkle"
	"github.com/tendermint/tendermint/crypto/tmhash"
)

type sh struct{
	value string
}

func (h sh) Hash() []byte {
	return tmhash.Sum([]byte(h.value))
}

func sliceDemo(){
	data := []merkle.Hasher{ &sh{"one"},&sh{"two"},&sh{"three"},&sh{"four"} }
	hash := merkle.SimpleHashFromHashers(data)
	fmt.Printf("root hash => %x\n",hash)
}

func mapDemo(){
	data := map[string]merkle.Hasher{
		"tom": &sh{"actor"},
		"mary":&sh{"teacher"},
		"linda":&sh{"scientist"},
		"luke":&sh{"fisher"}}
	hash := merkle.SimpleHashFromMap(data)
	fmt.Printf("root hash => %x\n",hash)
}

func main(){
	sliceDemo()
	mapDemo()
}
