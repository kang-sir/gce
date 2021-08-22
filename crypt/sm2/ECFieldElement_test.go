package sm2

import (
	"gce/constant"
	"math/big"
	"reflect"
	"testing"
)

func TestECFieldElement_Square(t *testing.T) {
	GxSquareNum, _ := new(big.Int).SetString("110765054780706081731929645201630847064343024009898629435113713897114513809343", 10)
	type fields struct {
		Num *big.Int
		Q   *big.Int
	}
	tests := []struct {
		name   string
		fields fields
		want   *ECFieldElement
	}{
		// TODO: Add test cases.
		{
			name: "GX Square",
			fields: fields{
				Num: sm2Curve.Gx,
				Q:   sm2Curve.P,
			},
			want: &ECFieldElement{
				Num: GxSquareNum,
				Q:   sm2Curve.P,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			field := &ECFieldElement{
				Num: tt.fields.Num,
				Q:   tt.fields.Q,
			}
			if got := field.Square(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Square() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestECFieldElement_Multiply(t *testing.T) {
	// Gx^2
	GxSquareNum, _ := new(big.Int).SetString("110765054780706081731929645201630847064343024009898629435113713897114513809343", 10)
	// 3Gx^2
	ThreeMulGxSquareNum, _ := new(big.Int).SetString("100710985921405747682948245176850755660528364045847505396498753824764171444031", 10)
	type fields struct {
		Num *big.Int
		Q   *big.Int
	}
	type args struct {
		fieldB *ECFieldElement
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *ECFieldElement
	}{
		// TODO: Add test cases.
		{
			name:   "MultiplyTest(3*Gx^2)",
			fields: fields{Num: constant.BigIntThree, Q: sm2Curve.P},
			args: args{
				fieldB: &ECFieldElement{
					Num: GxSquareNum,
					Q:   sm2Curve.P,
				},
			},
			want: &ECFieldElement{
				Num: ThreeMulGxSquareNum,
				Q:   sm2Curve.P,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			field := &ECFieldElement{
				Num: tt.fields.Num,
				Q:   tt.fields.Q,
			}
			if got := field.Multiply(tt.args.fieldB); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Multiply() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestECFieldElement_Add(t *testing.T) {
	// 3*Gx^2 + A
	ThreeMultyGxSquareAddA, _ := new(big.Int).SetString("100710985921405747682948245176850755660528364045847505396498753824764171444028", 10)
	// 3Gx^2
	ThreeMulGxSquareNum, _ := new(big.Int).SetString("100710985921405747682948245176850755660528364045847505396498753824764171444031", 10)

	type fields struct {
		Num *big.Int
		Q   *big.Int
	}
	type args struct {
		fieldB *ECFieldElement
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *ECFieldElement
	}{
		// TODO: Add test cases.
		{
			name: "AddTest(3*Gx^2 + A)",
			fields: fields{
				Num: ThreeMulGxSquareNum,
				Q:   sm2Curve.P,
			},
			args: args{
				fieldB: ECFieldElementA,
			},
			want: &ECFieldElement{
				Num: ThreeMultyGxSquareAddA,
				Q:   sm2Curve.P,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			field := &ECFieldElement{
				Num: tt.fields.Num,
				Q:   tt.fields.Q,
			}
			if got := field.Add(tt.args.fieldB); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Add() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestECFieldElement_Divide(t *testing.T) {
	// 2y
	twoGy := ECFieldElementGy.Multiply(ECFieldElementTwo)
	// 3*Gx^2 + A
	ThreeMultyGxSquareAddA, _ := new(big.Int).SetString("100710985921405747682948245176850755660528364045847505396498753824764171444028", 10)
	// (3*Gx^2 + A)/2y
	ThreeMultyGxSquareAddADivide2Gy, _ := new(big.Int).SetString("97089084182832288391840691297101644113875969642515613231745221902176049284016", 10)

	type fields struct {
		Num *big.Int
		Q   *big.Int
	}
	type args struct {
		fieldB *ECFieldElement
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *ECFieldElement
	}{
		// TODO: Add test cases.
		{
			name: "DivideTest[(3*Gx^2 + A)/(2Gy)]",
			fields: fields{
				Num: ThreeMultyGxSquareAddA,
				Q:   sm2Curve.P,
			},
			args: args{
				fieldB: twoGy,
			},
			want: &ECFieldElement{
				Num: ThreeMultyGxSquareAddADivide2Gy,
				Q:   sm2Curve.P,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			field := &ECFieldElement{
				Num: tt.fields.Num,
				Q:   tt.fields.Q,
			}
			if got := field.Divide(tt.args.fieldB); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Divide() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestECFieldElement_Subtract(t *testing.T) {
	// (3*Gx^2 + A)/2y
	ThreeMultyGxSquareAddADivide2Gy, _ := new(big.Int).SetString("97089084182832288391840691297101644113875969642515613231745221902176049284016", 10)
	// lambda^2-2x
	lambdaECFieldElement := ECFieldElement{Num: ThreeMultyGxSquareAddADivide2Gy, Q: sm2Curve.P}
	lambdaECFieldElementSquare := lambdaECFieldElement.Square()
	// 2Gx
	twoGx := ECFieldElementGx.Multiply(ECFieldElementTwo)

	G3x, _ := new(big.Int).SetString("39264624226210491828350299246801547995169179683885222662092227069188840275282", 10)
	G3xElement := &ECFieldElement{
		Num: G3x,
		Q:   sm2Curve.P,
	}

	type fields struct {
		Num *big.Int
		Q   *big.Int
	}
	type args struct {
		fieldB *ECFieldElement
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *ECFieldElement
	}{
		// TODO: Add test cases.
		{
			name: "SubtractTest[lambda^2-2Gx]",
			fields: fields{
				Num: lambdaECFieldElementSquare.Num,
				Q:   lambdaECFieldElementSquare.Q,
			},
			args: args{
				fieldB: twoGx,
			},
			want: G3xElement,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			field := &ECFieldElement{
				Num: tt.fields.Num,
				Q:   tt.fields.Q,
			}
			if got := field.Subtract(tt.args.fieldB); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Subtract() = %v, want %v", got, tt.want)
			}
		})
	}
}
