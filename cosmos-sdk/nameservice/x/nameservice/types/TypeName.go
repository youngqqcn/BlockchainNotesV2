package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Name struct {
	Creator sdk.AccAddress `json:"creator" yaml:"creator"`
	ID      string         `json:"id" yaml:"id"`
    Value string `json:"value" yaml:"value"`
    Price string `json:"price" yaml:"price"`
}