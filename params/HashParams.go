package params

import "gce/constant"

type HashAlg struct {
	// 摘要算法名称
	AlgName constant.HashAlgName

	// 摘要算法参数
	AlgParam map[string]interface{}
}
