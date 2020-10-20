package mytokenapp

import "github.com/tendermint/tendermint/crypto/merkle"

func (app *MyTokenApp) balanceMapToByteMap() map[string][]byte {

	////app.Accounts
	//newMap := make(map[string][]byte, len(app.Accounts))
	//for addr, balance := range app.Accounts {
	//	balanceBytes, err := codec.MarshalBinaryBare(balance)
	//	if err != nil {
	//		panic(err)
	//	}
	//	newMap[addr] = balanceBytes
	//}
	//
	//return newMap
	return map[string][]byte{}
}

func (app *MyTokenApp) getRootHash() []byte {
	bytesMap := app.balanceMapToByteMap()
	return merkle.SimpleHashFromMap(bytesMap)
}

func (app *MyTokenApp) getProofBytes(address string) []byte {
	bytesMap := app.balanceMapToByteMap()

	_, proofs, _ := merkle.SimpleProofsFromMap(bytesMap)

	bz, err := codec.MarshalBinaryBare(proofs[address])
	if err != nil {
		panic(err)
	}

	return bz
}

func (app *MyTokenApp) getProof(address string) *merkle.SimpleProof {
	bytesMap := app.balanceMapToByteMap()
	_, proofs, _ := merkle.SimpleProofsFromMap(bytesMap)
	return proofs[address]
}

func (app *MyTokenApp) getProofOp(address string) merkle.SimpleValueOp {
	bytesMap := app.balanceMapToByteMap()
	_, proofs, _ := merkle.SimpleProofsFromMap(bytesMap)
	return merkle.NewSimpleValueOp([]byte(address), proofs[address])
}
