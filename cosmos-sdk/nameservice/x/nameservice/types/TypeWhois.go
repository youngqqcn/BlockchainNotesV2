package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var MinNamePrice = sdk.Coins{sdk.NewInt64Coin("qtoken", 10)}

type Whois struct {
	Owner sdk.AccAddress `json:"owner" yaml:"owner"`
	Value string         `json:"value" yaml:"value"`
	Price sdk.Coins      `json:"price" yaml:"price"`
}

// NewWhois  返回一个 Whois对象 最小的价格
func NewWhois() Whois {
	return Whois{
		Price: MinNamePrice,
	}
}

// String  格式化为字符串
func (w Whois) String() string {
	return strings.TrimSpace(fmt.Sprintf(
		`Owner:%s
		Value:%s
		Price:%s`, w.Owner, w.Value, w.Price))
}
