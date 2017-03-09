package conf


// file extension -> content-type map
var ContentTypeExt = map[string]string{
	"application/vnd.ms-excel": "xls",
	"application/vnd.ms-powerpoint": "ppt",
	"application/vnd.ms-word": "doc",
	"application/vnd.oasis.chart": "odc",
	"application/vnd.oasis.database": "odb",
	"application/vnd.oasis.formula": "odf",
	"application/vnd.oasis.image": "odi",
	"application/vnd.oasis.opendocument.graphics": "odg",
	"application/vnd.oasis.opendocument.graphics-template": "otg",
	"application/vnd.oasis.opendocument.presentation": "odp",
	"application/vnd.oasis.opendocument.presentation-template": "otp",
	"application/vnd.oasis.opendocument.text": "odt",
	"application/vnd.oasis.opendocument.text-master": "odm",
	"application/vnd.oasis.opendocument.text-template": "ott",
	"application/vnd.oasis.opendocument.text-web": "oth",
	"application/vnd.oasis.spreadsheet": "ods",
	"application/vnd.oasis.spreadsheet-template": "ots",
	"application/vnd.openxmlformats-officedocument.presentationml.presentation": "pptx",
	"application/vnd.openxmlformats-officedocument.presentationml.slideshow": "ppsx",
	"application/vnd.openxmlformats-officedocument.presentationml.template": "potx",
	"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet": "xlsx",
	"application/vnd.openxmlformats-officedocument.spreadsheetml.template": "xltx",
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": "docx",
	"application/vnd.openxmlformats-officedocument.wordprocessingml.template": "dotx",
	"application/x-ole-storage": "msg",
	"image/gif": "gif",
	"image/jpeg": "jpg",
	"image/png": "png",
	"text/plain": "txt",
}


// file extension -> content-type map
var ExtContentType = map[string]string{
	"doc":  "application/vnd.ms-word",
	"docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
	"dotx": "application/vnd.openxmlformats-officedocument.wordprocessingml.template",
	"xls":  "application/vnd.ms-excel",
	"xlsx": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
	"xltx": "application/vnd.openxmlformats-officedocument.spreadsheetml.template",
	"ppt":  "application/vnd.ms-powerpoint",
	"pptx": "application/vnd.openxmlformats-officedocument.presentationml.presentation",
	"ppsx": "application/vnd.openxmlformats-officedocument.presentationml.slideshow",
	"potx": "application/vnd.openxmlformats-officedocument.presentationml.template",

	"odg": "application/vnd.oasis.opendocument.graphics",
	"otg": "application/vnd.oasis.opendocument.graphics-template",
	"otp": "application/vnd.oasis.opendocument.presentation-template",
	"odp": "application/vnd.oasis.opendocument.presentation",
	"odm": "application/vnd.oasis.opendocument.text-master",
	"odt": "application/vnd.oasis.opendocument.text",
	"oth": "application/vnd.oasis.opendocument.text-web",
	"ott": "application/vnd.oasis.opendocument.text-template",
	"ods": "application/vnd.oasis.spreadsheet",
	"ots": "application/vnd.oasis.spreadsheet-template",
	"odc": "application/vnd.oasis.chart",
	"odf": "application/vnd.oasis.formula",
	"odb": "application/vnd.oasis.database",
	"odi": "application/vnd.oasis.image",

	"txt": "text/plain",
	"msg": "application/x-ole-storage",

	"jpg":  "image/jpeg",
	"jpeg": "image/jpeg",
	"gif":  "image/gif",
	"png":  "image/png",
}