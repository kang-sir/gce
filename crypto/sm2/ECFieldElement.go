package sm2

import (
	"gce/constant"
	"math/big"
)

var ECFieldElementThree *ECFieldElement
var ECFieldElementTwo *ECFieldElement
var ECFieldElementA *ECFieldElement
var ECFieldElementGx *ECFieldElement
var ECFieldElementGy *ECFieldElement

type ECFieldElement struct {
	Num, Q *big.Int
}

func initECFieldElement(sm2Curve sm2P256Curve) {
	ECFieldElementThree = &ECFieldElement{Num: constant.BigIntThree, Q: sm2Curve.P}
	ECFieldElementTwo = &ECFieldElement{Num: constant.BigIntTwo, Q: sm2Curve.P}
	ECFieldElementA = &ECFieldElement{Num: sm2Curve.A, Q: sm2Curve.P}
	ECFieldElementGx = &ECFieldElement{Num: sm2Curve.Gx, Q: sm2Curve.P}
	ECFieldElementGy = &ECFieldElement{Num: sm2Curve.Gy, Q: sm2Curve.P}
}

func (field *ECFieldElement) equals(fieldB *ECFieldElement) bool {
	if field.Num.Cmp(fieldB.Num) == 0 {
		return true
	}
	return false
}

func (field *ECFieldElement) Add(fieldB *ECFieldElement) *ECFieldElement {
	addRes := new(big.Int)
	addRes.Add(field.Num, fieldB.Num)
	addRes.Mod(addRes, field.Q)
	return &ECFieldElement{Num: addRes, Q: field.Q}
}

func (field *ECFieldElement) Subtract(fieldB *ECFieldElement) *ECFieldElement {
	subRes := new(big.Int)
	subRes.Sub(field.Num, fieldB.Num)
	subRes.Mod(subRes, field.Q)
	return &ECFieldElement{Num: subRes, Q: field.Q}
}

func (field *ECFieldElement) Multiply(fieldB *ECFieldElement) *ECFieldElement {
	mulRes := new(big.Int)
	mulRes.Mul(field.Num, fieldB.Num)
	mulRes.Mod(mulRes, field.Q)
	return &ECFieldElement{Num: mulRes, Q: field.Q}
}

func (field *ECFieldElement) Divide(fieldB *ECFieldElement) *ECFieldElement {
	bModQInverse := new(big.Int)
	bModQInverse.ModInverse(fieldB.Num, field.Q)
	divRes := new(big.Int)
	divRes.Mul(field.Num, bModQInverse)
	divRes.Mod(divRes, field.Q)
	return &ECFieldElement{Num: divRes, Q: field.Q}
}

func (field *ECFieldElement) Square() *ECFieldElement {
	squareVal := new(big.Int)
	squareVal.Mul(field.Num, field.Num)
	squareVal = squareVal.Mod(squareVal, field.Q)
	return &ECFieldElement{Num: squareVal, Q: field.Q}
}
