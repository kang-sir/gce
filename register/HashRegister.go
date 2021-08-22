package register

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"errors"
	"gce/crypt/sm3"
	"hash"
)

var hashFuncMap = make(map[string]func() hash.Hash)

func init() {
	hashFuncMap["MD5"] = md5.New
	hashFuncMap["SHA1"] = sha1.New
	hashFuncMap["SHA256"] = sha256.New
	hashFuncMap["SHA512"] = sha512.New
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
	if hashFunc, ok := hashFuncMap[hashName]; ok {
		return hashFunc
	}
	panic("name=[" + hashName + "]的hash函数尚未注册")
}
