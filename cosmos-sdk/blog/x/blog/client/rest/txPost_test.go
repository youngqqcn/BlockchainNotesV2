package rest

import (
	"encoding/json"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
)

func TestCreatePost(t *testing.T) {

	t.Log("=====test create post")

	pr := createPostRequest{
		BaseReq: rest.NewBaseReq("from",
			"memo",
			"chaiid", "gas", "gasAdjustment", 123, 13,
			sdk.NewCoins(sdk.NewCoin("cosmos", sdk.NewInt(1))),
			sdk.NewDecCoins(sdk.NewDecCoin("cosmos", sdk.NewInt(1))),
			true),
		Creator: "cosmos1wt47yve6l29yjtxtsajhltr2vqhf7mpw5n6fx6",
		Title:   "this is title",
		Body:    "this is body",
	}

	// ctx := context.NewCLIContext().WithCodec(cdc)
	bz, err := json.Marshal(pr)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("data : %v\n", bz)
}
