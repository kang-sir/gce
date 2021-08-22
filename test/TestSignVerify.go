package main

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/hex"
	"fmt"
	"gce/constant"
	"gce/params"
	"gce/provider/soft"
)

func main() {

	provider := soft.GoSoftProvider{Name: "SOFT-GCE"}
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	//key, err := sm2.GenerateKey()
	if err != nil {
		fmt.Println(err)
	}
	hashData, err := provider.HashData([]byte("123"), params.HashAlg{
		AlgName: constant.SHA256,
	})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(hashData)

	signBytes, err := provider.SignHash(key, hashData)
	if err != nil {
		fmt.Println(err)
	}
	verRes, err := provider.VerifySignByHash(&key.PublicKey, hashData, signBytes)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(verRes)

	cipher, err := provider.PubKeyEncrypt(&key.PublicKey, []byte("123"))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("私钥加密密文：" + hex.EncodeToString(cipher))

	plain, err := provider.PriKeyDecrypt(key, cipher)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("私钥解密明文：" + string(plain))
}
