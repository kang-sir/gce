package provider

import (
	"crypto"
	"gce/params"
)

type GceProvider interface {
	GenRandom(len int32) (randBytes []byte, err error)

	HashData(oriDataBytes []byte, alg params.HashAlg) (hashBytes []byte, err error)

	SignHash(priKey crypto.PrivateKey, hashDataBytes []byte) (signBytes []byte, err error)

	VerifySignByHash(pubKey crypto.PublicKey, hashDataBytes []byte, signBytes []byte) (verRes bool, err error)

	PubKeyEncrypt(pubKey crypto.PublicKey, plainData []byte) (cipher []byte, err error)

	PriKeyDecrypt(priKey crypto.PrivateKey, cipherData []byte) (plain []byte, err error)

	SymEncrypt(symKey []byte, alg string, iv []byte, plainData []byte) (cipher []byte, err error)

	SymDecrypt(symKey []byte, alg string, iv []byte, cipherData []byte) (plain []byte, err error)
}
