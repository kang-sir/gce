package pad

import (
	"bytes"
	"encoding/binary"
	"errors"
	"gce/util/common"
	"strconv"
)

type pkcs5Padding struct {
	Name string
}

func NewPkcs5Padding() (padding Padding) {
	return &pkcs5Padding{Name: "PKCS5Padding"}
}

// 填充数据
func (padding *pkcs5Padding) PackData(oriData []byte, blockLen int) ([]byte, error) {
	if blockLen == 0 {
		return nil, errors.New("PackData error,wrong blockLen " + strconv.Itoa(blockLen))
	}
	oriDataLen := len(oriData)
	remainBlockLen := oriDataLen % blockLen
	fillNum := blockLen - remainBlockLen

	fillBytes := make([]byte, fillNum)
	for key, _ := range fillBytes {
		fillBytes[key] = uint8(fillNum)
	}
	finalData := common.BytesCombine(oriData, fillBytes)
	return finalData, nil
}

// 去掉填充数据
func (padding *pkcs5Padding) DePackData(packData []byte) ([]byte, error) {
	packDataLen := len(packData)
	if packDataLen == 0 {
		return nil, errors.New("DePackData error,wrong packDataLen " + strconv.Itoa(packDataLen))
	}
	fillByte := packData[packDataLen-1]
	fillBytes := []byte{0x00, 0x00, 0x00, fillByte}
	var fillNum int32
	binary.Read(bytes.NewBuffer(fillBytes), binary.BigEndian, &fillNum)
	oriData := packData[:(packDataLen - int(fillNum))]
	return oriData, nil
}
