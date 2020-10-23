package main

import (
	"fmt"
	"github.com/tendermint/tendermint/abci/server"
	"github.com/tendermint/tendermint/abci/types"
)


func main()  {

	app := types.NewBaseApplication()
	svr , err := server.NewServer(":26658", "socket", app)
	if err != nil {
		panic(err)
	}
	err = svr.Start()
	if err != nil {
		panic(err)
	}
	defer svr.Stop()
	fmt.Println("abci server started")

	select {}
}



