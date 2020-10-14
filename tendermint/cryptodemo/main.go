package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	cryptoamino "github.com/tendermint/tendermint/crypto/encoding/amino"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	//amino "github.com/tendermint/go-amino"
)

func testSecp256k1() {
	/*
		// tendermint  secp256k1  地址生成实现

		// Address returns a Bitcoin style addresses: RIPEMD160(SHA256(pubkey))
		func (pubKey PubKeySecp256k1) Address() crypto.Address {
			hasherSHA256 := sha256.New()
			hasherSHA256.Write(pubKey[:]) // does not error
			sha := hasherSHA256.Sum(nil)

			hasherRIPEMD160 := ripemd160.New()
			hasherRIPEMD160.Write(sha) // does not error
			return crypto.Address(hasherRIPEMD160.Sum(nil))
		}

	*/

	privKey := secp256k1.GenPrivKey()

	fmt.Println("private key: ", hex.EncodeToString(privKey.Bytes()) )

	pubKey := privKey.PubKey()
	fmt.Println("public key: ", hex.EncodeToString( pubKey.Bytes()) )

	address := privKey.PubKey().Address()
	fmt.Println("address : ", address)


}





func testEd25519()  {

	privKey := ed25519.GenPrivKey()
	fmt.Println("private key: ", hex.EncodeToString(privKey.Bytes()) )

	pubKey := privKey.PubKey()
	fmt.Println("public key: ", hex.EncodeToString( pubKey.Bytes()) )

	address := privKey.PubKey().Address()
	fmt.Println("address : ", address)
}

type Letter struct {
	Msg []byte
	Sig []byte
	PubKey []byte
}

func testSign() string {

	msg := []byte("hello")

	privKey := secp256k1.GenPrivKey()
	sig, _ := privKey.Sign(msg)


	pubKeyBytes := privKey.PubKey().Bytes()  // 使用了 amino 编码格式
	letter := Letter{ msg, sig,  pubKeyBytes }

	data , _ := json.Marshal( letter  )
	fmt.Println( hex.EncodeToString( data) )

	return hex.EncodeToString(data)
}

func testVerify( letter string ) bool  {

	rawdata , _ := hex.DecodeString(letter)
	lt := Letter{}
	err := json.Unmarshal( rawdata , &lt)
	if err != nil {
		fmt.Println(err)
		return false
	}

	//var _ crypto.PubKey = secp256k1.PubKeySecp256k1{}
	var _ crypto.PubKey = ed25519.PubKeyEd25519{}
	pubkey , _ := cryptoamino.PubKeyFromBytes(lt.PubKey)  // 使用amino进行法序列化

	fmt.Println( pubkey.Bytes() )

	return pubkey.VerifyBytes(lt.Msg, lt.Sig)
}



func main() {

	//testSecp256k1()

	//testEd25519()

	letter :=  testSign()
	if  testVerify(letter) {
		fmt.Println("verify successed")
	} else {
		fmt.Println("verify failed")
	}

}
