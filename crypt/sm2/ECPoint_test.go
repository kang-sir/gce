package sm2

import (
	"math/big"
	"reflect"
	"testing"
)

func TestECPoint_Twice(t *testing.T) {

	ECFieldElementGx := &ECFieldElement{Num: sm2Curve.Gx, Q: sm2Curve.P}
	ECFieldElementGy := &ECFieldElement{Num: sm2Curve.Gy, Q: sm2Curve.P}

	G2x, _ := new(big.Int).SetString("39264624226210491828350299246801547995169179683885222662092227069188840275282", 10)
	G2y, _ := new(big.Int).SetString("22488263119165336147085100012701851696717994160539376183636150600277949031363", 10)
	G2xField := &ECFieldElement{Num: G2x, Q: sm2Curve.P}
	G2yField := &ECFieldElement{Num: G2y, Q: sm2Curve.P}

	type fields struct {
		X *ECFieldElement
		Y *ECFieldElement
	}
	tests := []struct {
		name   string
		fields fields
		wantP2 *ECPoint
	}{
		// TODO: Add test cases.
		{
			name:   "GBasePointTwice",
			fields: fields{X: ECFieldElementGx, Y: ECFieldElementGy},
			wantP2: &ECPoint{X: G2xField, Y: G2yField},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p1 := &ECPoint{
				X: tt.fields.X,
				Y: tt.fields.Y,
			}
			if gotP2 := p1.Twice(); !reflect.DeepEqual(gotP2, tt.wantP2) {
				t.Errorf("Twice() = %v, want %v", gotP2, tt.wantP2)
			}
		})
	}
}

func Test_computeLambda(t *testing.T) {

	GbasePoint := &ECPoint{
		X: ECFieldElementGx,
		Y: ECFieldElementGy,
	}
	lambdaNum, _ := new(big.Int).SetString("97089084182832288391840691297101644113875969642515613231745221902176049284016", 10)
	wantLambda := &ECFieldElement{
		Num: lambdaNum,
		Q:   sm2Curve.P,
	}

	type args struct {
		p1 *ECPoint
		p2 *ECPoint
	}
	tests := []struct {
		name       string
		args       args
		wantLambda *ECFieldElement
	}{
		// TODO: Add test cases.
		{
			name: "GBasePointTwice",
			args: args{
				p1: GbasePoint,
				p2: GbasePoint,
			},
			wantLambda: wantLambda,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotLambda := computeLambda(tt.args.p1, tt.args.p2); !reflect.DeepEqual(gotLambda, tt.wantLambda) {
				t.Errorf("computeLambda() = %v, want %v", gotLambda, tt.wantLambda)
			}
		})
	}
}
