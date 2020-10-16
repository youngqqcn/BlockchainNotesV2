package merkletree

import (
	"encoding/hex"
	"fmt"
	"github.com/tendermint/tendermint/crypto/merkle"
	"github.com/tendermint/tendermint/crypto/tmhash"
	"testing"
)

type Data struct {
	value string
}

func (data *Data) Hash() []byte {
	return tmhash.Sum([]byte(data.value))
}

func TestMerkleHash(t *testing.T) {

	//data := []merkle.KVPair {
	//merkle.KVPair{[]byte("yqq"), []byte("1000", )},

	//}

	//data := []merkle.KVPair{
	//	merkle.KVPair{}
	//}

	//newSimpleMap
	//merkle.Key{}
	//merkle.KeyPath{}
	//merkle.SimpleHashFromByteSlices()

	m := map[string][]byte{
		"yqq":   []byte("1000"),
		"tom":   []byte("2000"),
		"alice": []byte("3000"),
		"jack":  []byte("4000"),
	}

	h := merkle.SimpleHashFromMap(m)
	fmt.Printf("hash of map : %v\n", hex.EncodeToString(h))

	root, proofs, keys := merkle.SimpleProofsFromMap(m)
	fmt.Printf("merkle root: %v \n", hex.EncodeToString(root))
	fmt.Printf("proof: %v\n", proofs)
	fmt.Printf("key : %v\n", keys)

	fmt.Printf("yqq aunts : %v\n", proofs["yqq"].Aunts)
	fmt.Printf("tom aunts : %v\n", proofs["tom"].Aunts)
	fmt.Printf("alice aunts : %v\n", proofs["alice"].Aunts)
	fmt.Printf("jack aunts : %v\n", proofs["jack"].Aunts)

	// 需要对 value 进行hash, 然后proof
	hleaf := merkle.KVPair{Key: []byte("yqq"), Value: tmhash.Sum(m["yqq"])}
	if err := proofs["yqq"].Verify(root, hleaf.Bytes()); err != nil {
		//fmt.Println("verify failed")
		fmt.Println(err)
	} else {
		fmt.Println("verify successed")
	}

	fmt.Println("=====================================")

	bzs := [][]byte{
		[]byte("hello"),
		[]byte("goood"),
		[]byte("apple"),
		[]byte("pine"),
	}

	hs := merkle.SimpleHashFromByteSlices(bzs) // 递归实现
	fmt.Printf(" h : %v \n", hex.EncodeToString(hs))

	r, p := merkle.SimpleProofsFromByteSlices(bzs)
	fmt.Printf("merkle root: %v \n", r)
	fmt.Printf("proof: %v\n", p)

	if err := p[0].Verify(r, bzs[0]); err != nil {
		fmt.Println("verify failed", err)
	} else {
		fmt.Println("verify successed")
	}

	ihs := merkle.SimpleHashFromByteSlicesIterative(bzs) // 迭代实现
	fmt.Printf(" h : %v \n", hex.EncodeToString(ihs))

}
