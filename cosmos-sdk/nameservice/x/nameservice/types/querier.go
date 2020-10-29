package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const QueryListWhois = "list-whois"
const QueryGetWhois = "get-whois"
const QueryResolveName = "resolve-name"

type QueryResResolve struct {
	Value string `json:"value"`
}

func (r QueryResResolve) String() string {
	return r.Value
}

type QueryResNames []string

func (r QueryResNames) String() string {
	return strings.Join(r[:], "\n")
}



// 自定义的查询, 而不使用 types.Whois, 是为了把name字段也加上
type QueryWhois struct {
	Owner sdk.AccAddress `json:"owner" `
	Value string         `json:"value"`
	Price sdk.Coins      `json:"price"`
	Name string         `json:"name"`
}

// String  格式化为字符串
func (w QueryWhois) String() string {
	return strings.TrimSpace(fmt.Sprintf(
		`Name:%s
		Owner:%s
		Value:%s
		Price:%s`, w.Name, w.Owner, w.Value, w.Price))
}

type QueryResNameses []string

func (r QueryResNameses) String() string {
	return strings.Join(r[:], "\n")
}
