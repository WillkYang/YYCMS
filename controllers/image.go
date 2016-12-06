package controllers

import (
	cnf "YYCMS/conf"
	m "YYCMS/models"
	"strings"
	"os"
	"strconv"
	"path"
	"YYCMS/utils"
	"errors"
	"runtime/debug"
)

type ImageController struct {
	BaseController
}

func (c *ImageController) Upload() {
	fileInfo := Server{}.Upload(c.Ctx.Request, nil, GenRealPathWithDir(cnf.UploadImagePath), false)
	var err error
	err = genImageThumbnail(fileInfo, 750, 600)
	err = genImageThumbnail(fileInfo, 1040, 640)
	debug.FreeOSMemory()

	if err != nil {
		c.finish(errors.New(strconv.Itoa(m.FileExtError)))
	} else {
		c.finish("/"+fileInfo.Path)
	}
}

func genImageThumbnail(fileInfo UploadFileInfo, width, height int) error {
	newFileNamePath := strings.Replace(fileInfo.Path, fileInfo.FileName, fileInfo.FileName + "-" + strconv.Itoa(width) + "*" + strconv.Itoa(height), -1)
	newImagePath := strings.Replace(newFileNamePath, cnf.UploadImagePath[2:], cnf.TransferImagePath, -1)
	os.MkdirAll(path.Dir(newImagePath),os.ModePerm)
	if err := utils.ImageThumbnail(fileInfo.Path,uint(width),uint(height),newImagePath); err != nil {
		return err
	} else {
		return nil
	}
}