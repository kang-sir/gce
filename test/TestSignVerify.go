package main

import (
	"fmt"
	"gce/constant"
	"gce/crypto/sm2"
	"gce/params"
	"gce/provider/soft"
)

func main() {

	provider := soft.GoSoftProvider{Name: "SOFT-GCE"}
	key, err := sm2.GenerateKey()
	if err != nil {
		fmt.Println(err)
	}
	hashData, err := provider.HashData([]byte("123"), params.HashAlg{
		AlgName: constant.SM3,
		AlgParam: params.SM3Param{
			PubKey: key.PublicKey,
			UserId: []byte("1234567812345678"),
		},
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
}
