package main

import (
	"fmt"
	"gce/crypto/sm2"
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
	fmt.Println(key.PublicKey.X)
	fmt.Println(key.PublicKey.Y)
	fmt.Println(key.D)
}
