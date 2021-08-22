package sm2

import (
	"gce/constant"
	"math/big"
)

type ecPoint struct {
	X, Y *ecFieldElement
}

// SM2的原点O
var InfinityPoint = ecPoint{X: nil, Y: nil}

// SM2曲线的基点G
var gBasePoint *ecPoint

func initecPoint(sm2Curve sm2P256Curve) {
	gBasePoint = &ecPoint{X: ecFieldElementGx, Y: ecFieldElementGy}
}

func (p1 *ecPoint) IsInfinity() bool {
	if p1.X == nil || p1.Y == nil || p1.X.Num.Cmp(constant.BigIntZero) == 0 || p1.Y.Num.Cmp(constant.BigIntZero) == 0 {
		return true
	}
	return false
}

func (p1 *ecPoint) Twice() (p2 *ecPoint) {
	p2 = p1.Add(p1)
	return
}

func (p1 *ecPoint) Add(p2 *ecPoint) (p3 *ecPoint) {
	if p1.IsInfinity() {
		p3 = p2
		return
	}
	if p2.IsInfinity() {
		p3 = p1
		return
	}
	p3 = new(ecPoint)
	// x3 = lambda^2-x2-x1   y3 = lambda(x1-x3)-y1
	lambda := computeLambda(p1, p2)
	p3.X = lambda.Square().Subtract(p2.X).Subtract(p1.X)
	p3.Y = lambda.Multiply(p1.X.Subtract(p3.X)).Subtract(p1.Y)
	return
}

func (p1 *ecPoint) Multiply(k *big.Int) (p3 *ecPoint) {
	pointArr := []*ecPoint{&InfinityPoint, p1}
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

func (p1 *ecPoint) IsOnCurve() bool {
	// y^2 = x^3 + ax + b
	ySquare := p1.Y.Square()
	xSquareMulX := p1.X.Square().Multiply(p1.X) //x^3
	res := xSquareMulX.Add(ecFieldElementA.Multiply(p1.X)).Add(ecFieldElementB)
	return ySquare.equals(res)
}

// 获取(x1,y1)(x2,y2)的斜率
func computeLambda(p1 *ecPoint, p2 *ecPoint) (lambda *ecFieldElement) {
	if p1.X.equals(p2.X) && p1.Y.equals(p2.Y) {
		// p1 == p2  lambda = (3x^2+a)/2y
		lambda = ecFieldElementThree.Multiply(p1.X.Square()).Add(ecFieldElementA).Divide(p1.Y.Multiply(ecFieldElementTwo))
	} else {
		// p1 != p2  lambda = (y2-y1)/(x2-x1)
		lambda = p2.Y.Subtract(p1.Y).Divide(p2.X.Subtract(p1.X))
	}
	return
}
