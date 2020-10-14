package main

import (
	"fmt"
	"github.com/tendermint/tendermint/rpc/client/http"
	"mytokenapp/mytokenapp"
	"time"
)



func main()  {

	c, err := http.New("http://127.0.0.1:26657", "/websocket")
	if err != nil {
		// handle error
		panic(err)
	}

	// call Start/Stop if you're subscribing to events
	err = c.Start()
	if err != nil {
		// handle error
		panic(err)
	}
	defer c.Stop()

	res, err := c.Status()
	if err != nil {
		// handle error
		panic(err)
	}

	fmt.Println(res)


	// release

	release(c)

	transfer(c)

}


func transfer(c *http.HTTP) {
	wallet := mytokenapp.LoadWalletFromFile("/home/yqq/BlockchainNotesV2/tendermint/mytokenapp/mytokenapp/superuser.wallet")
	tx :=  mytokenapp.NewTx( mytokenapp.NewTransferPayload(
			wallet.GetAddress("yqq"), wallet.GetAddress("superuser"),
		1000, time.Now().Unix(), "transfer test" ) )

	if err := tx.Sign( wallet.GetPrivKeyByLabel("superuser") ) ; err != nil {
		panic(err)
	}

	bztx, err := mytokenapp.MarshalBinaryBare(tx)
	if err != nil {
		panic(err)
	}

	ret , err := c.BroadcastTxCommit(  bztx)
	if err != nil {
		panic(err)
	}

	fmt.Println("broadcast response : %v", ret)
}


func release( c *http.HTTP) {
	wallet := mytokenapp.LoadWalletFromFile("/home/yqq/BlockchainNotesV2/tendermint/mytokenapp/mytokenapp/superuser.wallet")
	tx :=  mytokenapp.NewTx( mytokenapp.NewReleasePayload( wallet.GetAddress("superuser"),
		wallet.GetAddress("yqq"), 1000, time.Now().Unix(), "release" ) )

	if err := tx.Sign( wallet.GetPrivKeyByLabel("superuser") ) ; err != nil {
		panic(err)
	}

	bztx, err := mytokenapp.MarshalBinaryBare(tx)
	if err != nil {
		panic(err)
	}

	ret , err := c.BroadcastTxCommit(  bztx)
	if err != nil {
		panic(err)
	}

	fmt.Println("broadcast response : %v", ret)
}
