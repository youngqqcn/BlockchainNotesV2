package rest

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/youngqqcn/scavenge/x/scavenge/types"
)

// Used to not have an error if strconv is unused
var _ = strconv.Itoa(42)

type createScavengeRequest struct {
	BaseReq     rest.BaseReq `json:"base_req"`
	Creator     string       `json:"creator"`
	Description string       `json:"description"`
	// SolutionHash string       `json:"solutionHash"`
	Reward   string `json:"reward"`
	Solution string `json:"solution"`
	// Scavenger    string       `json:"scavenger"`
}

func createScavengeHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createScavengeRequest
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}
		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}
		creator, err := sdk.AccAddressFromBech32(req.Creator)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		parsedDescription := req.Description

		solutionHash := sha256.Sum256([]byte(req.Solution))
		solutionHashHex := hex.EncodeToString(solutionHash[:])

		parsedSolutionHash := solutionHashHex

		rewardStr := req.Reward
		if !strings.Contains(rewardStr, "token") {
			rewardStr += "token"
		}

		parsedReward, err := sdk.ParseCoins(rewardStr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		}

		parsedSolution := req.Solution

		// parsedScavenger := req.Scavenger

		msg := types.NewMsgCreateScavenge(
			creator,
			parsedDescription,
			parsedSolutionHash,
			parsedReward,
			parsedSolution,
		)

		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, baseReq, []sdk.Msg{msg})
	}
}

type createRevealRequest struct {
	BaseReq  rest.BaseReq `json:"base_req"`
	Creator  string       `json:"creator"`
	Solution string       `json:"solution"`
}

func createRevealHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		var req createScavengeRequest
		if !rest.ReadRESTReq(rw, r, cliCtx.Codec, &req) {
			rest.WriteErrorResponse(rw, http.StatusBadRequest, "failed to parse request")
			return
		}
		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(rw) {
			return
		}
		creator, err := sdk.AccAddressFromBech32(req.Creator)
		if err != nil {
			rest.WriteErrorResponse(rw, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgRevealSolution(creator, req.Solution)
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(rw, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(rw, cliCtx, baseReq, []sdk.Msg{msg})
	}

}
