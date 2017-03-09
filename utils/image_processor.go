package utils

import (
	"YYCMS/utils/YYLog"
	"image/jpeg"
	"image/png"
	"os"
	"path"
	"strings"

	"github.com/nfnt/resize"
)

type ImageProcessor struct {
	Width          int
	Height         int
	ImageProcessor string
}

const (
	HtmlImage    = "HtmlImage"
	ArticleThumb = "ArticleThumb"
)

var ConfigImageProcessor = map[string]ImageProcessor{
	HtmlImage:    ImageProcessor{600, 0, "processImage"}, //0 表示按照宽度等比缩放
	ArticleThumb: ImageProcessor{300, 200, "MustProcessImage"},
}

func MustProcessImage(imageName string) {
	defer func() {
		if r := recover(); r != nil {
			YYLog.Error(r)
		}
	}()

}

func ImageThumbnail(file string, width uint, height uint, to string) error {
	// 打开图片并解码
	file_origin, _ := os.Open(file)
	switch strings.ToLower(path.Ext(file)) {
	case ".jpeg", ".jpg":
		origin, err := jpeg.Decode(file_origin)
		if err != nil {
			YYLog.Error(err)
			//尝试另一种解码
			file_origin.Seek(0, 0)
			origin, err = png.Decode(file_origin)
			if err != nil {
				YYLog.Error(err)
				return err
			}
		}
		defer file_origin.Close()
		canvas := resize.Thumbnail(width, height, origin, resize.Lanczos3)
		file_out, err := os.Create(to)
		if err != nil {
			return err
		}
		defer file_out.Close()
		jpeg.Encode(file_out, canvas, &jpeg.Options{Quality: 80})
	case ".png":
		origin, err := png.Decode(file_origin)
		if err != nil {
			YYLog.Error(err)
			//尝试另一种解码
			file_origin.Seek(0, 0)
			origin, err = jpeg.Decode(file_origin)
			if err != nil {
				YYLog.Error(err)
				return err
			}
		}
		defer file_origin.Close()
		canvas := resize.Thumbnail(width, height, origin, resize.Lanczos3)
		file_out, err := os.Create(to)
		if err != nil {
			return err
		}
		defer file_out.Close()
		png.Encode(file_out, canvas)
	}
	return nil
}
