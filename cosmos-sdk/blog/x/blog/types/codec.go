package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// RegisterCodec registers concrete types on codec
func RegisterCodec(cdc *codec.Codec) {
	// this line is used by starport scaffolding # 1
		cdc.RegisterConcrete(MsgCreateComment{}, "blog/CreateComment", nil)
		cdc.RegisterConcrete(MsgSetComment{}, "blog/SetComment", nil)
		cdc.RegisterConcrete(MsgDeleteComment{}, "blog/DeleteComment", nil)
	cdc.RegisterConcrete(MsgCreatePost{}, "blog/CreatePost", nil)
}

// ModuleCdc defines the module codec
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
