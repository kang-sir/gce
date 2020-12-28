package sm2

import (
	"gce/constant"
	"math/big"
)

var ecFieldElementThree *ecFieldElement
var ecFieldElementTwo *ecFieldElement
var ecFieldElementA *ecFieldElement
var ecFieldElementB *ecFieldElement
var ecFieldElementGx *ecFieldElement
var ecFieldElementGy *ecFieldElement

type ecFieldElement struct {
	Num, Q *big.Int
}

func initECFieldElement(sm2Curve sm2P256Curve) {
	ecFieldElementThree = &ecFieldElement{Num: constant.BigIntThree, Q: sm2Curve.P}
	ecFieldElementTwo = &ecFieldElement{Num: constant.BigIntTwo, Q: sm2Curve.P}
	ecFieldElementA = &ecFieldElement{Num: sm2Curve.A, Q: sm2Curve.P}
	ecFieldElementB = &ecFieldElement{Num: sm2Curve.B, Q: sm2Curve.P}
	ecFieldElementGx = &ecFieldElement{Num: sm2Curve.Gx, Q: sm2Curve.P}
	ecFieldElementGy = &ecFieldElement{Num: sm2Curve.Gy, Q: sm2Curve.P}
}

func (field *ecFieldElement) GetNumBytes32() []byte {
	numBytes := field.Num.Bytes()
	numLen := len(numBytes)
	if numLen < 32 {
		numBytes = append(make([]byte, 32-numLen), numBytes...)
	}
	return numBytes
}

func (field *ecFieldElement) equals(fieldB *ecFieldElement) bool {
	if field.Num.Cmp(fieldB.Num) == 0 {
		return true
	}
	return false
}

func (field *ecFieldElement) Add(fieldB *ecFieldElement) *ecFieldElement {
	addRes := new(big.Int)
	addRes.Add(field.Num, fieldB.Num)
	addRes.Mod(addRes, field.Q)
	return &ecFieldElement{Num: addRes, Q: field.Q}
}

func (field *ecFieldElement) Subtract(fieldB *ecFieldElement) *ecFieldElement {
	subRes := new(big.Int)
	subRes.Sub(field.Num, fieldB.Num)
	subRes.Mod(subRes, field.Q)
	return &ecFieldElement{Num: subRes, Q: field.Q}
}

func (field *ecFieldElement) Multiply(fieldB *ecFieldElement) *ecFieldElement {
	mulRes := new(big.Int)
	mulRes.Mul(field.Num, fieldB.Num)
	mulRes.Mod(mulRes, field.Q)
	return &ecFieldElement{Num: mulRes, Q: field.Q}
}

func (field *ecFieldElement) Divide(fieldB *ecFieldElement) *ecFieldElement {
	bModQInverse := new(big.Int)
	bModQInverse.ModInverse(fieldB.Num, field.Q)
	divRes := new(big.Int)
	divRes.Mul(field.Num, bModQInverse)
	divRes.Mod(divRes, field.Q)
	return &ecFieldElement{Num: divRes, Q: field.Q}
}

func (field *ecFieldElement) Square() *ecFieldElement {
	squareVal := new(big.Int)
	squareVal.Mul(field.Num, field.Num)
	squareVal = squareVal.Mod(squareVal, field.Q)
	return &ecFieldElement{Num: squareVal, Q: field.Q}
}
