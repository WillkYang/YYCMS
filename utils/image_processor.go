package utils

import (
	"github.com/astaxie/beego"
	"os"
	"path"
	"image/jpeg"
	"github.com/nfnt/resize"
	"image/png"
)

type ImageProcessor struct {
	Width int
	Height int
	ImageProcessor string
}

const (
	HtmlImage = "HtmlImage"
	ArticleThumb = "ArticleThumb"
)

var ConfigImageProcessor = map[string]ImageProcessor{
	HtmlImage: ImageProcessor{600, 0, "processImage"},	//0 表示按照宽度等比缩放
	ArticleThumb: ImageProcessor{300, 200, "MustProcessImage"},
}

func MustProcessImage(imageName string) {
	defer func() {
		if r := recover(); r != nil {
			beego.Error(r)
		}
	}()


}

func ImageThumbnail(file string, width uint, height uint, to string) error {
	// 打开图片并解码
	file_origin, _ := os.Open(file)
	switch path.Ext(file) {
	case ".jpeg", ".jpg":
		origin, err := jpeg.Decode(file_origin)
		if err != nil {
			return err
		}
		defer file_origin.Close()
		canvas := resize.Thumbnail(width, height, origin, resize.Lanczos3)
		file_out, err := os.Create(to)
		if err != nil {
			return err
		}
		defer file_out.Close()
		jpeg.Encode(file_out, canvas, &jpeg.Options{80})
	case ".png":
		origin, err := png.Decode(file_origin)
		if err != nil {
			beego.Error(err)
			return err
		}
		defer file_origin.Close()
		canvas := resize.Thumbnail(width, height, origin, resize.Lanczos3)
		file_out, err := os.Create(to)
		if err != nil {
			return err
		}
		defer file_out.Close()
		png.Encode(file_out, canvas,)
	}
	return nil
}
