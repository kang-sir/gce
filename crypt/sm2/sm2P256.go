package sm2

import (
	"crypto/elliptic"
	"crypto/rand"
	"gce/constant"
	"gce/crypt"
	"gce/crypt/sm3"
	"gce/util/num"
	"io"
	"math/big"
)

var sm2Curve sm2P256Curve

type sm2P256Curve struct {
	// 曲线参数
	*elliptic.CurveParams
	// 曲线参数a，b
	A, B *big.Int
	// 曲线余因子(h=p/n)
	h *big.Int
}

// 初始化椭圆曲线参数信息（GM定义）
func init() {
	sm2Curve.CurveParams = &elliptic.CurveParams{Name: "SM2P256-Curve"}
	sm2Curve.P = crypt.P
	sm2Curve.N = crypt.N
	sm2Curve.Gx = crypt.Gx
	sm2Curve.Gy = crypt.Gy
	sm2Curve.A = crypt.A
	sm2Curve.B = crypt.B
	sm2Curve.BitSize = 256
	sm2Curve.h = new(big.Int).Div(sm2Curve.P, sm2Curve.N)
	// 初始化常量
	initECFieldElement(sm2Curve)
	initecPoint(sm2Curve)
}

func getRandomD() *big.Int {
	randBytes := make([]byte, sm2Curve.N.BitLen()/8)
	d := new(big.Int)
	for {
		io.ReadFull(rand.Reader, randBytes)
		d.SetBytes(randBytes)
		if d.Cmp(constant.BigIntOne) > 0 && d.Cmp(sm2Curve.N) < 0 {
			break
		}
	}
	return d
}

func GetCurve() sm2P256Curve {
	return sm2Curve
}

func KDF(Z []byte, msgLen int) (K []byte) {
	sm3HashFunc := sm3.New()
	K = make([]byte, msgLen)
	// 循环计算Ha[i]
	ct := uint32(1)
	blockCount := msgLen / sm3HashFunc.BlockSize()
	for i := 0; i < blockCount; i++ {
		tmpData := append(Z, num.Uint32ToBytes(ct)...)
		hashBlock := sm3HashFunc.Sum(tmpData)
		copy(K[i*sm3HashFunc.BlockSize():(i+1)*sm3HashFunc.BlockSize()], hashBlock)
		ct++
	}
	// 处理最后一个Ha[v]
	tmpData := append(Z, num.Uint32ToBytes(ct)...)
	hashBlock := sm3HashFunc.Sum(tmpData)
	copy(K[blockCount*sm3HashFunc.BlockSize():], hashBlock)
	return
}

func XOR(x []byte, y []byte) (z []byte) {
	xLen := len(x)
	yLen := len(y)
	var zLen int
	if xLen <= yLen {
		zLen = xLen
	} else {
		zLen = yLen
	}
	z = make([]byte, zLen)
	for i := 0; i < zLen; i++ {
		z[i] = x[i] ^ y[i]
	}
	return
}

func CreatePoint(xBytes []byte, yBytes []byte) ecPoint {
	x := new(big.Int).SetBytes(xBytes)
	y := new(big.Int).SetBytes(yBytes)

	return ecPoint{
		X: &ecFieldElement{
			Num: x,
			Q:   sm2Curve.P,
		},
		Y: &ecFieldElement{
			Num: y,
			Q:   sm2Curve.P,
		},
	}
}
