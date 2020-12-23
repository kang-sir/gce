package sm2

import (
	"crypto/rand"
	"gce/constant"
	"io"
	"math/big"
)

type PublicKey struct {
	X, Y *big.Int
}
type PrivateKey struct {
	PublicKey
	D *big.Int
}

func GenerateKey() (priKey *PrivateKey, err error) {
	// 产生随机数D，作为私钥
	randBytes := make([]byte, sm2Curve.N.BitLen()/8)
	d := new(big.Int)
	for {
		io.ReadFull(rand.Reader, randBytes)
		d.SetBytes(randBytes)
		if d.Cmp(constant.BigIntOne) > 0 && d.Cmp(sm2Curve.N) < 0 {
			break
		}
	}
	// 计算 P(x,y)=[k]G
	P := GBasePoint.Multiply(d)
	priKey = &PrivateKey{D: d, PublicKey: PublicKey{X: P.X.Num, Y: P.Y.Num}}
	return
}
