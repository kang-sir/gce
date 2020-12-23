package sm2

import (
	"math/big"
)

type ECPoint struct {
	X, Y *ECFieldElement
}

// SM2的原点O
var InfinityPoint = ECPoint{X: nil, Y: nil}

// SM2曲线的基点G
var GBasePoint *ECPoint

func initECPoint(sm2Curve sm2P256Curve) {
	GBasePoint = &ECPoint{X: ECFieldElementGx, Y: ECFieldElementGy}
}

func (p1 *ECPoint) IsInfinity() bool {
	if p1.X == nil || p1.Y == nil {
		return true
	}
	return false
}

func (p1 *ECPoint) Twice() (p2 *ECPoint) {
	p2 = p1.Add(p1)
	return
}

func (p1 *ECPoint) Add(p2 *ECPoint) (p3 *ECPoint) {
	if p1.IsInfinity() {
		p3 = p2
		return
	}
	if p2.IsInfinity() {
		p3 = p1
		return
	}
	p3 = new(ECPoint)
	// x3 = lambda^2-x2-x1   y3 = lambda(x1-x3)-y1
	lambda := computeLambda(p1, p2)
	p3.X = lambda.Square().Subtract(p2.X).Subtract(p1.X)
	p3.Y = lambda.Multiply(p1.X.Subtract(p3.X)).Subtract(p1.Y)
	return
}

func (p1 *ECPoint) Multiply(k *big.Int) (p3 *ECPoint) {
	pointArr := []*ECPoint{&InfinityPoint, p1}
	// 获取K的二进制位长度
	bitLen := k.BitLen()
	// 循环判断二进制位进行加和
	for i := 0; i < bitLen; i++ {
		addIndex := k.Bit(i)
		modifyIndex := 1 - addIndex
		twice := pointArr[modifyIndex].Twice()
		pointArr[modifyIndex] = twice.Add(pointArr[addIndex])
	}
	p3 = pointArr[0]
	return
}

// 获取(x1,y1)(x2,y2)的斜率
func computeLambda(p1 *ECPoint, p2 *ECPoint) (lambda *ECFieldElement) {
	if p1.X.equals(p2.X) && p1.Y.equals(p2.Y) {
		// p1 == p2  lambda = (3x^2+a)/2y
		lambda = ECFieldElementThree.Multiply(p1.X.Square()).Add(ECFieldElementA).Divide(p1.Y.Multiply(ECFieldElementTwo))
	} else {
		// p1 != p2  lambda = (y2-y1)/(x2-x1)
		lambda = p2.Y.Subtract(p1.Y).Divide(p2.X.Subtract(p1.X))
	}
	return
}
