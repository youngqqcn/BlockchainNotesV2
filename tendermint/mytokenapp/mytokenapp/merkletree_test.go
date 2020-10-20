package mytokenapp

import (
	"testing"
)

func Test_balanceMapToByteMap(t *testing.T) {

	//app := MyTokenApp{map[string]int64{
	//	"yqq":  10000,
	//	"tom":  20000,
	//	"hans": 30000,
	//}}
	//
	//bm := app.balanceMapToByteMap()
	//for k, v := range bm {
	//	fmt.Printf("%s => %v\n", k, v)
	//}
	//
	//root := app.getRootHash()
	//
	//const TestAddr = "yqq"
	//
	//fmt.Printf("root hash: %v\n", hex.EncodeToString(root))
	////bzProofBytes := app.getProofBytes(TestAddr)
	//r, p, _ := merkle.SimpleProofsFromMap(bm)
	//require.Equal(t, r, root, "merkle root not matched")
	//
	//leaf := merkle.KVPair{Key: []byte("yqq"), Value: tmhash.Sum(bm["yqq"])}
	//err := p[TestAddr].Verify(root, leaf.Bytes())
	//require.NoError(t, err, "proof error : %v", err)
	//
	//bzProofBytes := app.getProofBytes("yqq")
	//
	//var sp merkle.SimpleProof
	//err = codec.UnmarshalBinaryBare(bzProofBytes, &sp)
	//require.NoError(t, err, "error : %v", err)
	//
	//require.Equal(t, sp.String(), app.getProof("yqq").String(), "err : proof not equal")
	//
	//err = sp.Verify(r, leaf.Bytes())
	//require.NoError(t, err, "error : %v", err)
}
