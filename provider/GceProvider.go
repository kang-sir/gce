package provider

import (
	"crypto"
	"gce/params"
)

type GceProvider interface {
	HashData(oriDataBytes []byte, alg params.HashAlg) (hashBytes []byte, err error)

	SignHash(priKey crypto.PrivateKey, hashDataBytes []byte) (signBytes []byte, err error)

	VerifySignByHash(pubKey crypto.PublicKey, hashDataBytes []byte, signBytes []byte) (verRes bool, err error)
}
