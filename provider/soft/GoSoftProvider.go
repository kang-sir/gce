package soft

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"gce/crypto/sm2"
	"gce/crypto/sm3"
	"gce/params"
	"gce/register"
	"reflect"
)

type GoSoftProvider struct {
	Name string
}

func (provider *GoSoftProvider) HashData(oriDataBytes []byte, alg params.HashAlg) (hashBytes []byte, err error) {
	hashFunc := register.GetHashFunc(string(alg.AlgName))
	hash := hashFunc()
	if param, ok := alg.AlgParam.(params.SM3Param); ok {
		if sm3Func, ok := hash.(*sm3.SM3); ok {
			sm3Func.AddId(*param.PubKey.X, *param.PubKey.Y, param.UserId)
		}
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
