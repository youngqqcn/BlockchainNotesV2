package rest

import (
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
)

// RegisterRoutes registers blog-related REST handlers to a router
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	// this line is used by starport scaffolding # 1
		r.HandleFunc("/blog/comment", createCommentHandler(cliCtx)).Methods("POST")
		r.HandleFunc("/blog/comment", listCommentHandler(cliCtx, "blog")).Methods("GET")
		r.HandleFunc("/blog/comment/{key}", getCommentHandler(cliCtx, "blog")).Methods("GET")
		r.HandleFunc("/blog/comment", setCommentHandler(cliCtx)).Methods("PUT")
		r.HandleFunc("/blog/comment", deleteCommentHandler(cliCtx)).Methods("DELETE")

		
	//   r.HandleFunc("/blog/comment", listCommentHandler(cliCtx, "blog")).Methods("GET")
	//   r.HandleFunc("/blog/comment", createCommentHandler(cliCtx)).Methods("POST")
	r.HandleFunc("/blog/post", listPostHandler(cliCtx, "blog")).Methods("GET")
	r.HandleFunc("/blog/post", createPostHandler(cliCtx)).Methods("POST")
}
