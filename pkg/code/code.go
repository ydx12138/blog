package code

// ErrorCode 错误码
type ErrorCode struct {
	HttpCode int    // HTTP 状态码
	BizCode  int    // 业务错误码（前端可根据这个码做不同处理）
	Message  string // 默认提示信息
}

// ========== 通用错误 ==========
var (
	Success       = ErrorCode{200, 0, "success"}
	BadRequest    = ErrorCode{200, 1001, "参数错误"}
	Unauthorized  = ErrorCode{200, 1002, "请先登录"}
	Forbidden     = ErrorCode{200, 1003, "无权限"}
	NotFound      = ErrorCode{200, 1004, "资源不存在"}
	InternalError = ErrorCode{200, 1005, "服务器内部错误"} //所有dao层err，全返回这个
)

// ========== 业务错误（按模块） ==========

// 用户模块
var (
	ErrUserExist    = ErrorCode{200, 2001, "用户已存在"}
	ErrUserNotFound = ErrorCode{200, 2002, "用户不存在"}
	ErrPassword     = ErrorCode{200, 2003, "密码错误"}
)

// 订单模块（示例）
var (
	ErrOrderNotFound = ErrorCode{200, 3001, "订单不存在"}
	ErrOrderPaid     = ErrorCode{200, 3002, "订单已支付"}
)
