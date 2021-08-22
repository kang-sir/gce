package pad

type Padding interface {
	PackData(oriData []byte, blockLen int) ([]byte, error)

	DePackData(packData []byte) ([]byte, error)
}
