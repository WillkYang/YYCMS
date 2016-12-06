package models

//SystemError 系统内部错误
const SystemError = -1

const (
	NoError = iota
	ParamsMissError //1
	ParamsTypeError
	LoginError
	LoginAccountError
	LoginPwdError
	AccountExistError
	OldPwdError
	PwdWeakError
	CookieError
	NoPermissonError //10
	FileExtError
	EmailFormatError
	EmailNoneExistError
	AccountUnuseError
	IDCardNumberError
	PhoneNumberError
	EmailExistError
	NetworkPlatformInfoExistError
	JSONUnMarsalError
	AccountExistAndPleaseRefreshPageError //20
	InfoNotExistError
	CanNotDelCateError
	PwdRepeatError
	CaptchaError
	VideoFormatError
	DataBaseGetError
)

var ErrInfo = map[int]string{
	SystemError:                           "系统内部错误",
	NoError:                               "无错误",
	ParamsMissError:                       "参数不完整",
	ParamsTypeError:                       "参数类型错误",
	LoginError:                            "账号或者密码错误",
	LoginAccountError:                     "账号错误",
	LoginPwdError:                         "密码错误",
	AccountExistError:                     "用户名已经存在",
	OldPwdError:                           "旧密码输入错误",
	PwdWeakError:                          "密码强度太弱,请输入6位以上数字和字母的组合",
	CookieError:                           "登录过期,请重新登录",
	NoPermissonError:                      "用户权限不够",
	FileExtError:                          "文件格式不正确",
	EmailFormatError:                      "邮箱格式不正确",
	EmailNoneExistError:                   "邮箱不存在",
	AccountUnuseError:                     "账号已经被冻结",
	IDCardNumberError:                     "身份证号码格式不正确",
	PhoneNumberError:                      "手机号码格式不正确",
	EmailExistError:                       "邮箱已经存在",
	NetworkPlatformInfoExistError:         "网络平台信息已经存在",
	JSONUnMarsalError:                     "JSON解析错误",
	AccountExistAndPleaseRefreshPageError: "账号已经存在,请刷新页面重试",
	InfoNotExistError:                     "信息不存在",
	CanNotDelCateError:                    "该栏目不允许删除",
	PwdRepeatError:                        "新密码不能和旧密码一致",
	CaptchaError:                          "验证码错误",
	VideoFormatError:                      "视频文件格式错误，目前只支持mp4,flv,mov,rmvb,avi格式文件",
	DataBaseGetError:                            "数据库获取异常",
}
var ErrCode = map[string]int{
	"系统内部错误":                 SystemError,
	"无错误":                    NoError,
	"参数不完整":                  ParamsMissError,
	"参数类型错误":                 ParamsTypeError,
	"账号或者密码错误":               LoginError,
	"账号错误":                   LoginAccountError,
	"密码错误":                   LoginPwdError,
	"用户名已经存在":                AccountExistError,
	"旧密码输入错误":                OldPwdError,
	"密码强度太弱,请输入6位以上数字和字母的组合": PwdWeakError,
	"登录过期,请重新登录":             CookieError,
	"用户权限不够":                 NoPermissonError,
	"文件格式不正确":                FileExtError,
	"邮箱格式不正确":                EmailFormatError,
	"邮箱不存在":                  EmailNoneExistError,
	"账号已经被冻结":                AccountUnuseError,
	"身份证号码格式不正确":             IDCardNumberError,
	"手机号码格式不正确":              PhoneNumberError,
	"邮箱已经存在":                 EmailExistError,
	"网络平台信息已经存在":             NetworkPlatformInfoExistError,
	"JSON解析错误":               JSONUnMarsalError,
	"账号已经存在,请刷新页面重试":         AccountExistAndPleaseRefreshPageError,
	"信息不存在":                  InfoNotExistError,
	"该栏目不允许删除":               CanNotDelCateError,
	"新密码不能和旧密码一致":            PwdRepeatError,
	"验证码错误":                  CaptchaError,
	"视频文件格式错误，目前只支持mp4,flv,mov,rmvb,avi格式文件": VideoFormatError,
	"数据库获取异常":                            DataBaseGetError,
}
