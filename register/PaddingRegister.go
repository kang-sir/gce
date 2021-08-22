package register

import (
	"errors"
	"gce/crypt/pad"
)

var paddingMap = make(map[string]pad.Padding)

func init() {
	paddingMap["PKCS5Padding"] = pad.NewPkcs5Padding()
}

func PaddingRegister(name string, padding pad.Padding) error {
	if _, ok := paddingMap[name]; ok {
		return errors.New("name=[" + name + "]的Padding已被注册")
	}
	paddingMap[name] = padding
	return nil
}

func GetPadding(name string) (padding pad.Padding, err error) {
	if padding, ok := paddingMap[name]; ok {
		return padding, nil
	}
	panic("name=[" + name + "]的Padding方式尚未注册")
}
