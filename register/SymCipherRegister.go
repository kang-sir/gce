package register

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"errors"
	"gce/crypt/sm4"
)

var blockFuncMap = make(map[string]func([]byte) (cipher.Block, error))

func init() {
	blockFuncMap["SM4"] = sm4.NewCipher
	blockFuncMap["TDES"] = des.NewTripleDESCipher
	blockFuncMap["AES"] = aes.NewCipher
}

func SymCipherFuncRegister(name string, cipherFunc func([]byte) (cipher.Block, error)) (err error) {
	if _, ok := blockFuncMap[name]; ok {
		err = errors.New("name=[" + name + "]的SymCipher函数已被注册")
		return
	}
	blockFuncMap[name] = cipherFunc
	return
}

func GetSymCipherBlock(name string) func([]byte) (cipher.Block, error) {
	if cipherBlock, ok := blockFuncMap[name]; ok {
		return cipherBlock
	}
	panic("name=[" + name + "]的SymCipher函数尚未注册")
}
