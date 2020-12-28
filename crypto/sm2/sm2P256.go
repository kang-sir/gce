package sm2

import (
	"crypto/elliptic"
	"crypto/rand"
	"gce/constant"
	"gce/crypto/sm3"
	"gce/util/num"
	"io"
	"math/big"
)

var sm2Curve sm2P256Curve

var userIdDef = []byte("1234567812345678")

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
	sm2Curve.P, _ = new(big.Int).SetString("FFFFFFFEFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF00000000FFFFFFFFFFFFFFFF", 16)
	sm2Curve.N, _ = new(big.Int).SetString("FFFFFFFEFFFFFFFFFFFFFFFFFFFFFFFF7203DF6B21C6052B53BBF40939D54123", 16)
	sm2Curve.Gx, _ = new(big.Int).SetString("32c4ae2c1f1981195f9904466a39c9948fe30bbff2660be1715a4589334c74c7", 16)
	sm2Curve.Gy, _ = new(big.Int).SetString("bc3736a2f4f6779c59bdcee36b692153d0a9877cc62a474002df32e52139f0a0", 16)
	sm2Curve.A, _ = new(big.Int).SetString("FFFFFFFEFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF00000000FFFFFFFFFFFFFFFC", 16)
	sm2Curve.B, _ = new(big.Int).SetString("28E9FA9E9D9F5E344D5A9E4BCF6509A7F39789F515AB8F92DDBCBD414D940E93", 16)
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

func getZA(pubKey *PublicKey, userId []byte) []byte {
	if userId == nil || len(userId) == 0 {
		userId = userIdDef
	}
	IDBitLen := len(userId) * 8
	entLenBytes := num.Uint16T0Bytes(uint16(IDBitLen))
	var finalOriBytes []byte
	// ENTLa
	finalOriBytes = append(finalOriBytes, entLenBytes...)
	// IDa
	finalOriBytes = append(finalOriBytes, userId...)
	// a,b
	finalOriBytes = append(finalOriBytes, sm2Curve.A.Bytes()...)
	finalOriBytes = append(finalOriBytes, sm2Curve.B.Bytes()...)
	// Gx,Gy
	finalOriBytes = append(finalOriBytes, sm2Curve.Gx.Bytes()...)
	finalOriBytes = append(finalOriBytes, sm2Curve.Gy.Bytes()...)
	// Ax,Ay(防止Ax,Ay的高位为0的情况)
	xZero := make([]byte, 32-len(pubKey.X.Bytes()))
	finalOriBytes = append(finalOriBytes, xZero...)
	finalOriBytes = append(finalOriBytes, pubKey.X.Bytes()...)
	yZero := make([]byte, 32-len(pubKey.Y.Bytes()))
	finalOriBytes = append(finalOriBytes, yZero...)
	finalOriBytes = append(finalOriBytes, pubKey.Y.Bytes()...)
	// 对数据进行Hash，返回ZA
	sm3HashFunc := sm3.New()
	return sm3HashFunc.Sum(finalOriBytes)
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
