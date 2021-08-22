package main

import (
	"encoding/hex"
	"fmt"
	"gce/provider/soft"
)

func main() {

	provider := soft.GoSoftProvider{Name: "SOFT-GCE"}

	symKey, _ := provider.GenRandom(24)

	cipherData, err := provider.SymEncrypt(symKey, "TDES/ECB/PKCS5Padding", nil, []byte("12345678123456781234567812345678"))

	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(hex.EncodeToString(cipherData))

	plainData, err := provider.SymDecrypt(symKey, "TDES/ECB/PKCS5Padding", nil, cipherData)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(plainData))
}
