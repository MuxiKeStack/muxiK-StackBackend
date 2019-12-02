package service

import (
	"context"
	"errors"
	"github.com/qiniu/api.v7/v7/auth/qbox"
	"github.com/qiniu/api.v7/v7/storage"
	"github.com/spf13/viper"
	"io"
	"strings"

	"strconv"
	"time"
)

var (
	accessKey    string
	secretKey    string
	bucketName   string
	domainName   string
	upToken      string
	setTimeEpoch int64
	typeMap      map[string]bool
)

func getType(filename string) (string, error) {
	i := strings.LastIndex(filename, ".")
	fileType := filename[i:]
	if !typeMap[fileType] {
		return "", errors.New("the file type is not allowed")
	}
	return fileType, nil
}

func getToken() {
	var maxInt uint64 = 1 << 32
	putPolicy := storage.PutPolicy{
		Scope:   bucketName,
		Expires: maxInt,
	}
	mac := qbox.NewMac(accessKey, secretKey)
	upToken = putPolicy.UploadToken(mac)
}

func getObjectName(filename string, id uint32) (string, error) {
	if setTimeEpoch == 0 {
		setTimeEpoch = time.Now().Unix()
	}
	fileType, err := getType(filename)
	if err != nil {
		return "", err
	}
	timeEpochNow := time.Now().Unix()
	objectName := strconv.FormatUint(uint64(id), 10) + "-" + strconv.FormatInt(timeEpochNow%setTimeEpoch, 10) + fileType
	return objectName, nil
}

func UploadImage(filename string, id uint32, r io.ReaderAt, dataLen int64) (string, error) {
	// 先初始化一些信息
	accessKey = viper.GetString("oss.access_key")
	secretKey = viper.GetString("oss.secret_key")
	bucketName = viper.GetString("oss.bucket_name")
	domainName = viper.GetString("oss.domain_name")
	typeMap = map[string]bool{".jpg": true, ".png": true, ".bmp": true, "jpeg": true, "gif": true}

	if upToken == "" {
		getToken()
	}

	objectName, err := getObjectName(filename, id)
	if err != nil {
		return "", err
	}

	// 下面是七牛云的oss所需信息，objectName对应key是文件上传路径
	cfg := storage.Config{Zone: &storage.ZoneHuanan, UseHTTPS: false, UseCdnDomains: true}
	formUploader := storage.NewResumeUploader(&cfg)
	ret := storage.PutRet{}
	putExtra := storage.RputExtra{Params: map[string]string{"x:name": "STACK"}}
	err = formUploader.Put(context.Background(), &ret, upToken, objectName, r, dataLen, &putExtra)
	//err = formUploader.PutFile(context.Background(), &ret, upToken, objectName, "/home/bowser/Pictures/maogai/1.jpg", &putExtra)
	if err != nil {
		return "", err
	}
	url := domainName + objectName
	return url, nil
}
