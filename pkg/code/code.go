package code

// ErrorCode 错误码
type ErrorCode struct {
	HttpCode int    // HTTP 状态码
	BizCode  int    // 业务错误码
	Message  string // 默认提示信息
}

// ========== 通用错误 ==========
var (
	Success       = ErrorCode{200, 0, "success"}
	BadRequest    = ErrorCode{200, 1001, "参数错误"}
	Unauthorized  = ErrorCode{200, 1002, "请先登录"}
	Forbidden     = ErrorCode{200, 1003, "无权限"}
	NotFound      = ErrorCode{200, 1004, "资源不存在"}
	InternalError = ErrorCode{200, 1005, "服务器内部错误"}
)

// ========== 用户模块 ==========
var (
	ErrUserExist    = ErrorCode{200, 2001, "用户已存在"}
	ErrUserNotFound = ErrorCode{200, 2002, "用户不存在"}
	ErrPassword     = ErrorCode{200, 2003, "密码错误"}
	EmailRepeat     = ErrorCode{200, 2004, "邮箱已经存在"}
)

// ========== 文章模块 ==========
var (
	ErrArticleNotFound = ErrorCode{200, 3001, "文章不存在"}
)

// ========== 评论模块 ==========
var (
	ErrCommentNotFound = ErrorCode{200, 4001, "评论不存在"}
)

// ========== 上传模块 ==========
var (
	ErrFileType     = ErrorCode{200, 5001, "不支持的文件类型"}
	ErrFileTooLarge = ErrorCode{200, 5002, "文件大小超过限制"}
)

// ========== 权限 ==========
var (
	ErrNoPermission = ErrorCode{200, 6001, "无操作权限"}
)
