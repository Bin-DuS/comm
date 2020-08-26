package ucerror

const (
	Module_General       = 0
	Module_MsgServer     = 1
	Module_UniformServer = 2
)

// 定义错误码
var (
	// 通用部分
	S_Error            = NewCode(Module_General, -1)
	S_Ok               = NewCode(Module_General, 0)
	S_IllegalParam     = NewCode(Module_General, 10100)
	S_PermissionDenied = NewCode(Module_General, 10101)

	// MsgServer使用
	S_MessageCustom1 = NewCode(Module_MsgServer, 1)

	// UniformServer使用
	S_UniformCustom1 = NewCode(Module_UniformServer, 1)
)

// 定义错误码信息
var errMap = map[int]string{
	S_IllegalParam:     "parse input params error",
	S_PermissionDenied: "this user no right access api",
}
