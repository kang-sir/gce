package sm4

const blockSize = 16

type SM4 struct {
	keyBytes []byte
}

func NewCipher(keyBytes []byte) *SM4 {
	if len(keyBytes) != blockSize {
		panic("new SM4 Cipher error, key len is not 16")
	}
	return &SM4{keyBytes: keyBytes}
}

func (sm4 *SM4) BlockSize() int {
	return blockSize
}

func (sm4 *SM4) Encrypt(dst, src []byte) {
	if len(src) != blockSize {
		panic("SM4 encrypt error, plain len is not 16")
	}
	cipherBlock := ProcessBlock(src, sm4.keyBytes, true)
	copy(dst, cipherBlock)
}

func (sm4 *SM4) Decrypt(dst, src []byte) {
	if len(src) != blockSize {
		panic("SM4 encrypt error, cipher block len is not 16")
	}
	cipherBlock := ProcessBlock(src, sm4.keyBytes, false)
	copy(dst, cipherBlock)
}
