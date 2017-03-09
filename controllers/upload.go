package controllers

import (
	cnf "YYCMS/conf"
	"YYCMS/utils/YYLog"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Unknwon/com"
	"github.com/agelinazf/egb"
)

type Server struct {
	fileDir string
}

type UploadFileInfo struct {
	FileName   string //文件名		eg - "3lk4jlk34jk2jk5434524235"
	Ext        string //文件后缀 		eg: "mp4"
	Path       string //文件保存路径  eg: "/x/xx/xxx.mp4"
	OriginName string //原文件名 		eg: "hehe"
}

// 文件上传接口
// allowExt 若不判断文件后缀,则传nil - ".mp4", ".mov", ...
func (cc Server) Upload(request *http.Request, allowExt []string, uploadDir string, isNeedTransCode bool) UploadFileInfo {
	defer func() {
		if r := recover(); r != nil {
			YYLog.Error(r)
		}
	}()

	//1.判断请求
	//c.MustPost()

	_, fileHeader, err := request.FormFile("file")
	handleErr(err)

	fileName := fileHeader.Filename

	//检查文件名
	if !strings.Contains(fileName, ".") {
		//二进制数据
		mime := fileHeader.Header.Get("Content-Type")
		if mime != "" {
			if ext := cnf.ContentTypeExt[mime]; ext != "" {
				fileName = fileName + "." + ext
			}
		}
	}

	fileExt := filepath.Ext(fileName)

	var newFilePath string
	var chunksUploadDir string

	//2.后缀检查
	if allowExt != nil && len(allowExt) > 0 && !checkFileExt(fileExt, allowExt) {
		panic(fileName + "不符合的文件类型")
	}

	inputFile, err := fileHeader.Open()
	defer inputFile.Close()
	handleErr(err)

	newFileName := Md5(fileName)
	if len(uploadDir) <= 0 {
		uploadDir = GenRealPathWithDir(cnf.UploadFileDir)
	}
	YYLog.Debug(uploadDir)
	os.Mkdir(path.Dir(uploadDir), os.ModePerm)

	//3.根据上传策略修改路径和文件名
	if len(request.FormValue("chunks")) > 0 {
		chunksUploadDir = fmt.Sprintf("%s%s/", uploadDir, newFileName)
		_, err := os.Stat(chunksUploadDir)
		if err != nil && !os.IsExist(err) {
			os.MkdirAll(chunksUploadDir, os.ModePerm)
		}
		chunkFileName := fmt.Sprintf("%s-%s", newFileName, request.FormValue("chunk"))
		newFilePath = filepath.Join(chunksUploadDir, chunkFileName+fileExt)
	} else {
		//由文件生成md5值作为文件名
		md5h := md5.New()
		_, err := io.Copy(md5h, inputFile)
		handleErr(err)
		newFileName = fmt.Sprintf("%x%d", md5h.Sum(nil), time.Now().Unix())
		newFilePath = filepath.Join(uploadDir, newFileName+fileExt)
	}

	//4.文件上传
	outputFile, err := os.OpenFile(newFilePath, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	//outputFile,err := os.Create(newFilePath,)
	defer outputFile.Close()
	handleErr(err)

	inputFile.Seek(0, 0)
	_, err = io.Copy(outputFile, inputFile)
	handleErr(err)

	//5.根据上传策略修改文件分片计数,合并
	if len(request.FormValue("chunks")) > 0 {
		chunks, _ := strconv.Atoi(request.FormValue("chunks"))
		fileChunks := addFileCount(chunksUploadDir + "chunks")
		if chunks == fileChunks {
			//合并文件
			newFilePath = MergeFile(chunksUploadDir, uploadDir)
			//转码视频
			saveTranscodeFilePath := checkTranscodeDir() + newFileName + ".mp4"
			go TranscodeToMp4(newFilePath, saveTranscodeFilePath)
			YYLog.Info("转码完成:", saveTranscodeFilePath)
			newFilePath = saveTranscodeFilePath
		}
	}

	return UploadFileInfo{
		Ext:        fileExt,
		OriginName: fileName,
		FileName:   newFileName,
		Path:       newFilePath,
	}
}

/*
	传入视频分片目录路径和输出文件目录
	输出合并后的文件路径
	同时删除源文件
*/
func MergeFile(mergeDir, outputDir string) string {

	fileNames, _ := com.StatDir(mergeDir)
	newfilename := filepath.Base(mergeDir)
	newfilePath := outputDir + newfilename + filepath.Ext(fileNames[0])

	os.RemoveAll(newfilePath)
	outputfile, _ := os.OpenFile(newfilePath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0777)
	defer outputfile.Close()

	if len(fileNames) > 1 {
		for i := 0; i < len(fileNames)-1; i++ {
			inputfile, err := os.OpenFile(fmt.Sprintf("%s%s-%d%s", mergeDir, filepath.Base(mergeDir), i, filepath.Ext(fileNames[0])), os.O_RDONLY, 0777)
			handleErr(err)
			defer inputfile.Close()
			_, err = io.Copy(outputfile, inputfile)
			handleErr(err)
		}
	}
	////由文件生成md5值作为文件名
	//md5h := md5.New()
	//_,err := io.Copy(md5h, outputfile)
	//handleErr(err)

	//newFileFullName := fmt.Sprintf("%x%d%s", md5h.Sum(nil), time.Now().Unix(), filepath.Ext(fileNames[0]))
	//err = os.Rename(newfilePath,outputDir + newFileFullName)
	//handleErr(err)

	os.RemoveAll(mergeDir)

	return newfilePath
}

var mutex sync.Mutex

//添加文件已到达分片数
func addFileCount(path string) int {

	mutex.Lock()
	defer mutex.Unlock()

	count := 0

	fi := &os.File{}
	defer fi.Close()

	_, err := os.Stat(path)
	if err != nil && !os.IsExist(err) {
		fi, err = os.Create(path)
	} else {
		fi, err = os.Open(path)
		fd, err := ioutil.ReadAll(fi)
		handleErr(err)

		count, err = strconv.Atoi(string(fd))
		handleErr(err)
	}

	count++
	err = ioutil.WriteFile(path, []byte(strconv.Itoa(count)), os.ModePerm)
	handleErr(err)
	return count
}

func checkFileExt(fileExt string, allowExt []string) bool {
	for _, ext := range allowExt {
		if fileExt == ext {
			return true
		}
	}
	return false
}

func handleErr(err error) {
	if err != nil {
		YYLog.Error(err)
		panic(err)
	}
}

// returns filename from path.
// Given "xxx.mp4" it returns "xxx".
func fileNamePrefix(fileName string) string {
	strArr := strings.Split(fileName, ".")
	strLength := len(strArr)
	if strLength >= 2 {
		return strArr[len(strArr)-2]
	}
	return fileName
}

//在目录后加上年月
func GenRealPathWithDir(dir string) string {
	now := time.Now()
	newDir := fmt.Sprintf("%s/%d/%d/%d/", dir, now.Year(), now.Month(), now.Day())
	_, err := os.Stat(newDir)
	if err != nil && !os.IsExist(err) {
		os.MkdirAll(newDir, os.ModePerm)
	}
	return newDir
}

func Md5(value string) string {
	h := md5.New()
	h.Write([]byte(value))
	return fmt.Sprintf("%s", hex.EncodeToString(h.Sum(nil)))
}

func TranscodeToMp4(src, dst string) {
	YYLog.Warning(src)
	YYLog.Warning(dst)
	command := fmt.Sprintf("ffmpeg -y -i %s %s %s", src, "-vcodec h264 -strict -2", dst)
	cmd := exec.Command("/bin/sh", "-c", command)

	_, err := cmd.StdoutPipe()
	if err != nil {
		YYLog.Error("Error:can not obtain the stdout pipe for command No.0: %s", err)
		return
	}
	err = cmd.Start()

	if err != nil {
		YYLog.Error(err)
	}
	//_, err = cmd.Process.Wait()
	err = cmd.Wait()

	if err != nil {
		YYLog.Error("Command finished with error: %v", err)
	}
}

func checkTranscodeDir() string {
	transcodeVideoRealPath := cnf.TranscodeVideoPath + egb.TimeYear() + "/" + egb.TimeMonth() + "/"
	if !com.IsExist(transcodeVideoRealPath) {
		os.MkdirAll(transcodeVideoRealPath, os.ModePerm)
	}
	return transcodeVideoRealPath
}
