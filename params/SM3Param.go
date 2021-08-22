package params

import (
	"gce/crypt/sm2"
)

type SM3Param struct {
	PubKey sm2.PublicKey
	UserId []byte
}
