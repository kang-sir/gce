package common

import (
	"bytes"
	"encoding/binary"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
)

//BytesCombine 多个[]byte数组合并成一个[]byte
func BytesCombine(pBytes ...[]byte) []byte {
	return bytes.Join(pBytes, []byte(""))
}

// 将GBK的字节转为UTF8的字符串
func ConvertGBKBytes2UTF8(GBKBytes []byte) ([]byte, error) {
	bytesReader := bytes.NewReader(GBKBytes)
	transReader := transform.NewReader(bytesReader, simplifiedchinese.GBK.NewDecoder())
	utf8Bytes, err := ioutil.ReadAll(transReader)
	if err != nil {
		return nil, err
	}
	return utf8Bytes, nil
}

func XOR(dataBytes []byte) []byte {
	for index, byte := range dataBytes {
		dataBytes[index] = byte ^ 0xff
	}
	return dataBytes
}
func BooleanToByte(in bool) []byte {
	outDataByte := make([]byte, 1)
	if in {
		outDataByte[0] = 1
	} else {
		outDataByte[0] = 0
	}
	return outDataByte
}

func ByteToInt(in []byte) int {
	bytebuff := bytes.NewBuffer(in)
	var data int64
	binary.Read(bytebuff, binary.BigEndian, &data)
	return int(data)
}
