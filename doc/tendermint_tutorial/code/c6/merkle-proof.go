package main

import (
	"fmt"
	"github.com/tendermint/tendermint/crypto/tmhash"
	"github.com/tendermint/tendermint/crypto/merkle"
)

type sh struct{
	value string
}

func (h sh) Hash() []byte {
	return tmhash.Sum([]byte(h.value))
}

func sliceDemo(){
	data := []merkle.Hasher{ &sh{"one"},&sh{"two"},&sh{"three"},&sh{"four"} }
	root,proofs := merkle.SimpleProofsFromHashers(data)
	fmt.Printf("root hash => %x\n",root)
	fmt.Printf("proof for one => %+v\n",proofs[0])
	valid := proofs[0].Verify(0,4,data[0].Hash(),root)
	fmt.Printf("data[0] is valid? => %t\n",valid)
}

func mapDemo(){
	data := map[string]merkle.Hasher{
		"tom": &sh{"actor"},
		"mary":&sh{"teacher"},
		"linda":&sh{"scientist"},
		"luke":&sh{"fisher"}}
	root,proofs,keys := merkle.SimpleProofsFromMap(data)
	fmt.Printf("root hash => %x\n",root)
	fmt.Printf("proof for tom => %+v\n",proofs["tom"])
	fmt.Printf("keys sorted => %v\n",keys)
	kvpair := merkle.KVPair{Key:[]byte("tom"),Value: data["tom"].Hash()}
	valid := proofs["tom"].Verify(3,4,kvpair.Hash(),root)
	fmt.Printf("data[\"tom\"] is valid? => %t\n",valid)
}

func main(){
	sliceDemo()
	mapDemo()
}
