package soft

import (
	"crypto"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"gce/crypt"
	"gce/crypt/sm2"
	"gce/crypt/sm3"
	"gce/params"
	"gce/register"
	"reflect"
	"strings"
)

type GoSoftProvider struct {
	Name string
}

func (provider *GoSoftProvider) GenRandom(len int32) (randBytes []byte, err error) {
	randBytes = make([]byte, len)
	rand.Read(randBytes)
	return
}

func (provider *GoSoftProvider) HashData(oriDataBytes []byte, alg params.HashAlg) (hashBytes []byte, err error) {
	hashFunc := register.GetHashFunc(string(alg.AlgName))
	hash := hashFunc()
	if param, ok := alg.AlgParam.(params.SM3Param); ok {
		ZA := sm3.GetZA(*param.PubKey.X, *param.PubKey.Y, param.UserId)
		hash.Write(ZA)
	}
	hashBytes = hash.Sum(oriDataBytes)
	return
}

func (provider *GoSoftProvider) SignHash(priKey crypto.PrivateKey, hashDataBytes []byte) (signBytes []byte, err error) {
	if priKeyObj, ok := priKey.(*sm2.PrivateKey); ok {
		signBytes, err = sm2.SignHashPkcs1(priKeyObj, hashDataBytes)
	} else if priKeyObj, ok := priKey.(*rsa.PrivateKey); ok {
		signBytes, err = rsa.SignPKCS1v15(rand.Reader, priKeyObj, 0, hashDataBytes)
	} else {
		err = errors.New("not supported private key type " + reflect.TypeOf(priKey).String())
	}
	return
}

func (provider *GoSoftProvider) VerifySignByHash(pubKey crypto.PublicKey, hashDataBytes []byte, signBytes []byte) (verRes bool, err error) {
	if pubKeyObj, ok := pubKey.(*sm2.PublicKey); ok {
		verRes, err = sm2.VerifySignByHash(pubKeyObj, hashDataBytes, signBytes)
	} else if pubKeyObj, ok := pubKey.(*rsa.PublicKey); ok {
		err = rsa.VerifyPKCS1v15(pubKeyObj, 0, hashDataBytes, signBytes)
		verRes = err == nil
	} else {
		err = errors.New("not supported public key type " + reflect.TypeOf(pubKey).String())
	}
	return
}

func (provider *GoSoftProvider) PubKeyEncrypt(pubKey crypto.PublicKey, plainData []byte) (cipher []byte, err error) {
	if pubKeyObj, ok := pubKey.(*sm2.PublicKey); ok {
		cipher, err = sm2.Encrypt(pubKeyObj, plainData)
	} else if pubKeyObj, ok := pubKey.(*rsa.PublicKey); ok {
		cipher, err = rsa.EncryptPKCS1v15(rand.Reader, pubKeyObj, plainData)
	} else {
		err = errors.New("not supported public key type " + reflect.TypeOf(pubKey).String())
	}
	return
}

func (provider *GoSoftProvider) PriKeyDecrypt(priKey crypto.PrivateKey, cipherData []byte) (plain []byte, err error) {
	if priKeyObj, ok := priKey.(*sm2.PrivateKey); ok {
		plain, err = sm2.Decrypt(priKeyObj, cipherData)
	} else if priKeyObj, ok := priKey.(*rsa.PrivateKey); ok {
		plain, err = rsa.DecryptPKCS1v15(rand.Reader, priKeyObj, cipherData)
	} else {
		err = errors.New("not supported private key type " + reflect.TypeOf(priKey).String())
	}
	return
}

func (provider *GoSoftProvider) SymEncrypt(symKey []byte, alg string, iv []byte, plainData []byte) (cipherData []byte, err error) {
	algAttrs := strings.Split(alg, "/")
	if len(algAttrs) != 3 {
		err = errors.New("err format sym alg " + alg)
		return
	}
	// 获取对称加密的block函数
	cipherBlockFunc := register.GetSymCipherBlock(algAttrs[0])
	block, err := cipherBlockFunc(symKey)
	// 获取padding方式对象
	padding, err := register.GetPadding(algAttrs[2])
	if err != nil {
		return
	}
	packData, err := padding.PackData(plainData, block.BlockSize())
	if err != nil {
		return
	}
	var blockMode cipher.BlockMode
	switch algAttrs[1] {
	case "ECB":
		blockMode = crypt.NewECBEncrypter(block)
		break
	case "CBC":
		blockMode = cipher.NewCBCEncrypter(block, iv)
		break
	default:
		panic("not supported encrypt mode = " + algAttrs[1])
	}
	cipherData = make([]byte, len(packData))
	blockMode.CryptBlocks(cipherData, packData)
	return
}

func (provider *GoSoftProvider) SymDecrypt(symKey []byte, alg string, iv []byte, cipherData []byte) (plain []byte, err error) {
	algAttrs := strings.Split(alg, "/")
	if len(algAttrs) != 3 {
		err = errors.New("err format sym alg " + alg)
		return
	}
	// 获取对称加密的block函数
	cipherBlockFunc := register.GetSymCipherBlock(algAttrs[0])
	block, err := cipherBlockFunc(symKey)
	if err != nil {
		return
	}
	// 获取padding方式对象
	padding, err := register.GetPadding(algAttrs[2])
	if err != nil {
		return
	}
	var blockMode cipher.BlockMode
	switch algAttrs[1] {
	case "ECB":
		blockMode = crypt.NewECBDecrypter(block)
		break
	case "CBC":
		blockMode = cipher.NewCBCDecrypter(block, iv)
		break
	default:
		panic("not supported decrypt mode = " + algAttrs[1])
	}
	plain = make([]byte, len(cipherData))
	blockMode.CryptBlocks(plain, cipherData)
	plain, err = padding.DePackData(plain)
	return
}
