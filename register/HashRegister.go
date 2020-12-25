package register

import (
	"crypto"
	"errors"
	"gce/crypto/sm3"
	"hash"
)

var hashFuncMap = make(map[string]func() hash.Hash)

func init() {
	hashFuncMap["MD4"] = crypto.MD4.New
	hashFuncMap["MD5"] = crypto.MD5.New
	hashFuncMap["SHA1"] = crypto.SHA1.New
	hashFuncMap["SHA224"] = crypto.SHA224.New
	hashFuncMap["SHA256"] = crypto.SHA256.New
	hashFuncMap["SHA384"] = crypto.SHA384.New
	hashFuncMap["SHA512"] = crypto.SHA512.New
	hashFuncMap["SM3"] = sm3.New
}

func HashFuncRegister(name string, hashFunc func() hash.Hash) (err error) {
	if _, ok := hashFuncMap[name]; ok {
		err = errors.New("name=[" + name + "]的hash函数已被注册")
		return
	}
	hashFuncMap[name] = hashFunc
	return
}

func GetHashFunc(hashName string) func() hash.Hash {
	return hashFuncMap[hashName]
}
