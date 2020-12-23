package provider

import (
	"crypto"
	"gce/params"
)

type GceProvider interface {
	HashData(oriDataBytes []byte, alg params.HashAlg)

	SignHash(hashDataBytes []byte, priKey crypto.PrivateKey)

	VerifySignByHash(hashDataBytes []byte, pubKey crypto.PublicKey)
}
