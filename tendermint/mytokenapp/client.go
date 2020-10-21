package main

import (
	"encoding/hex"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/crypto/merkle"
	"github.com/tendermint/tendermint/libs/os"
	"strings"

	//"github.com/tendermint/tendermint/libs/os"
	"github.com/tendermint/tendermint/rpc/client/http"
	myapp "mytokenapp/mytokenapp"
	"time"
)

var cli, _ = http.New("http://127.0.0.1:26657", "/websocket")

var subcmd string
var rootCmd = &cobra.Command{
	Use: "client",
	Run: func(cmd *cobra.Command, args []string) {
		if len(subcmd) == 0 {
			cmd.Help()
			return
		}
	},
}

func init() {
	//rootCmd.Flags().StringVarP(&subcmd, "initwallet", "i", "", "init wallet for test")
	//rootCmd.Flags().StringVarP(&subcmd, "release", "r", "", "release token")
	//rootCmd.Flags().StringVarP(&subcmd, "transfer", "t", "default", "transfer token ")
	//rootCmd.Flags().StringVarP(&subcmd, "initwallet", "i", "", "init wallet for test")
}

func main() {

	initWalletCmd := cobra.Command{
		Use:   "initwallet",
		Short: "init wallet for test",
		Run: func(cmd *cobra.Command, args []string) {

			filepath, err := cmd.Flags().GetString("filepath")
			if err != nil {
				cmd.Help()
				return
			}

			labels, err := cmd.Flags().GetStringArray("labels")
			if err != nil {
				cmd.Help()
				return
			}

			if (len(filepath) == 0) || (len(labels) == 0) {
				cmd.Help()
				return
			}

			initWallet(filepath, labels)
		},
	}

	initWalletCmd.Flags().StringP("filepath", "f", "wallet.dat", " filepath to save wallet file ")
	// e.g:  ./bin/client initwallet  --filepath wallet.dat -lyqq -lsuperuser -ltom -ljack
	initWalletCmd.Flags().StringArrayP("labels", "l", []string{"yqq", "superuser", "tom"}, "labels of wallet accounts")

	releaseCmd := cobra.Command{
		Use:   "release",
		Short: "release token to 'superuser'",
		Run: func(cmd *cobra.Command, args []string) {

			walletFilePath, _ := cmd.Flags().GetString("walletfile")
			toLabel, _ := cmd.Flags().GetString("tolabel")
			valueRelease, _ := cmd.Flags().GetInt64("value")

			if os.FileExists(walletFilePath) && len(toLabel) > 0 && valueRelease > 0 {
				release(walletFilePath, toLabel, valueRelease)
				return
			}
			cmd.Help()
			return
		},
	}
	releaseCmd.Flags().StringP("walletfile", "w", "wallet.dat", " path of wallet file ")
	releaseCmd.Flags().StringP("tolabel", "t", "yqq", "the account that release token to")
	releaseCmd.Flags().Int64P("value", "v", 1000, "the value to release")

	transferCmd := cobra.Command{
		Use:   "transfer",
		Short: "transfer token test",
		Run: func(cmd *cobra.Command, args []string) {
			fromLabel, _ := cmd.Flags().GetString("fromlabel")
			walletFilePath, _ := cmd.Flags().GetString("walletfile")
			toLabel, _ := cmd.Flags().GetString("tolabel")
			valueTransfer, _ := cmd.Flags().GetInt64("value")

			if os.FileExists(walletFilePath) && len(toLabel) > 0 && len(fromLabel) > 0 && valueTransfer > 0 {
				transfer(walletFilePath, fromLabel, toLabel, valueTransfer)
				return
			}

			cmd.Help()
			return
		},
	}

	transferCmd.Flags().StringP("walletfile", "w", "wallet.dat", " path of wallet file ")
	transferCmd.Flags().StringP("fromlabel", "f", "yqq", "the from account ")
	transferCmd.Flags().StringP("tolabel", "t", "tom", "the account that release token to")
	transferCmd.Flags().Int64P("value", "v", 1000, "the value to transfer")

	queryBalanceCmd := cobra.Command{
		Use: "querybalance",
		Run: func(cmd *cobra.Command, args []string) {
			filepath, _ := cmd.Flags().GetString("walletfile")
			label, _ := cmd.Flags().GetString("label")
			queryBalance(filepath, label)
		},
	}
	queryBalanceCmd.Flags().StringP("walletfile", "w", "wallet.dat", " path of wallet file ")
	queryBalanceCmd.Flags().StringP("label", "l", "yqq", "the account label to query")

	queryTxCmd := cobra.Command{
		Use: "querytx",
		Run: func(cmd *cobra.Command, args []string) {
			txhash, _ := cmd.Flags().GetString("hash")
			_, err := hex.DecodeString(txhash)
			if len(txhash) > 0 && err == nil {
				queryTx(txhash)
				return
			}
			cmd.Help()
			return
		},
	}
	queryTxCmd.Flags().StringP("hash", "x", "", " tx hash")

	rootCmd.AddCommand(&initWalletCmd)
	rootCmd.AddCommand(&releaseCmd)
	rootCmd.AddCommand(&transferCmd)
	rootCmd.AddCommand(&queryBalanceCmd)
	rootCmd.AddCommand(&queryTxCmd)
	//
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}

}

func initWallet(filepath string, labels []string) error {

	if len(labels) == 0 {
		return fmt.Errorf("labels is empty")
	}
	if len(filepath) == 0 {
		return fmt.Errorf("invalid filepath")
	}

	nw := myapp.NewWallet()
	for _, label := range labels {
		nw.GenNewPrivKey(label)
	}
	return nw.Save(filepath)
}

// 查询交易
func queryTx(txHash string) {

	txhash, err := hex.DecodeString(txHash)
	if err != nil {
		panic(err)
	}

	tx, err := cli.Tx(txhash, true)
	if err != nil {
		panic(err)
	}

	fmt.Printf("txhash: %v\n", hex.EncodeToString(tx.Hash))
	fmt.Printf("height: %v\n", tx.Height)
	//fmt.Printf("tx: %v\n", tx.Tx.String())

	tx.Tx.String()

	var trx myapp.Tx
	//hexTx, _ := hex.DecodeString( tx.Tx.String())
	err = myapp.UnMarshalBinaryBare(tx.Tx, &trx)
	if err != nil {
		panic(err)
	}

	//fmt.Println("from:", trx.Payload.(* (myapp.TransferPayload) ).FromAddress)

	if _, ok := trx.Payload.(*myapp.ReleasePayload); ok {
		p, _ := trx.Payload.(*myapp.ReleasePayload)
		fmt.Println("from:", p.FromAddress)
		fmt.Println("to:", p.ToAddress)
		fmt.Println("value:", p.Value)
		fmt.Println("sequence:", p.Sequence)
		fmt.Println("memo:", p.Memo)
	}
	if _, ok := trx.Payload.(*myapp.TransferPayload); ok {
		p, _ := trx.Payload.(*myapp.TransferPayload)
		fmt.Println("from:", p.FromAddress)
		fmt.Println("to:", p.ToAddress)
		fmt.Println("value:", p.Value)
		fmt.Println("sequence:", p.Sequence)
		fmt.Println("memo:", p.Memo)
	}

	//
	//jsonTx ,err := myapp.MarshalJSON(trx)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println("======================")
	//fmt.Println(jsonTx)

	// OK
	// 参考: /home/yqq/go/pkg/mod/github.com/orientwalt/tendermint@v90.0.7+incompatible/lite/proxy/query_test.go
	// 第一种方式   Validate 内部调用了  Proof.Verify
	key := tx.Tx.Hash()
	keyHash := merkle.SimpleHashFromByteSlices([][]byte{key})
	if err := tx.Proof.Validate(keyHash); err != nil {
		panic(err)
	}

	// 第二种方式
	// 参考: /home/yqq/go/pkg/mod/github.com/tendermint/tendermint@v0.33.8/rpc/client/rpc_test.go
	proof := tx.Proof
	if err := proof.Proof.Verify(proof.RootHash, txhash); err != nil {
		panic(err)
	}

	fmt.Println()

}

// 查询余额
func queryBalance(filepath, label string) {
	wallet := myapp.LoadWalletFromFile(filepath)
	address := wallet.GetAddress(label)

	rsp, err := cli.ABCIQuery("", []byte(strings.ToUpper(hex.EncodeToString(address))))
	if err != nil {
		panic(err)
	}

	fmt.Printf("label:%v , address:%v, balance : %+v\n", label, address.String(), string(rsp.Response.Value))
}

func transfer(filepath, fromLabel, toLabel string, value int64) {
	wallet := myapp.LoadWalletFromFile(filepath)
	tx := myapp.NewTx(myapp.NewTransferPayload(
		wallet.GetAddress(fromLabel), wallet.GetAddress(toLabel),
		value, time.Now().Unix(), "transfer test"))

	if err := tx.Sign(wallet.GetPrivKeyByLabel(fromLabel)); err != nil {
		panic(err)
	}

	bztx, err := myapp.MarshalBinaryBare(tx)
	if err != nil {
		panic(err)
	}

	//fmt.Println("hex: ", hex.EncodeToString(bztx))

	ret, err := cli.BroadcastTxCommit(bztx)
	if err != nil {
		panic(err)
	}

	fmt.Printf("broadcast response : %v", ret)
	fmt.Printf("hash: %v", ret.Hash)
	fmt.Printf("height: %v", ret.Height)
	fmt.Printf("tx: %v", ret.DeliverTx.String())

}

func release(filepath, toLabel string, value int64) {
	wallet := myapp.LoadWalletFromFile(filepath)
	tx := myapp.NewTx(myapp.NewReleasePayload(wallet.GetAddress("superuser"),
		wallet.GetAddress(toLabel), value, time.Now().Unix(), "release"))

	if err := tx.Sign(wallet.GetPrivKeyByLabel("superuser")); err != nil {
		panic(err)
	}

	bztx, err := myapp.MarshalBinaryBare(tx)
	if err != nil {
		panic(err)
	}

	ret, err := cli.BroadcastTxCommit(bztx)
	if err != nil {
		panic(err)
	}

	fmt.Printf("broadcast response : %v", ret)
	fmt.Printf("hash: %v", ret.Hash)
	fmt.Printf("height: %v", ret.Height)
	fmt.Printf("tx: %v", ret.DeliverTx.String())
}
