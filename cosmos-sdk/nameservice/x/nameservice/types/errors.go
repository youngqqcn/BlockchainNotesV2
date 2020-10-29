package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// 这里可以自定义一些错误

var (
	ErrInvalid = sdkerrors.Register(ModuleName, 1, "custom error message")
)
