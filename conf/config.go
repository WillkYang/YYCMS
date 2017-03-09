package conf

const (
	AdminAuthCookieKey = "aep3HRwtagqczu79"
	SuperAdminRoleId   = -1

	TopCateId = 0

	DefaultPageSize = 8

	//CachePath 缓存文件地址
	CachePath = "cache/"
	//ModelCachePath 模型缓存地址
	ModelCachePath = "cache/model/"
	//CategoryCachePath 栏目缓存地址
	CategoryCachePath = "cache/category/"

	UploadFileDir = "file/upload"

	TranscodeVideoPath = "./upload/transcode/video/"

	DefaultAdminPassword = "PTpGWn4E"

	//图片上传路径
	UploadImagePath = "./upload/private/uploadimg/"
	//公开图片保存路径
	UploadPublicImagePath = "./upload/public/uploadimg/"
	//缩略图路径
	TransferThumbImagePath = "./transfer/img/thumb/"
	TransferHtmlImagePath  = "./transfer/img/html/"
	TransferImagePath      = "./transfer/img/"
)

//配置
const (
	ModelConfigFilePath = "conf/models.json"
	AllActionsPath      = "conf/actions.json"
)
