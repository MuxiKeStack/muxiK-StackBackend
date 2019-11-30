package service

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/spf13/viper"
	"io"
	"strconv"
	"time"
)

var (
	endPoint        = viper.GetString("oss.endpoint")
	accessKeyId     = viper.GetString("oss.access_key_id")
	accessKeySecret = viper.GetString("oss.access_key_secret")
	bucketName      = viper.GetString("oss.bucket_name")
	maxImageNumber  = viper.GetInt64("oss.max_image_number_per_person")
	domainName      = viper.GetString("domain_name")
)

func UploadImage(id uint32, r io.Reader) (string, error) {
	client, err := oss.New(endPoint, accessKeyId, accessKeySecret)
	if err != nil {
		return "", err
	}
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return "", err
	}
	timeEpochNow := time.Now().Unix()
	objectName := strconv.Itoa(int(id)) + "-" + strconv.Itoa(int(timeEpochNow%maxImageNumber))
	err = bucket.PutObject(objectName, r)
	if err != nil {
		return "", err
	}
	url := domainName + objectName
	return url, nil
}
