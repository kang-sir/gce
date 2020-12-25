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

}
