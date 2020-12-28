package sm2

import (
	"bytes"
	"errors"
	"gce/constant"
	"gce/crypto/sm3"
	"math/big"
)

type PublicKey struct {
	X, Y *big.Int
}
type PrivateKey struct {
	PublicKey
	D *big.Int
}

func (pubKey *PublicKey) createPoint() *ecPoint {
	return &ecPoint{
		X: &ecFieldElement{
			Num: pubKey.X,
			Q:   sm2Curve.P,
		},
		Y: &ecFieldElement{
			Num: pubKey.Y,
			Q:   sm2Curve.P,
		},
	}
}

func GenerateKey() (priKey *PrivateKey, err error) {
	// 产生随机数D，作为私钥
	d := getRandomD()
	// 计算 P(x,y)=[k]G
	P := gBasePoint.Multiply(d)
	priKey = &PrivateKey{D: d, PublicKey: PublicKey{X: P.X.Num, Y: P.Y.Num}}
	return
}

func SignHashPkcs1(priKey *PrivateKey, hashBytes []byte) (signBytes []byte, err error) {
	var R, S = new(big.Int), new(big.Int)
	for {
		// 产生随机数K
		k := getRandomD()
		P := gBasePoint.Multiply(k)
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
	// 将R和S转为字节数组(防止R/S高位为0的情况)
	RBytes := R.Bytes()
	RLen := len(RBytes)
	SBytes := S.Bytes()
	SLen := len(SBytes)
	signBytes = make([]byte, 64)
	copy(signBytes[32-RLen:32], RBytes)
	copy(signBytes[64-SLen:], SBytes)
	return
}

func VerifySignByHash(pubKey *PublicKey, hashBytes []byte, signBytes []byte) (pass bool, err error) {
	R := new(big.Int).SetBytes(signBytes[:32])
	S := new(big.Int).SetBytes(signBytes[32:])
	// 判断签名R、S的范围是否正确，否则验签失败
	if R.Cmp(constant.BigIntOne) < 0 || R.Cmp(sm2Curve.N) >= 0 ||
		S.Cmp(constant.BigIntOne) < 0 || S.Cmp(sm2Curve.N) >= 0 {
		pass = false
		return
	}
	// 如果t = [(R+S) mod n] == 0,验签失败
	RaddS := new(big.Int).Add(R, S)
	if RaddS.Mod(RaddS, sm2Curve.N).Cmp(constant.BigIntZero) == 0 {
		pass = false
		return
	}
	// 计算(x,y) = [S]G + [t]Pa
	Pa := pubKey.createPoint()
	Point := gBasePoint.Multiply(S).Add(Pa.Multiply(RaddS))
	// 计算R'和R是否相等
	EFromHash := new(big.Int).SetBytes(hashBytes)
	Rget := EFromHash.Add(EFromHash, Point.X.Num) //R'=(e+x) mod n
	Rget = Rget.Mod(Rget, sm2Curve.N)
	if Rget.Cmp(R) == 0 {
		pass = true
	}
	return
}

func Encrypt(pubKey *PublicKey, plainData []byte) (cipherData []byte, err error) {
	// 生成随机数k
	k := getRandomD()
	// 计算C1(x1,y1)
	c1Point := gBasePoint.Multiply(k)
	SPoint := pubKey.createPoint().Multiply(sm2Curve.h)
	if SPoint.IsInfinity() {
		err = errors.New("compute S=[h]Pb result is zero")
		return
	}
	C1 := append(c1Point.X.GetNumBytes32(), c1Point.Y.GetNumBytes32()...)
	// 计算C2(x2,y2)
	c2Point := pubKey.createPoint().Multiply(k)
	x2Bytes := c2Point.X.GetNumBytes32()
	y2Bytes := c2Point.Y.GetNumBytes32()
	// KDF密钥派生函数，计算C2
	Z := append(x2Bytes, y2Bytes...)
	t := KDF(Z, len(plainData))
	C2 := XOR(plainData, t)
	// 计算C3
	sm3Hash := sm3.New()
	x2My2 := append(x2Bytes, plainData...)
	x2My2 = append(x2My2, y2Bytes...)
	C3 := sm3Hash.Sum(x2My2)
	// 输出密文C1||C3||C2
	cipherData = append(cipherData, C1...)
	cipherData = append(cipherData, C3...)
	cipherData = append(cipherData, C2...)
	return
}

//Decrypt SM2私钥解密接口
func Decrypt(priKey *PrivateKey, cipherData []byte) (plainData []byte, err error) {
	// 明文数据长度
	plainDataLen := len(cipherData) - 96
	if plainDataLen <= 0 {
		err = errors.New("Decrypt error, cipherData len is less than 96 ")
		return
	}
	// 获取C1, C2, C3
	C1x, C1y, C3, C2 := make([]byte, 32), make([]byte, 32), make([]byte, 32), make([]byte, plainDataLen)
	copy(C1x, cipherData[0:32])
	copy(C1y, cipherData[32:64])
	copy(C3, cipherData[64:96])
	copy(C2, cipherData[96:])
	// 判断C1是否在椭圆曲线上
	c1Point := CreatePoint(C1x, C1y)
	if !c1Point.IsOnCurve() {
		err = errors.New("Decrypt error, C1 is not on sm2Curve ")
		return
	}
	// 使用余因子求S
	S := c1Point.Multiply(sm2Curve.h)
	if S.IsInfinity() {
		err = errors.New("Decrypt error, S=[h]C1 is zero ")
		return
	}
	// (x2,y2)=[dB]C1
	C2Point := c1Point.Multiply(priKey.D)
	x2Bytes := C2Point.X.GetNumBytes32()
	y2Bytes := C2Point.Y.GetNumBytes32()
	Z := append(x2Bytes, y2Bytes...)
	t := KDF(Z, plainDataLen)
	plainData = XOR(t, C2)

	// 计算明文摘要是否正确
	oriHash := append(x2Bytes, plainData...)
	oriHash = append(oriHash, y2Bytes...)
	hashRes := sm3.New().Sum(oriHash)
	if !bytes.Equal(hashRes, C3) {
		err = errors.New("Decrypt error, hash is not matched ")
	}
	return
}
