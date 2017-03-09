package controllers

import (
	"YYCMS/utils/YYLog"
	"encoding/json"
	"io/ioutil"
	"os"
	"regexp"

	"github.com/agelinazf/egb"
	"github.com/astaxie/beego"
)

type UeditorController struct {
	beego.Controller
}

type UploadimageUE struct {
	url      string
	title    string
	original string
	state    string
}

type List struct {
	Url string `json:"url"`
}

type Listimage struct {
	State string `json:"state"`
	List  []List `json:"list"`
	Start int    `json:"start"`
	Total int    `json:"total"`
}

type ListCatch struct {
	Url    string `json:"url"`
	Source string `json:"source"`
	State  string `json:"state"`
}

type Catchimage struct {
	State string      `json:"state"`
	List  []ListCatch `json:"list"`
}

func (c *UeditorController) ControllerUE() {
	c.EnableRender = false
	op := c.Input().Get("action")
	switch op {
	case "config": //这里是conf/config.json
		file, err := os.Open("./YYCMS_Data/ueditor/config.json")
		if err != nil {
			YYLog.Error(err.Error())
		}
		defer file.Close()
		fd, _ := ioutil.ReadAll(file)
		src := string(fd)
		re, _ := regexp.Compile("\\/\\*[\\S\\s]+?\\*\\/") //参考php的$CONFIG = json_decode(preg_replace("/\/\*[\s\S]+?\*\//", "", file_get_contents("config.json")), true);
		//将php中的正则移植到go中，需要将/ \/\*[\s\S]+?\*\/  /去掉前后的/，然后将\改成2个\\
		//参考//去除所有尖括号内的HTML代码，并换成换行符
		// re, _ = regexp.Compile("\\<[\\S\\s]+?\\>")
		// src = re.ReplaceAllString(src, "\n")
		//当把<和>换成/*和*\时，斜杠/和*之间加双斜杠\\才行。
		src = re.ReplaceAllString(src, "")
		tt := []byte(src)
		var r interface{}
		json.Unmarshal(tt, &r) //这个byte要解码
		c.Data["json"] = r
		c.ServeJSON()
	case "uploadimage", "uploadfile", "uploadvideo":
		//保存上传的图片
		//获取上传的文件，直接可以获取表单名称对应的文件名，不用另外提取
		_, h, err := c.GetFile("upfile")
		if err != nil {
			YYLog.Error(err)
		}
		date := egb.TimeYear() + "/" + egb.TimeMonth() + "/" + egb.TimeDay() + "/"
		onlytime := egb.TimeNowUnix()
		os.MkdirAll("./upload/public/ueditor/"+date, os.ModePerm)
		path := "./upload/public/ueditor/" + date + onlytime + h.Filename
		err = c.SaveToFile("upfile", path) //.Join("attachment", attachment)) //存文件    WaterMark(path)    //给文件加水印
		if err != nil {
			YYLog.Error(err)
		}
		c.Data["json"] = map[string]interface{}{"state": "SUCCESS", "url": "/upload/public/ueditor/" + date + onlytime + h.Filename, "title": onlytime + h.Filename, "original": onlytime + h.Filename}
		c.ServeJSON()
	}
}
