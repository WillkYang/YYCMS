package controllers

import (
	"YYCMS/utils/YYLog"
	"bytes"
	"encoding/base64"
	"image/jpeg"
	"strings"
	"sync"
	"time"

	"github.com/Unknwon/com"
	"github.com/hanguofeng/gocaptcha"
)

type AuthCodeController struct {
	LoginController
}

var captcha *gocaptcha.Captcha
var once sync.Once

//获取验证码
func (this *AuthCodeController) GetCaptcha() {
	key, dist := makeCaptcha()
	info := make(map[string]string, 0)
	info["key"] = key
	info["info"] = "data:image/png;base64," + strings.Replace(string(dist), "\u0000", "", -1)
	this.AjaxMsg(info, 0, "", "")
}

func makeCaptcha() (string, []byte) {
	createCaptcha()
	key, err1 := captcha.GetKey(4)
	if err1 != nil {
		YYLog.Error(err1.Error())
	}
	img, err := captcha.GetImage(key)
	if err != nil {
		YYLog.Error(err.Error())
	}
	emptyBuff := bytes.NewBuffer(nil)                  //开辟一个新的空buff
	jpeg.Encode(emptyBuff, img, nil)                   //img写入到buff
	dist := make([]byte, 10000)                        //开辟存储空间
	base64.StdEncoding.Encode(dist, emptyBuff.Bytes()) //buff转成base64
	//fmt.Println(string(dist))                          //输出图片base64(type = []byte)
	return key, dist
}

func VerifyCaptcha(key, word string) (bool, string) {
	return captcha.Verify(key, word)
}

//验证验证码
func (this *AuthCodeController) VerifyCaptcha() {
	key := this.MustStr("key")
	word := this.MustStr("word")
	isSuc, reason := captcha.Verify(key, word)
	if isSuc {
		this.AjaxMsg(nil, 0, "", "成功")
	} else {
		this.AjaxMsg(nil, 1, reason, "")
	}
}

func createCaptcha() {
	if captcha == nil {
		YYLog.Debug("创建Captcha")
		once.Do(func() {
			wordDict, captchaConfig, imageConfig, filterConfig, storeConfig := loadConfig()
			wordmgr, err := gocaptcha.CreateWordManagerFromDataFile(wordDict)
			if err != nil {
				YYLog.Error(err)
			}
			captcha, err = gocaptcha.CreateCaptcha(wordmgr, captchaConfig, imageConfig, filterConfig, storeConfig)
			if err != nil {
				YYLog.Error(err)
			}
		})
	}
}
func loadConfig() (string, *gocaptcha.CaptchaConfig, *gocaptcha.ImageConfig, *gocaptcha.FilterConfig, *gocaptcha.StoreConfig) {
	fileNames, _ := com.StatDir(".")
	YYLog.Warning(fileNames)

	data_path := "./static/authcode_data/"

	wordDict := data_path + "en_char"

	captchaConfig := new(gocaptcha.CaptchaConfig)
	captchaConfig.LifeTime = 10 * time.Minute

	imageConfig := new(gocaptcha.ImageConfig)
	imageConfig.FontFiles = []string{data_path + "zpix.ttf"}
	imageConfig.FontSize = 26
	imageConfig.Height = 40
	imageConfig.Width = 120

	filterConfig := new(gocaptcha.FilterConfig)
	filterConfig.Init()
	filterConfig.Filters = []string{"ImageFilterNoiseLine", "ImageFilterNoisePoint", "ImageFilterStrike"}

	var filterConfigGroup *gocaptcha.FilterConfigGroup
	filterConfigGroup = new(gocaptcha.FilterConfigGroup)
	filterConfigGroup.Init()
	filterConfigGroup.SetItem("Num", "1")
	filterConfig.SetGroup("ImageFilterNoiseLine", filterConfigGroup)
	filterConfigGroup = new(gocaptcha.FilterConfigGroup)
	filterConfigGroup.Init()
	filterConfigGroup.SetItem("Num", "1")
	filterConfig.SetGroup("ImageFilterNoisePoint", filterConfigGroup)
	filterConfigGroup = new(gocaptcha.FilterConfigGroup)
	filterConfigGroup.Init()
	filterConfigGroup.SetItem("Num", "1")
	filterConfig.SetGroup("ImageFilterStrike", filterConfigGroup)

	storeConfig := new(gocaptcha.StoreConfig)
	storeConfig.Engine = gocaptcha.STORE_ENGINE_BUILDIN
	storeConfig.GcDivisor = 100
	storeConfig.GcProbability = 1

	return wordDict, captchaConfig, imageConfig, filterConfig, storeConfig
}
