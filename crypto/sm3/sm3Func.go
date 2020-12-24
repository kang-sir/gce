package sm3

import (
	"encoding/binary"
	"errors"
)

// SM3布尔函数
func FF(index int, x, y, z uint32) uint32 {
	if index >= 0 && index <= 15 {
		return x ^ y ^ z
	} else {
		return (x & y) | (x & z) | (y & z)
	}
}

func GG(index int, x, y, z uint32) uint32 {
	if index >= 0 && index <= 15 {
		return x ^ y ^ z
	} else {
		return (x & y) | (^x & z)
	}
}

// 获取常量T
func T(index int) uint32 {
	if index < 16 {
		return 0x79cc4519
	} else {
		return 0x7a879d8a
	}
}

// SM3置换函数P0
func leftRotate(x uint32, i uint32) uint32 {
	return x<<(i%32) | x>>(32-i%32)
}

func P0(x uint32) uint32 {
	return x ^ leftRotate(x, 9) ^ leftRotate(x, 17)
}

func P1(x uint32) uint32 {
	return x ^ leftRotate(x, 15) ^ leftRotate(x, 23)
}

// SM3压缩函数
func CF(IV [8]uint32, blockBytes []byte) (nexIv [8]uint32, err error) {
	blockLen := len(blockBytes)
	if blockLen != 64 {
		err = errors.New("SM3 error, CF param blockBytes length not 64")
		return
	}
	// 开始数据压缩过程
	wj, wjExt, tmpErr := MsgExtend(blockBytes)
	if tmpErr != nil {
		err = tmpErr
		return
	}

	A, B, C, D, E, F, G, H := IV[0], IV[1], IV[2], IV[3], IV[4], IV[5], IV[6], IV[7]
	for i := 0; i < 64; i++ {
		SS1 := leftRotate(leftRotate(A, 12)+E+leftRotate(T(i), uint32(i)), 7)
		SS2 := SS1 ^ leftRotate(A, 12)
		TT1 := FF(i, A, B, C) + D + SS2 + wjExt[i]
		TT2 := GG(i, E, F, G) + H + SS1 + wj[i]
		D = C
		C = leftRotate(B, 9)
		B = A
		A = TT1
		H = G
		G = leftRotate(F, 19)
		F = E
		E = P0(TT2)
	}
	nexIv[0] = A ^ IV[0]
	nexIv[1] = B ^ IV[1]
	nexIv[2] = C ^ IV[2]
	nexIv[3] = D ^ IV[3]
	nexIv[4] = E ^ IV[4]
	nexIv[5] = F ^ IV[5]
	nexIv[6] = G ^ IV[6]
	nexIv[7] = H ^ IV[7]
	return
}

// SM3的消息扩展函数
func MsgExtend(blockBytes []byte) (wj [68]uint32, wjExt [64]uint32, err error) {

	// 计算Wj
	for i := 0; i < 16; i++ {
		wj[i] = binary.BigEndian.Uint32(blockBytes[i*4 : (i+1)*4])
	}
	for i := 16; i < 68; i++ {
		p1Param := wj[i-16] ^ wj[i-9] ^ leftRotate(wj[i-3], 15)
		wj[i] = P1(p1Param) ^ leftRotate(wj[i-13], 7) ^ wj[i-6]
	}
	// 计算Wj'
	for i := 0; i < 64; i++ {
		wjExt[i] = wj[i] ^ wj[i+4]
	}
	return
}

// SM3消息补位函数，SM3补位规则 msg + 1000....000 + ([8]byte=msgLen) = 64byte的整数倍
func Padding(msg []byte, msgTotalBitLen uint64) (paddingDealMsg []byte) {
	// 先填充1字节的10000000
	paddingMsg := append(msg, 0x80)
	// 继续填充0,直到paddingMsgLen%64=56
	currentLen := len(paddingMsg) % 64
	paddingZeroLen := 0
	if currentLen <= 56 {
		paddingZeroLen = 56 - currentLen
	} else {
		paddingZeroLen = 64 - (currentLen - 56)
	}
	paddingZeroBytes := make([]byte, paddingZeroLen)
	msgTotalLenBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(msgTotalLenBytes, msgTotalBitLen)
	paddingDealMsg = append(paddingMsg, append(paddingZeroBytes, msgTotalLenBytes...)...)
	return
}
