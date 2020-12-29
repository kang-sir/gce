package sm4

import (
	"reflect"
	"testing"
)

func TestSM4(t *testing.T) {

	plain := []byte{(byte)(0x01), (byte)(0x23), (byte)(0x45), (byte)(0x67), (byte)(0x89), (byte)(0xab), (byte)(0xcd), (byte)(0xef), (byte)(0xfe), (byte)(0xdc), (byte)(0xba), (byte)(0x98), (byte)(0x76), (byte)(0x54), (byte)(0x32), (byte)(0x10)}
	symKey := plain
	cipherBlock := []byte{(byte)(0x68), (byte)(0x1e), (byte)(0xdf), (byte)(0x34), (byte)(0xd2), (byte)(0x06), (byte)(0x96), (byte)(0x5e), (byte)(0x86), (byte)(0xb3), (byte)(0xe9), (byte)(0x4f), (byte)(0x53), (byte)(0x6e), (byte)(0x42), (byte)(0x46)}

	t.Run("SM4EncryptTest", func(t *testing.T) {
		// SM4加密
		sm4 := NewCipher(symKey)
		dst := make([]byte, 16)
		sm4.Encrypt(dst, plain)

		if !reflect.DeepEqual(cipherBlock, dst) {
			t.Errorf("Encrypt() = %v, want %v", dst, cipherBlock)
		}
	})

	t.Run("SM4DecryptTest", func(t *testing.T) {
		// SM4加密
		sm4 := NewCipher(symKey)
		dst := make([]byte, 16)
		sm4.Decrypt(dst, cipherBlock)

		if !reflect.DeepEqual(plain, dst) {
			t.Errorf("Decrypt() = %v, want %v", dst, plain)
		}
	})

}
