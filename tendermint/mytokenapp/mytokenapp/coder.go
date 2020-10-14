package mytokenapp

import (
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

var codec = amino.NewCodec()


func init() {
	codec.RegisterInterface((*crypto.PrivKey)(nil), nil)
	codec.RegisterInterface((*crypto.PubKey)(nil), nil)
	codec.RegisterConcrete(&secp256k1.PrivKeySecp256k1{}, "secp256k1/privkey", nil)
	codec.RegisterConcrete(&secp256k1.PubKeySecp256k1{}, "secp256k1/pubkey", nil)
}


func MarshalBinaryBare(o interface{}) ([]byte , error) {
	return codec.MarshalBinaryBare(o)
}

func UnMarshalBinaryBare(bz []byte, ptr interface{}) error{
	return codec.UnmarshalBinaryBare(bz, ptr )
}

func MarshalJSON(o interface{}) ([]byte , error) {
	return codec.MarshalJSON(o)
}

func UnMarshalJSON(bz []byte, ptr interface{}) error {
	return codec.UnmarshalJSON(bz, ptr)
}
