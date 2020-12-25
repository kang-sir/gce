package num

import (
	"reflect"
	"testing"
)

func TestUint16T0Bytes(t *testing.T) {
	type args struct {
		number uint16
	}
	tests := []struct {
		name         string
		args         args
		wantNumBytes [2]byte
	}{
		// TODO: Add test cases.
		{
			name: "Uint16ToBytesTest(12)",
			args: args{
				number: 12,
			},
			wantNumBytes: [2]byte{0, 12},
		},
		{
			name: "Uint16ToBytesTest(256)",
			args: args{
				number: 256,
			},
			wantNumBytes: [2]byte{1, 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotNumBytes := Uint16T0Bytes(tt.args.number); !reflect.DeepEqual(gotNumBytes, tt.wantNumBytes) {
				t.Errorf("Uint16T0Bytes() = %v, want %v", gotNumBytes, tt.wantNumBytes)
			}
		})
	}
}
