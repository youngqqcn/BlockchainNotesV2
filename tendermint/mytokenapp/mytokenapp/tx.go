package mytokenapp

import "crypto"

type Tx struct {
	Payload 		Payload				// 交易内容
	Signature 		[]byte				// 签名
	PubKey			crypto.PublicKey	// 公钥
}









