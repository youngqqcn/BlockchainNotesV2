package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Student struct {
	Creator sdk.AccAddress `json:"creator" yaml:"creator"`
	ID      string         `json:"id" yaml:"id"`
    Stuname string `json:"stuname" yaml:"stuname"`
    Age int32 `json:"age" yaml:"age"`
    Gender bool `json:"gender" yaml:"gender"`
    Homeaddr string `json:"homeaddr" yaml:"homeaddr"`
}