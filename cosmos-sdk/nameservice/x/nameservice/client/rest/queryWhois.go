package rest

import (
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
	"github.com/youngqqcn/nameservice/x/nameservice/types"
)

func listWhoisHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/list-whois", storeName), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func getWhoisHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		key := vars["key"]

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/get-whois/%s", storeName, key), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func resolveNameHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", storeName, types.QueryResolveName, name), nil)
		if err != nil {
			rest.WriteErrorResponse(rw, http.StatusNotFound, err.Error())
			return
		}
		rest.PostProcessResponse(rw, cliCtx, res)
	}
}
