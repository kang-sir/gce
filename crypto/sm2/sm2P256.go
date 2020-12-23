package sm2

import (
	"crypto/elliptic"
	"math/big"
	"sync"
)

var initonce sync.Once
var sm2Curve sm2P256Curve

type sm2P256Curve struct {
	// 曲线参数
	*elliptic.CurveParams
	// 曲线参数a，b
	A, B *big.Int
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
	// 初始化常量
	initECFieldElement(sm2Curve)
	initECPoint(sm2Curve)
}
