package main

import (
	"encoding/base64"
	"fmt"
	"gce/crypto/sm2"
	"gce/crypto/sm3"
	"math/big"
)

type MyName struct {
	*Address
	Name string
	Age  *big.Int
}

type Address struct {
	Email string
	Phone string
}

var myName MyName

func main() {

	key, err := sm2.GenerateKey()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(base64.StdEncoding.EncodeToString(key.X.Bytes()))
	fmt.Println(base64.StdEncoding.EncodeToString(key.Y.Bytes()))
	hash := sm3.New().Sum([]byte("abc"))
	pkcs1, err := sm2.SignHashPkcs1(key, hash)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(base64.StdEncoding.EncodeToString(pkcs1))

	hash = sm3.New().Sum([]byte("abc"))
	res, err := sm2.VerifySignByHash(&key.PublicKey, hash, pkcs1)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("验签结果：", res)

	encrypt, err := sm2.Encrypt(&key.PublicKey, []byte("123测试数据"))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("加密结果：", base64.StdEncoding.EncodeToString(encrypt))

	decrypt, err := sm2.Decrypt(key, encrypt)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("解密结果：", string(decrypt))

}
