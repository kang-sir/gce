package sm2

import (
	"gce/constant"
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
	d := getRandomD()
	// 计算 P(x,y)=[k]G
	P := GBasePoint.Multiply(d)
	priKey = &PrivateKey{D: d, PublicKey: PublicKey{X: P.X.Num, Y: P.Y.Num}}
	return
}

func SignHashPkcs1(priKey *PrivateKey, hashBytes []byte) (signBytes []byte, err error) {
	var R, S = new(big.Int), new(big.Int)
	for {
		// 产生随机数K
		k := getRandomD()
		P := GBasePoint.Multiply(k)
		e := new(big.Int).SetBytes(hashBytes)
		R = e.Add(e, P.X.Num)
		R = R.Mod(R, sm2Curve.N)
		if R.Sign() == 0 {
			continue
		}
		var rAddK big.Int
		if rAddK.Add(R, k).Cmp(sm2Curve.N) == 0 {
			continue
		}
		// 计算S
		S = S.Add(constant.BigIntOne, priKey.D) // 1+dA
		S = S.ModInverse(S, sm2Curve.N)         // 1+dA相对于N的逆元

		rMulD := new(big.Int).Mul(R, priKey.D) // R*dA
		kSubRMulD := k.Sub(k, rMulD)           // k-R*dA

		S = S.Mul(S, kSubRMulD)  //1+dA相对于N的逆元 * (k-R*dA)
		S = S.Mod(S, sm2Curve.N) // mod n
		if S.Sign() == 0 {
			continue
		}
		break
	}
	// 将R和S转为字节数组
	signBytes = make([]byte, 64)
	copy(signBytes[:32], R.Bytes())
	copy(signBytes[32:], S.Bytes())
	return
}
