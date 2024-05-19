package model

import (
	"crypto/rand"
	"fmt"
	"github.com/tjfoc/gmsm/sm2"
	"os"
)

var PublicKeyX string
var PublicKeyY string
var PrivateKey *sm2.PrivateKey

func init() {
	// 生成 SM2 密钥对
	var err error
	PrivateKey, err = sm2.GenerateKey(rand.Reader)
	if err != nil {
		fmt.Println("Failed to generate SM2 key pair:", err)
		os.Exit(1)
	}
	publicKey := PrivateKey.PublicKey
	PublicKeyX = publicKey.X.String()
	PublicKeyY = publicKey.Y.String()

	fmt.Println("sm2公钥序列化成功...")
}
