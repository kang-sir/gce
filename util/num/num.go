package num

func BytesToUint32(b []byte) uint32 {
	_ = b[3] // bounds check hint to compiler; see golang.org/issue/14808
	return uint32(b[0])<<24 | uint32(b[1])<<16 | uint32(b[2])<<8 | uint32(b[3])
}

func Uint32ToBytes(number uint32) (numBytes []byte) {
	numBytes = append(numBytes, byte((number>>24)&0xff))
	numBytes = append(numBytes, byte((number>>16)&0xff))
	numBytes = append(numBytes, byte((number>>8)&0xff))
	numBytes = append(numBytes, byte(number&0xff))
	return
}

func Uint16T0Bytes(number uint16) (numBytes []byte) {
	numBytes = append(numBytes, byte((number>>8)&0xff))
	numBytes = append(numBytes, byte(number&0xff))
	return
}
