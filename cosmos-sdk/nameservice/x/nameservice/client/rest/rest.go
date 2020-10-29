package rest

import (
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
)

// RegisterRoutes registers nameservice-related REST handlers to a router
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	// this line is used by starport scaffolding # 1
	r.HandleFunc("/nameservice/whois", buyNameHandler(cliCtx)).Methods("POST")
	r.HandleFunc("/nameservice/whois", listWhoisHandler(cliCtx, "nameservice")).Methods("GET")
	r.HandleFunc("/nameservice/whois/{key}", getWhoisHandler(cliCtx, "nameservice")).Methods("GET")
	r.HandleFunc("/nameservice/whois/resolve-name/{name}", resolveNameHandler(cliCtx, "nameservice")).Methods("GET")
	r.HandleFunc("/nameservice/whois", setNameHandler(cliCtx)).Methods("PUT")
	r.HandleFunc("/nameservice/whois", deleteNameHandler(cliCtx)).Methods("DELETE")
}
