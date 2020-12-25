package sm3

import (
	"encoding/binary"
	"hash"
	"log"
)

type SM3 struct {
	// 存储已计算的摘要
	digest [8]uint32
	// 数据长度
	msgBitLen uint64
	// 剩余的未处理的字节
	unHandledMsg []byte
}

func New() hash.Hash {
	sm3 := new(SM3)
	sm3.Reset()
	return sm3
}

func (sm3 *SM3) Write(msg []byte) (n int, err error) {
	// 计算写入消息长度，并更新sm3对象数据长度(Bit长度)
	writeMsgLen := len(msg) * 8
	sm3.msgBitLen += uint64(writeMsgLen)
	// 将消息整合，添加之前未处理的消息
	allMsg := append(sm3.unHandledMsg, msg...)
	// 计算块数目
	allMsgLen := len(allMsg)
	remainMsgLen := allMsgLen % sm3.BlockSize()
	needDealMsgLen := allMsgLen - remainMsgLen
	// 更新未处理的消息数据
	sm3.unHandledMsg = allMsg[needDealMsgLen:]
	// 对数据段进行处理
	for i := 0; i < needDealMsgLen; i += sm3.BlockSize() {
		tmpBlock := allMsg[i : i+sm3.BlockSize()]
		// 压缩数据，计算nextIV
		sm3.digest, err = CF(sm3.digest, tmpBlock)
		if err != nil {
			return
		}
	}
	n = needDealMsgLen
	return
}

func (sm3 *SM3) Sum(msg []byte) []byte {
	sm3.Write(msg)
	// 处理未处理的消息
	var err error
	paddingDealMsg := Padding(sm3.unHandledMsg, sm3.msgBitLen)
	sm3.unHandledMsg = []byte{}
	// 进行最终数据的处理
	needDealMsgLen := len(paddingDealMsg)
	// 对数据段进行处理
	for i := 0; i < needDealMsgLen; i += sm3.BlockSize() {
		tmpBlock := paddingDealMsg[i : i+sm3.BlockSize()]
		// 压缩数据，计算nextIV
		sm3.digest, err = CF(sm3.digest, tmpBlock)
		if err != nil {
			log.Fatalf("SM3 error，compute paddingMsg error:%s", err)
			return []byte{}
		}
	}
	// 返回摘要
	finalDigest := make([]byte, sm3.Size())
	for i := 0; i < 8; i++ {
		binary.BigEndian.PutUint32(finalDigest[i*4:], sm3.digest[i])
	}
	return finalDigest
}

func (sm3 *SM3) Reset() {
	// 初始IV值
	sm3.digest[0] = 0x7380166f
	sm3.digest[1] = 0x4914b2b9
	sm3.digest[2] = 0x172442d7
	sm3.digest[3] = 0xda8a0600
	sm3.digest[4] = 0xa96f30bc
	sm3.digest[5] = 0x163138aa
	sm3.digest[6] = 0xe38dee4d
	sm3.digest[7] = 0xb0fb0e4e
	// 初始化数据长度
	sm3.msgBitLen = 0
	// 未处理数据长度为0
	sm3.unHandledMsg = []byte{}
}

func (sm3 *SM3) Size() int {
	return 32
}

func (sm3 *SM3) BlockSize() int {
	return 64
}
