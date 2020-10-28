package rest

import (
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
)

// RegisterRoutes registers nameservice-related REST handlers to a router
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
  // this line is used by starport scaffolding # 1
		r.HandleFunc("/nameservice/name", createNameHandler(cliCtx)).Methods("POST")
		r.HandleFunc("/nameservice/name", listNameHandler(cliCtx, "nameservice")).Methods("GET")
		r.HandleFunc("/nameservice/name/{key}", getNameHandler(cliCtx, "nameservice")).Methods("GET")
		r.HandleFunc("/nameservice/name", setNameHandler(cliCtx)).Methods("PUT")
		r.HandleFunc("/nameservice/name", deleteNameHandler(cliCtx)).Methods("DELETE")

		
}
