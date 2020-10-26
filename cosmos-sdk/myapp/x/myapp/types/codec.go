package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// RegisterCodec registers concrete types on codec
func RegisterCodec(cdc *codec.Codec) {
  // this line is used by starport scaffolding # 1
		cdc.RegisterConcrete(MsgCreateStudent{}, "myapp/CreateStudent", nil)
		cdc.RegisterConcrete(MsgSetStudent{}, "myapp/SetStudent", nil)
		cdc.RegisterConcrete(MsgDeleteStudent{}, "myapp/DeleteStudent", nil)
		cdc.RegisterConcrete(MsgCreatePost{}, "myapp/CreatePost", nil)
		cdc.RegisterConcrete(MsgSetPost{}, "myapp/SetPost", nil)
		cdc.RegisterConcrete(MsgDeletePost{}, "myapp/DeletePost", nil)
}

// ModuleCdc defines the module codec
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
