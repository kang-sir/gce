package num

func BytesToUint32(b []byte) uint32 {
	_ = b[3] // bounds check hint to compiler; see golang.org/issue/14808
	return uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24
}

func Uint16T0Bytes(number uint16) (numBytes []byte) {
	numBytes = append(numBytes, byte((number>>8)&0xff))
	numBytes = append(numBytes, byte(number&0xff))
	return
}
