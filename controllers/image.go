package controllers

import (
	cnf "YYCMS/conf"
	m "YYCMS/models"
	"YYCMS/utils"
	"YYCMS/utils/YYLog"
	"errors"
	"os"
	"path"
	"runtime/debug"
	"strconv"
	"strings"
)

type ImageController struct {
	BaseController
}

func (c *ImageController) Upload() {
	var fileInfo UploadFileInfo
	var uploadPath string
	if c.Bool("isPrivate") {
		uploadPath = cnf.UploadImagePath
	} else {
		uploadPath = cnf.UploadPublicImagePath
	}
	fileInfo = Server{}.Upload(c.Ctx.Request, nil, GenRealPathWithDir(uploadPath), false)

	m.CreateOneKvdb(fileInfo.FileName+fileInfo.Ext, fileInfo.OriginName)
	var err error
	err = genImageThumbnail(fileInfo, uploadPath, cnf.TransferThumbImagePath, 600, 400)
	err = genImageThumbnail(fileInfo, uploadPath, cnf.TransferHtmlImagePath, 1040, 640)

	debug.FreeOSMemory()
	if err != nil {
		YYLog.Error("上传图片失败:", err.Error())
		c.finish(errors.New(strconv.Itoa(m.FileExtError)))
	} else {
		dbImagtePath := "/" + strings.Replace(fileInfo.Path, cnf.UploadImagePath[2:], "", -1)
		YYLog.Warning(dbImagtePath)
		c.finish(dbImagtePath)
	}
}

func genImageThumbnail(fileInfo UploadFileInfo, uploadPath, transferPath string, width, height int) error {
	//newFileNamePath := strings.Replace(fileInfo.Path, fileInfo.FileName, fileInfo.FileName + "-" + strconv.Itoa(width) + "*" + strconv.Itoa(height), -1)
	newImagePath := strings.Replace(fileInfo.Path, uploadPath[2:], transferPath, -1)
	os.MkdirAll(path.Dir(newImagePath), os.ModePerm)
	if err := utils.ImageThumbnail(fileInfo.Path, uint(width), uint(height), newImagePath); err != nil {
		YYLog.Error(err)
		return err
	} else {
		return nil
	}
}
