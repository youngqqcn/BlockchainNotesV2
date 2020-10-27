package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Post 博客的文章类型
type Post struct {
	// 消息的发送者
	Creator sdk.AccAddress `json:"creator" yaml:"creator"`
	// 全局唯一的标识
	ID    string `json:"id" yaml:"id"`
	Title string `json:"title" yaml:"title"`
	Body  string `json:"body" yaml:"body"`
}
