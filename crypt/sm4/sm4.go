package sm4

import "crypto/cipher"

const blockSize = 16

type SM4 struct {
	keyBytes []byte
}

func NewCipher(keyBytes []byte) (cipher.Block, error) {
	if len(keyBytes) != blockSize {
		panic("new SM4 Cipher error, key len is not 16")
	}
	return &SM4{keyBytes: keyBytes}, nil
}

func (sm4 *SM4) BlockSize() int {
	return blockSize
}

func (sm4 *SM4) Encrypt(dst, src []byte) {
	if len(src) != blockSize {
		panic("SM4 encrypt error, plain len is not 16")
	}
	cipherBlock := processBlock(src, sm4.keyBytes, true)
	copy(dst, cipherBlock)
}

func (sm4 *SM4) Decrypt(dst, src []byte) {
	if len(src) != blockSize {
		panic("SM4 encrypt error, cipher block len is not 16")
	}
	cipherBlock := processBlock(src, sm4.keyBytes, false)
	copy(dst, cipherBlock)
}
