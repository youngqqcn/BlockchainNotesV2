package rest

import (
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
)

// RegisterRoutes registers myapp-related REST handlers to a router
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
  // this line is used by starport scaffolding # 1
		r.HandleFunc("/myapp/student", createStudentHandler(cliCtx)).Methods("POST")
		r.HandleFunc("/myapp/student", listStudentHandler(cliCtx, "myapp")).Methods("GET")
		r.HandleFunc("/myapp/student/{key}", getStudentHandler(cliCtx, "myapp")).Methods("GET")
		r.HandleFunc("/myapp/student", setStudentHandler(cliCtx)).Methods("PUT")
		r.HandleFunc("/myapp/student", deleteStudentHandler(cliCtx)).Methods("DELETE")

		
		r.HandleFunc("/myapp/post", createPostHandler(cliCtx)).Methods("POST")
		r.HandleFunc("/myapp/post", listPostHandler(cliCtx, "myapp")).Methods("GET")
		r.HandleFunc("/myapp/post/{key}", getPostHandler(cliCtx, "myapp")).Methods("GET")
		r.HandleFunc("/myapp/post", setPostHandler(cliCtx)).Methods("PUT")
		r.HandleFunc("/myapp/post", deletePostHandler(cliCtx)).Methods("DELETE")

		
}
