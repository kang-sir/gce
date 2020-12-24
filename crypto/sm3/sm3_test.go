package sm3

import (
	"encoding/base64"
	"reflect"
	"testing"
)

func TestSM3_Sum(t *testing.T) {
	abcHash, _ := base64.StdEncoding.DecodeString("Zsfw9GLu7dnR8tRr3BDk4kFnxIdc8veiKX2gK49LqOA=")

	byte90Hash, _ := base64.StdEncoding.DecodeString("1yeFqHLkHIKxqUklAJXlmjn1n3fMssgpOj9dUIBrayU=")
	data := make([]byte, 90)
	for i := 0; i < 90; i++ {
		data[i] = 'a'
	}

	type fields struct {
		digest       [8]uint32
		msgBitLen    uint64
		unHandledMsg []byte
	}
	type args struct {
		msg []byte
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []byte
	}{
		// TODO: Add test cases.
		{
			name: "SM3Test(abc)",
			fields: fields{
				digest:       [8]uint32{0x7380166f, 0x4914b2b9, 0x172442d7, 0xda8a0600, 0xa96f30bc, 0x163138aa, 0xe38dee4d, 0xb0fb0e4e},
				msgBitLen:    0,
				unHandledMsg: []byte{},
			},
			args: args{
				msg: []byte("abc"),
			},
			want: abcHash,
		},
		{
			name: "SM3Test(90bytes)",
			fields: fields{
				digest:       [8]uint32{0x7380166f, 0x4914b2b9, 0x172442d7, 0xda8a0600, 0xa96f30bc, 0x163138aa, 0xe38dee4d, 0xb0fb0e4e},
				msgBitLen:    0,
				unHandledMsg: []byte{},
			},
			args: args{
				msg: data,
			},
			want: byte90Hash,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sm3 := &SM3{
				digest:       tt.fields.digest,
				msgBitLen:    tt.fields.msgBitLen,
				unHandledMsg: tt.fields.unHandledMsg,
			}
			if got := sm3.Sum(tt.args.msg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Sum() = %v, want %v", got, tt.want)
			}
		})
	}
}
