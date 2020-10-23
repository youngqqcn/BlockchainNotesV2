package lib

import (
	"github.com/tendermint/tendermint/crypto/tmhash"
	"github.com/tendermint/tendermint/crypto/merkle"
)

type Balance int

func (b Balance) Hash() []byte{
	v,_ := codec.MarshalBinaryBare(b)
	return tmhash.Sum(v)
}

func (app *TokenApp) stateToHasherMap() map[string]merkle.Hasher {
	hashers := map[string]merkle.Hasher{}
	for addr,val := range app.Accounts {
		balance := Balance(val)
		hashers[addr] = &balance
	}
	return hashers
}

func (app *TokenApp) getRootHash() []byte {
	hashers := app.stateToHasherMap()
	return merkle.SimpleHashFromMap(hashers)
}

func (app *TokenApp) getProofBytes(addr string) []byte {
	hashers := app.stateToHasherMap()
	_,proofs,_ := merkle.SimpleProofsFromMap(hashers)
	bz,err := codec.MarshalBinaryBare(proofs[addr])
	if err != nil  { return  []byte{} }
	return bz
}
